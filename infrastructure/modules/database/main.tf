resource "aws_db_subnet_group" "main" {
  name       = "${var.prefix}-db-subnet-group"
  subnet_ids = var.subnet_ids

  tags = merge(var.tags, { Name = "${var.prefix}-db-subnet-group" })
}

resource "aws_security_group" "rds" {
  name        = "${var.prefix}-rds-sg"
  description = "Allow Postgres from within the VPC"
  vpc_id      = var.vpc_id

  ingress {
    description = "Postgres from VPC"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, { Name = "${var.prefix}-rds-sg" })
}

resource "aws_db_instance" "main" {
  identifier        = "${var.prefix}-postgres"
  engine            = "postgres"
  engine_version    = "16"
  instance_class    = "db.t3.micro" # free tier
  allocated_storage = 20
  storage_type      = "gp2"

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  publicly_accessible    = false
  multi_az               = false
  storage_encrypted      = true

  backup_retention_period = 0    # free tier does not support automated backups
  skip_final_snapshot     = true

  tags = merge(var.tags, { Name = "${var.prefix}-postgres" })
}
