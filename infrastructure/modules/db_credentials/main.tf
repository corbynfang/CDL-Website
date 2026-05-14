locals {
  database_url = "postgresql://${var.db_username}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=require"
}

resource "aws_secretsmanager_secret" "db_url" {
  name        = "${var.prefix}/database-url"
  description = "PostgreSQL connection string for the CDL API"

  tags = merge(var.tags, { Name = "${var.prefix}/database-url" })
}

resource "aws_secretsmanager_secret_version" "db_url" {
  secret_id     = aws_secretsmanager_secret.db_url.id
  secret_string = local.database_url
}
