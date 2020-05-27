terraform {
    backend s3 {
        bucket = "sentinel-terraform-remote-state"
        key = "terraform.tfstate"
        region = "us-east-1"
    }
}