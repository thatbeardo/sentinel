resource aws_internet_gateway gw {
  vpc_id = var.vpc_id

  tags = {
    Name = "dev-ig"
  }
}

resource aws_eip eip {
  vpc = true
  depends_on = [aws_internet_gateway.gw]
}