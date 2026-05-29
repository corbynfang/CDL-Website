package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var rosterTestDB *gorm.DB

func TestMain(m *testing.M) {
	os.Exit(runRosterTests(m))
}

func runRosterTests(m *testing.M) int {
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
		&models.Season{},
		&models.Team{},
		&models.Player{},
		&models.Tournament{},
		&models.Match{},
		&models.MatchMap{},
		&models.PlayerMapStats{},
		&models.TeamRoster{},
	); err != nil {
		log.Println("gorm: automigrate failed:", err)
		return m.Run()
	}

	rosterTestDB = db
	return m.Run()
}

func rosterTx(t *testing.T) *gorm.DB {
	t.Helper()
	if rosterTestDB == nil {
		t.Skip("postgres test container unavailable (is Docker running?)")
	}
	tx := rosterTestDB.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { tx.Rollback() }) //nolint:errcheck
	return tx
}

func mkSeason(t *testing.T, db *gorm.DB, id uint, code string) {
	t.Helper()
	require.NoError(t, db.Create(&models.Season{
		ID: id, Name: code + " Season", GameTitle: code, GameCode: code, StartDate: time.Now(),
	}).Error)
}

func mkTeam(t *testing.T, db *gorm.DB, id uint, name, abbr string, franchiseID *uint) {
	t.Helper()
	require.NoError(t, db.Create(&models.Team{
		ID: id, Name: name, Abbreviation: abbr, FranchiseID: franchiseID,
	}).Error)
}

func mkFranchise(t *testing.T, db *gorm.DB, id uint, key string) {
	t.Helper()
	require.NoError(t, db.Create(&models.Franchise{
		ID: id, FranchiseKey: key, Name: key,
	}).Error)
}

func mkTournament(t *testing.T, db *gorm.DB, id, seasonID uint) {
	t.Helper()
	require.NoError(t, db.Create(&models.Tournament{
		ID: id, SeasonID: seasonID, Name: fmt.Sprintf("Event %d", id),
		Slug: fmt.Sprintf("event-%d", id), TournamentType: "major", StartDate: time.Now(),
	}).Error)
}

func mkPlayer(t *testing.T, db *gorm.DB, id uint, tag string) {
	t.Helper()
	require.NoError(t, db.Create(&models.Player{ID: id, Gamertag: tag}).Error)
}

func mkMatch(t *testing.T, db *gorm.DB, id, tournamentID, team1, team2 uint, date time.Time) {
	t.Helper()
	require.NoError(t, db.Create(&models.Match{
		ID: id, TournamentID: tournamentID, Team1ID: team1, Team2ID: team2, MatchDate: date,
	}).Error)
}

func mkMap(t *testing.T, db *gorm.DB, matchID uint, mapNumber int, played bool) {
	t.Helper()
	mm := models.MatchMap{MatchID: matchID, MapNumber: mapNumber, Played: played}
	require.NoError(t, db.Create(&mm).Error)
	// MatchMap.Played has gorm `default:true`, so a zero-value false is omitted on
	// insert and the DB default wins. Force the column when we actually want a DNP map.
	if !played {
		require.NoError(t, db.Model(&models.MatchMap{}).Where("id = ?", mm.ID).
			Update("played", false).Error)
	}
}

func mkMapStat(t *testing.T, db *gorm.DB, matchID uint, mapNumber int, playerID, teamID uint) {
	t.Helper()
	require.NoError(t, db.Create(&models.PlayerMapStats{
		MatchID: matchID, MapNumber: mapNumber, PlayerID: playerID, TeamID: teamID,
	}).Error)
}

func findStint(stints []rosterStint, playerID, teamID, seasonID uint) *rosterStint {
	for i := range stints {
		s := &stints[i]
		if s.PlayerID == playerID && s.TeamID == teamID && s.SeasonID == seasonID {
			return s
		}
	}
	return nil
}

func TestInferRosterStints_PlayerOnTwoTeamsAcrossSeasons(t *testing.T) {
	db := rosterTx(t)

	mkSeason(t, db, 1, "BO6")
	mkSeason(t, db, 2, "MW3")
	mkTeam(t, db, 1, "OpTic Texas", "OTX", nil)
	mkTeam(t, db, 2, "Atlanta FaZe", "ATL", nil)
	mkTournament(t, db, 1, 1)
	mkTournament(t, db, 2, 2)
	mkPlayer(t, db, 1, "Dashy")

	mkMatch(t, db, 1, 1, 1, 2, time.Now())
	mkMap(t, db, 1, 1, true)
	mkMapStat(t, db, 1, 1, 1, 1) // player1 on team1, season1

	mkMatch(t, db, 2, 2, 1, 2, time.Now())
	mkMap(t, db, 2, 1, true)
	mkMapStat(t, db, 2, 1, 1, 2) // player1 on team2, season2

	stints, err := inferRosterStints(db)
	require.NoError(t, err)
	require.Len(t, stints, 2)
	require.NotNil(t, findStint(stints, 1, 1, 1), "expected player1/team1/season1 stint")
	require.NotNil(t, findStint(stints, 1, 2, 2), "expected player1/team2/season2 stint")
}

