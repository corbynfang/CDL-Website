package handlers

// stats_test.go — API contract tests for GetTopKDPlayers and GetAllPlayersKDStats.
// Each test checks required fields and computed values, not fragile row counts.

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// statsEnvelope is the shape both stats handlers return.
type statsEnvelope struct {
	Timestamp int64            `json:"timestamp"`
	Players   []map[string]any `json:"players"`
	Count     int              `json:"count"`
}

// statPlayerCols are the columns returned by the aggregated stats query.
var statPlayerCols = []string{
	"player_id", "gamertag", "avatar_url", "team_abbr",
	"season_kills", "season_deaths", "season_assists",
}

// ── GetTopKDPlayers ───────────────────────────────────────────────────────────

func TestGetTopKDPlayers_ResponseShape(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).WillReturnRows(
		sqlmock.NewRows(statPlayerCols).AddRow(1, "Scump", "", "OTX", 500, 400, 50))

	c, w := newCtx(nil, "")
	GetTopKDPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.NotZero(t, body.Timestamp, "timestamp required")
	assert.Equal(t, 1, body.Count, "count must equal len(players)")
	require.Len(t, body.Players, 1)

	p := body.Players[0]
	for _, field := range []string{"player_id", "gamertag", "season_kills", "season_deaths", "season_kd"} {
		assert.Contains(t, p, field, "GetTopKDPlayers player must contain %s", field)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopKDPlayers_KDCalculated(t *testing.T) {
	mock := setupMockDB(t)
	// 100 kills / 50 deaths = 2.00
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).WillReturnRows(
		sqlmock.NewRows(statPlayerCols).AddRow(2, "Cellium", "", "ATL", 100, 50, 10))

	c, w := newCtx(nil, "")
	GetTopKDPlayers(c)

	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Len(t, body.Players, 1)
	assert.InDelta(t, 2.0, body.Players[0]["season_kd"], 0.001,
		"season_kd must be kills/deaths = 100/50 = 2.00")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopKDPlayers_ZeroDeathsReturnsZeroNotPanic(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).WillReturnRows(
		sqlmock.NewRows(statPlayerCols).AddRow(3, "Ghost", "", "ATL", 100, 0, 0))

	c, w := newCtx(nil, "")
	GetTopKDPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	// calculateKD returns 0 when deaths == 0 — documents the current contract (not Inf, not panic)
	assert.Equal(t, float64(0), body.Players[0]["season_kd"],
		"0 deaths → season_kd = 0 (not Inf, not crash)")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopKDPlayers_EmptyResults(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).
		WillReturnRows(sqlmock.NewRows(statPlayerCols))

	c, w := newCtx(nil, "")
	GetTopKDPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, 0, body.Count)
	assert.NotNil(t, body.Players, "players must be [] not null")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetAllPlayersKDStats ──────────────────────────────────────────────────────

func TestGetAllPlayersKDStats_ResponseShape(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).WillReturnRows(
		sqlmock.NewRows(statPlayerCols).AddRow(5, "Shotzzy", "", "OTX", 200, 150, 20))

	c, w := newCtx(nil, "")
	GetAllPlayersKDStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Len(t, body.Players, 1)

	p := body.Players[0]
	assert.Contains(t, p, "season_kd", "season_kd required")
	// GetAllPlayersKDStats exposes season_kd_plus_minus; GetTopKDPlayers does not.
	assert.Contains(t, p, "season_kd_plus_minus",
		"season_kd_plus_minus is exclusive to GetAllPlayersKDStats")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPlayersKDStats_KDPlusMinus(t *testing.T) {
	mock := setupMockDB(t)
	// 150/100 = 1.50 → plus_minus = 0.50
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).WillReturnRows(
		sqlmock.NewRows(statPlayerCols).AddRow(6, "Simp", "", "ATL", 150, 100, 5))

	c, w := newCtx(nil, "")
	GetAllPlayersKDStats(c)

	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Len(t, body.Players, 1)
	assert.InDelta(t, 0.5, body.Players[0]["season_kd_plus_minus"], 0.001,
		"season_kd_plus_minus = kd - 1.0 = 1.5 - 1.0 = 0.5")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPlayersKDStats_SeasonFilter(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).WillReturnRows(
		sqlmock.NewRows(statPlayerCols).AddRow(7, "aBeZy", "", "ATL", 80, 60, 8))

	// ?season_id=3 should be accepted without error
	c, w := newCtx(nil, "season_id=3")
	GetAllPlayersKDStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, 1, body.Count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPlayersKDStats_EmptyResults(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).
		WillReturnRows(sqlmock.NewRows(statPlayerCols))

	c, w := newCtx(nil, "")
	GetAllPlayersKDStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, 0, body.Count)
	assert.NotNil(t, body.Players, "players must be [] not null")
	assert.NoError(t, mock.ExpectationsWereMet())
}
