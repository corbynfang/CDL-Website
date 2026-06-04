package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type statsEnvelope struct {
	Timestamp int64            `json:"timestamp"`
	Players   []map[string]any `json:"players"`
	Count     int              `json:"count"`
}

var statPlayerCols = []string{
	"player_id", "gamertag", "avatar_url", "team_abbr",
	"season_kills", "season_deaths", "season_assists",
}

func TestGetTopKDPlayers_ResponseShape(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).WillReturnRows(
		sqlmock.NewRows(statPlayerCols).AddRow(1, "Scump", "", "OTX", 500, 400, 50))

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetTopKDPlayers(c)

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

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetTopKDPlayers(c)

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

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetTopKDPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(0), body.Players[0]["season_kd"],
		"0 deaths → season_kd = 0 (not Inf, not crash)")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopKDPlayers_EmptyResults(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).
		WillReturnRows(sqlmock.NewRows(statPlayerCols))

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetTopKDPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, 0, body.Count)
	assert.NotNil(t, body.Players, "players must be [] not null")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPlayersKDStats_ResponseShape(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).WillReturnRows(
		sqlmock.NewRows(statPlayerCols).AddRow(5, "Shotzzy", "", "OTX", 200, 150, 20))

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetAllPlayersKDStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Len(t, body.Players, 1)

	p := body.Players[0]
	assert.Contains(t, p, "season_kd", "season_kd required")
	assert.Contains(t, p, "season_kd_plus_minus",
		"season_kd_plus_minus is exclusive to GetAllPlayersKDStats")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPlayersKDStats_KDPlusMinus(t *testing.T) {
	mock := setupMockDB(t)
	// 150/100 = 1.50 → plus_minus = 0.50
	mock.ExpectQuery(`SELECT.*FROM player_tournament_stats`).WillReturnRows(
		sqlmock.NewRows(statPlayerCols).AddRow(6, "Simp", "", "ATL", 150, 100, 5))

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetAllPlayersKDStats(c)

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

	h := newTestHandler(t)
	c, w := newCtx(nil, "season_id=3")
	h.GetAllPlayersKDStats(c)

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

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetAllPlayersKDStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body statsEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, 0, body.Count)
	assert.NotNil(t, body.Players, "players must be [] not null")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTeamStats_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "x"}}, "")
	h.GetTeamStats(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid team ID", errBody(t, w.Body.Bytes()))
}

func TestGetTeamStats_Empty(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "999"}}, "")
	h.GetTeamStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var stats []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &stats))
	assert.Empty(t, stats)
}

func TestGetTeamStats_ReturnsRows(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t) // season 1, tournament 1, teams 1 & 2
	placement := 1
	require.NoError(t, database.DB.Create(&models.TeamTournamentStats{
		TournamentID: 1, TeamID: 1, Placement: &placement,
		MatchesPlayed: 5, MatchesWon: 4, MatchesLost: 1,
	}).Error)

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetTeamStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var stats []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &stats))
	require.Len(t, stats, 1)
	assert.Equal(t, float64(5), stats[0]["matches_played"])
}

func TestGetPlayerStats_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "nope"}}, "")
	h.GetPlayerStats(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid player ID", errBody(t, w.Body.Bytes()))
}

func TestGetPlayerStats_Empty(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetPlayerStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var stats []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &stats))
	assert.Empty(t, stats)
}
