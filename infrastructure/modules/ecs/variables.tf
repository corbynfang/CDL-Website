variable "prefix" {
  type = string
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "ecr_repository_url" {
  type = string
}

variable "db_secret_arn" {
  description = "Secrets Manager ARN for DATABASE_URL"
  type        = string
}

variable "jwt_secret_arn" {
  description = "Secrets Manager ARN for SUPABASE_JWT_SECRET"
  type        = string
}

variable "supabase_url" {
  description = "Supabase project URL for JWKS endpoint"
  type        = string
}

variable "cloud_map_service_arn" {
  description = "Cloud Map service ARN — ECS registers task IPs here for API Gateway routing"
  type        = string
}

variable "vpc_link_security_group_id" {
  description = "API Gateway VPC Link security group — ECS tasks allow inbound from this SG only"
  type        = string
}

variable "container_port" {
  type    = number
  default = 8080
}

variable "tags" {
  type    = map(string)
  default = {}
}
