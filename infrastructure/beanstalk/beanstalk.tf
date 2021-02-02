provider aws {
    region = var.region
}

module vpc {
    source     = "git::https://github.com/cloudposse/terraform-aws-vpc.git?ref=master"
    namespace  = var.namespace
    stage      = var.stage
    name       = var.name
    cidr_block = "172.16.0.0/16"
}

module subnets {
    source               = "git::https://github.com/cloudposse/terraform-aws-dynamic-subnets.git?ref=master"
    availability_zones   = ["us-east-1a", "us-east-1b"]
    namespace            = var.namespace
    stage                = var.stage
    name                 = var.name
    vpc_id               = module.vpc.vpc_id
    igw_id               = module.vpc.igw_id
    cidr_block           = module.vpc.vpc_cidr_block
    nat_gateway_enabled  = false
    nat_instance_enabled = false
}

module elastic_beanstalk_application {
    source      = "git::https://github.com/cloudposse/terraform-aws-elastic-beanstalk-application.git?ref=master"
    namespace   = var.namespace
    stage       = var.stage
    name        = var.name
    description = "Sentinel Staging"
}

module elastic_beanstalk_environment {
    source                             = "git::https://github.com/cloudposse/terraform-aws-elastic-beanstalk-environment.git?ref=master"
    namespace                          = var.namespace
    stage                              = var.stage
    name                               = var.name
    description                        = "Test elastic_beanstalk_environment "
    region                             = var.region
    availability_zone_selector         = "Any 2"
    dns_zone_id                        = var.dns_zone_id
    elastic_beanstalk_application_name = module.elastic_beanstalk_application.elastic_beanstalk_application_name

    instance_type           = "t2.micro"
    autoscale_min           = 1
    autoscale_max           = 1
    updating_min_in_service = 0
    updating_max_batch      = 1

    loadbalancer_type             = "application"
    vpc_id                        = module.vpc.vpc_id
    loadbalancer_subnets          = module.subnets.public_subnet_ids
    loadbalancer_certificate_arn  = var.acm_certificate_arn
    loadbalancer_ssl_policy       = "ELBSecurityPolicy-2016-08"
    application_subnets           = module.subnets.public_subnet_ids
    allowed_security_groups       = [module.vpc.vpc_default_security_group_id]
    healthcheck_url               = var.healthcheck_url
    associate_public_ip_address   = true

    // https://docs.aws.amazon.com/elasticbeanstalk/latest/platforms/platforms-supported.html
    // https://docs.aws.amazon.com/elasticbeanstalk/latest/platforms/platforms-supported.html#platforms-supported.docker
    solution_stack_name = "64bit Amazon Linux 2 v3.1.1 running Go 1"

    additional_settings = [
      {
        namespace = "aws:elasticbeanstalk:application:environment"
        name      = "DB_URI"
        value     = var.db_uri
      },
      {
        namespace = "aws:elasticbeanstalk:application:environment"
        name      = "USERNAME"
        value     = var.username
      },
      {
        namespace = "aws:elasticbeanstalk:application:environment"
        name      = "PASSWORD"
        value     = var.password
      },
      {
        namespace = "aws:elasticbeanstalk:application:environment"
        name      = "HOST"
        value     = var.host
      },
      {
        namespace = "aws:elasticbeanstalk:application:environment"
        name      = "PORT"
        value     = "5000"
      },
      {
        namespace = "aws:elasticbeanstalk:application:environment"
        name      = "NEW_RELIC_LICENSE"
        value     = var.newrelicLicense
      }
    ]
  }
