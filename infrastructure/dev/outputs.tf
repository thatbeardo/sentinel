output alb_hostname {
  value = aws_alb.main.dns_name
}

output website_bucket_id {
  value = aws_s3_bucket.website_bucket.id
}