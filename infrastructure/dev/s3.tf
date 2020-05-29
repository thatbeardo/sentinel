resource aws_s3_bucket website_bucket {
  bucket = "guardmy.app"
  acl    = "public-read"
  policy = file("templates/s3/policy.json")

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}