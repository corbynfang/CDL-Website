package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

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
