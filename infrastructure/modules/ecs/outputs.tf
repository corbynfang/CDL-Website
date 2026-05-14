output "cluster_name" {
  value = aws_ecs_cluster.main.name
}

output "service_name" {
  value = aws_ecs_service.api.name
}

output "task_definition_family" {
  value = aws_ecs_task_definition.api.family
}

output "task_security_group_id" {
  description = "Security group attached to ECS tasks — used by run-task for the seeder"
  value       = aws_security_group.tasks.id
}
