resource aws_instance neo4j {
  ami             = "ami-081a719c938b93687"
  instance_type   = "t2.micro"
  subnet_id       = element(aws_subnet.public.*.id, 0)
  security_groups = [aws_security_group.db.id]

  depends_on      = [aws_security_group.db]
}