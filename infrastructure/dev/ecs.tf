resource aws_ecs_cluster main {
  name = "sentinel-community-cluster"
}

data template_file sentinel_app {
  template = file("./templates/ecs/sentinel_app.json.tpl")

  vars = {
    app_image      = var.app_image
    app_port       = var.app_port
    fargate_cpu    = var.fargate_cpu
    fargate_memory = var.fargate_memory
    aws_region     = var.aws_region
    db_uri         = aws_instance.neo4j.public_dns
    host           = aws_alb.main.dns_name
    username       = var.username
    password       = aws_instance.neo4j.id
  }
  depends_on = [aws_instance.neo4j, aws_alb.main]
}

resource aws_ecs_task_definition app {
  family                   = "sentinel-app-task"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.fargate_cpu
  memory                   = var.fargate_memory
  container_definitions    = data.template_file.sentinel_app.rendered
}

resource aws_ecs_service main {
  name            = "sentinel-community-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = var.app_count
  launch_type     = "FARGATE"

  network_configuration {
    security_groups  = [aws_security_group.ecs_tasks.id]
    subnets          = aws_subnet.private.*.id
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_alb_target_group.app.id
    container_name   = "sentinel-app"
    container_port   = var.app_port
  }

  depends_on = [aws_alb_listener.front_end, aws_iam_role_policy_attachment.ecs_task_execution_role_attachment]
}