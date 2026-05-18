package handlers

// ratelimit.go — per-IP rate limiting middleware using a token bucket algorithm.
//
// Each IP gets a bucket of 10 tokens that refills at 1 token/second.
// This allows short bursts (clicking through pages quickly) while blocking
// sustained high-rate traffic like scrapers or runaway scripts.
//
// A background goroutine cleans up limiters for IPs not seen in 5 minutes
// so the map doesn't grow forever.

import (
	"net/http"
	"strings"
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
		// 1 token per second sustained, burst up to 10 tokens.
		v = &visitor{limiter: rate.NewLimiter(rate.Every(time.Second), 10)}
		visitors[ip] = v
	}
	v.lastSeen = time.Now()
	return v.limiter
}

// RateLimit is a Gin middleware that returns 429 when an IP exceeds the limit.
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// CloudFront sets X-Forwarded-For as "client-ip, cloudfront-ip".
		// Take only the first entry — that's the real client address.
		ip := c.ClientIP()
		if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
			ip = strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
		}

		if !getVisitor(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down.",
			})
			return
		}
		c.Next()
	}
}
