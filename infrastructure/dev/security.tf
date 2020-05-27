# ALB Security Group: Edit to restrict access to the application
resource aws_security_group lb {
  name        = "sentinel-community-load-balancer-security-group"
  description = "controls access to the ALB"
  vpc_id      = aws_vpc.main.id

  ingress {
    protocol    = "tcp"
    from_port   = var.app_port
    to_port     = var.app_port
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Traffic to the ECS cluster should only come from the ALB
resource aws_security_group ecs_tasks {
  name        = "sentinel-community-ecs-tasks-security-group"
  description = "allow inbound access from the ALB only"
  vpc_id      = aws_vpc.main.id

  ingress {
    protocol        = "tcp"
    from_port       = var.app_port
    to_port         = var.app_port
    security_groups = [aws_security_group.lb.id]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Neo4j Security group opens ports 7474 (HTTP), 7687 (Bolt), and 22 (SSH) 
resource aws_security_group db {
  name        = "sentinel-db-tasks-security-group"
  description = "allow inbound access on specific ports only"
  vpc_id      = aws_vpc.main.id

  ingress {
    protocol        = "tcp"
    from_port       = var.db_http_port
    to_port         = var.db_http_port
    cidr_blocks     = ["0.0.0.0/0"]
  }

  ingress {
    protocol        = "tcp"
    from_port       = var.db_bolt_port
    to_port         = var.db_bolt_port
    cidr_blocks     = ["0.0.0.0/0"]
  }

  ingress {
    protocol        = "tcp"
    from_port       = var.db_ssh_port
    to_port         = var.db_ssh_port
    cidr_blocks     = ["0.0.0.0/0"]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

