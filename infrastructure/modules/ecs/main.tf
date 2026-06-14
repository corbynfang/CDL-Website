# --- IAM: Task Execution Role ---
# ECS uses this role (not your code) to pull images from ECR and read secrets.

data "aws_iam_policy_document" "ecs_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "execution" {
  name               = "${var.prefix}-ecs-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_assume.json
  tags               = var.tags
}

resource "aws_iam_role_policy_attachment" "ecs_execution_policy" {
  role       = aws_iam_role.execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy" "read_secret" {
  name = "${var.prefix}-read-db-secret"
  role = aws_iam_role.execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["secretsmanager:GetSecretValue"]
      Resource = [var.db_secret_arn, var.jwt_secret_arn]
    }]
  })
}

# --- CloudWatch Logs ---
resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/${var.prefix}-api"
  retention_in_days = 14

  tags = var.tags
}

# --- Security Group for ECS Tasks ---
resource "aws_security_group" "tasks" {
  name        = "${var.prefix}-ecs-tasks-sg"
  description = "Allow inbound from ALB only" # description is immutable in AWS; kept as-is to avoid SG replacement
  vpc_id      = var.vpc_id

  ingress {
    description     = "API port from VPC Link"
    from_port       = var.container_port
    to_port         = var.container_port
    protocol        = "tcp"
    security_groups = [var.vpc_link_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, { Name = "${var.prefix}-ecs-tasks-sg" })
}

# --- ECS Cluster ---
resource "aws_ecs_cluster" "main" {
  name = "${var.prefix}-cluster"

  setting {
    name  = "containerInsights"
    value = "disabled"
  }

  tags = var.tags
}

# --- Task Definition ---
# Blueprint for the container: image, CPU, RAM, env vars, secrets, logging.
resource "aws_ecs_task_definition" "api" {
  family                   = "${var.prefix}-api"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc" # required for Fargate
  cpu                      = 256      # 0.25 vCPU
  memory                   = 512      # MB
  execution_role_arn       = aws_iam_role.execution.arn

  container_definitions = jsonencode([{
    name      = "${var.prefix}-api"
    image     = "${var.ecr_repository_url}:latest"
    essential = true

    portMappings = [{
      containerPort = var.container_port
      protocol      = "tcp"
    }]

    environment = [
      {
        name  = "SUPABASE_URL"
        value = var.supabase_url
      }
    ]

    secrets = [
      {
        name      = "DATABASE_URL"
        valueFrom = var.db_secret_arn
      },
      {
        name      = "SUPABASE_JWT_SECRET"
        valueFrom = var.jwt_secret_arn
      }
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.ecs.name
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "ecs"
      }
    }
  }])

  tags = var.tags
}

# --- ECS Service ---
# Keeps the task running. Restarts it on crash. Registers IP with ALB target group.
resource "aws_ecs_service" "api" {
  name            = "${var.prefix}-api"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count = 1

  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 1
  }

  force_new_deployment = true

  network_configuration {
    subnets          = var.subnet_ids
    security_groups  = [aws_security_group.tasks.id]
    assign_public_ip = true # Fargate in public subnet needs a public IP to pull from ECR
  }

  service_registries {
    registry_arn   = var.cloud_map_service_arn
    container_name = "${var.prefix}-api"
    container_port = var.container_port
  }

  deployment_minimum_healthy_percent = 50
  deployment_maximum_percent         = 200

  tags = var.tags
}
