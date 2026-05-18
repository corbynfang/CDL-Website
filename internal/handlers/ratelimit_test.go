package handlers

// ratelimit_test.go — tests for the per-IP rate limiting middleware.
//
// Key concepts shown here:
//   - Testing middleware directly: we call RateLimit() and pass it a fake Gin
//     context instead of wiring it into a full router. Same result, much simpler.
//   - Global state: the visitors map is package-level. Each test uses a unique
//     fake IP so tests don't interfere with each other.
//   - Sentinel error: assert.AnError is a built-in test error value — use it
//     whenever you need "some error" without caring about the message.

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// fireRequest simulates one API call from the given IP through the rate limiter.
// Returns the HTTP status code — either 200 (allowed) or 429 (blocked).
func fireRequest(ip string) int {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/players", nil)
	c.Request.Header.Set("X-Forwarded-For", ip)

	handler := RateLimit()
	handler(c)

	// If the middleware didn't abort (i.e., it allowed the request through),
	// Gin leaves the status at 200. If it called AbortWithStatusJSON(429,...),
	// the recorder already has 429.
	if !c.IsAborted() {
		w.WriteHeader(http.StatusOK)
	}
	return w.Code
}

func TestRateLimit_AllowsFirstTenRequests(t *testing.T) {
	// The bucket starts with 10 tokens (burst size). All 10 should succeed.
	// Use a unique IP so this test doesn't share a bucket with other tests.
	ip := "192.0.2.10"

	for i := 1; i <= 10; i++ {
		code := fireRequest(ip)
		assert.Equal(t, http.StatusOK, code, "request %d should be allowed", i)
	}
}

func TestRateLimit_BlocksAfterBurst(t *testing.T) {
	// After exhausting the 10-token burst, the 11th request must get 429.
	ip := "192.0.2.11"

	for i := 0; i < 10; i++ {
		fireRequest(ip)
	}
	// 11th request — bucket is empty
	code := fireRequest(ip)
	assert.Equal(t, http.StatusTooManyRequests, code)

	var body map[string]string
	// The response recorder has the JSON body even on 429
	_ = body
}

func TestRateLimit_DifferentIPsHaveIndependentBuckets(t *testing.T) {
	// Each IP gets its own bucket. Exhausting one IP's bucket must not
	// affect a different IP.
	ipA := "192.0.2.20"
	ipB := "192.0.2.21"

	// Drain ipA's bucket completely
	for i := 0; i < 11; i++ {
		fireRequest(ipA)
	}

	// ipB should still have a full bucket
	code := fireRequest(ipB)
	assert.Equal(t, http.StatusOK, code, "ipB should not be affected by ipA's limit")
}

func TestRateLimit_ExtractsFirstIPFromXForwardedFor(t *testing.T) {
	// CloudFront sets X-Forwarded-For: <client-ip>, <cloudfront-ip>
	// The rate limiter must use only the first IP so the same real client
	// is correctly tracked regardless of which CloudFront edge handled it.
	clientIP := "192.0.2.30"
	edge1 := fmt.Sprintf("%s, 13.32.0.1",  clientIP)  // same client, edge node 1
	edge2 := fmt.Sprintf("%s, 13.32.0.99", clientIP)  // same client, edge node 2

	// Send 10 requests through edge node 1 — drains the bucket for clientIP
	for i := 0; i < 10; i++ {
		fireRequest(edge1)
	}

	// 11th request through a different edge node must still be blocked —
	// because it's the same underlying client IP
	code := fireRequest(edge2)
	assert.Equal(t, http.StatusTooManyRequests, code,
		"same client through different CloudFront edge must share the same bucket")
}

func TestRateLimit_UniqueIPsPerTest(t *testing.T) {
	// Sanity check: 5 different IPs all get fresh buckets
	for i := 0; i < 5; i++ {
		ip := fmt.Sprintf("10.0.%d.1", i)
		code := fireRequest(ip)
		assert.Equal(t, http.StatusOK, code, "fresh IP %s should be allowed", ip)
	}
}
