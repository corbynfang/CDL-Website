output "ecr_repository_url" {
  description = "Push your Docker image here: docker push <this_url>:latest"
  value       = module.ecr.repository_url
}

output "cloudfront_domain" {
  description = "Your site URL (before DNS is configured)"
  value       = module.frontend.cloudfront_domain
}

output "route53_nameservers" {
  description = "IMPORTANT: Set these as nameservers at your domain registrar for cdlytics.com"
  value       = module.frontend.route53_nameservers
}

output "s3_bucket_name" {
  description = "Upload your Svelte build here: aws s3 sync frontend/dist s3://<this>"
  value       = module.frontend.s3_bucket_name
}

output "cloudfront_distribution_id" {
  description = "Used by deploy.sh to invalidate the CDN cache after uploading new frontend files"
  value       = module.frontend.cloudfront_distribution_id
}

output "ecs_cluster" {
  value = module.ecs.cluster_name
}

output "ecs_service" {
  value = module.ecs.service_name
}

output "ecs_task_security_group_id" {
  value = module.ecs.task_security_group_id
}

output "public_subnet_ids" {
  description = "Public subnet IDs — needed by run-task for the seeder"
  value       = module.network.public_subnet_ids
}
