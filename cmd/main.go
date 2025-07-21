package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	database.ConnectDatabase()
	// Enable AutoMigrate to ensure GORM can map the tables properly
	database.AutoMigrate()
	defer database.CloseDatabase()

	router := setupRouter()

	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)

	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

func setupRouter() *gin.Engine {
	if getEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(corsMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "CDL Stats API is running",
		})
	})

	api := router.Group("/api/v1")
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
		api.GET("/players/:id/kd-stats", handlers.GetPlayerKDStats)
		api.GET("/stats/top-kd", handlers.GetTopKDPlayers)
		api.GET("/stats/top-kd-new", handlers.GetTopKDPlayersNew)

		// Tournament routes
		api.GET("/tournaments", handlers.GetTournaments)
		api.GET("/tournaments/:id", handlers.GetTournament)

		// KD stats for all players and all majors
		api.GET("/stats/all-kd-by-tournament", handlers.GetAllPlayersKDStats)
	}

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
