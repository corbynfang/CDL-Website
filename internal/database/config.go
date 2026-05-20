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

// ConnectDatabase initializes the database connection with exponential backoff.
// Retrying with backoff instead of immediately exiting prevents ECS crash-loops
// from hammering the Supabase pooler and triggering its circuit breaker.
func ConnectDatabase() {
	_ = godotenv.Load()
	_ = godotenv.Load(".env.railway")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	const maxAttempts = 5
	backoff := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := tryConnect(dsn)
		if err == nil {
			log.Println("Connected to PostgreSQL database")
			return
		}

		log.Printf("DB connection attempt %d/%d failed: %v", attempt, maxAttempts, err)
		if attempt == maxAttempts {
			log.Fatal("Failed to connect to database after max retries")
		}
		time.Sleep(backoff)
		if backoff < 30*time.Second {
			backoff *= 2
		}
	}
}

func tryConnect(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	DB = db
	return nil
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

	// Composite index for the GetTeams inner subquery which always filters on
	// both season_id and tournament_type. season_id alone has a single-column index
	// but tournament_type has none — this covers both filter columns in one seek.
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_tournaments_season_type
		ON tournaments (season_id, tournament_type)`)

	log.Println("Database migration completed")
}
