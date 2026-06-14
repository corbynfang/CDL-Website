output "api_gateway_domain" {
  description = "Hostname of the HTTP API — use as CloudFront origin domain_name (no https:// prefix)"
  value       = "${aws_apigatewayv2_api.main.id}.execute-api.${var.aws_region}.amazonaws.com"
}

output "invoke_url" {
  description = "Full invoke URL including stage path"
  value       = aws_apigatewayv2_stage.default.invoke_url
}

output "cloud_map_service_arn" {
  description = "ARN of the Cloud Map service — pass to ECS service_registries"
  value       = aws_service_discovery_service.api.arn
}

output "vpc_link_security_group_id" {
  description = "Security group on the VPC Link ENIs — ECS task SG allows inbound from this"
  value       = aws_security_group.vpc_link.id
}
