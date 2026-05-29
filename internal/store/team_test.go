package store

import (
	"context"
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

var storeTestDB *gorm.DB

func TestMain(m *testing.M) {
	os.Exit(runStoreTests(m))
}

func runStoreTests(m *testing.M) int {
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
		&models.TeamRoster{},
		&models.Tournament{},
		&models.Match{},
		&models.MatchMap{},
		&models.PlayerMapStats{},
	); err != nil {
		log.Println("gorm: automigrate failed:", err)
		return m.Run()
	}

	storeTestDB = db
	return m.Run()
}

func storeTx(t *testing.T) *gorm.DB {
	t.Helper()
	if storeTestDB == nil {
		t.Skip("postgres test container unavailable (is Docker running?)")
	}
	tx := storeTestDB.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { tx.Rollback() }) //nolint:errcheck
	return tx
}

// gamertags is a small helper to extract names for assertions.
func gamertags(players []models.Player) []string {
	out := make([]string, len(players))
	for i, p := range players {
		out[i] = p.Gamertag
	}
	return out
}

func mkSeasonAt(t *testing.T, db *gorm.DB, id uint, code string, start time.Time) {
	t.Helper()
	require.NoError(t, db.Create(&models.Season{
		ID: id, Name: code, GameTitle: code, GameCode: code, StartDate: start,
	}).Error)
}

func mkTeamRow(t *testing.T, db *gorm.DB, id uint, name, abbr string) {
	t.Helper()
	require.NoError(t, db.Create(&models.Team{ID: id, Name: name, Abbreviation: abbr}).Error)
}

func mkPlayerRow(t *testing.T, db *gorm.DB, id uint, tag string) {
	t.Helper()
	require.NoError(t, db.Create(&models.Player{ID: id, Gamertag: tag}).Error)
}

func mkTour(t *testing.T, db *gorm.DB, id, seasonID uint, name string) {
	t.Helper()
	require.NoError(t, db.Create(&models.Tournament{
		ID: id, SeasonID: seasonID, Name: name, Slug: name, TournamentType: "major", StartDate: time.Now(),
	}).Error)
}

func mkMatchRow(t *testing.T, db *gorm.DB, id, tournamentID, team1, team2 uint, date time.Time) {
	t.Helper()
	require.NoError(t, db.Create(&models.Match{
		ID: id, TournamentID: tournamentID, Team1ID: team1, Team2ID: team2, MatchDate: date,
	}).Error)
}

func mkMapRow(t *testing.T, db *gorm.DB, matchID uint, mapNumber int, played bool) {
	t.Helper()
	mm := models.MatchMap{MatchID: matchID, MapNumber: mapNumber, Played: played}
	require.NoError(t, db.Create(&mm).Error)
	if !played { // gorm default:true is gone, but force it defensively in fixtures
		require.NoError(t, db.Model(&models.MatchMap{}).Where("id = ?", mm.ID).
			Update("played", false).Error)
	}
}

func mkPMSRow(t *testing.T, db *gorm.DB, matchID uint, mapNumber int, playerID, teamID uint) {
	t.Helper()
	require.NoError(t, db.Create(&models.PlayerMapStats{
		MatchID: matchID, MapNumber: mapNumber, PlayerID: playerID, TeamID: teamID,
	}).Error)
}

