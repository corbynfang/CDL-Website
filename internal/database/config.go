package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	// Load .env file if it exists
	_ = godotenv.Load()
	_ = godotenv.Load(".env.railway")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to PostgreSQL database")

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}

// AutoMigrate creates tables using GORM
func AutoMigrate() {
	err := DB.AutoMigrate(
		// Franchise must come before Team (FK dependency)
		&Franchise{},
		&Season{},
		&Team{},
		&Player{},
		&TeamRoster{},
		&Tournament{},
		&Match{},
		// MatchMap and PlayerMapStats depend on Match and Team
		&MatchMap{},
		&PlayerMapStats{},
		&PlayerMatchStats{},
		&PlayerTournamentStats{},
		&TeamTournamentStats{},
		&Coach{},
		&PlayerTransfer{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// pg_trgm enables trigram indexes, which make ILIKE '%scump%' fast even with
	// a leading wildcard. A regular B-tree index can't help with leading wildcards —
	// it can only seek from the left side of the string. Trigrams split every word
	// into overlapping 3-character chunks and index all of them, so any substring
	// match is fast regardless of where in the string it appears.
	DB.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm")
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_players_gamertag_trgm
		ON players USING gin (gamertag gin_trgm_ops)`)

	log.Println("Database migration completed")
}
