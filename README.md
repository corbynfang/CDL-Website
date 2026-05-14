# CDLytics

Call of Duty League statistics and analytics platform covering seasons from Black Ops Cold War through Black Ops 6.

## Stack

**Frontend** — React 19, TypeScript, Vite, Tailwind CSS, React Router

**Backend** — Go (Gin), GORM, PostgreSQL

**Infrastructure** — AWS (ECS Fargate, RDS, CloudFront, S3, ECR, ALB, Secrets Manager), Terraform

## Features

- K/D leaderboards filterable by season (Cold War → BO6)
- Team rosters and stats per season
- Player profiles with match history
- Tournament brackets
- Transfer history

## Project Structure

```
├── cmd/
│   ├── main.go          # API server entry point
│   └── seed/main.go     # One-time database seeder (reads CSV data)
├── internal/
│   ├── database/        # GORM models and DB connection
│   └── handlers/        # Gin route handlers
├── frontend-react/      # React + TypeScript frontend
├── infrastructure/      # Terraform modules (AWS)
│   └── modules/
│       ├── network/     # VPC, subnets
│       ├── ecr/         # Docker image registry
│       ├── database/    # RDS PostgreSQL
│       ├── alb/         # Application Load Balancer
│       ├── ecs/         # Fargate cluster + service
│       └── frontend/    # S3 + CloudFront + ACM + Route53
├── database/            # CSV source data files
└── deploy/deploy.sh     # Full deploy script
```

## Local Development

```bash
# Backend
go run cmd/main.go

# Frontend
cd frontend-react
npm install
npm run dev
```

The frontend dev server proxies `/api/*` to CloudFront by default. To hit a local backend instead, set `VITE_API_URL=http://localhost:8080/api/v1` in a `.env` file.

## Deploying

Prerequisites: AWS CLI configured, Terraform >= 1.9, Docker, jq

```bash
# First deploy — seeds database from CSV files
SEED=true ./deploy/deploy.sh

# All subsequent deploys
./deploy/deploy.sh
```

The script runs `terraform apply`, builds and pushes the Docker image to ECR, updates the ECS service, builds the React frontend, and syncs it to S3.

## Data

Season data sourced from CDL match records:

| Season | Game |
|--------|------|
| 2024-25 | Black Ops 6 |
| 2023-24 | Modern Warfare III |
| 2022-23 | Modern Warfare II |
| 2021-22 | Vanguard |
| 2020-21 | Black Ops Cold War |
