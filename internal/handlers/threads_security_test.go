package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetThread_EmptyForNewMatch(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 1)

	r := newTestRouter(New(database.DB))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/matches/1/thread", nil))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetThread_InvalidMatchID(t *testing.T) {
	setupPGTx(t)
	r := newTestRouter(New(database.DB))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/matches/abc/thread", nil))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePost_RequiresAuth(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 1)
	t.Setenv("SUPABASE_JWT_SECRET", testJWTSecret)

	r := newTestRouter(New(database.DB))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/matches/1/thread/posts", jsonBody(t, map[string]string{"body": "hello"})))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreatePost_RequiresProfileSetup(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 1)

	token := signJWT(t, "uid-no-profile")
	r := newTestRouter(New(database.DB))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/matches/1/thread/posts", jsonBody(t, map[string]string{"body": "hello"}), token))
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestEditPost_OnlyOwnerCanEdit(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 1)

	ownerToken := signJWT(t, "uid-edit-owner")
	otherToken := signJWT(t, "uid-edit-other")
	r := newTestRouter(New(database.DB))

	for _, tc := range []struct{ token, username string }{{ownerToken, "EditOwner"}, {otherToken, "EditOther"}} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/auth/profile", jsonBody(t, map[string]string{"username": tc.username}), tc.token))
		require.Equal(t, http.StatusOK, w.Code)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/matches/1/thread/posts", jsonBody(t, map[string]string{"body": "original"}), ownerToken))
	require.Equal(t, http.StatusCreated, w.Code)

	var created map[string]any
	require.NoError(t, decodeJSON(w.Body.Bytes(), &created))
	postID := int(created["id"].(float64))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPut, postPath(postID), jsonBody(t, map[string]string{"body": "hijacked"}), otherToken))
	assert.Equal(t, http.StatusForbidden, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPut, postPath(postID), jsonBody(t, map[string]string{"body": "updated"}), ownerToken))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeletePost_OnlyOwnerCanDelete(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 1)

	ownerToken := signJWT(t, "uid-del-owner")
	otherToken := signJWT(t, "uid-del-other")
	r := newTestRouter(New(database.DB))

	for _, tc := range []struct{ token, username string }{{ownerToken, "DelOwner"}, {otherToken, "DelOther"}} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/auth/profile", jsonBody(t, map[string]string{"username": tc.username}), tc.token))
		require.Equal(t, http.StatusOK, w.Code)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/matches/1/thread/posts", jsonBody(t, map[string]string{"body": "to be deleted"}), ownerToken))
	require.Equal(t, http.StatusCreated, w.Code)

	var created map[string]any
	require.NoError(t, decodeJSON(w.Body.Bytes(), &created))
	postID := int(created["id"].(float64))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodDelete, postPath(postID), nil, otherToken))
	assert.Equal(t, http.StatusForbidden, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodDelete, postPath(postID), nil, ownerToken))
	assert.Equal(t, http.StatusOK, w.Code)
}
