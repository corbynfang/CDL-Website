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
	dbURL := getEnv("DATABASE_URL", "")
	if dbURL != "" {
		dsn = dbURL
		log.Println("Using DATABASE_URL from environment")
	} else {
		// Build connection string from individual parts
		dsn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=UTC",
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_NAME", "cdl_stats"),
			getEnv("DB_SSLMODE", "disable"),
		)
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
