# CDLytics

Full-stack Call of Duty League analytics platform covering five competitive seasons (Black Ops Cold War through Black Ops 6). Built end-to-end as a solo project — backend API, React frontend, and full AWS infrastructure provisioned with Terraform.

Live at **[cdlytics.com](https://cdlytics.com)**

## Stack

| Layer | Tech |
|-------|------|
| Frontend | React 19, TypeScript, Vite, Tailwind CSS, React Router |
| Backend | Go (Gin), GORM, PostgreSQL |
| Auth | Supabase |
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
- **Match stats** — per-map breakdown with player K/D, kills, deaths, damage, and mode-specific stats (hill time, plants, defuses, first bloods)
- **Transfer history** — chronological player movement across all five seasons
- **User accounts** — Supabase-backed registration and sign-in with JWT auth validated on the Go backend
- **Match discussion threads** — per-match comment threads for signed-in users
- **Rate limiting** with sliding-window logic and `X-Forwarded-For` parsing behind CloudFront
- **Live event strip** surfacing in-progress events on the home page

## Screenshots

### Home

![CDLytics home page](docs/screenshots/Screenshot%202026-06-09%20at%201.35.39%20PM.png)

The landing page is built around a single global search that resolves players, teams, and events from one box, with quick-pick chips (`Simp`, `Shotzzy`, `Scrap`, `OpTic Texas`, `EWC 2025`) for common queries. Four category cards route into Players, Teams, Events, and Stats, and a **Featured Events** strip surfaces marquee tournaments. The navbar carries Sign In / Sign Out for authenticated users.

### Player profile

![Player profile with per-mode K/D](docs/screenshots/Screenshot%202026-06-09%20at%201.35.59%20PM.png)

Each player page pairs an identity card (gamertag, active status, avatar) with a **K/D Statistics** panel that breaks performance down by game mode — Overall, Hardpoint, Search & Destroy, and Control — each value color-coded above/below 1.0 with a relative bar. The mode splits are computed live by aggregating per-map kills/deaths joined to `match_maps.mode`, so they're accurate for every season (including eras that ship no pre-aggregated stats). Below, tabs for **Last 5 / Matches / Event Stats / Events / Career** drive a match log showing per-series result, K/D, kills, deaths, and date.

### Team page — era switching, rosters & franchise history

![Atlanta FaZe team page with era selector](docs/screenshots/Screenshot%202026-06-09%20at%201.36.12%20PM.png)

Teams are modeled per era (one row per franchise per game), so a team page carries an **ERA** dropdown that re-scopes the whole view to any season the franchise played. The roster shows the selected era's lineup as player cards with a **Current Roster / Players Used** toggle, while the **Franchise History** rail on the right lists every era with the active one marked. Because eras stay linked to one franchise, rebrands and relocations remain connected rather than fragmented.

### K/D rankings

![K/D rankings leaderboard](docs/screenshots/Screenshot%202026-06-09%20at%201.36.35%20PM.png)

The Stats page is a server-paginated, season-filterable leaderboard (the season dropdown scopes to any single era). Rows rank players by K/D with kills and deaths alongside, and the K/D column is color-graded so the top of the board reads at a glance.

### Transfers

![Transfers / roster moves](docs/screenshots/Screenshot%202026-06-09%20at%201.36.45%20PM.png)

The Transfers page is a chronological feed of roster moves across all five seasons. Each entry tags the move type (**SIGNING** / **RELEASE**), the from → to teams, the player's role (AR / SMG), the season, and the date.

### Tournament overview

![CDL 2021 Stage 1 Major overview](docs/screenshots/Screenshot%202026-06-09%20at%201.36.55%20PM.png)

Full event pages open on an **Overview** tab showing teams, prize pool, format, and dates, with all participating team logos as a quick visual index. Tabs for Bracket, Matches, Teams, and Stats let you drill into any angle of the tournament.

### Tournament bracket

![Tournament bracket view](docs/screenshots/Screenshot%202026-06-09%20at%201.37.05%20PM.png)

A custom canvas-based bracket renderer adapts layout per event format — winners bracket, elimination bracket, and grand finals are all laid out in one view with round-by-round navigation. Round filter chips (`Winners Round 1`, `Elimination Finals`, etc.) let you jump to any stage.

### Match stats

![Per-map match stat breakdown](docs/screenshots/Screenshot%202026-06-09%20at%201.37.18%20PM.png)

Each match page shows a per-map breakdown: winning team, score, and a full stat table with player K/D, kills, deaths, damage, and mode-specific columns — hill time for Hardpoint, plants/defuses/first bloods for Search & Destroy. Both sides are shown side by side for each map played.

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
│       ├── auth.go          # JWT validation middleware (Supabase)
│       └── threads.go       # Match discussion thread endpoints
├── frontend-react/          # React + TypeScript SPA
│   └── src/
│       ├── components/      # Page and feature components
│       │   ├── auth/        # AuthModal (sign-in / sign-up flows)
│       │   ├── events/      # Full events feature (bracket, group stage, hero, tabs…)
│       │   └── threads/     # MatchThread — per-match discussion component
│       ├── context/         # AuthContext — global auth state via Supabase session
│       ├── hooks/           # Custom React hooks
│       ├── lib/             # Supabase client, bracket layout utilities
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
