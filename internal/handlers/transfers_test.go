package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
