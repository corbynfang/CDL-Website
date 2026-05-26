package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func fireRequest(ip string) int {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/players", nil)
	c.Request.Header.Set("X-Forwarded-For", ip)

	handler := RateLimit()
	handler(c)

	if !c.IsAborted() {
		w.WriteHeader(http.StatusOK)
	}
	return w.Code
}

func TestRateLimit_AllowsFirstTwentyRequests(t *testing.T) {
	ip := "192.0.2.10"

	for i := 1; i <= 20; i++ {
		code := fireRequest(ip)
		assert.Equal(t, http.StatusOK, code, "request %d should be allowed", i)
	}
}

func TestRateLimit_BlocksAfterBurst(t *testing.T) {
	ip := "192.0.2.11"

	for i := 0; i < 20; i++ {
		fireRequest(ip)
	}

	code := fireRequest(ip)
	assert.Equal(t, http.StatusTooManyRequests, code)
}

func TestRateLimit_DifferentIPsHaveIndependentBuckets(t *testing.T) {

	ipA := "192.0.2.20"
	ipB := "192.0.2.21"

	for i := 0; i < 21; i++ {
		fireRequest(ipA)
	}

	code := fireRequest(ipB)
	assert.Equal(t, http.StatusOK, code, "ipB should not be affected by ipA's limit")
}

func TestRateLimit_ExtractsFirstIPFromXForwardedFor(t *testing.T) {
	clientIP := "192.0.2.30"
	edge1 := fmt.Sprintf("%s, 13.32.0.1",  clientIP)  // same client, edge node 1
	edge2 := fmt.Sprintf("%s, 13.32.0.99", clientIP)  // same client, edge node 2

	for i := 0; i < 20; i++ {
		fireRequest(edge1)
	}

	// 21st request through a different edge node must still be blocked —
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
