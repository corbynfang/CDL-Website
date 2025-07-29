package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/handlers"
	"github.com/gin-gonic/gin"
)

// Rate limiting map
var requestCounts = make(map[string]int)
var lastReset = time.Now()

// Rate limiting middleware
func rateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Reset counters every minute
		if now.Sub(lastReset) > time.Minute {
			requestCounts = make(map[string]int)
			lastReset = now
		}

		// Check rate limit (100 requests per minute per IP)
		if requestCounts[clientIP] > 100 {
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		requestCounts[clientIP]++
		c.Next()
	}
}

// Input validation middleware
func validateInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate ID parameters
		if id := c.Param("id"); id != "" {
			if _, err := strconv.Atoi(id); err != nil {
				c.JSON(400, gin.H{"error": "Invalid ID parameter"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func main() {
	// Connect to database
	database.ConnectDatabase()
	defer database.CloseDatabase()

	// Auto-migrate database tables
	database.AutoMigrate()

	// Create Gin router
	r := gin.Default()

	// Security middleware
	r.Use(func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:;")

		// CORS - restrict to your domain only
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := []string{"https://cdlytics.me", "http://localhost:3000", "http://localhost:5173"}

		for _, allowed := range allowedOrigins {
			if origin == allowed {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Apply rate limiting and input validation
	r.Use(rateLimit())
	r.Use(validateInput())

	// API routes
	api := r.Group("/api/v1")
	{
		// Team routes
		api.GET("/teams", handlers.GetTeams)
		api.GET("/teams/:id", handlers.GetTeam)
		api.GET("/teams/:id/players", handlers.GetTeamPlayers)
		api.GET("/teams/:id/stats", handlers.GetTeamStats)

		// Player routes
		api.GET("/players", handlers.GetPlayers)
		api.GET("/players/:id", handlers.GetPlayer)
		api.GET("/players/:id/stats", handlers.GetPlayerStats)
		api.GET("/players/:id/kd", handlers.GetPlayerKDStats)
		api.GET("/players/top-kd", handlers.GetTopKDPlayers)
		api.GET("/players/top-kd-new", handlers.GetTopKDPlayersNew)
		api.GET("/players/all-kd-stats-tournament", handlers.GetAllPlayersKDStats)
		api.GET("/stats/all-kd-by-tournament", handlers.GetAllPlayersKDStats)

		// Tournament routes
		api.GET("/tournaments", handlers.GetTournaments)
		api.GET("/tournaments/:id", handlers.GetTournament)

		// Transfers routes
		api.GET("/transfers", handlers.GetTransfers)
	}

	// Serve static files from frontend/dist/assets
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")

	// Serve index.html for all non-API routes (SPA catch-all)
	r.NoRoute(func(c *gin.Context) {
		// Don't serve index.html for API routes
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}

		// Check if file exists in dist folder
		filePath := filepath.Join("./frontend/dist", c.Request.URL.Path)
		if _, err := os.Stat(filePath); err == nil {
			c.File(filePath)
			return
		}

		// Serve index.html for all other routes (SPA routing)
		c.File("./frontend/dist/index.html")
	})

	// Get port from environment variable (Railway sets this)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port for local development
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
