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

type paginatedBody struct {
	Data       []map[string]any `json:"data"`
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	} `json:"pagination"`
}

func TestGetPlayers_DefaultPagination(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "players"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	playerRows := sqlmock.NewRows([]string{"id", "gamertag", "first_name", "last_name",
		"country", "role", "is_active", "liquipedia_url", "twitter_handle", "avatar_url"}).
		AddRow(1, "Scump", "Seth", "Abner", "US", "flex", true, "", "", "").
		AddRow(2, "Crimsix", "Ian", "Porter", "US", "AR", true, "", "", "")
	mock.ExpectQuery(`SELECT \* FROM "players"`).
		WillReturnRows(playerRows)

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var body paginatedBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.Len(t, body.Data, 2)
	assert.Equal(t, "Scump", body.Data[0]["gamertag"])

	assert.Equal(t, 1, body.Pagination.Page)
	assert.Equal(t, 25, body.Pagination.Limit)
	assert.Equal(t, 2, body.Pagination.Total)
	assert.Equal(t, 1, body.Pagination.TotalPages)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlayers_SecondPage(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "players"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(30))

	playerRows := sqlmock.NewRows([]string{"id", "gamertag", "first_name", "last_name",
		"country", "role", "is_active", "liquipedia_url", "twitter_handle", "avatar_url"}).
		AddRow(26, "Temp", "Anthony", "Terrell", "US", "flex", true, "", "", "")
	mock.ExpectQuery(`SELECT \* FROM "players"`).
		WillReturnRows(playerRows)

	h := newTestHandler(t)
	c, w := newCtx(nil, "page=2&limit=25")
	h.GetPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body paginatedBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.Equal(t, 2, body.Pagination.Page)
	assert.Equal(t, 30, body.Pagination.Total)
	assert.Equal(t, 2, body.Pagination.TotalPages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlayers_WithSearch(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "players" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	playerRows := sqlmock.NewRows([]string{"id", "gamertag", "first_name", "last_name",
		"country", "role", "is_active", "liquipedia_url", "twitter_handle", "avatar_url"}).
		AddRow(7, "Scump", "Seth", "Abner", "US", "flex", true, "", "", "")
	mock.ExpectQuery(`SELECT \* FROM "players" WHERE`).
		WillReturnRows(playerRows)

	h := newTestHandler(t)
	c, w := newCtx(nil, "search=scump")
	h.GetPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body paginatedBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.Len(t, body.Data, 1)
	assert.Equal(t, "Scump", body.Data[0]["gamertag"])
	assert.Equal(t, 1, body.Pagination.Total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlayers_EmptyResults(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "players"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT \* FROM "players"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "gamertag"}))

	h := newTestHandler(t)
	c, w := newCtx(nil, "search=zzznobody")
	h.GetPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body paginatedBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.Len(t, body.Data, 0)
	assert.Equal(t, 0, body.Pagination.Total)
	assert.Equal(t, 1, body.Pagination.TotalPages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlayers_CountDBError(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "players"`).
		WillReturnError(assert.AnError)

	h := newTestHandler(t)
	c, w := newCtx(nil, "")
	h.GetPlayers(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Contains(t, body["error"], "Failed to fetch players")
}

// ── GetPlayerMatches ──────────────────────────────────────────────────────────

func TestGetPlayerMatches_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "notanumber"}}, "")
	h.GetPlayerMatches(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPlayerMatches_EmptyResponse(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT .+ FROM "player_match_stats" JOIN matches`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "match_id", "player_id", "team_id",
			"maps_played", "total_kills", "total_deaths", "total_assists",
			"total_damage", "kd_ratio", "kda_ratio", "adr",
			"created_at", "updated_at",
		}))

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetPlayerMatches(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	events, _ := body["events"].([]interface{})
	assert.Empty(t, events)
	assert.Equal(t, float64(0), body["total"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlayerKDStats_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "notanumber"}}, "")
	h.GetPlayerKDStats(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}


type playerKDBody struct {
	PlayerID       float64          `json:"player_id"`
	Gamertag       string           `json:"gamertag"`
	TotalKills     float64          `json:"total_kills"`
	TotalDeaths    float64          `json:"total_deaths"`
	TotalAssists   float64          `json:"total_assists"`
	AvgKD          float64          `json:"avg_kd"`
	HpKDRatio      float64          `json:"hp_kd_ratio"`
	SndKDRatio     float64          `json:"snd_kd_ratio"`
	ControlKDRatio float64          `json:"control_kd_ratio"`
	TournamentStats []map[string]any `json:"tournament_stats"`
}

var playerTournamentStatCols = []string{
	"id", "player_id", "team_id", "tournament_id",
	"total_kills", "total_deaths", "total_assists", "total_damage",
	"kd_ratio", "kda_ratio", "overall_maps", "overall_plus_minus",
	"hp_kills", "hp_deaths", "hp_kd_ratio", "hp_maps",
	"snd_kills", "snd_deaths", "snd_kd_ratio", "snd_maps",
	"control_kd_ratio", "control_maps",
}


func TestGetPlayerMatches_ResponseShape(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 10)

	require.NoError(t, database.DB.Create(&models.Player{
		ID: 1, Gamertag: "Shotzzy", IsActive: true,
	}).Error)
	require.NoError(t, database.DB.Create(&models.PlayerMatchStats{
		MatchID: 10, PlayerID: 1, TeamID: 1,
		MapsPlayed: 3, TotalKills: 24, TotalDeaths: 16, TotalAssists: 3,
		KDRatio: 1.50, KDARatio: 1.69, ADR: 850,
	}).Error)

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetPlayerMatches(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.Contains(t, body, "player_id")
	assert.Contains(t, body, "events")
	assert.Contains(t, body, "total")

	events, _ := body["events"].([]any)
	require.Len(t, events, 1, "one event (tournament) expected")

	event, _ := events[0].(map[string]any)
	for _, field := range []string{"event", "year", "tournament_id", "matches"} {
		assert.Contains(t, event, field, "event must contain %s", field)
	}

	matches, _ := event["matches"].([]any)
	require.Len(t, matches, 1)

	match, _ := matches[0].(map[string]any)
	for _, field := range []string{"match_id", "date", "opponent", "opponent_abbr", "result", "kd", "kills", "deaths"} {
		assert.Contains(t, match, field, "match entry must contain %s", field)
	}
}

func TestGetPlayer_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "bad"}}, "")
	h.GetPlayer(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "Invalid player ID", body["error"])
}

func TestGetPlayer_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "players"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "gamertag"}))

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetPlayer(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlayer_Success(t *testing.T) {
	mock := setupMockDB(t)
	rows := sqlmock.NewRows([]string{"id", "gamertag", "first_name", "last_name",
		"country", "role", "is_active"}).
		AddRow(7, "Scump", "Seth", "Abner", "US", "flex", true)
	mock.ExpectQuery(`SELECT \* FROM "players"`).
		WillReturnRows(rows)

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "7"}}, "")
	h.GetPlayer(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var player map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &player))
	assert.Equal(t, "Scump", player["gamertag"])
	assert.Equal(t, "US", player["country"])
	assert.NoError(t, mock.ExpectationsWereMet())
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
