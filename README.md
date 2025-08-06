# 🎮 CDL Stats - Competitive Call of Duty League Statistics Platform

[![Go](https://img.shields.io/badge/Go-1.24.5-blue.svg)](https://golang.org/)
[![React](https://img.shields.io/badge/React-19.1.0-61dafb.svg)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.8.3-blue.svg)](https://www.typescriptlang.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.4.17-38b2ac.svg)](https://tailwindcss.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://www.postgresql.org/)

A comprehensive full-stack web application for exploring and analyzing Call of Duty League (CDL) player and team statistics. Built with modern technologies to provide esports fans, analysts, and gaming enthusiasts with detailed insights into competitive Call of Duty performance data.

**🌐 Live Demo: [https://cdlytics.me](https://cdlytics.me)**

## ✨ Features

### 🏆 **Team Analytics**
- Complete team rosters and player details
- Team performance statistics and rankings
- Historical team data across multiple majors
- Interactive team logos and branding

### 👤 **Player Statistics**
- Individual player performance metrics
- Kill/Death ratios and detailed stats
- Player avatars and profile information
- Career statistics and trends

### 📊 **Data Visualization**
- Interactive KD statistics dashboard
- Player comparison tools
- Performance tracking across events
- Real-time data updates

### 🎨 **Modern UI/UX**
- Responsive design for all devices
- Dark mode support
- Smooth animations and transitions
- Intuitive navigation and search

### 🔄 **Transfer Tracking**
- Player transfer history
- Team roster changes
- Transfer window analytics
- Historical movement data

## 🛠 Tech Stack

### **Frontend**
- **React 19** - Modern UI framework with hooks
- **TypeScript** - Type-safe development
- **Vite** - Lightning-fast build tool
- **Tailwind CSS** - Utility-first styling
- **React Router** - Client-side routing
- **Axios** - HTTP client for API calls

### **Backend**
- **Go 1.24** - High-performance server
- **Gin** - Fast HTTP web framework
- **GORM** - Database ORM
- **PostgreSQL** - Reliable data storage
- **RESTful API** - Clean API design

### **Development Tools**
- **ESLint** - Code quality and consistency
- **PostCSS** - CSS processing
- **Docker** - Containerization support

## 📁 Project Structure

```
cdl-website/
├── 📂 frontend/                 # React + Vite frontend
│   ├── 📂 src/
│   │   ├── 📂 components/      # UI components
│   │   │   ├── Home.tsx        # Landing page
│   │   │   ├── Teams.tsx       # Team listings
│   │   │   ├── Players.tsx     # Player listings
│   │   │   ├── PlayerDetails.tsx # Player profiles
│   │   │   ├── TeamDetails.tsx # Team profiles
│   │   │   ├── KDStats.tsx     # Statistics dashboard
│   │   │   ├── Transfers.tsx   # Transfer tracking
│   │   │   └── Layout.tsx      # App layout
│   │   ├── 📂 services/        # API integration
│   │   ├── 📂 types/           # TypeScript definitions
│   │   └── 📂 assets/          # Images and static files
│   └── package.json
├── 📂 internal/                 # Go backend
│   ├── 📂 handlers/            # API route handlers
│   └── 📂 database/            # Database models & config
├── 📂 database/                # Database migrations
├── 📂 datacsv/                 # CSV data files
│   ├── Major-1-players-stats.csv
│   ├── CDL_Major2_All_48_Player_Stats.csv
│   ├── CDL_Major3_All_48_Players.csv
│   ├── Major4data.csv
│   ├── CodChamps2025.csv
│   └── 2025_Team_stats.csv
├── 📂 cmd/                     # Go application entry point
├── Dockerfile                  # Container configuration
└── go.mod                      # Go dependencies
```

## 🚀 Quick Start

### Prerequisites

- **Node.js** 18+ and npm
- **Go** 1.24+
- **PostgreSQL** 15+

### 1. Clone the Repository

```bash
git clone https://github.com/corbynfang/CDL-Website
cd cdl-website
```

### 2. Backend Setup

```bash
# Install Go dependencies
go mod tidy

# Set up environment variables
cp .env.example .env
# Edit .env with your database credentials

# Run database migrations
# (See database/ directory for schema)

# Start the Go server
go run cmd/main.go
```

### 3. Frontend Setup

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

Visit [http://localhost:3000](http://localhost:3000) to view the application.

## 🌐 API Endpoints

### Teams
- `GET /api/v1/teams` - List all teams
- `GET /api/v1/teams/:id` - Get team details
- `GET /api/v1/teams/:id/players` - Get team players

### Players
- `GET /api/v1/players` - List all players
- `GET /api/v1/players/:id` - Get player details
- `GET /api/v1/players/:id/stats` - Get player statistics

### Statistics
- `GET /api/v1/stats/kd` - KD statistics
- `GET /api/v1/stats/majors` - Major event statistics

## 📊 Data Sources

The application includes comprehensive CDL data from:
- **Major 1** - Player statistics and team performance
- **Major 2** - All 48 player stats
- **Major 3** - Complete player roster
- **Major 4** - Latest event data
- **Call of Duty Champs 2025** - Championship statistics
- **2025 Team Stats** - Current season team data

## 🐳 Docker Deployment

```bash
# Build and run with Docker
docker build -t cdl-website .
docker run -p 8080:8080 cdl-website
```

## 🤝 Contributing

Contributions are welcome! This is a learning project, so feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow TypeScript best practices
- Add proper error handling
- Include TypeScript types for new features
- Test your changes locally before submitting

## 📝 License

This project is open source and part of the Call of Duty League application ecosystem. Feel free to use and modify for educational purposes.

## 🙏 Acknowledgments

- Call of Duty League for the competitive gaming platform
- The CDL community for inspiration and feedback
- Open source contributors and the gaming community

---

**Built with ❤️ for the Call of Duty League community**

*For questions or support, please open an issue on GitHub.* # Updated Tue Aug  5 00:02:04 CDT 2025
# Updated Wed Aug  6 17:32:55 CDT 2025
