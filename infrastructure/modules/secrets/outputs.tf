output "db_password" {
  description = "Generated database password"
  value       = random_password.db.result
  sensitive   = true
}
