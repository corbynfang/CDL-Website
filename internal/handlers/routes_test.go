package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	RegisterRoutes(api, h)
	r.NoRoute(func(c *gin.Context) { c.JSON(http.StatusNotFound, gin.H{"error": "not found"}) })
	return r
}

func TestRegisterRoutes_AllRoutesPresent(t *testing.T) {
	r := newTestRouter(New(nil)) // registering routes needs no DB
	got := map[string]bool{}
	for _, ri := range r.Routes() {
		got[ri.Method+" "+ri.Path] = true
	}

	want := []string{
		"GET /api/v1/seasons",
		"GET /api/v1/seasons/:id",
		"GET /api/v1/seasons/active",
		"GET /api/v1/teams",
		"GET /api/v1/teams/:id",
		"GET /api/v1/teams/:id/players",
		"GET /api/v1/teams/:id/stats",
		"GET /api/v1/players",
		"GET /api/v1/players/:id",
		"GET /api/v1/players/:id/stats",
		"GET /api/v1/players/:id/kd",
		"GET /api/v1/players/:id/matches",
		"GET /api/v1/players/:id/franchise-career",
		"GET /api/v1/players/top-kd",
		"GET /api/v1/stats/all-kd-by-tournament",
		"GET /api/v1/matches/:id",
		"GET /api/v1/franchises",
		"GET /api/v1/franchises/:key",
		"GET /api/v1/tournaments",
		"GET /api/v1/tournaments/slug/:slug",
		"GET /api/v1/tournaments/:id",
		"GET /api/v1/tournaments/:id/bracket",
		"GET /api/v1/tournaments/:id/matches",
		"GET /api/v1/tournaments/:id/teams",
		"GET /api/v1/tournaments/:id/stats",
		"GET /api/v1/transfers",
	}

	for _, w := range want {
		assert.True(t, got[w], "route not registered: %s", w)
	}
	assert.Len(t, r.Routes(), len(want), "unexpected number of routes registered")
}

func TestRouter_UnknownRouteReturns404(t *testing.T) {
	r := newTestRouter(New(nil))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/does-not-exist", nil))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouter_E2E_ListSeasons(t *testing.T) {
	setupPGTx(t)
	require.NoError(t, database.DB.Create(&models.Season{
		ID: 1, Name: "BO6 2025", GameTitle: "Black Ops 6", GameCode: "BO6", StartDate: time.Now(),
	}).Error)

	r := newTestRouter(New(database.DB))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/seasons", nil))

	assert.Equal(t, http.StatusOK, w.Code)
	var body any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body), "response must be valid JSON")
}
