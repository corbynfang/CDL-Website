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

	"github.com/DATA-DOG/go-sqlmock"
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

