# CDL Stats - Competitive Call of Duty League Statistics Platform

[![Go](https://img.shields.io/badge/Go-1.24.5-blue.svg)](https://golang.org/)
[![React](https://img.shields.io/badge/React-19.1.0-61dafb.svg)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.8.3-blue.svg)](https://www.typescriptlang.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.4.17-38b2ac.svg)](https://tailwindcss.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://www.postgresql.org/)

A comprehensive full-stack web application for exploring and analyzing Call of Duty League (CDL) player and team statistics. Built with modern technologies to provide esports fans, analysts, and gaming enthusiasts with detailed insights into competitive Call of Duty performance data.

**Live Demo: [https://cdlytics.me](https://cdlytics.me)**

## Features

### Team Analytics
- Complete team rosters and player details
- Team performance statistics and rankings
- Historical team data across multiple majors
- Interactive team logos and branding

### Player Statistics
- Individual player performance metrics
- Kill/Death ratios and detailed stats
- Player avatars and profile information
- Career statistics and trends
- Tournament-specific KD statistics
- Player match history

### Data Visualization
- Interactive KD statistics dashboard
- Player comparison tools
- Performance tracking across events
- Tournament-based statistics
- Top KD players rankings

### Modern UI/UX
- Responsive design for all devices
- Dark mode support
- Smooth animations and transitions
- Intuitive navigation and search

### Transfer Tracking
- Player transfer history
- Team roster changes
- Transfer window analytics
- Historical movement data

### Security Features
- Rate limiting to prevent abuse
- Input validation and sanitization
- Enhanced security headers
- HTTPS redirect middleware
- Request logging for security monitoring
- SQL injection prevention via prepared statements

## Tech Stack

### Frontend
- **React 19** - Modern UI framework with hooks
- **TypeScript** - Type-safe development
- **Vite 7** - Lightning-fast build tool
- **Tailwind CSS** - Utility-first styling
- **React Router v7** - Client-side routing
- **Axios** - HTTP client for API calls

### Backend
- **Go 1.24.5** - High-performance server
- **Gin** - Fast HTTP web framework
- **GORM** - Database ORM
- **PostgreSQL** - Reliable data storage
- **RESTful API** - Clean API design

### Development Tools
- **ESLint** - Code quality and consistency
- **PostCSS** - CSS processing
- **Docker** - Containerization support
- **Railway** - Deployment platform

## Project Structure

```
cdl-website/
├── frontend/                 # React + Vite frontend
│   ├── src/
│   │   ├── components/      # UI components
│   │   │   ├── Home.tsx        # Landing page
│   │   │   ├── Teams.tsx       # Team listings
│   │   │   ├── Players.tsx     # Player listings
│   │   │   ├── PlayerDetail.tsx # Player profiles
│   │   │   ├── TeamDetail.tsx # Team profiles
│   │   │   ├── Stats.tsx     # Statistics dashboard
│   │   │   ├── Transfers.tsx   # Transfer tracking
│   │   │   └── Layout.tsx      # App layout
│   │   ├── services/        # API integration
│   │   ├── hooks/           # Custom React hooks
│   │   ├── types/           # TypeScript definitions
│   │   ├── config/          # Configuration files
│   │   └── assets/          # Images and static files
│   └── package.json
├── internal/                 # Go backend
│   ├── handlers/            # API route handlers
│   └── database/            # Database models & config
├── cmd/                     # Go application entry point
│   └── main.go             # Server setup and routing
├── database/                # Database files
│   └── season_stats.csv    # CSV data files
├── Dockerfile              # Container configuration
├── railway.json           # Railway deployment config
└── go.mod                 # Go dependencies
```

## API Endpoints

### Teams
- `GET /api/v1/teams` - List all active teams
- `GET /api/v1/teams/:id` - Get team details by ID
- `GET /api/v1/teams/:id/players` - Get players for a specific team
- `GET /api/v1/teams/:id/stats` - Get team statistics

### Players
- `GET /api/v1/players` - List all players
- `GET /api/v1/players/:id` - Get player details by ID
- `GET /api/v1/players/:id/stats` - Get player statistics
- `GET /api/v1/players/:id/kd` - Get player KD statistics across tournaments
- `GET /api/v1/players/:id/matches` - Get player match history
- `GET /api/v1/players/top-kd` - Get top KD players
- `GET /api/v1/players/top-kd-new` - Get top KD players (new format)
- `GET /api/v1/players/all-kd-stats-tournament` - Get all players KD stats by tournament

### Tournaments
- `GET /api/v1/tournaments` - List all tournaments
- `GET /api/v1/tournaments/:id` - Get tournament details by ID

### Statistics
- `GET /api/v1/stats/all-kd-by-tournament` - Get all players KD statistics by tournament

### Transfers
- `GET /api/v1/transfers` - Get all player transfers

## Database Models

The application uses the following main database models:

- **Season** - CDL seasons and game titles
- **Team** - Team information, logos, colors
- **Player** - Player profiles, gamertags, avatars
- **TeamRoster** - Many-to-many relationship between teams and players
- **Tournament** - Tournament information and details
- **Match** - Match results and scores
- **PlayerMatchStats** - Individual player performance per match
- **PlayerTournamentStats** - Aggregated player stats per tournament
- **TeamTournamentStats** - Team performance per tournament
- **PlayerTransfer** - Player transfer history
- **Coach** - Team coaching staff

## Security Features

The application implements multiple layers of security:

### Rate Limiting
- Request rate limiting per IP address
- Automatic reset intervals

### Input Validation
- ID parameter validation
- Query parameter sanitization
- SQL injection prevention via prepared statements

### Security Headers
- Content Security Policy
- Frame protection
- XSS protection
- Content type protection
- HSTS for secure connections
- Referrer and permissions policies

### CORS
- Restricted to allowed origins only
- Configured for production and development environments

### HTTPS
- Automatic HTTPS redirect in production
- Modern TLS configuration
- Secure cipher suites

## Data Sources

The application includes comprehensive CDL data from:
- Major 1 - Player statistics and team performance
- Major 2 - All 48 player stats
- Major 3 - Complete player roster
- Major 4 - Latest event data
- Call of Duty Champs 2025 - Championship statistics
- 2025 Team Stats - Current season team data

Data is stored in PostgreSQL and can be imported from CSV files in the `database/` directory.

## License

This project is open source and part of the Call of Duty League application ecosystem. Feel free to use and modify for educational purposes.

## Acknowledgments

- Call of Duty League for the competitive gaming platform
- Breakingpoint.gg for most of the data that I scraped
- The CDL community for inspiration and feedback
- Open source contributors and the gaming community

---

**Built for the Call of Duty League community**

*For questions or support, please open an issue on GitHub.*
