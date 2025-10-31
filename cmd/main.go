package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func httpsRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Railway sets X-Forwarded-Proto header
		// Only redirect if not localhost and using HTTP
		proto := c.Request.Header.Get("X-Forwarded-Proto")
		host := c.Request.Host

		// Don't redirect localhost or if already HTTPS
		if proto == "http" && host != "localhost:8080" && host != "localhost:3000" {
			httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(301, httpsURL)
			c.Abort()
			return
		}
		c.Next()
	}
}

// Security middleware with enhanced headers
func securityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Enhanced security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Enhanced Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none';")

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
	}
}

// Request logging middleware
func requestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log security-relevant information
		return fmt.Sprintf("[SECURITY] %s | %d | %s | %s | %s | %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Method,
			param.Path,
			param.ClientIP,
			param.Request.UserAgent(),
		)
	})
}

func main() {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Connect to database
	database.ConnectDatabase()
	defer database.CloseDatabase()

	// Auto-migrate database tables (skip for faster startup since tables exist)
	// database.AutoMigrate()

	// Create Gin router
	r := gin.New() // Use gin.New() instead of gin.Default() for more control

	r.Use(requestLogger())
	r.Use(gin.Recovery())
	r.Use(httpsRedirect())
	r.Use(securityMiddleware())
	r.Use(rateLimit())
	r.Use(validateInput())

	// Add cache-control headers to prevent caching of HTML (must be before routes)
	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/" || (!strings.HasPrefix(c.Request.URL.Path, "/api") && !strings.HasPrefix(c.Request.URL.Path, "/assets")) {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
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
		api.GET("/players/:id/matches", handlers.GetPlayerMatches)
		api.GET("/players/top-kd", handlers.GetTopKDPlayers)
		api.GET("/players/top-kd-new", handlers.GetTopKDPlayersNew)
		api.GET("/players/all-kd-stats-tournament", handlers.GetAllPlayersKDStats)
		api.GET("/stats/all-kd-by-tournament", handlers.GetAllPlayersKDStats)

		// Tournament routes
		api.GET("/tournaments", handlers.GetTournaments)
		api.GET("/tournaments/:id", handlers.GetTournament)

		// Transfers routes
		api.GET("/transfers", handlers.GetTransfers)

		// Debug/Validation routes
		api.GET("/debug/validation", handlers.GetDatabaseValidation)
	}

	// Serve static files from frontend/dist/assets with cache-busting headers
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

	// Create HTTP server with security configurations
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		// TLS configuration for HTTPS (Railway handles this)
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(server.ListenAndServe())
}
