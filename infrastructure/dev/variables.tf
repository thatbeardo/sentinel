# variables.tf

variable aws_region {
  description = "The AWS region things are created in"
  default     = "us-east-1"
}

variable ecs_task_execution_role_name {
  description = "ECS task execution role name"
  default = "sentinelCommunityEcsTaskExecutionRole"
}

variable az_count {
  description = "Number of AZs to cover in a given region"
  default     = "2"
}

variable app_image {
  description = "Docker image to run in the ECS cluster. Passed from the CD workflow"
  default = "dummy-image"
}

variable app_port {
  description = "Port exposed by the docker image to redirect traffic to"
  default     = 8080
}

variable app_count {
  description = "Number of docker containers to run"
  default     = 3
}

variable health_check_path {
  default = "/docs"
}

variable fargate_cpu {
  description = "Fargate instance CPU units to provision (1 vCPU = 1024 CPU units)"
  default     = "1024"
}

variable fargate_memory {
  description = "Fargate instance memory to provision (in MiB)"
  default     = "2048"
}

variable username {
  description = "Database username"
  default     = "neo4j"
}

variable db_bolt_port {
  description = "Ingress port 7474 to establish http - insecure connections remotely"
  default     = 7474
}

variable db_http_port {
  description = "Ingress port 7687 to establish bolt connections"
  default     = 7687
}

variable db_ssh_port {
  description = "Ingress port 22 to establish ssh connections remotely"
  default     = 22
}