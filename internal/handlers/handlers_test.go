package handlers

import (
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
