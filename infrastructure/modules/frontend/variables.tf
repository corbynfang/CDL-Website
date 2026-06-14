variable "prefix" {
  type = string
}

variable "api_gateway_domain" {
  description = "Hostname of the API Gateway HTTP API — CloudFront uses this as the API origin"
  type        = string
}

variable "domain_name" {
  description = "Root domain (e.g. cdlytics.com)"
  type        = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
