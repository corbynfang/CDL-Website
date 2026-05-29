package handlers

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var pgDB *gorm.DB

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	ctx := context.Background()

	ctr, err := tcpg.Run(ctx,
		"postgres:16-alpine",
		tcpg.WithDatabase("cdltest"),
		tcpg.WithUsername("test"),
		tcpg.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		log.Println("testcontainers: container start failed:", err)
		return m.Run()
	}
	defer ctr.Terminate(ctx) //nolint:errcheck

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Println("testcontainers: connection string failed:", err)
		return m.Run()
	}

	db, err := gorm.Open(gormpg.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Println("gorm: open failed:", err)
		return m.Run()
	}

	if err = db.AutoMigrate(
		&models.Franchise{},
		&models.Player{},
		&models.Season{},
		&models.Team{},
		&models.TeamRoster{},
		&models.Tournament{},
		&models.Match{},
		&models.MatchMap{},
		&models.PlayerMapStats{},
		&models.PlayerMatchStats{},
		&models.PlayerTournamentStats{},
		&models.TeamTournamentStats{},
		&models.Coach{},
		&models.PlayerTransfer{},
	); err != nil {
		log.Println("gorm: automigrate failed:", err)
		return m.Run()
	}

	pgDB = db
	return m.Run()
}

func setupPGTx(t *testing.T) {
	t.Helper()
	if pgDB == nil {
		t.Skip("postgres test container unavailable (is Docker running?)")
	}
	tx := pgDB.Begin()
	require.NoError(t, tx.Error)

	old := database.DB
	database.DB = tx
	t.Cleanup(func() {
		tx.Rollback() //nolint:errcheck
		database.DB = old
	})
}

func pgSeason(t *testing.T) {
	t.Helper()
	require.NoError(t, database.DB.Create(&models.Season{
		ID: 1, Name: "BO6 Season 2025", GameTitle: "Black Ops 6", GameCode: "BO6",
		StartDate: time.Now(),
	}).Error)
}

func pgTournament(t *testing.T) {
	t.Helper()
	require.NoError(t, database.DB.Create(&models.Tournament{
		ID: 1, SeasonID: 1, Name: "CDL Major 1 2025", Slug: "cdl-major-1-2025",
		TournamentType: "major", StartDate: time.Now(),
	}).Error)
}

func pgTeams(t *testing.T) {
	t.Helper()
	require.NoError(t, database.DB.Create(&models.Team{
		ID: 1, Name: "OpTic Texas", Abbreviation: "OTX",
		PrimaryColor: "#000", SecondaryColor: "#0f0",
	}).Error)
	require.NoError(t, database.DB.Create(&models.Team{
		ID: 2, Name: "Atlanta FaZe", Abbreviation: "ATL",
		PrimaryColor: "#000", SecondaryColor: "#f00",
	}).Error)
}

func pgMatchEnv(t *testing.T) {
	t.Helper()
	pgSeason(t)
	pgTournament(t)
	pgTeams(t)
}

func pgMatch(t *testing.T, matchID uint) {
	t.Helper()
	winnerID := uint(1)
	require.NoError(t, database.DB.Create(&models.Match{
		ID:              matchID,
		TournamentID:    1,
		Team1ID:         1,
		Team2ID:         2,
		WinnerID:        &winnerID,
		Team1Score:      3,
		Team2Score:      1,
		MatchDate:       time.Now(),
		Format:          "BO5",
		MatchType:       "winners_r1",
		BracketRound:    "winners_r1",
		BracketPosition: 1,
	}).Error)
}
