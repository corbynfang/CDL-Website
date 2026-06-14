locals {
  prefix = "cdl"
  tags = {
    Project     = "cdl-website"
    Environment = "production"
    ManagedBy   = "terraform"
  }
}

module "network" {
  source             = "./modules/network"
  prefix             = local.prefix
  availability_zones = ["us-east-1a", "us-east-1b"]
  tags               = local.tags
}

module "ecr" {
  source = "./modules/ecr"
  prefix = local.prefix
  tags   = local.tags
}

module "db_credentials" {
  source            = "./modules/db_credentials"
  prefix            = local.prefix
  connection_string = var.database_url
  tags              = local.tags
}

data "aws_secretsmanager_secret" "jwt" {
  name = "${local.prefix}/jwt-secret"
}


module "api_gateway" {
  source         = "./modules/api_gateway"
  prefix         = local.prefix
  aws_region     = var.aws_region
  vpc_id         = module.network.vpc_id
  subnet_ids     = module.network.public_subnet_ids
  container_port = var.container_port
  tags           = local.tags
}

module "ecs" {
  source                     = "./modules/ecs"
  prefix                     = local.prefix
  aws_region                 = var.aws_region
  vpc_id                     = module.network.vpc_id
  subnet_ids                 = module.network.public_subnet_ids
  ecr_repository_url         = module.ecr.repository_url
  db_secret_arn              = module.db_credentials.secret_arn
  jwt_secret_arn             = data.aws_secretsmanager_secret.jwt.arn
  supabase_url               = var.supabase_url
  cloud_map_service_arn      = module.api_gateway.cloud_map_service_arn
  vpc_link_security_group_id = module.api_gateway.vpc_link_security_group_id
  container_port             = var.container_port
  tags                       = local.tags
}

module "frontend" {
  source             = "./modules/frontend"
  prefix             = local.prefix
  api_gateway_domain = module.api_gateway.api_gateway_domain
  domain_name        = "cdlytics.com"
  tags               = local.tags
}
