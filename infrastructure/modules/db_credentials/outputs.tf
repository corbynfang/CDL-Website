output "secret_arn" {
  description = "ARN of the Secrets Manager secret — ECS task definition references this to inject DATABASE_URL"
  value       = aws_secretsmanager_secret.db_url.arn
}
