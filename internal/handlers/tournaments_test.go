package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTournaments_Empty(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetTournaments(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var ts []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &ts))
	assert.Empty(t, ts)
}

func TestGetTournaments_ReturnsSeeded(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetTournaments(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var ts []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &ts))
	assert.Len(t, ts, 1)
}

func TestGetTournamentBySlug_NotFound(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "slug", Value: "nope"}}, "")
	h.GetTournamentBySlug(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTournamentBySlug_Success(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "slug", Value: "cdl-major-1-2025"}}, "")
	h.GetTournamentBySlug(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var d struct {
		Tournament map[string]any `json:"tournament"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &d))
	assert.Equal(t, "cdl-major-1-2025", d.Tournament["slug"])
}

func TestGetTournament_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "abc"}}, "")
	h.GetTournament(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid tournament ID", errBody(t, w.Body.Bytes()))
}

func TestGetTournament_NotFound(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "999"}}, "")
	h.GetTournament(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTournament_Success(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetTournament(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var tour map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &tour))
	assert.Equal(t, "CDL Major 1 2025", tour["name"])
}

// --- GET /tournaments/:id/bracket ------------------------------------------

func TestGetTournamentBracket_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "x"}}, "")
	h.GetTournamentBracket(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid tournament ID", errBody(t, w.Body.Bytes()))
}

func TestGetTournamentBracket_NotFound(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "999"}}, "")
	h.GetTournamentBracket(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTournamentBracket_SuccessEmpty(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetTournamentBracket(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var br map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &br))
	assert.Equal(t, float64(0), br["total_matches"])
}

// --- GET /tournaments/:id/matches ------------------------------------------

func TestGetTournamentMatches_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "bad"}}, "")
	h.GetTournamentMatches(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid tournament ID", errBody(t, w.Body.Bytes()))
}

func TestGetTournamentMatches_Empty(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "999"}}, "")
	h.GetTournamentMatches(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var ms []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &ms))
	assert.Empty(t, ms)
}

func TestGetTournamentMatches_ReturnsSeeded(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 1)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetTournamentMatches(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var ms []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &ms))
	assert.Len(t, ms, 1)
}

// --- GET /tournaments/:id/teams --------------------------------------------

func TestGetTournamentTeams_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "nope"}}, "")
	h.GetTournamentTeams(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid tournament ID", errBody(t, w.Body.Bytes()))
}

func TestGetTournamentTeams_Empty(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "999"}}, "")
	h.GetTournamentTeams(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var teams []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &teams))
	assert.Empty(t, teams)
}

func TestGetTournamentStats_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "xyz"}}, "")
	h.GetTournamentStats(c)
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

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "99"}}, "")
	h.GetTournamentStats(c)

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

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetTournamentStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Len(t, body, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}
