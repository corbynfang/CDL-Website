# CDL Stats Frontend

A modern React TypeScript frontend for the CDL Stats API, built with Vite and Tailwind CSS.

## Features

- **Modern React**: Built with React 18 and TypeScript for type safety
- **Tailwind CSS**: Utility-first CSS framework for rapid UI development
- **React Router**: Client-side routing for seamless navigation
- **Axios**: HTTP client for API communication
- **Responsive Design**: Mobile-first responsive layout
- **Dark Theme**: Modern dark theme optimized for gaming aesthetics

## Tech Stack

- **React 18** - UI library
- **TypeScript** - Type safety and better developer experience
- **Vite** - Fast build tool and development server
- **Tailwind CSS** - Utility-first CSS framework
- **React Router** - Client-side routing
- **Axios** - HTTP client for API calls

## Getting Started

### Prerequisites

- Node.js 16+ 
- npm or yarn
- Go backend running on localhost:8080

### Installation

1. Install dependencies:
```bash
npm install
```

2. Create environment file:
```bash
echo "VITE_API_URL=http://localhost:8080/api/v1" > .env
```

3. Start the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:3000`

### Building for Production

```bash
npm run build
```

The built files will be in the `dist` directory.

## Project Structure

```
src/
├── components/          # React components
│   ├── Layout.tsx      # Main layout with navigation
│   ├── Home.tsx        # Home page
│   ├── Teams.tsx       # Teams listing page
│   └── Players.tsx     # Players listing page
├── services/           # API services
│   └── api.ts          # Axios configuration and API calls
├── types/              # TypeScript type definitions
│   └── index.ts        # Interface definitions
├── App.tsx             # Main app component with routing
├── main.tsx            # Application entry point
└── index.css           # Global styles with Tailwind
```

## API Integration

The frontend communicates with the Go backend through the following endpoints:

- `GET /api/v1/teams` - Get all teams
- `GET /api/v1/teams/:id` - Get specific team
- `GET /api/v1/teams/:id/players` - Get team players
- `GET /api/v1/players` - Get all players
- `GET /api/v1/players/:id` - Get specific player
- `GET /api/v1/players/:id/stats` - Get player statistics

## Development

### Adding New Components

1. Create a new component in `src/components/`
2. Add the route in `src/App.tsx`
3. Update the navigation in `src/components/Layout.tsx`

### Styling

The project uses Tailwind CSS for styling. Custom styles can be added in `src/index.css` using the `@layer` directive.

### Type Safety

All API responses are typed using TypeScript interfaces defined in `src/types/index.ts`. When adding new API endpoints, make sure to:

1. Add the interface to `src/types/index.ts`
2. Add the API function to `src/services/api.ts`
3. Use proper typing in components

## Environment Variables

- `VITE_API_URL`: Base URL for the API (default: http://localhost:8080/api/v1)

## Contributing

1. Follow the existing code style
2. Add proper TypeScript types for new features
3. Test the application thoroughly
4. Update documentation as needed

## License

This project is part of the CDL Stats application.
