data aws_ami neo4j_image {
  most_recent = true

  filter {
    name   = "name"
    values = ["neo4j-community-bolt-enabled"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["self"] 
}

resource aws_instance neo4j {
  ami             = data.aws_ami.neo4j_image.id
  instance_type   = "t2.micro"
  subnet_id       = element(aws_subnet.public.*.id, 0)
  security_groups = [aws_security_group.db.id]

  depends_on      = [aws_security_group.db]
}