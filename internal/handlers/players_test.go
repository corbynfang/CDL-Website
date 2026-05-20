package handlers

// players_test.go — tests for GetPlayers and GetPlayer handlers.
//
// Key concepts shown here:
//   - Mock database: we never touch a real database. sqlmock intercepts the
//     SQL that GORM generates and returns rows we control.
//   - ExpectQuery vs ExpectExec: SELECT → ExpectQuery, INSERT/UPDATE/DELETE → ExpectExec.
//   - GetPlayers runs TWO queries (COUNT then SELECT), so we set up two expectations
//     in the exact order they fire.
//   - Response shape: GetPlayers returns { data: [...], pagination: {...} }, so we
//     unmarshal into a struct that matches that shape.

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// paginatedBody is the shape GetPlayers returns.
// Unexported fields use any because the exact player fields don't matter for
// these tests — we just care that the envelope is correct.
type paginatedBody struct {
	Data       []map[string]any `json:"data"`
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	} `json:"pagination"`
}

// ── GetPlayers ───────────────────────────────────────────────────────────────

func TestGetPlayers_DefaultPagination(t *testing.T) {
	mock := setupMockDB(t)

	// GetPlayers runs COUNT first, then SELECT.
	// We must register expectations in the same order the handler fires them.

	// Expectation 1: COUNT query — returns 2 total players.
	// sqlmock.NewRows defines column names; AddRow adds one result row.
	mock.ExpectQuery(`SELECT count\(\*\) FROM "players"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// Expectation 2: paginated SELECT — returns the two player rows.
	playerRows := sqlmock.NewRows([]string{"id", "gamertag", "first_name", "last_name",
		"country", "role", "is_active", "liquipedia_url", "twitter_handle", "avatar_url"}).
		AddRow(1, "Scump", "Seth", "Abner", "US", "flex", true, "", "", "").
		AddRow(2, "Crimsix", "Ian", "Porter", "US", "AR", true, "", "", "")
	mock.ExpectQuery(`SELECT \* FROM "players"`).
		WillReturnRows(playerRows)

	c, w := newCtx(nil, "")
	GetPlayers(c)

	// Check HTTP status
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response into our typed struct
	var body paginatedBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	// The data array should have 2 players
	assert.Len(t, body.Data, 2)
	assert.Equal(t, "Scump", body.Data[0]["gamertag"])

	// Pagination metadata should reflect defaults and the COUNT result
	assert.Equal(t, 1,  body.Pagination.Page)
	assert.Equal(t, 25, body.Pagination.Limit)
	assert.Equal(t, 2,  body.Pagination.Total)
	assert.Equal(t, 1,  body.Pagination.TotalPages)

	// Verify sqlmock got exactly the queries we expected — no more, no fewer
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlayers_SecondPage(t *testing.T) {
	mock := setupMockDB(t)

	// 30 total players, page 2, limit 25 → only 5 players on page 2
	mock.ExpectQuery(`SELECT count\(\*\) FROM "players"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(30))

	playerRows := sqlmock.NewRows([]string{"id", "gamertag", "first_name", "last_name",
		"country", "role", "is_active", "liquipedia_url", "twitter_handle", "avatar_url"}).
		AddRow(26, "Temp", "Anthony", "Terrell", "US", "flex", true, "", "", "")
	mock.ExpectQuery(`SELECT \* FROM "players"`).
		WillReturnRows(playerRows)

	// ?page=2 is passed via the query string
	c, w := newCtx(nil, "page=2&limit=25")
	GetPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body paginatedBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	// Page 2 of 2 (ceil(30/25) = 2)
	assert.Equal(t, 2,  body.Pagination.Page)
	assert.Equal(t, 30, body.Pagination.Total)
	assert.Equal(t, 2,  body.Pagination.TotalPages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlayers_WithSearch(t *testing.T) {
	mock := setupMockDB(t)

	// When search is provided, both COUNT and SELECT include WHERE gamertag ILIKE
	mock.ExpectQuery(`SELECT count\(\*\) FROM "players" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	playerRows := sqlmock.NewRows([]string{"id", "gamertag", "first_name", "last_name",
		"country", "role", "is_active", "liquipedia_url", "twitter_handle", "avatar_url"}).
		AddRow(7, "Scump", "Seth", "Abner", "US", "flex", true, "", "", "")
	mock.ExpectQuery(`SELECT \* FROM "players" WHERE`).
		WillReturnRows(playerRows)

	c, w := newCtx(nil, "search=scump")
	GetPlayers(c)

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

	// Search that matches nobody — count=0, data=[]
	mock.ExpectQuery(`SELECT count\(\*\) FROM "players"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT \* FROM "players"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "gamertag"}))

	c, w := newCtx(nil, "search=zzznobody")
	GetPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body paginatedBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	// Empty data, but pagination still present
	assert.Len(t, body.Data, 0)
	assert.Equal(t, 0, body.Pagination.Total)
	// buildMeta guarantees TotalPages is never 0
	assert.Equal(t, 1, body.Pagination.TotalPages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlayers_CountDBError(t *testing.T) {
	mock := setupMockDB(t)

	// Simulate a database failure on the COUNT query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "players"`).
		WillReturnError(assert.AnError) // assert.AnError is a generic sentinel error

	c, w := newCtx(nil, "")
	GetPlayers(c)

	// Handler should return 500 when the DB fails
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Contains(t, body["error"], "Failed to fetch players")
}

// ── GetPlayerMatches ──────────────────────────────────────────────────────────

func TestGetPlayerMatches_InvalidID(t *testing.T) {
	c, w := newCtx(gin.Params{{Key: "id", Value: "notanumber"}}, "")
	GetPlayerMatches(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPlayerMatches_EmptyResponse(t *testing.T) {
	mock := setupMockDB(t)

	// GORM fires one SELECT on player_match_stats; 0 rows means no preload queries.
	mock.ExpectQuery(`SELECT \* FROM "player_match_stats"`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "match_id", "player_id", "team_id",
			"maps_played", "total_kills", "total_deaths", "total_assists",
			"total_damage", "kd_ratio", "kda_ratio", "adr",
			"created_at", "updated_at",
		}))

	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	GetPlayerMatches(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	events, _ := body["events"].([]interface{})
	assert.Empty(t, events)
	assert.Equal(t, float64(0), body["total"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── sortEventsByRecentMatch ───────────────────────────────────────────────────

// ── GetPlayerKDStats ──────────────────────────────────────────────────────────

func TestGetPlayerKDStats_InvalidID(t *testing.T) {
	c, w := newCtx(gin.Params{{Key: "id", Value: "notanumber"}}, "")
	GetPlayerKDStats(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPlayerKDStats_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "players"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "gamertag"}))

	c, w := newCtx(gin.Params{{Key: "id", Value: "99"}}, "")
	GetPlayerKDStats(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// playerKDBody mirrors the shape GetPlayerKDStats returns.
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

// playerTournamentStatCols lists all columns GetPlayerKDStats reads from player_tournament_stats.
var playerTournamentStatCols = []string{
	"id", "player_id", "team_id", "tournament_id",
	"total_kills", "total_deaths", "total_assists", "total_damage",
	"kd_ratio", "kda_ratio", "overall_maps", "overall_plus_minus",
	"hp_kills", "hp_deaths", "hp_kd_ratio", "hp_maps",
	"snd_kills", "snd_deaths", "snd_kd_ratio", "snd_maps",
	"control_kd_ratio", "control_maps",
}

func TestGetPlayerKDStats_ResponseShape(t *testing.T) {
	mock := setupMockDB(t)
	now := time.Now()

	// 1. Player lookup
	mock.ExpectQuery(`SELECT \* FROM "players"`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "gamertag", "avatar_url", "is_active", "created_at", "updated_at"}).
			AddRow(1, "Shotzzy", "", true, now, now))

	// 2. Tournament stats (one row with mode-specific columns)
	mock.ExpectQuery(`SELECT \* FROM "player_tournament_stats"`).WillReturnRows(
		sqlmock.NewRows(playerTournamentStatCols).AddRow(
			1, 1, 1, 5, // id, player_id, team_id, tournament_id
			120, 80, 10, 50000, // kills, deaths, assists, damage
			1.50, 1.63, 10, 40, // kd_ratio, kda_ratio, overall_maps, plus_minus
			60, 40, 1.50, 5, // hp_kills, hp_deaths, hp_kd_ratio, hp_maps
			30, 20, 1.50, 3, // snd_kills, snd_deaths, snd_kd_ratio, snd_maps
			1.20, 2, // control_kd_ratio, control_maps
		))

	// 3. Tournament preload
	mock.ExpectQuery(`SELECT \* FROM "tournaments"`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "slug", "tournament_type", "start_date", "is_lan", "created_at", "updated_at"}).
			AddRow(5, "CDL Major 1 2025", "cdl-major-1-2025", "major", now, true, now, now))

	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	GetPlayerKDStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body playerKDBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.Equal(t, "Shotzzy", body.Gamertag, "gamertag required")
	assert.InDelta(t, 1.0, body.AvgKD, 5.0, "avg_kd must be a finite number")
	assert.NotNil(t, body.TournamentStats, "tournament_stats required")
	require.Len(t, body.TournamentStats, 1)

	ts := body.TournamentStats[0]
	for _, field := range []string{"tournament_id", "tournament_name", "kills", "deaths", "kd_ratio", "maps_played"} {
		assert.Contains(t, ts, field, "tournament_stats entry must contain %s", field)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlayerKDStats_ControlKDZeroWhenNoControlMaps(t *testing.T) {
	mock := setupMockDB(t)
	now := time.Now()

	mock.ExpectQuery(`SELECT \* FROM "players"`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "gamertag", "avatar_url", "is_active", "created_at", "updated_at"}).
			AddRow(2, "Simp", "", true, now, now))

	// control_maps = 0 → handler skips adding to ctlKDSum, ctlMapsTotal stays 0
	mock.ExpectQuery(`SELECT \* FROM "player_tournament_stats"`).WillReturnRows(
		sqlmock.NewRows(playerTournamentStatCols).AddRow(
			2, 2, 1, 5,
			100, 70, 5, 30000,
			1.43, 1.50, 8, 30,
			50, 30, 1.67, 4,
			25, 18, 1.39, 2,
			0.0, 0, // control_kd_ratio=0, control_maps=0
		))

	mock.ExpectQuery(`SELECT \* FROM "tournaments"`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "slug", "tournament_type", "start_date", "is_lan", "created_at", "updated_at"}).
			AddRow(5, "CDL Major 1 2025", "cdl-major-1-2025", "major", now, true, now, now))

	c, w := newCtx(gin.Params{{Key: "id", Value: "2"}}, "")
	GetPlayerKDStats(c)

	var body playerKDBody
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	// Documents current contract: control_kd_ratio is 0, not null, when no control maps played.
	// The frontend must distinguish 0 from "no data" — currently it does not.
	assert.Equal(t, float64(0), body.ControlKDRatio,
		"control_kd_ratio = 0 when no control maps (not null)")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetPlayerMatches — with data ──────────────────────────────────────────────

func TestGetPlayerMatches_ResponseShape(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 10)

	require.NoError(t, database.DB.Create(&database.Player{
		ID: 1, Gamertag: "Shotzzy", IsActive: true,
	}).Error)
	require.NoError(t, database.DB.Create(&database.PlayerMatchStats{
		MatchID: 10, PlayerID: 1, TeamID: 1,
		MapsPlayed: 3, TotalKills: 24, TotalDeaths: 16, TotalAssists: 3,
		KDRatio: 1.50, KDARatio: 1.69, ADR: 850,
	}).Error)

	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	GetPlayerMatches(c)

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

// ── sortEventsByRecentMatch ───────────────────────────────────────────────────

func TestSortEventsByRecentMatch_OrdersByRecentMatchFirst(t *testing.T) {
	// Event B (match_id 5) was built before event A (match_id 10) — simulates map iteration.
	events := []gin.H{
		{"event": "B", "matches": []gin.H{{"match_id": uint(5)}}},
		{"event": "A", "matches": []gin.H{{"match_id": uint(10)}}},
	}
	sortEventsByRecentMatch(events)
	assert.Equal(t, "A", events[0]["event"])
	assert.Equal(t, "B", events[1]["event"])
}

func TestSortEventsByRecentMatch_EmptyMatchesLastRight(t *testing.T) {
	events := []gin.H{
		{"event": "empty", "matches": []gin.H{}},
		{"event": "has-match", "matches": []gin.H{{"match_id": uint(1)}}},
	}
	sortEventsByRecentMatch(events)
	assert.Equal(t, "has-match", events[0]["event"])
	assert.Equal(t, "empty", events[1]["event"])
}

func TestSortEventsByRecentMatch_StableOnSingleEvent(t *testing.T) {
	events := []gin.H{
		{"event": "only", "matches": []gin.H{{"match_id": uint(7)}}},
	}
	sortEventsByRecentMatch(events)
	assert.Equal(t, "only", events[0]["event"])
}

