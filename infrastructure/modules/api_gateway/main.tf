resource "aws_service_discovery_private_dns_namespace" "main" {
  name = "${var.prefix}.local"
  vpc  = var.vpc_id
  tags = var.tags
}

resource "aws_service_discovery_service" "api" {
  name          = "${var.prefix}-api"
  force_destroy = true

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.main.id

    dns_records {
      ttl  = 10
      type = "SRV"
    }

    routing_policy = "WEIGHTED"
  }

  tags = var.tags
}

# Security group for the VPC Link's ENIs (the network interfaces API Gateway creates in the VPC)
resource "aws_security_group" "vpc_link" {
  name        = "${var.prefix}-vpc-link-sg"
  description = "API Gateway VPC Link outbound to ECS tasks"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, { Name = "${var.prefix}-vpc-link-sg" })
}

# VPC Link — bridges API Gateway to private VPC resources via ENIs in the subnets
resource "aws_apigatewayv2_vpc_link" "main" {
  name               = "${var.prefix}-vpc-link"
  security_group_ids = [aws_security_group.vpc_link.id]
  subnet_ids         = var.subnet_ids
  tags               = var.tags
}

resource "aws_apigatewayv2_api" "main" {
  name          = "${var.prefix}-http-api"
  protocol_type = "HTTP"
  tags          = var.tags
}

resource "aws_apigatewayv2_integration" "api" {
  api_id             = aws_apigatewayv2_api.main.id
  integration_type   = "HTTP_PROXY"
  integration_uri    = aws_service_discovery_service.api.arn
  integration_method = "ANY"
  connection_type    = "VPC_LINK"
  connection_id      = aws_apigatewayv2_vpc_link.main.id
}

resource "aws_apigatewayv2_route" "default" {
  api_id    = aws_apigatewayv2_api.main.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.api.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.main.id
  name        = "$default"
  auto_deploy = true
  tags        = var.tags
}
