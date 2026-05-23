package handlers

// ratelimit.go — per-IP rate limiting middleware using a token bucket algorithm.
//
// Each IP gets a bucket of 20 tokens that refills at 2 tokens/second.
// An event detail page fires ~5 requests on first open; 20-token burst covers
// ~4 rapid page loads before the bucket empties. Scrapers that sustain >2 req/s
// will still be blocked within a few seconds.
//
// A background goroutine cleans up limiters for IPs not seen in 5 minutes
// so the map doesn't grow forever.

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors   = make(map[string]*visitor)
	visitorsMu sync.Mutex
)

func init() {
	go cleanupVisitors()
}

// cleanupVisitors removes IPs that haven't made a request in 5 minutes.
// Without this the map would grow by one entry per unique IP forever.
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		visitorsMu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(visitors, ip)
			}
		}
		visitorsMu.Unlock()
	}
}

func getVisitor(ip string) *rate.Limiter {
	visitorsMu.Lock()
	defer visitorsMu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		// 2 tokens per second sustained, burst up to 20 tokens.
		v = &visitor{limiter: rate.NewLimiter(rate.Every(500*time.Millisecond), 20)}
		visitors[ip] = v
	}
	v.lastSeen = time.Now()
	return v.limiter
}

// RateLimit is a Gin middleware that returns 429 when an IP exceeds the limit.
// Uses c.ClientIP() which respects the trusted-proxy list set in main.go —
// that makes XFF safe to read because only trusted infrastructure (the ALB)
// can inject the header. Never parse XFF manually: a caller who hits the ALB
// directly could forge any IP and rotate through unlimited fresh buckets.
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !getVisitor(ip).Allow() {
			c.Header("Retry-After", "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down.",
			})
			return
		}
		c.Next()
	}
}
