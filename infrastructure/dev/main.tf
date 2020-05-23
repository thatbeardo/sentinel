module sentinel-vpc {
  source = "../modules/vpc"
  vpc_cidr = "10.0.0.0/28"
  vpc_tenancy = "default"
}

module sentinel-subnet-private-east-1 {
  source = "../modules/vpc/subnet"
  subnet_cidr_block = "10.0.0.0/28"
  subnet_name = "private-east-1"
  availability_zone = "us-east-1"
  map_public_ip_on_launch = false
  vpc_id = module.sentinel-vpc.vpc_id
}

module sentinel-subnet-private-east-2 {
  source = "../modules/vpc/subnet"
  subnet_cidr_block = "10.0.0.16/28"
  subnet_name = "private-east-2"
  availability_zone = "us-east-2"
  map_public_ip_on_launch = false
  vpc_id = module.sentinel-vpc.vpc_id
}

module sentinel-subnet-public-east-1 {
  source = "../modules/vpc/subnet"
  subnet_cidr_block = "10.0.0.32/12"
  subnet_name = "public-east-1"
  availability_zone = "us-east-1"
  map_public_ip_on_launch = true
  vpc_id = module.sentinel-vpc.vpc_id
}

module sentinel-subnet-public-east-2 {
  source = "../modules/vpc/subnet"
  subnet_cidr_block = "10.0.0.48/12"
  subnet_name = "public-east-2"
  availability_zone = "us-east-2"
  map_public_ip_on_launch = true
  vpc_id = module.sentinel-vpc.vpc_id
}

module sentinel-dev-gateway {
  source = "../modules/internet-gateway"
  vpc_id = module.sentinel-vpc.vpc_id
}

module private_route_table {
  source = "../modules/route-table"
  vpc_id = module.sentinel-vpc.vpc_id
  route_table_name = "private-route-table"
}

module private_route_table_association_subnet_1 {
  source = "../modules/route-table-association"
  subnet_id = module.sentinel-subnet-private-east-1.subnet_id
  route_table_id = module.private_route_table.route_table_id
}

module private_route_table_association_subnet_2 {
  source = "../modules/route-table-association"
  subnet_id = module.sentinel-subnet-private-east-2.subnet_id
  route_table_id = module.private_route_table.route_table_id
}

module nat_gateway {
  source = "../modules/nat"
  nat_name = "nat_gateway_east-1"
  subnet_id = module.sentinel-subnet-public-east-1.subnet_id
}

resource aws_route private_internet_access {
  route_table_id         = module.public_route_table.route_table_id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id = module.nat_gateway.nat_id
}

module public_route_table {
  source = "../modules/route-table"
  vpc_id = module.sentinel-vpc.vpc_id
  route_table_name = "public-route-table"
}

module public_route_table_association_subnet_1 {
  source = "../modules/route-table-association"
  subnet_id = module.sentinel-subnet-public-east-1.subnet_id
  route_table_id = module.public_route_table.route_table_id
}

module public_route_table_association_subnet_2 {
  source = "../modules/route-table-association"
  subnet_id = module.sentinel-subnet-public-east-2.subnet_id
  route_table_id = module.public_route_table.route_table_id
}

resource aws_route public_internet_access {
  route_table_id         = module.public_route_table.route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = module.sentinel-dev-gateway.ig_id
}
