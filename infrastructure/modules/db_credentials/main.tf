resource "aws_secretsmanager_secret" "db_url" {
  name        = "${var.prefix}/database-url"
  description = "PostgreSQL connection string for the CDL API"

  tags = merge(var.tags, { Name = "${var.prefix}/database-url" })
}

resource "aws_secretsmanager_secret_version" "db_url" {
  secret_id     = aws_secretsmanager_secret.db_url.id
  secret_string = var.connection_string

  # Terraform creates the initial value but never overwrites it afterward.
  # Rotate the password directly in AWS Secrets Manager (console or CLI).
  # Updating terraform.tfvars will not change the live secret.
  lifecycle {
    ignore_changes = [secret_string]
  }
}
