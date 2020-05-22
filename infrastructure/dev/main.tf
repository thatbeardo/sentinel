module "sentinel-vpc" {
  source = "../modules/vpc"
  vpc_cidr = "10.0.0.0/28"
  vpc_tenancy = "default"
}

module "sentinel-subnet-private-1a" {
  source = "../modules/vpc/subnet"
  subnet_cidr_block = "10.0.0.0/12"
  subnet_name = "private-1a"
  vpc_id = module.sentinel-vpc.vpc_id
}