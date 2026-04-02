package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// UserTier represents the subscription tier of a user.
type UserTier int

const (
	TierFree      UserTier = iota // Free tier: 1 concurrent connection
	TierPro                       // Pro tier: 5 concurrent connections
	TierEnterprise                // Enterprise tier: unlimited (capped at 100)
)

// ConcurrencyConfig holds settings for per-user concurrent connection limits.
type ConcurrencyConfig struct {
	FreeLimit      int32 // Max concurrent connections for free users
	ProLimit       int32 // Max concurrent connections for pro users
	EnterpriseCap  int32 // Hard cap even for enterprise users
	TierResolver   func(*gin.Context) UserTier // Resolve user tier from request
}

// DefaultConcurrencyConfig returns the default concurrency configuration.
func DefaultConcurrencyConfig() ConcurrencyConfig {
	return ConcurrencyConfig{
		FreeLimit:     1,
		ProLimit:      5,
		EnterpriseCap: 100,
		TierResolver:  defaultTierResolver,
	}
}

// defaultTierResolver returns Free tier by default.
// Replace this with actual tier lookup (e.g. from DB or JWT claims).
func defaultTierResolver(c *gin.Context) UserTier {
	// Check if tier is set in context (e.g. by auth middleware)
	if tier, exists := c.Get("userTier"); exists {
		if t, ok := tier.(UserTier); ok {
			return t
		}
	}
	return TierFree
}

// concurrentCounter tracks per-user connection counts.
type concurrentCounter struct {
	counters sync.Map // map[string]*atomic.Int32
}

// concurrentManager is the global concurrent connection manager.
var concurrentMgr = &concurrentCounter{}

// ConcurrencyLimiter returns a middleware that limits concurrent connections per user.
func ConcurrencyLimiter(cfg ConcurrencyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := extractUserID(c)
		if userID == "" {
			// No user identity — use IP as fallback
			userID = c.ClientIP()
		}

		tier := cfg.TierResolver(c)
		limit := tierLimit(tier, cfg)

		// Get or create counter for this user
		val, _ := concurrentMgr.counters.LoadOrStore(userID, &atomic.Int32{})
		counter := val.(*atomic.Int32)

		current := counter.Add(1)
		defer counter.Add(-1)

		if current > limit {
			slog.Warn("concurrent connection limit exceeded",
				"user_id", userID,
				"current", current,
				"limit", limit,
				"tier", tier,
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "并发连接数已达上限，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// tierLimit returns the concurrent connection limit for a given tier.
func tierLimit(tier UserTier, cfg ConcurrencyConfig) int32 {
	switch tier {
	case TierFree:
		return cfg.FreeLimit
	case TierPro:
		return cfg.ProLimit
	case TierEnterprise:
		return cfg.EnterpriseCap
	default:
		return cfg.FreeLimit
	}
}

// extractUserID gets the user ID from the Gin context.
func extractUserID(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		return fmtUserID(userID)
	}
	return ""
}

// ActiveConnections returns the current number of active connections for a user.
func ActiveConnections(userID string) int32 {
	val, ok := concurrentMgr.counters.Load(userID)
	if !ok {
		return 0
	}
	return val.(*atomic.Int32).Load()
}
