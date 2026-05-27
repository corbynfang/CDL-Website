# CDLytics

Full-stack Call of Duty League analytics platform covering five competitive seasons (Black Ops Cold War through Black Ops 6). Built end-to-end as a solo project ; backend API, React frontend, and full AWS infrastructure provisioned with Terraform.

Live at **[cdlytics.com](https://cdlytics.com)**

## Stack

| Layer | Tech |
|-------|------|
| Frontend | React 19, TypeScript, Vite, Tailwind CSS, React Router |
| Backend | Go (Gin), GORM, PostgreSQL |
| Infrastructure | AWS (ECS Fargate, RDS, CloudFront, S3, ECR, ALB, Secrets Manager) |
| IaC | Terraform |
| CI/CD | GitHub Actions |
| Containerization | Docker |

## Features

- **K/D leaderboards** filterable by season and map type, with server-side pagination and search
- **Player profiles** with cross-season stat history and match logs
- **Team rosters** with season-by-season breakdowns
- **Tournament brackets** — custom canvas-based renderer that adapts layout per event format (single-elim, group stage, EWC format)
- **Events system** — full event pages with hero, standings, match results, bracket view, and team/stat tabs
- **Transfer history** — chronological player movement across all five seasons
- **Rate limiting** with sliding-window logic and `X-Forwarded-For` parsing behind CloudFront
- **Live event strip** surfacing in-progress events on the home page

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  CloudFront CDN                                             │
│    ├── /api/*  → ALB → ECS Fargate (Go API, Docker)        │
│    └── /*      → S3 (React SPA, static assets)             │
└─────────────────────────────────────────────────────────────┘
         │
         ▼
    RDS PostgreSQL (private subnet, not internet-exposed)
```

The API and frontend are decoupled — the Go server handles all `/api/v1/*` routes and the React app is served as a static build from S3. CloudFront routes between them so both share a single domain with no CORS overhead.

Infrastructure is split into Terraform modules (network, ECR, RDS, ALB, ECS, frontend), each independently plannable.

## Project Structure

```
├── cmd/
│   ├── main.go              # API server entry point
│   └── seed/main.go         # One-time database seeder (reads CSV data)
├── internal/
│   ├── database/            # GORM models and DB connection
│   └── handlers/            # Gin route handlers + tests
├── frontend-react/          # React + TypeScript SPA
│   └── src/
│       ├── components/      # Page and feature components
│       │   └── events/      # Full events feature (bracket, group stage, hero, tabs…)
│       ├── hooks/           # Custom React hooks
│       ├── services/        # API client (api.ts)
│       └── types/           # Shared TypeScript types
├── infrastructure/          # Terraform modules (AWS)
│   └── modules/
│       ├── network/         # VPC, subnets, security groups
│       ├── ecr/             # Docker image registry
│       ├── database/        # RDS PostgreSQL
│       ├── alb/             # Application Load Balancer
│       ├── ecs/             # Fargate cluster + service definition
│       └── frontend/        # S3 + CloudFront + ACM + Route53
├── database/                # CSV source data (CDL match records)
└── deploy/deploy.sh         # End-to-end deploy script
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

The Vite dev server proxies `/api/*` to CloudFront by default. To point at a local backend:

```bash
# frontend-react/.env
VITE_API_URL=http://localhost:8080/api/v1
```

## Deploying

Prerequisites: AWS CLI configured, Terraform >= 1.9, Docker, jq

```bash
# First deploy — provisions infrastructure and seeds the database from CSV
SEED=true ./deploy/deploy.sh

# All subsequent deploys
./deploy/deploy.sh
```

The script runs `terraform apply`, builds and pushes the Docker image to ECR, forces a new ECS task deployment, builds the React frontend, and syncs it to S3 with cache invalidation.

## Testing

```bash
# Go unit + integration tests
go test ./...

# Frontend component tests (Vitest)
cd frontend-react && npm test
```

## Data

Season data sourced from CDL match records:

| Season | Game |
|--------|------|
| 2024-25 | Black Ops 6 |
| 2023-24 | Modern Warfare III |
| 2022-23 | Modern Warfare II |
| 2021-22 | Vanguard |
| 2020-21 | Black Ops Cold War |