func TestGetPlayers_NoSeasonDefaultsToLatestByStartDate(t *testing.T) {
	db := storeTx(t)
	ctx := context.Background()

	require.NoError(t, db.Create(&models.Season{ID: 1, Name: "BO6", GameTitle: "BO6", GameCode: "BO6",
		StartDate: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)}).Error) // newest, low id
	require.NoError(t, db.Create(&models.Season{ID: 2, Name: "MW3", GameTitle: "MW3", GameCode: "MW3",
		StartDate: time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)}).Error) // older, high id
	require.NoError(t, db.Create(&models.Team{ID: 1, Name: "OpTic Texas", Abbreviation: "OTX"}).Error)
	require.NoError(t, db.Create(&models.Player{ID: 1, Gamertag: "OldGuy"}).Error)
	require.NoError(t, db.Create(&models.Player{ID: 2, Gamertag: "NewGuy"}).Error)

	require.NoError(t, db.Create(&models.TeamRoster{TeamID: 1, PlayerID: 1, SeasonID: 2}).Error) // older season
	require.NoError(t, db.Create(&models.TeamRoster{TeamID: 1, PlayerID: 2, SeasonID: 1}).Error) // newer season

	st := NewGormTeamStore(db)
	players, err := st.GetPlayers(ctx, 1, "")
	require.NoError(t, err)
	require.Equal(t, []string{"NewGuy"}, gamertags(players),
		"no season_id should return only the latest season (BO6) roster")
}

// An explicit season_id returns that season's roster regardless of recency.
func TestGetPlayers_ExplicitSeasonOverridesDefault(t *testing.T) {
	db := storeTx(t)
	ctx := context.Background()

	require.NoError(t, db.Create(&models.Season{ID: 1, Name: "BO6", GameTitle: "BO6", GameCode: "BO6",
		StartDate: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)}).Error)
	require.NoError(t, db.Create(&models.Season{ID: 2, Name: "MW3", GameTitle: "MW3", GameCode: "MW3",
		StartDate: time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)}).Error)
	require.NoError(t, db.Create(&models.Team{ID: 1, Name: "OpTic Texas", Abbreviation: "OTX"}).Error)
	require.NoError(t, db.Create(&models.Player{ID: 1, Gamertag: "OldGuy"}).Error)
	require.NoError(t, db.Create(&models.Player{ID: 2, Gamertag: "NewGuy"}).Error)
	require.NoError(t, db.Create(&models.TeamRoster{TeamID: 1, PlayerID: 1, SeasonID: 2}).Error)
	require.NoError(t, db.Create(&models.TeamRoster{TeamID: 1, PlayerID: 2, SeasonID: 1}).Error)

	st := NewGormTeamStore(db)
	players, err := st.GetPlayers(ctx, 1, "2") // ask for the older season explicitly
	require.NoError(t, err)
	require.Equal(t, []string{"OldGuy"}, gamertags(players))
}

// Duplicate roster rows for the same player must not produce duplicate players.
func TestGetPlayers_DistinctPlayers(t *testing.T) {
	db := storeTx(t)
	ctx := context.Background()

	require.NoError(t, db.Create(&models.Season{ID: 1, Name: "BO6", GameTitle: "BO6", GameCode: "BO6",
		StartDate: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)}).Error)
	require.NoError(t, db.Create(&models.Team{ID: 1, Name: "OpTic Texas", Abbreviation: "OTX"}).Error)
	require.NoError(t, db.Create(&models.Player{ID: 1, Gamertag: "Dashy"}).Error)
	require.NoError(t, db.Create(&models.TeamRoster{TeamID: 1, PlayerID: 1, SeasonID: 1}).Error)
	require.NoError(t, db.Create(&models.TeamRoster{TeamID: 1, PlayerID: 1, SeasonID: 1}).Error) // dup

	st := NewGormTeamStore(db)
	players, err := st.GetPlayers(ctx, 1, "")
	require.NoError(t, err)
	require.Equal(t, []string{"Dashy"}, gamertags(players))
}

