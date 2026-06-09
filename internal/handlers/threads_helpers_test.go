package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "test-jwt-secret-for-cdlytics"

func signJWT(t *testing.T, uid string) string {
	t.Helper()
	t.Setenv("SUPABASE_JWT_SECRET", testJWTSecret)
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  uid,
		"role": "authenticated",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	signed, err := tok.SignedString([]byte(testJWTSecret))
	require.NoError(t, err)
	return signed
}

func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

func authReq(method, path string, body *bytes.Buffer, token string) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, body)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func decodeJSON(b []byte, v any) error {
	return json.Unmarshal(b, v)
}

func postPath(id int) string {
	return fmt.Sprintf("/api/v1/thread/posts/%d", id)
}
