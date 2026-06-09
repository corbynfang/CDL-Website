package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncProfile_CreatesUser(t *testing.T) {
	setupPGTx(t)
	token := signJWT(t, "uid-sync")
	r := newTestRouter(New(database.DB))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/auth/profile", jsonBody(t, map[string]string{"username": "TestUser"}), token))

	require.Equal(t, http.StatusOK, w.Code)
	var user models.User
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &user))
	assert.Equal(t, "TestUser", user.Username)
	assert.Equal(t, "uid-sync", user.SupabaseUID)
}

func TestSyncProfile_Idempotent(t *testing.T) {
	setupPGTx(t)
	token := signJWT(t, "uid-idempotent")
	r := newTestRouter(New(database.DB))

	for range 2 {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/auth/profile", jsonBody(t, map[string]string{"username": "SameUser"}), token))
		assert.Equal(t, http.StatusOK, w.Code)
	}

	var count int64
	require.NoError(t, database.DB.Model(&models.User{}).Where("supabase_uid = ?", "uid-idempotent").Count(&count).Error)
	assert.EqualValues(t, 1, count)
}

func TestGetMe_NotFound_WithoutProfile(t *testing.T) {
	setupPGTx(t)
	token := signJWT(t, "uid-no-profile")
	r := newTestRouter(New(database.DB))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodGet, "/api/v1/auth/me", nil, token))
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestGetMe_ReturnsProfile(t *testing.T) {
	setupPGTx(t)
	token := signJWT(t, "uid-getme")
	r := newTestRouter(New(database.DB))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/auth/profile", jsonBody(t, map[string]string{"username": "GetMeUser"}), token))
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodGet, "/api/v1/auth/me", nil, token))
	require.Equal(t, http.StatusOK, w.Code)

	var user models.User
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &user))
	assert.Equal(t, "GetMeUser", user.Username)
}

func TestDeleteMe_RemovesAccount(t *testing.T) {
	setupPGTx(t)
	token := signJWT(t, "uid-delete-me")
	r := newTestRouter(New(database.DB))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/auth/profile", jsonBody(t, map[string]string{"username": "GoneUser"}), token))
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodDelete, "/api/v1/auth/me", nil, token))
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodGet, "/api/v1/auth/me", nil, token))
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
