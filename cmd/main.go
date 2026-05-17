package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	database.ConnectDatabase()
	defer database.CloseDatabase()
	database.AutoMigrate()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// CORS middleware — frontend is served from CloudFront/S3, not this server
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := []string{
			"https://cdlytics.com",
			"https://www.cdlytics.com",
			"http://localhost:5173",
			"http://localhost:5174",
		}

		for _, allowed := range allowedOrigins {
			if origin == allowed {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Cache-Control")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api/v1")
	{
		api.GET("/seasons", handlers.GetSeasons)
		api.GET("/seasons/:id", handlers.GetSeason)
		api.GET("/seasons/active", handlers.GetActiveSeason)

		api.GET("/teams", handlers.GetTeams)
		api.GET("/teams/:id", handlers.GetTeam)
		api.GET("/teams/:id/players", handlers.GetTeamPlayers)
		api.GET("/teams/:id/stats", handlers.GetTeamStats)

		api.GET("/players", handlers.GetPlayers)
		api.GET("/players/:id", handlers.GetPlayer)
		api.GET("/players/:id/stats", handlers.GetPlayerStats)
		api.GET("/players/:id/kd", handlers.GetPlayerKDStats)
		api.GET("/players/:id/matches", handlers.GetPlayerMatches)
		api.GET("/players/:id/franchise-career", handlers.GetPlayerFranchiseCareer)
		api.GET("/players/top-kd", handlers.GetTopKDPlayers)
		api.GET("/players/top-kd-new", handlers.GetTopKDPlayersNew)

		api.GET("/stats/all-kd-by-tournament", handlers.GetAllPlayersKDStats)
		api.GET("/players/all-kd-stats-tournament", handlers.GetAllPlayersKDStats)

		api.GET("/matches/:id", handlers.GetMatch)

		api.GET("/franchises", handlers.GetFranchises)
		api.GET("/franchises/:key", handlers.GetFranchise)

		api.GET("/tournaments", handlers.GetTournaments)
		api.GET("/tournaments/:id", handlers.GetTournament)
		api.GET("/tournaments/:id/bracket", handlers.GetTournamentBracket)

		api.GET("/transfers", handlers.GetTransfers)

		api.GET("/debug/validation", handlers.GetDatabaseValidation)
	}

	// 404 for unknown API routes — CloudFront handles all non-API routes before they get here
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(server.ListenAndServe())
}