func TestInferRosterStints_ExcludesDNPMaps(t *testing.T) {
	db := rosterTx(t)

	mkSeason(t, db, 1, "BO6")
	mkTeam(t, db, 1, "OpTic Texas", "OTX", nil)
	mkTeam(t, db, 2, "Atlanta FaZe", "ATL", nil)
	mkTournament(t, db, 1, 1)
	mkPlayer(t, db, 1, "Dashy")
	mkPlayer(t, db, 2, "Ghosty") // appears only on a DNP map

	mkMatch(t, db, 1, 1, 1, 2, time.Now())
	mkMap(t, db, 1, 1, true)  // played
	mkMap(t, db, 1, 2, false) // DNP
	mkMapStat(t, db, 1, 1, 1, 1)
	mkMapStat(t, db, 1, 2, 2, 1) // stat on the DNP map

	stints, err := inferRosterStints(db)
	require.NoError(t, err)
	require.Len(t, stints, 1)
	require.NotNil(t, findStint(stints, 1, 1, 1), "played-map player should be rostered")
	require.Nil(t, findStint(stints, 2, 1, 1), "DNP-only player must not be rostered")
}

func TestInferRosterStints_DatesAreMinMaxMatchDate(t *testing.T) {
	db := rosterTx(t)

	mkSeason(t, db, 1, "BO6")
	mkTeam(t, db, 1, "OpTic Texas", "OTX", nil)
	mkTeam(t, db, 2, "Atlanta FaZe", "ATL", nil)
	mkTournament(t, db, 1, 1)
	mkPlayer(t, db, 1, "Dashy")

	early := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	late := time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC)
	sentinel := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

	mkMatch(t, db, 1, 1, 1, 2, early)
	mkMap(t, db, 1, 1, true)
	mkMapStat(t, db, 1, 1, 1, 1)

	mkMatch(t, db, 2, 1, 1, 2, late)
	mkMap(t, db, 2, 1, true)
	mkMapStat(t, db, 2, 1, 1, 1)

	mkMatch(t, db, 3, 1, 1, 2, sentinel) // undated match — ignored for bounds
	mkMap(t, db, 3, 1, true)
	mkMapStat(t, db, 3, 1, 1, 1)

	stints, err := inferRosterStints(db)
	require.NoError(t, err)
	require.Len(t, stints, 1)
	st := findStint(stints, 1, 1, 1)
	require.NotNil(t, st)
	require.True(t, st.StartDate.Valid)
	require.True(t, st.EndDate.Valid)
	require.Equal(t, early.UTC(), st.StartDate.Time.UTC())
	require.Equal(t, late.UTC(), st.EndDate.Time.UTC())
	require.Equal(t, 3, st.MapCount, "all played maps count toward the appearance")
}

func TestInferRosterStints_SeparateHistoricalTeams(t *testing.T) {
	db := rosterTx(t)

	franchiseID := uint(1)
	mkFranchise(t, db, franchiseID, "surge")
	mkSeason(t, db, 1, "MW3")
	mkSeason(t, db, 2, "BO6")
	mkTeam(t, db, 1, "Seattle Surge", "SEA", &franchiseID)
	mkTeam(t, db, 2, "Vancouver Surge", "VAN", &franchiseID)
	mkTournament(t, db, 1, 1)
	mkTournament(t, db, 2, 2)
	mkPlayer(t, db, 1, "Pred")

	mkMatch(t, db, 1, 1, 1, 2, time.Now())
	mkMap(t, db, 1, 1, true)
	mkMapStat(t, db, 1, 1, 1, 1) // Seattle Surge, season1

	mkMatch(t, db, 2, 2, 1, 2, time.Now())
	mkMap(t, db, 2, 1, true)
	mkMapStat(t, db, 2, 1, 1, 2) // Vancouver Surge, season2

	stints, err := inferRosterStints(db)
	require.NoError(t, err)
	require.Len(t, stints, 2)
	seattle := findStint(stints, 1, 1, 1)
	vancouver := findStint(stints, 1, 2, 2)
	require.NotNil(t, seattle, "expected Seattle Surge stint")
	require.NotNil(t, vancouver, "expected Vancouver Surge stint")
	require.NotEqual(t, seattle.TeamID, vancouver.TeamID, "shared franchise must keep distinct team rows")
}
