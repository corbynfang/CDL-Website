resource "aws_secretsmanager_secret" "db_url" {
  name        = "${var.prefix}/database-url"
  description = "PostgreSQL connection string for the CDL API"

  tags = merge(var.tags, { Name = "${var.prefix}/database-url" })
}

resource "aws_secretsmanager_secret_version" "db_url" {
  secret_id     = aws_secretsmanager_secret.db_url.id
  secret_string = var.connection_string
}
