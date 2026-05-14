variable "prefix" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "container_port" {
  type    = number
  default = 8080
}

variable "tags" {
  type    = map(string)
  default = {}
}
