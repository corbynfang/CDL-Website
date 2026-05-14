output "cloudfront_domain" {
  description = "CloudFront domain name (e.g. d1abc.cloudfront.net)"
  value       = aws_cloudfront_distribution.main.domain_name
}

output "s3_bucket_name" {
  description = "Name of the S3 bucket holding the frontend build"
  value       = aws_s3_bucket.frontend.bucket
}

output "s3_bucket_id" {
  value = aws_s3_bucket.frontend.id
}

output "route53_nameservers" {
  description = "Update your domain registrar to use these nameservers"
  value       = aws_route53_zone.main.name_servers
}

output "cloudfront_distribution_id" {
  description = "Used by the deploy script to invalidate the CDN cache after uploading new files"
  value       = aws_cloudfront_distribution.main.id
}
