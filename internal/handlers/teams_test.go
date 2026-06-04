package handlers

// teams_test.go — handler-level tests for GetTeamPlayers scope/season_id parsing
// and store-method routing. Uses a fake TeamStore so no SQL/DB is involved.

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeTeamStore records which roster method the handler routed to and the
// season_id it forwarded. Non-roster methods are unused no-ops.
type fakeTeamStore struct {
	calledMethod string
	calledSeason string
	players      []models.Player
}

func (f *fakeTeamStore) GetPlayers(_ context.Context, _ int, seasonID string) ([]models.Player, error) {
	f.calledMethod = "GetPlayers"
	f.calledSeason = seasonID
	return f.players, nil
}

func (f *fakeTeamStore) GetLatestMatchRoster(_ context.Context, _ int, seasonID string) ([]models.Player, error) {
	f.calledMethod = "GetLatestMatchRoster"
	f.calledSeason = seasonID
	return f.players, nil
}

func (f *fakeTeamStore) ListActiveCDL(context.Context) ([]models.Team, error) { return nil, nil }
func (f *fakeTeamStore) ListForSeason(context.Context, string, string) ([]models.Team, error) {
	return nil, nil
}
func (f *fakeTeamStore) ListAll(context.Context) ([]models.Team, error)          { return nil, nil }
func (f *fakeTeamStore) GetByID(context.Context, int) (*models.Team, error)       { return nil, nil }
func (f *fakeTeamStore) GetStats(context.Context, int) ([]models.TeamTournamentStats, error) {
	return nil, nil
}

func handlerWithFakeTeams(f *fakeTeamStore) *Handler {
	return &Handler{teams: services.NewTeamService(f, nil)}
}

func errBody(t *testing.T, body []byte) string {
	t.Helper()
	var m map[string]string
	require.NoError(t, json.Unmarshal(body, &m))
	return m["error"]
}

func TestGetTeamPlayers_InvalidTeamID(t *testing.T) {
	f := &fakeTeamStore{}
	h := handlerWithFakeTeams(f)
	c, w := newCtx(gin.Params{{Key: "id", Value: "abc"}}, "")
	h.GetTeamPlayers(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid team ID", errBody(t, w.Body.Bytes()))
	assert.Empty(t, f.calledMethod, "store must not be called on invalid team id")
}

func TestGetTeamPlayers_InvalidSeasonID(t *testing.T) {
	f := &fakeTeamStore{}
	h := handlerWithFakeTeams(f)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "season_id=abc")
	h.GetTeamPlayers(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid season_id", errBody(t, w.Body.Bytes()))
	assert.Empty(t, f.calledMethod, "store must not be called on invalid season_id")
}

func TestGetTeamPlayers_InvalidScope(t *testing.T) {
	f := &fakeTeamStore{}
	h := handlerWithFakeTeams(f)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "scope=bogus")
	h.GetTeamPlayers(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid scope", errBody(t, w.Body.Bytes()))
	assert.Empty(t, f.calledMethod, "store must not be called on invalid scope")
}

func TestGetTeamPlayers_DefaultScopeUsesCurrentRoster(t *testing.T) {
	f := &fakeTeamStore{players: []models.Player{{ID: 1, Gamertag: "Dashy"}}}
	h := handlerWithFakeTeams(f)
	c, w := newCtx(gin.Params{{Key: "id", Value: "4"}}, "")
	h.GetTeamPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "GetLatestMatchRoster", f.calledMethod, "no scope should default to current roster")
	assert.Equal(t, "", f.calledSeason)

	var players []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &players))
	require.Len(t, players, 1)
	assert.Equal(t, "Dashy", players[0]["gamertag"])
}

func TestGetTeamPlayers_ScopeCurrentForwardsSeasonID(t *testing.T) {
	f := &fakeTeamStore{}
	h := handlerWithFakeTeams(f)
	c, w := newCtx(gin.Params{{Key: "id", Value: "4"}}, "season_id=2&scope=current")
	h.GetTeamPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "GetLatestMatchRoster", f.calledMethod)
	assert.Equal(t, "2", f.calledSeason, "season_id must be forwarded to the store")
}

func TestGetTeamPlayers_ScopeAllUsesStints(t *testing.T) {
	f := &fakeTeamStore{}
	h := handlerWithFakeTeams(f)
	c, w := newCtx(gin.Params{{Key: "id", Value: "22"}}, "season_id=2&scope=all")
	h.GetTeamPlayers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "GetPlayers", f.calledMethod, "scope=all must use the stint union path")
	assert.Equal(t, "2", f.calledSeason)
}

func TestGetTeam_InvalidID(t *testing.T) {
	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "notanumber"}}, "")
	h.GetTeam(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "Invalid team ID", body["error"])
}

func TestGetTeam_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "teams"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "999"}}, "")
	h.GetTeam(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTeam_Success(t *testing.T) {
	mock := setupMockDB(t)
	rows := sqlmock.NewRows([]string{"id", "name", "abbreviation", "city", "logo_url",
		"primary_color", "secondary_color", "is_active"}).
		AddRow(1, "Atlanta FaZe", "ATL", "Atlanta", "", "#000", "#f00", true)
	mock.ExpectQuery(`SELECT \* FROM "teams"`).
		WillReturnRows(rows)

	h := newTestHandler(t)
	c, w := newCtx(gin.Params{{Key: "id", Value: "1"}}, "")
	h.GetTeam(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var team map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &team))
	assert.Equal(t, "Atlanta FaZe", team["name"])
	assert.Equal(t, "ATL", team["abbreviation"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTeams_ScopeAll_Success(t *testing.T) {
	mock := setupMockDB(t)
	rows := sqlmock.NewRows([]string{"id", "name", "abbreviation", "is_active"}).
		AddRow(1, "Atlanta FaZe", "ATL", true).
		AddRow(2, "Boston Breach", "BOS", true)
	mock.ExpectQuery(`SELECT \* FROM "teams"`).
		WillReturnRows(rows)

	h := newTestHandler(t)
	c, w := newCtx(nil, "scope=all")
	h.GetTeams(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var teams []map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &teams))
	assert.Len(t, teams, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTeams_ScopeAll_Empty(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "teams"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	h := newTestHandler(t)
	c, w := newCtx(nil, "scope=all")
	h.GetTeams(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
