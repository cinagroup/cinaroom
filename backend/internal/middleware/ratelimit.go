package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

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

// TokenBucketConfig holds settings for the per-user token bucket limiter.
type TokenBucketConfig struct {
	Rate       rate.Limit   // tokens added per second
	Burst      int          // bucket capacity (max burst)
	KeyFunc    func(*gin.Context) string // extract rate-limit key from request
	CleanupInterval time.Duration        // how often to purge stale entries
	MaxAge          time.Duration        // max age before entry is removed
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
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// TokenBucketLimiter returns a per-user token bucket rate limiting middleware.
// Default: 100 requests/minute per user (≈ 1.67 req/s, burst 20).
func TokenBucketLimiter(cfg TokenBucketConfig) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	// Default key function: use authenticated user ID, fallback to IP.
	if cfg.KeyFunc == nil {
		cfg.KeyFunc = func(c *gin.Context) string {
			if userID, exists := c.Get("userID"); exists {
				return fmtUserID(userID)
			}
			return c.ClientIP()
		}
	}

	// Default cleanup settings
	if cfg.CleanupInterval == 0 {
		cfg.CleanupInterval = 5 * time.Minute
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 10 * time.Minute
	}

	// Background cleanup
	go func() {
		for {
			time.Sleep(cfg.CleanupInterval)
			mu.Lock()
			for key, v := range visitors {
				if time.Since(v.lastSeen) > cfg.MaxAge {
					delete(visitors, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		key := cfg.KeyFunc(c)

		mu.Lock()
		v, exists := visitors[key]
		if !exists {
			v = &visitor{limiter: rate.NewLimiter(cfg.Rate, cfg.Burst)}
			visitors[key] = v
		}
		v.lastSeen = time.Now()
		mu.Unlock()

		if !v.limiter.Allow() {
			slog.Warn("token bucket rate limit exceeded", "key", key, "path", c.Request.URL.Path)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// DefaultTokenBucketConfig returns the default per-user rate limit:
// 100 req/min (≈ 1.67 tokens/s) with burst of 20.
func DefaultTokenBucketConfig() TokenBucketConfig {
	return TokenBucketConfig{
		Rate:  rate.Limit(100.0 / 60.0), // 100 per minute
		Burst: 20,
	}
}

// fmtUserID safely formats a user ID value.
func fmtUserID(v interface{}) string {
	switch id := v.(type) {
	case uint:
		return fmtUint(id)
	case int:
		return fmtInt(id)
	case string:
		return id
	default:
		return ""
	}
}

func fmtUint(v uint) string {
	return formatPositive(int64(v))
}

func fmtInt(v int) string {
	return formatPositive(int64(v))
}

func formatPositive(v int64) string {
	if v < 0 {
		return "0"
	}
	// Fast path for small numbers
	if v < 10 {
		return string(rune('0' + v))
	}
	var buf [20]byte
	pos := len(buf)
	for v > 0 {
		pos--
		buf[pos] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[pos:])
}
