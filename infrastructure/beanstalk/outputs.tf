output application_name {
  value = module.elastic_beanstalk_application.elastic_beanstalk_application_name
}

output environment_name {
  value = module.elastic_beanstalk_environment.name
}

output website_bucket_id {
  value = aws_s3_bucket.website_bucket.id
}