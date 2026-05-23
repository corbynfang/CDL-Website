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

	// Trust only the VPC CIDR (10.0.0.0/16) — the ALB sits there and is the
	// only machine that should be forwarding requests. With this set, c.ClientIP()
	// walks X-Forwarded-For and returns the first IP that is NOT in the trusted
	// range, which is the real client IP added by CloudFront.
	if err := r.SetTrustedProxies([]string{"10.0.0.0/16"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Security headers on every response — defence-in-depth layer.
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
		c.Next()
	})

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

	// Lightweight health check — use this as the ALB target group health check
	// path instead of /api/v1/teams to avoid running a five-table join every 30s.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	api.Use(handlers.RateLimit())
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

		api.GET("/stats/all-kd-by-tournament", handlers.GetAllPlayersKDStats)

		api.GET("/matches/:id", handlers.GetMatch)

		api.GET("/franchises", handlers.GetFranchises)
		api.GET("/franchises/:key", handlers.GetFranchise)

		api.GET("/tournaments", handlers.GetTournaments)
		api.GET("/tournaments/slug/:slug", handlers.GetTournamentBySlug)
		api.GET("/tournaments/:id", handlers.GetTournament)
		api.GET("/tournaments/:id/bracket", handlers.GetTournamentBracket)
		api.GET("/tournaments/:id/matches", handlers.GetTournamentMatches)
		api.GET("/tournaments/:id/teams", handlers.GetTournamentTeams)
		api.GET("/tournaments/:id/stats", handlers.GetTournamentStats)

		api.GET("/transfers", handlers.GetTransfers)
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
