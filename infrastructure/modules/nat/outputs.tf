output nat_eip {
    value = aws_eip.gw.id
}

output nat_id {
    value = aws_nat_gateway.gw.id
}