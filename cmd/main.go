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

	r.GET("/robots.txt", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.Header("Cache-Control", "public, max-age=86400")
		c.String(http.StatusOK, "User-agent: *\nAllow: /\nSitemap: https://cdlytics.com/sitemap.xml\n")
	})

	h := handlers.New(database.DB)
	r.GET("/sitemap.xml", h.GetSitemap)

	api := r.Group("/api/v1")
	api.Use(middleware.RateLimit())
	handlers.RegisterRoutes(api, h)

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
