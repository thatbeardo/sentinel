resource aws_eip gw {
  vpc        = true
}

resource aws_nat_gateway gw {
  allocation_id = aws_eip.gw.id
  subnet_id     = var.subnet_id

  tags = {
    Name = var.nat_name
  }
}