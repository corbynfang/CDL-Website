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

func TestGetFranchises_Empty(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetFranchises(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var fr []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &fr))
	assert.Empty(t, fr)
}

func TestGetFranchises_ReturnsSeeded(t *testing.T) {
	setupPGTx(t)
	require.NoError(t, database.DB.Create(&models.Franchise{ID: 1, FranchiseKey: "optic", Name: "OpTic"}).Error)

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetFranchises(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var fr []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &fr))
	require.Len(t, fr, 1)
	assert.Equal(t, "optic", fr[0]["franchise_key"])
}

func TestGetFranchise_MissingKey(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "key", Value: ""}}, "")
	h.GetFranchise(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Missing franchise key", errBody(t, w.Body.Bytes()))
}

func TestGetFranchise_NotFound(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "key", Value: "nope"}}, "")
	h.GetFranchise(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetFranchise_SuccessWithEras(t *testing.T) {
	setupPGTx(t)
	fid := uint(1)
	require.NoError(t, database.DB.Create(&models.Franchise{ID: fid, FranchiseKey: "surge", Name: "Surge"}).Error)
	require.NoError(t, database.DB.Create(&models.Team{ID: 1, Name: "Seattle Surge", Abbreviation: "SEA", FranchiseID: &fid}).Error)
	require.NoError(t, database.DB.Create(&models.Team{ID: 2, Name: "Vancouver Surge", Abbreviation: "VAN", FranchiseID: &fid}).Error)
	require.NoError(t, database.DB.Create(&models.Team{ID: 3, Name: "Vancouver Surge MW3", Abbreviation: "VAN", FranchiseID: &fid}).Error)
	pgSeason(t)
	pgTournament(t)
	require.NoError(t, database.DB.Create(&models.Player{ID: 1, Gamertag: "Sib"}).Error)
	require.NoError(t, database.DB.Create(&models.Player{ID: 2, Gamertag: "Pred"}).Error)
	require.NoError(t, database.DB.Create(&models.Match{ID: 1, TournamentID: 1, Team1ID: 1, Team2ID: 2, MatchDate: time.Now()}).Error)
	require.NoError(t, database.DB.Create(&models.PlayerMatchStats{ID: 1, MatchID: 1, PlayerID: 1, TeamID: 1}).Error)
	require.NoError(t, database.DB.Create(&models.PlayerMatchStats{ID: 2, MatchID: 1, PlayerID: 2, TeamID: 2}).Error)

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "key", Value: "surge"}}, "")
	h.GetFranchise(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var detail struct {
		Franchise map[string]any   `json:"franchise"`
		Eras      []map[string]any `json:"eras"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &detail))
	assert.Equal(t, "surge", detail.Franchise["franchise_key"])

	names := make([]string, len(detail.Eras))
	for i, e := range detail.Eras {
		names[i] = e["name"].(string)
	}
	assert.ElementsMatch(t, []string{"Seattle Surge", "Vancouver Surge"}, names,
		"data-bearing eras are returned; the empty MW3 era is filtered out")
}
