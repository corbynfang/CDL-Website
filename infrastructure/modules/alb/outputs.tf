output "alb_dns_name" {
  description = "DNS name of the ALB — used as CloudFront's API origin"
  value       = aws_lb.main.dns_name
}

output "target_group_arn" {
  description = "ARN of the target group — ECS service registers container IPs here"
  value       = aws_lb_target_group.api.arn
}

output "security_group_id" {
  description = "ALB security group ID — ECS task SG allows inbound from this"
  value       = aws_security_group.alb.id
}
