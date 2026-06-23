package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corbynfang/CDL-Website/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountCreation_PostLinkedToUser(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 1)

	const uid = "supabase-uid-corbyn"
	token := signJWT(t, uid)
	r := newTestRouter(New(database.DB))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/auth/profile", jsonBody(t, map[string]string{"username": "Corbyn"}), token))
	require.Equal(t, http.StatusOK, w.Code, "profile creation failed: %s", w.Body.String())

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/matches/1/thread/posts", jsonBody(t, map[string]string{"body": "OpTic looking great today"}), token))
	require.Equal(t, http.StatusCreated, w.Code, "post creation failed: %s", w.Body.String())

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/matches/1/thread", nil))
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	posts, ok := resp["data"].([]any)
	require.True(t, ok)
	require.Len(t, posts, 1)

	post := posts[0].(map[string]any)
	assert.Equal(t, "OpTic looking great today", post["body"])

	user := post["user"].(map[string]any)
	assert.Equal(t, "Corbyn", user["username"])
	assert.NotContains(t, user, "supabase_uid")
}

func TestDeleteAccount_CascadesToPosts(t *testing.T) {
	setupPGTx(t)
	pgMatchEnv(t)
	pgMatch(t, 1)

	const uid = "supabase-uid-delete-me"
	token := signJWT(t, uid)
	r := newTestRouter(New(database.DB))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/auth/profile", jsonBody(t, map[string]string{"username": "DeleteMe"}), token))
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodPost, "/api/v1/matches/1/thread/posts", jsonBody(t, map[string]string{"body": "this should disappear"}), token))
	require.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authReq(http.MethodDelete, "/api/v1/auth/me", nil, token))
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/matches/1/thread", nil))
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	posts, _ := resp["data"].([]any)
	assert.Empty(t, posts, "posts should be soft-deleted when account is deleted")
}
