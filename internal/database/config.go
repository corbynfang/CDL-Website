package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	var err error

	// Load environment variables (don't fail if .env doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	var dsn string
	var dbURL string

	dbURL = getEnv("DATABASE_URL", "")

	if dbURL != "" {
		dsn = dbURL
		log.Println("Using DATABASE_URL from environment")
	} else {
		// Build connection string from individual parts
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "")
		dbname := getEnv("DB_NAME", "cdl_stats")
		sslmode := getEnv("DB_SSLMODE", "disable")

		// Log the values being used (without password)
		log.Printf("DB connection details - Host: %s, Port: %s, User: %s, DB: %s, SSL: %s",
			host, port, user, dbname, sslmode)

		dsn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=UTC",
			user, password, host, port, dbname, sslmode)
		log.Println("Using individual DB environment variables")
	}

	// Connect to database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
			NoLowerCase:   false,
			NameReplacer:  nil,
		},
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to PostgreSQL database!")

	// Test the connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// AutoMigrate creates tables using GORM
func AutoMigrate() {
	err := DB.AutoMigrate(
		&Season{},
		&Team{},
		&Player{},
		&TeamRoster{},
		&Tournament{},
		&Match{},
		&PlayerMatchStats{},
		&TeamTournamentStats{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully!")
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting database instance: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connection closed successfully")
	}
}
