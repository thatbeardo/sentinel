[
  {
    "name": "sentinel-app",
    "image": "${app_image}",
    "cpu": ${fargate_cpu},
    "memory": ${fargate_memory},
    "networkMode": "awsvpc",
    "environment": [
      {"name": "DB_URI", "value": "bolt://${db_uri}:7687"},
      {"name": "USERNAME", "value": "${username}"},
      {"name": "PASSWORD", "value": "${password}"},
      {"name": "HOST", "value": "${host}"},
      {"name": "PORT", "value": "${app_port}"}
    ],
    "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/sentinel-app",
          "awslogs-region": "${aws_region}",
          "awslogs-stream-prefix": "ecs"
        }
    },
    "portMappings": [
      {
        "containerPort": ${app_port}
      }
    ]
  }
]