variable "prefix" {
  type = string
}

variable "alb_dns_name" {
  description = "DNS name of the ALB — CloudFront uses this as the API origin"
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