func TestGetLatestMatchRoster_SeasonReturnsLatestMatchNotUnion(t *testing.T) {
	db := storeTx(t)
	ctx := context.Background()

	mkSeasonAt(t, db, 2, "MW3", time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC))
	mkTeamRow(t, db, 1, "Boston Breach", "BOS")
	mkTeamRow(t, db, 9, "Opponent", "OPP")
	mkTour(t, db, 1, 2, "MW3 Major")
	for id, tag := range map[uint]string{1: "Snoopy", 2: "Owakening", 3: "Purj", 4: "Beans", 5: "Cammy", 6: "SeanyBench"} {
		mkPlayerRow(t, db, id, tag)
		require.NoError(t, db.Create(&models.TeamRoster{TeamID: 1, PlayerID: id, SeasonID: 2}).Error) // season union
	}

	// earlier match: Snoopy/Owakening/Purj/Beans
	mkMatchRow(t, db, 100, 1, 1, 9, time.Date(2024, 4, 13, 0, 0, 0, 0, time.UTC))
	mkMapRow(t, db, 100, 1, true)
	for _, pid := range []uint{1, 2, 3, 4} {
		mkPMSRow(t, db, 100, 1, pid, 1)
	}
	// latest match (EWC): Beans out, Cammy in
	mkMatchRow(t, db, 200, 1, 1, 9, time.Date(2024, 8, 17, 0, 0, 0, 0, time.UTC))
	mkMapRow(t, db, 200, 1, true)
	for _, pid := range []uint{1, 2, 3, 5} {
		mkPMSRow(t, db, 200, 1, pid, 1)
	}

	st := NewGormTeamStore(db)
	players, err := st.GetLatestMatchRoster(ctx, 1, "2")
	require.NoError(t, err)
	require.Equal(t, []string{"Cammy", "Owakening", "Purj", "Snoopy"}, gamertags(players),
		"selected season should return the latest-match 4, not the union of 6")
}

func TestGetLatestMatchRoster_NoSeasonResolvesLatestSeasonByMatchDate(t *testing.T) {
	db := storeTx(t)
	ctx := context.Background()

	mkSeasonAt(t, db, 1, "BO6", time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)) // newest, low id
	mkSeasonAt(t, db, 2, "MW3", time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)) // older, high id
	mkTeamRow(t, db, 1, "OpTic Texas", "OTX")
	mkTeamRow(t, db, 9, "Opponent", "OPP")
	mkTour(t, db, 1, 1, "BO6 Major")
	mkTour(t, db, 2, 2, "MW3 Major")
	for id, tag := range map[uint]string{1: "Mw3A", 2: "Mw3B", 3: "Bo6A", 4: "Bo6B"} {
		mkPlayerRow(t, db, id, tag)
	}
	mkMatchRow(t, db, 100, 2, 1, 9, time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)) // MW3 (older)
	mkMapRow(t, db, 100, 1, true)
	mkPMSRow(t, db, 100, 1, 1, 1)
	mkPMSRow(t, db, 100, 1, 2, 1)
	mkMatchRow(t, db, 200, 1, 1, 9, time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)) // BO6 (newer)
	mkMapRow(t, db, 200, 1, true)
	mkPMSRow(t, db, 200, 1, 3, 1)
	mkPMSRow(t, db, 200, 1, 4, 1)

	st := NewGormTeamStore(db)
	players, err := st.GetLatestMatchRoster(ctx, 1, "")
	require.NoError(t, err)
	require.Equal(t, []string{"Bo6A", "Bo6B"}, gamertags(players),
		"no season_id should resolve latest season (BO6) by match date, then its latest-match roster")
}

