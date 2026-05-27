package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	database.DB = gormDB
	t.Cleanup(func() { sqlDB.Close() })
	return mock
}

// newTestHandler creates a Handler wired to database.DB.
// Must be called after setupMockDB or setupPGTx if DB access is needed.
func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	return New(database.DB)
}

func newCtx(params gin.Params, rawQuery string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?"+rawQuery, nil)
	c.Params = params
	return c, w
}

func TestCalculateKD(t *testing.T) {
	tests := []struct {
		name   string
		kills  int
		deaths int
		want   float64
	}{
		{"normal ratio", 10, 5, 2.0},
		{"equal kd", 5, 5, 1.0},
		{"below 1", 3, 6, 0.5},
		{"zero deaths returns 0", 10, 0, 0},
		{"both zero", 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, services.CalculateKD(tt.kills, tt.deaths))
		})
	}
}

func TestValidateID(t *testing.T) {
	t.Run("valid number", func(t *testing.T) {
		id, err := validateID("42")
		assert.NoError(t, err)
		assert.Equal(t, 42, id)
	})
	t.Run("letters return error", func(t *testing.T) {
		_, err := validateID("abc")
		assert.Error(t, err)
	})
	t.Run("empty string returns error", func(t *testing.T) {
		_, err := validateID("")
		assert.Error(t, err)
	})
	t.Run("negative number is valid", func(t *testing.T) {
		id, err := validateID("-1")
		assert.NoError(t, err)
		assert.Equal(t, -1, id)
	})
	t.Run("float-like string returns error", func(t *testing.T) {
		_, err := validateID("1.5")
		assert.Error(t, err)
	})
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

func TestGetTeams_ScopeAll_Success(t *testing.T) {
	// scope=all triggers a simple GORM query (SELECT * FROM teams ORDER BY name ASC)
	// that sqlmock can match with a straightforward regex.
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
