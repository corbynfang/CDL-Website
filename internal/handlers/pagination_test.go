package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newCtxWithQuery(t *testing.T, rawQuery string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?"+rawQuery, nil)
	return c, w
}

func TestParsePagination(t *testing.T) {

	tests := []struct {
		name           string
		query          string
		wantPage       int
		wantLimit      int
		wantOffset     int
	}{
		// Defaults when no params supplied
		{
			name: "no params uses defaults",
			query: "",
			wantPage: 1, wantLimit: 25, wantOffset: 0,
		},
		// Page param
		{
			name: "page 2 gives offset 25",
			query: "page=2",
			wantPage: 2, wantLimit: 25, wantOffset: 25,
		},
		{
			name: "page 3 with limit 10 gives offset 20",
			query: "page=3&limit=10",
			wantPage: 3, wantLimit: 10, wantOffset: 20,
		},
		// Invalid page values — all fall back to default (page=1)
		{
			name: "page=0 falls back to 1",
			query: "page=0",
			wantPage: 1, wantLimit: 25, wantOffset: 0,
		},
		{
			name: "page=-1 falls back to 1",
			query: "page=-1",
			wantPage: 1, wantLimit: 25, wantOffset: 0,
		},
		{
			name: "page=abc falls back to 1",
			query: "page=abc",
			wantPage: 1, wantLimit: 25, wantOffset: 0,
		},
		// Limit param
		{
			name: "custom limit=50",
			query: "limit=50",
			wantPage: 1, wantLimit: 50, wantOffset: 0,
		},
		{
			name: "limit at max boundary (100) is allowed",
			query: "limit=100",
			wantPage: 1, wantLimit: 100, wantOffset: 0,
		},
		// Invalid limit values — all fall back to default (limit=25)
		{
			name: "limit=101 exceeds cap, falls back to 25",
			query: "limit=101",
			wantPage: 1, wantLimit: 25, wantOffset: 0,
		},
		{
			name: "limit=0 falls back to 25",
			query: "limit=0",
			wantPage: 1, wantLimit: 25, wantOffset: 0,
		},
		{
			name: "limit=xyz falls back to 25",
			query: "limit=xyz",
			wantPage: 1, wantLimit: 25, wantOffset: 0,
		},
	}

	for _, tt := range tests {
		// t.Run creates a subtest — if it fails you'll see which scenario broke.
		t.Run(tt.name, func(t *testing.T) {
			c, _ := newCtxWithQuery(t, tt.query)
			page, limit, offset := parsePagination(c)
			assert.Equal(t, tt.wantPage,   page,   "page")
			assert.Equal(t, tt.wantLimit,  limit,  "limit")
			assert.Equal(t, tt.wantOffset, offset, "offset")
		})
	}
}

func TestBuildMeta(t *testing.T) {
	// buildMeta calculates TotalPages = ceil(total / limit).
	// Edge cases: 0 total should still give TotalPages=1 (never return 0 pages).
	tests := []struct {
		name       string
		page       int
		limit      int
		total      int
		wantPages  int
	}{
		{"0 total gives 1 page", 1, 25, 0, 1},
		{"exact fit: 25/25=1", 1, 25, 25, 1},
		{"one over: 26/25 rounds up to 2", 1, 25, 26, 2},
		{"exact 100/25=4", 1, 25, 100, 4},
		{"101/25 rounds up to 5", 1, 25, 101, 5},
		{"1 result, large limit = 1 page", 1, 100, 1, 1},
		{"page number is passed through", 3, 25, 100, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := buildMeta(tt.page, tt.limit, tt.total)
			assert.Equal(t, tt.page,      meta.Page)
			assert.Equal(t, tt.limit,     meta.Limit)
			assert.Equal(t, tt.total,     meta.Total)
			assert.Equal(t, tt.wantPages, meta.TotalPages)
		})
	}
}
