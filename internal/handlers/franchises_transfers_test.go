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
	assert.Len(t, detail.Eras, 2, "both Surge eras share the franchise")
}

func TestGetTransfers_Empty(t *testing.T) {
	setupPGTx(t)
	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetTransfers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var env struct {
		Transfers []map[string]any `json:"transfers"`
		Count     int              `json:"count"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &env))
	assert.Equal(t, 0, env.Count)
	assert.Empty(t, env.Transfers)
}

func TestGetTransfers_ReturnsSeeded(t *testing.T) {
	setupPGTx(t)
	require.NoError(t, database.DB.Create(&models.Player{ID: 1, Gamertag: "Vivid"}).Error)
	require.NoError(t, database.DB.Create(&models.PlayerTransfer{
		ID: 1, PlayerID: 1, TransferDate: time.Now(), TransferType: "join",
		GameCode: "BO6", Season: "BO6",
	}).Error)

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetTransfers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var env struct {
		Transfers []map[string]any `json:"transfers"`
		Count     int              `json:"count"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &env))
	assert.Equal(t, 1, env.Count)
	require.Len(t, env.Transfers, 1)
}

func TestGetTransfers_PlayerFilter(t *testing.T) {
	setupPGTx(t)
	require.NoError(t, database.DB.Create(&models.Player{ID: 1, Gamertag: "Vivid"}).Error)
	require.NoError(t, database.DB.Create(&models.Player{ID: 2, Gamertag: "Scump"}).Error)
	require.NoError(t, database.DB.Create(&models.PlayerTransfer{ID: 1, PlayerID: 1, TransferDate: time.Now()}).Error)
	require.NoError(t, database.DB.Create(&models.PlayerTransfer{ID: 2, PlayerID: 2, TransferDate: time.Now()}).Error)

	h := newTestHandler(t)
	c, w := newCtx(nil, "player_id=1")
	h.GetTransfers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var env struct {
		Count int `json:"count"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &env))
	assert.Equal(t, 1, env.Count, "player_id filter should narrow to one transfer")
}
