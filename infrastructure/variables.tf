variable "aws_region" {
  description = "AWS region to deploy into"
  type        = string
  default     = "us-east-1"
}

variable "container_port" {
  description = "Port your Go API listens on"
  type        = number
  default     = 8080
}

variable "db_name" {
  description = "PostgreSQL database name"
  type        = string
  default     = "cdlwebsite"
}

variable "db_username" {
  description = "PostgreSQL master username"
  type        = string
  default     = "cdladmin"
}

variable "database_url" {
  description = "Full PostgreSQL connection string stored in Secrets Manager and injected into ECS"
  type        = string
  sensitive   = true
}
