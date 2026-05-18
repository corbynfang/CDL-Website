#!/usr/bin/env bash
set -euo pipefail

# ── Config ────────────────────────────────────────────────────────────────────
REGION="us-east-1"
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
INFRA_DIR="$PROJECT_ROOT/infrastructure"
FRONTEND_DIR="$PROJECT_ROOT/frontend-react"
GIT_SHA=$(git rev-parse --short HEAD)
SEED="${SEED:-false}" # set SEED=true on first deploy to load CSV data into RDS

echo "==> Deploying commit $GIT_SHA"

# ── Step 1: Apply Terraform ───────────────────────────────────────────────────
echo "==> Running terraform apply..."
cd "$INFRA_DIR"
terraform init -input=false
terraform apply -input=false -auto-approve

# Read outputs into variables
ECR_URL=$(terraform output -raw ecr_repository_url)
S3_BUCKET=$(terraform output -raw s3_bucket_name)
CF_DIST_ID=$(terraform output -raw cloudfront_distribution_id 2>/dev/null || echo "")
ECS_CLUSTER=$(terraform output -raw ecs_cluster)
ECS_SERVICE=$(terraform output -raw ecs_service)
TASK_FAMILY="cdl-api"
TASK_SG=$(terraform output -raw ecs_task_security_group_id)
# Convert Terraform list output ["subnet-a","subnet-b"] → subnet-a,subnet-b
SUBNETS=$(terraform output -json public_subnet_ids | jq -r 'join(",")')

echo "==> ECR: $ECR_URL"
echo "==> S3:  $S3_BUCKET"

# ── Step 2: Build and push Docker image ──────────────────────────────────────
echo "==> Authenticating with ECR..."
aws ecr get-login-password --region "$REGION" | \
  docker login --username AWS --password-stdin "$ECR_URL"

echo "==> Building Docker image..."
cd "$PROJECT_ROOT"
docker build --platform linux/amd64 -t "${ECR_URL}:${GIT_SHA}" -t "${ECR_URL}:latest" .

echo "==> Pushing image to ECR..."
docker push "${ECR_URL}:${GIT_SHA}"
docker push "${ECR_URL}:latest"

# ── Step 3: Update ECS task definition and service ───────────────────────────
echo "==> Registering new ECS task definition..."

CURRENT_TASK=$(aws ecs describe-task-definition \
  --task-definition "$TASK_FAMILY" \
  --region "$REGION" \
  --query "taskDefinition" \
  --output json)

NEW_TASK=$(echo "$CURRENT_TASK" | jq \
  --arg IMAGE "${ECR_URL}:${GIT_SHA}" \
  '.containerDefinitions[0].image = $IMAGE
   | del(.taskDefinitionArn, .revision, .status, .requiresAttributes,
         .placementConstraints, .compatibilities, .registeredAt, .registeredBy)')

NEW_TASK_ARN=$(aws ecs register-task-definition \
  --region "$REGION" \
  --cli-input-json "$NEW_TASK" \
  --query "taskDefinition.taskDefinitionArn" \
  --output text)

echo "==> New task definition: $NEW_TASK_ARN"

echo "==> Updating ECS service..."
aws ecs update-service \
  --region "$REGION" \
  --cluster "$ECS_CLUSTER" \
  --service "$ECS_SERVICE" \
  --task-definition "$NEW_TASK_ARN" \
  --force-new-deployment \
  --output json > /dev/null

# ── Step 3b: Seed database (first deploy only) ────────────────────────────────
# Run the seeder as a one-off Fargate task inside the VPC so it can reach
# the private RDS instance. The container image already has the CSV files baked in.
# Usage: SEED=true ./deploy/deploy.sh
if [ "$SEED" = "true" ]; then
  echo "==> Running database seeder (one-off ECS task)..."
  SEED_TASK_ARN=$(aws ecs run-task \
    --region "$REGION" \
    --cluster "$ECS_CLUSTER" \
    --task-definition "$NEW_TASK_ARN" \
    --launch-type FARGATE \
    --network-configuration "awsvpcConfiguration={subnets=[${SUBNETS}],securityGroups=[${TASK_SG}],assignPublicIp=ENABLED}" \
    --overrides "{\"containerOverrides\":[{\"name\":\"cdl-api\",\"command\":[\"./seeder\"]}]}" \
    --query "tasks[0].taskArn" \
    --output text)

  echo "==> Seeder task started: $SEED_TASK_ARN"
  echo "==> Waiting for seeder to finish (this may take a few minutes)..."

  aws ecs wait tasks-stopped \
    --region "$REGION" \
    --cluster "$ECS_CLUSTER" \
    --tasks "$SEED_TASK_ARN"

  EXIT_CODE=$(aws ecs describe-tasks \
    --region "$REGION" \
    --cluster "$ECS_CLUSTER" \
    --tasks "$SEED_TASK_ARN" \
    --query "tasks[0].containers[0].exitCode" \
    --output text)

  if [ "$EXIT_CODE" != "0" ]; then
    echo "ERROR: Seeder exited with code $EXIT_CODE — check CloudWatch logs at /ecs/cdl-api"
    exit 1
  fi

  echo "==> Seeder finished successfully."
fi

# ── Step 4: Build and upload frontend ────────────────────────────────────────
echo "==> Building Svelte frontend..."
cd "$FRONTEND_DIR"
npm ci
npm run build

echo "==> Syncing frontend to S3..."
aws s3 sync dist/ "s3://${S3_BUCKET}/" \
  --delete \
  --cache-control "public, max-age=31536000, immutable" \
  --exclude "index.html"

aws s3 cp dist/index.html "s3://${S3_BUCKET}/index.html" \
  --cache-control "no-cache, no-store, must-revalidate"

# ── Step 5: Invalidate CloudFront cache ───────────────────────────────────────
if [ -n "$CF_DIST_ID" ]; then
  echo "==> Invalidating CloudFront cache..."
  aws cloudfront create-invalidation \
    --distribution-id "$CF_DIST_ID" \
    --paths "/*" \
    --output json > /dev/null
fi

echo ""
echo "✓ Deploy complete!"
echo "  Site:    https://cdlytics.com"
echo "  Commit:  $GIT_SHA"
