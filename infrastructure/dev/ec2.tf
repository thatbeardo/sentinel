resource aws_instance neo4j {
  ami                     = "ami-08910cfbd20298e4e"
  instance_type           = "t2.micro"
  subnet_id               = element(aws_subnet.public.*.id, 0)
  vpc_security_group_ids  = [aws_security_group.db.id]

  depends_on              = [aws_security_group.db]
}