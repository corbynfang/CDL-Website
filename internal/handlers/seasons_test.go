package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestGetActiveSeason_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "seasons"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetActiveSeason(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "No active season found", body["error"])
	assert.NoError(t, mock.ExpectationsWereMet())
}
