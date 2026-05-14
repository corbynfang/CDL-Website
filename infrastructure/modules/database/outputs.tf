output "db_endpoint" {
  description = "RDS hostname (without port)"
  value       = aws_db_instance.main.address
}

output "db_port" {
  description = "RDS port"
  value       = aws_db_instance.main.port
}
