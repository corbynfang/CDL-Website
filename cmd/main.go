package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	database.ConnectDatabase()
	defer database.CloseDatabase()

	// Auto-migrate database tables
	database.AutoMigrate()

	// Create Gin router
	r := gin.Default()

	// CORS middleware (if needed)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

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

		// Transfer routes
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
