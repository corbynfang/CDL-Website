# CDL Stats – Competitive Call of Duty League Stats Website

Welcome to CDL Stats! This project is a full-stack web application for exploring and visualizing player and team statistics from the Call of Duty League (CDL). It’s designed for esports fans, analysts, and anyone interested in competitive gaming data.

## 🚀 Features

- Modern, responsive web interface built with React, TypeScript, and Tailwind CSS
- Fast, RESTful API backend powered by Go (Gin) and PostgreSQL
- Browse teams, players, and detailed match statistics
- Search and filter by player, team, or event
- Dark mode and mobile-friendly design
- Easy local development and deployment

## 🛠 Tech Stack

**Frontend:**
- React 19 + TypeScript
- Vite (build tool)
- Tailwind CSS (utility-first styling)
- React Router (client-side routing)
- Axios (API requests)

**Backend:**
- Go 1.24 (Gin web framework)
- GORM (ORM for PostgreSQL)
- PostgreSQL (data storage)
- RESTful API design

## 📁 Project Structure

```
cdl-website/
├── frontend/         # React + Vite frontend
│   ├── src/components/   # UI components (Home, Teams, Players, etc.)
│   ├── services/         # API calls
│   └── types/            # TypeScript types
├── internal/
│   ├── handlers/         # Go API route handlers
│   └── database/         # DB models and config
├── database/        # SQL schema and migration files
├── cmd/             # Go entrypoint
└── main             # Go main file
```

## 🏁 Getting Started

### Prerequisites

- Node.js 16+ and npm (for frontend)
- Go 1.24+ (for backend)
- PostgreSQL (database)

### 1. Clone the repository

```sh
git clone https://github.com/corbynfang/CDL-Website
cd cdl-website
```

### 2. Set up the backend

- Configure your PostgreSQL database (see `database/schema.sql` for structure)
- Set environment variables as needed (e.g., database URL)
- Run the Go server:

```sh
go run cmd/main.go
```

### 3. Set up the frontend

```sh
cd frontend
npm install
npm run dev
```

Visit [http://localhost:3000](http://localhost:3000) to view the app. (Will be deploying once I can figure out React code)

## 🌐 API Endpoints

- `GET /api/v1/teams` – List all teams
- `GET /api/v1/players` – List all players
- `GET /api/v1/players/:id` – Player details and stats
- ...and more (see backend source for full list)

## 🧑‍💻 Contributing (You don't have too, I am new to really big projects like this.)

- Fork the repo and create a feature branch
- Use clear commit messages and follow code style
- Add TypeScript types for new features
- Test your changes before submitting a PR

## 📄 License

This project is open source and part of the Call Of Duty League application. 