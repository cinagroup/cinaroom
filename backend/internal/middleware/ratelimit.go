package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"cinaroom-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// visitor stores rate-limit state for a single client.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiterConfig holds settings for the rate limiter.
type RateLimiterConfig struct {
	Rate  rate.Limit // tokens per second
	Burst int        // maximum burst size
}

// RateLimiter returns a per-IP rate-limiting middleware.
func RateLimiter(cfg RateLimiterConfig) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	// Periodically clean up stale entries.
	go func() {
		for {
			time.Sleep(3 * time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 5*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			v = &visitor{limiter: rate.NewLimiter(cfg.Rate, cfg.Burst)}
			visitors[ip] = v
		}
		v.lastSeen = time.Now()
		mu.Unlock()

		if !v.limiter.Allow() {
			slog.Warn("rate limit exceeded", "ip", ip, "path", c.Request.URL.Path)
			response.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
