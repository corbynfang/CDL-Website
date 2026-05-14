output "repository_url" {
  description = "Full URI of the ECR repository (used in docker push and ECS task definition)"
  value       = aws_ecr_repository.api.repository_url
}

output "repository_name" {
  description = "Name of the ECR repository"
  value       = aws_ecr_repository.api.name
}
