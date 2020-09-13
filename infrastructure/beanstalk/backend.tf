terraform {
    backend s3 {
        bucket = "sentinel-terraform-remote-state"
        key = "terraform.staging.tfstate"
        region = "us-east-1"
    }
}