// DNP maps must not contribute players to the latest-match roster.
// Regression for the "Carolina Royal Ravens MW3 era shows the BO6 roster" bug.
// The root cause was a single team row shared across MW3 and BO6; once each
// game-era is its own team row, the no-season default must scope each era's
// roster to that era's own season — an MW3-era team_id never resolves to BO6.
func TestGetLatestMatchRoster_PerEraTeamsScopeToOwnSeason(t *testing.T) {
	db := storeTx(t)
	ctx := context.Background()

	mkSeasonAt(t, db, 1, "BO6", time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)) // newer
	mkSeasonAt(t, db, 2, "MW3", time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)) // older
	// Two distinct era rows for the same franchise name.
	mkTeamRow(t, db, 1, "Carolina Royal Ravens", "CRR") // MW3 era
	mkTeamRow(t, db, 2, "Carolina Royal Ravens", "CRR") // BO6 era
	mkTeamRow(t, db, 9, "Opponent", "OPP")
	mkTour(t, db, 1, 1, "BO6 Major")
	mkTour(t, db, 2, 2, "MW3 Major")
	for id, tag := range map[uint]string{1: "Gwinn", 2: "Owakening", 3: "Lyly", 4: "Vortex"} {
		mkPlayerRow(t, db, id, tag)
	}

	// MW3-era team (id=1) plays only in the MW3 season.
	mkMatchRow(t, db, 100, 2, 1, 9, time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC))
	mkMapRow(t, db, 100, 1, true)
	mkPMSRow(t, db, 100, 1, 1, 1)
	mkPMSRow(t, db, 100, 1, 2, 1)

	// BO6-era team (id=2) plays only in the BO6 season.
	mkMatchRow(t, db, 200, 1, 2, 9, time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC))
	mkMapRow(t, db, 200, 1, true)
	mkPMSRow(t, db, 200, 1, 3, 2)
	mkPMSRow(t, db, 200, 1, 4, 2)

	st := NewGormTeamStore(db)

	mw3, err := st.GetLatestMatchRoster(ctx, 1, "")
	require.NoError(t, err)
	require.Equal(t, []string{"Gwinn", "Owakening"}, gamertags(mw3),
		"MW3-era team must return its MW3 roster, not the newer BO6 season")

	bo6, err := st.GetLatestMatchRoster(ctx, 2, "")
	require.NoError(t, err)
	require.Equal(t, []string{"Lyly", "Vortex"}, gamertags(bo6),
		"BO6-era team must return its own BO6 roster")
}

func TestGetLatestMatchRoster_DNPMapsExcluded(t *testing.T) {
	db := storeTx(t)
	ctx := context.Background()

	mkSeasonAt(t, db, 1, "BO6", time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC))
	mkTeamRow(t, db, 1, "OpTic Texas", "OTX")
	mkTeamRow(t, db, 9, "Opponent", "OPP")
	mkTour(t, db, 1, 1, "BO6 Major")
	for id, tag := range map[uint]string{1: "Dashy", 2: "Kenny", 3: "Pred", 4: "Shotzzy", 5: "GhostSub"} {
		mkPlayerRow(t, db, id, tag)
	}
	mkMatchRow(t, db, 100, 1, 1, 9, time.Now())
	mkMapRow(t, db, 100, 1, true)  // played
	mkMapRow(t, db, 100, 2, false) // DNP
	for _, pid := range []uint{1, 2, 3, 4} {
		mkPMSRow(t, db, 100, 1, pid, 1)
	}
	mkPMSRow(t, db, 100, 2, 5, 1) // stray stat on the DNP map

	st := NewGormTeamStore(db)
	players, err := st.GetLatestMatchRoster(ctx, 1, "")
	require.NoError(t, err)
	require.Equal(t, []string{"Dashy", "Kenny", "Pred", "Shotzzy"}, gamertags(players),
		"DNP-map player must be excluded from the latest-match roster")
}

func TestGetLatestMatchRoster_FallsBackToStintsWhenNoMatchData(t *testing.T) {
	db := storeTx(t)
	ctx := context.Background()

	mkSeasonAt(t, db, 1, "BO6", time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC))
	mkTeamRow(t, db, 1, "OpTic Texas", "OTX")
	mkPlayerRow(t, db, 1, "Dashy")
	mkPlayerRow(t, db, 2, "Kenny")
	require.NoError(t, db.Create(&models.TeamRoster{TeamID: 1, PlayerID: 1, SeasonID: 1}).Error)
	require.NoError(t, db.Create(&models.TeamRoster{TeamID: 1, PlayerID: 2, SeasonID: 1}).Error)

	st := NewGormTeamStore(db)
	players, err := st.GetLatestMatchRoster(ctx, 1, "1")
	require.NoError(t, err)
	require.Equal(t, []string{"Dashy", "Kenny"}, gamertags(players),
		"no match data should fall back to the season's roster stints")
}
