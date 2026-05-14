variable "prefix" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "vpc_cidr" {
  description = "VPC CIDR block — RDS allows inbound Postgres from anything inside the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "db_name" {
  type = string
}

variable "db_username" {
  type = string
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "tags" {
  type    = map(string)
  default = {}
}
