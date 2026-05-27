package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/handlers"
	"github.com/corbynfang/CDL-Website/internal/middleware"
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

	if err := r.SetTrustedProxies([]string{"10.0.0.0/16"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	h := handlers.New(database.DB)

	api := r.Group("/api/v1")
	api.Use(middleware.RateLimit())
	{
		api.GET("/seasons", h.GetSeasons)
		api.GET("/seasons/:id", h.GetSeason)
		api.GET("/seasons/active", h.GetActiveSeason)

		api.GET("/teams", h.GetTeams)
		api.GET("/teams/:id", h.GetTeam)
		api.GET("/teams/:id/players", h.GetTeamPlayers)
		api.GET("/teams/:id/stats", h.GetTeamStats)

		api.GET("/players", h.GetPlayers)
		api.GET("/players/:id", h.GetPlayer)
		api.GET("/players/:id/stats", h.GetPlayerStats)
		api.GET("/players/:id/kd", h.GetPlayerKDStats)
		api.GET("/players/:id/matches", h.GetPlayerMatches)
		api.GET("/players/:id/franchise-career", h.GetPlayerFranchiseCareer)
		api.GET("/players/top-kd", h.GetTopKDPlayers)

		api.GET("/stats/all-kd-by-tournament", h.GetAllPlayersKDStats)

		api.GET("/matches/:id", h.GetMatch)

		api.GET("/franchises", h.GetFranchises)
		api.GET("/franchises/:key", h.GetFranchise)

		api.GET("/tournaments", h.GetTournaments)
		api.GET("/tournaments/slug/:slug", h.GetTournamentBySlug)
		api.GET("/tournaments/:id", h.GetTournament)
		api.GET("/tournaments/:id/bracket", h.GetTournamentBracket)
		api.GET("/tournaments/:id/matches", h.GetTournamentMatches)
		api.GET("/tournaments/:id/teams", h.GetTournamentTeams)
		api.GET("/tournaments/:id/stats", h.GetTournamentStats)

		api.GET("/transfers", h.GetTransfers)
	}

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
