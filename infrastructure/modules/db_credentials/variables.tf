variable "prefix" {
  type = string
}

variable "connection_string" {
  description = "Full PostgreSQL connection string to store in Secrets Manager"
  type        = string
  sensitive   = true
}

variable "tags" {
  type    = map(string)
  default = {}
}
