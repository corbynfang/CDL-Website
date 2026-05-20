package handlers

// matches_test.go — API contract tests for GetMatch and GetTournamentStats.
// GetMatch tests use a real Postgres testcontainer (via setupPGTx) so they are
// not brittle to GORM preload query ordering.
// GetTournamentStats tests use sqlmock because they involve no nested preloads.

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// matchDetailBody mirrors the top-level shape GetMatch returns.
type matchDetailBody struct {
	Match map[string]any   `json:"match"`
	Maps  []map[string]any `json:"maps"`
}

// ── GetMatch ──────────────────────────────────────────────────────────────────

func TestGetMatch_InvalidID(t *testing.T) {
	c, w := newCtx(gin.Params{{Key: "id", Value: "bad"}}, "")
	GetMatch(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMatch_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "matches"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	c, w := newCtx(gin.Params{{Key: "id", Value: "999"}}, "")
	GetMatch(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetMatch_ResponseShape verifies all required top-level match fields are present.
func TestGetMatch_ResponseShape(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 42)

	c, w := newCtx(gin.Params{{Key: "id", Value: "42"}}, "")
	GetMatch(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body matchDetailBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	m := body.Match
	require.NotNil(t, m, "response must contain 'match' key")
	for _, field := range []string{
		"id", "tournament_name", "team1_id", "team2_id",
		"team1_name", "team2_name", "team1_abbr", "team2_abbr",
		"team1_score", "team2_score", "match_date",
	} {
		assert.Contains(t, m, field, "match.%s required", field)
	}
	assert.NotNil(t, body.Maps, "response must contain 'maps' key")
}

// TestGetMatch_MapShape verifies required fields inside a map entry.
func TestGetMatch_MapShape(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 5)
	require.NoError(t, database.DB.Create(&database.MatchMap{
		MatchID: 5, MapNumber: 1, MapName: "Skyline", Mode: "Hardpoint",
		Score1: 250, Score2: 200, Played: true, DurationSec: 480,
	}).Error)

	c, w := newCtx(gin.Params{{Key: "id", Value: "5"}}, "")
	GetMatch(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body matchDetailBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	require.Len(t, body.Maps, 1, "one map should be returned")
	mp := body.Maps[0]
	for _, field := range []string{
		"map_number", "map_name", "mode",
		"score_1", "score_2", "played",
		"team1_stats", "team2_stats",
	} {
		assert.Contains(t, mp, field, "map.%s required", field)
	}
}

// TestGetMatch_EmptyStatArrays verifies that team1_stats and team2_stats are
// always arrays, never null, even when no player stats exist for a map.
func TestGetMatch_EmptyStatArrays(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 7)
	require.NoError(t, database.DB.Create(&database.MatchMap{
		MatchID: 7, MapNumber: 1, MapName: "Rewind", Mode: "Search and Destroy",
		Score1: 6, Score2: 4, Played: true,
	}).Error)

	c, w := newCtx(gin.Params{{Key: "id", Value: "7"}}, "")
	GetMatch(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body matchDetailBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	require.Len(t, body.Maps, 1)
	mp := body.Maps[0]
	team1Stats, _ := mp["team1_stats"].([]any)
	team2Stats, _ := mp["team2_stats"].([]any)
	assert.NotNil(t, team1Stats, "team1_stats must be an array, not null")
	assert.NotNil(t, team2Stats, "team2_stats must be an array, not null")
	assert.Len(t, team1Stats, 0)
	assert.Len(t, team2Stats, 0)
}

// ── GetTournamentStats ────────────────────────────────────────────────────────

func TestGetTournamentStats_InvalidID(t *testing.T) {
	c, w := newCtx(gin.Params{{Key: "id", Value: "xyz"}}, "")
	GetTournamentStats(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTournamentStats_ResponseShape(t *testing.T) {
	mock := setupMockDB(t)
	now := time.Now()

	statCols := []string{
		"id", "player_id", "team_id", "tournament_id",
		"total_kills", "total_deaths", "total_assists", "total_damage",
		"kd_ratio", "kda_ratio", "overall_maps", "overall_plus_minus",
	}
	mock.ExpectQuery(`SELECT \* FROM "player_tournament_stats"`).WillReturnRows(
		sqlmock.NewRows(statCols).
			AddRow(1, 10, 1, 99, 120, 80, 15, 60000, 1.50, 1.69, 10, 40))

	mock.ExpectQuery(`SELECT \* FROM "players"`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "gamertag", "is_active", "created_at", "updated_at"}).
			AddRow(10, "Shotzzy", true, now, now))

	mock.ExpectQuery(`SELECT \* FROM "teams"`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "abbreviation", "is_active", "is_cdl_franchise", "created_at", "updated_at"}).
			AddRow(1, "OpTic Texas", "OTX", true, true, now, now))

	c, w := newCtx(gin.Params{{Key: "id", Value: "99"}}, "")
	GetTournamentStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Len(t, body, 1)

	s := body[0]
	for _, field := range []string{"player_id", "total_kills", "total_deaths", "kd_ratio", "overall_maps"} {
		assert.Contains(t, s, field, "tournament stat must contain %s", field)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTournamentStats_EmptyResults(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery(`SELECT \* FROM "player_tournament_stats"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "player_id"}))

	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	GetTournamentStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Len(t, body, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}
