package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSeasons_Empty(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetSeasons(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var seasons []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &seasons))
	assert.Empty(t, seasons)
}

func TestGetSeasons_ReturnsSeeded(t *testing.T) {
	setupPGTx(t)
	require.NoError(t, database.DB.Create(&models.Season{ID: 1, Name: "BO6", GameTitle: "Black Ops 6", GameCode: "BO6", StartDate: time.Now()}).Error)
	require.NoError(t, database.DB.Create(&models.Season{ID: 2, Name: "MW3", GameTitle: "Modern Warfare III", GameCode: "MW3", StartDate: time.Now()}).Error)

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetSeasons(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var seasons []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &seasons))
	assert.Len(t, seasons, 2)
}

func TestGetSeason_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "abc"}}, "")
	h.GetSeason(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid season ID", errBody(t, w.Body.Bytes()))
}

func TestGetSeason_NotFound(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "999"}}, "")
	h.GetSeason(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Season not found", errBody(t, w.Body.Bytes()))
}

func TestGetSeason_Success(t *testing.T) {
	setupPGTx(t)
	require.NoError(t, database.DB.Create(&models.Season{ID: 1, Name: "BO6", GameTitle: "Black Ops 6", GameCode: "BO6", StartDate: time.Now()}).Error)

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetSeason(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var season map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &season))
	assert.Equal(t, "BO6", season["name"])
	assert.Equal(t, "BO6", season["game_code"])
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

func TestGetPlayerFranchiseCareer_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "bad"}}, "")
	h.GetPlayerFranchiseCareer(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid player ID", errBody(t, w.Body.Bytes()))
}

func TestGetPlayerFranchiseCareer_NotFound(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "999"}}, "")
	h.GetPlayerFranchiseCareer(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPlayerFranchiseCareer_EmptyForExistingPlayer(t *testing.T) {
	setupPGTx(t)
	require.NoError(t, database.DB.Create(&models.Player{ID: 1, Gamertag: "Scump"}).Error)

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetPlayerFranchiseCareer(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var career map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &career))
	assert.Equal(t, "Scump", career["gamertag"])
}
