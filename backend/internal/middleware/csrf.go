package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CSRFConfig holds settings for CSRF protection.
type CSRFConfig struct {
	// TokenLength is the number of random bytes in the CSRF token.
	TokenLength int
	// TokenLookup defines where the token is read from: "header:X-CSRF-Token" or "form:_csrf".
	TokenLookup string
	// CookieName is the name of the CSRF cookie.
	CookieName string
	// CookieDomain is the domain for the CSRF cookie.
	CookieDomain string
	// CookiePath is the path for the CSRF cookie.
	CookiePath string
	// CookieSecure marks the cookie as Secure (HTTPS only).
	CookieSecure bool
	// CookieHTTPOnly marks the cookie as HttpOnly.
	CookieHTTPOnly bool
	// CookieSameSite controls the SameSite attribute: "Strict", "Lax", or "None".
	CookieSameSite string
	// TokenExpiry is how long a CSRF token remains valid.
	TokenExpiry time.Duration
	// ExemptPaths are paths that skip CSRF checks (e.g. API paths with Bearer auth).
	ExemptPaths []string
	// ExemptPrefixes are path prefixes that skip CSRF checks.
	ExemptPrefixes []string
}

// DefaultCSRFConfig returns a secure default CSRF configuration.
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		TokenLength:    32,
		TokenLookup:    "header:X-CSRF-Token",
		CookieName:     "cinaseek_csrf",
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: false, // Must be readable by JS for SPA
		CookieSameSite: "Lax",
		TokenExpiry:    1 * time.Hour,
		ExemptPrefixes: []string{"/api/", "/health"},
	}
}

// csrfTokenStore holds issued CSRF tokens with their expiry.
type csrfTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]time.Time // token → expiry
}

var csrfStore = &csrfTokenStore{
	tokens: make(map[string]time.Time),
}

// generateCSRFToken creates a cryptographically random hex token.
func generateCSRFToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CSRF returns a CSRF protection middleware.
// It sets SameSite cookies and validates CSRF tokens on state-changing requests.
func CSRF(cfg CSRFConfig) gin.HandlerFunc {
	// Start cleanup goroutine
	go csrfStore.cleanupLoop()

	return func(c *gin.Context) {
		// Skip exempt paths
		if isCSRFExempt(c.Request.URL.Path, cfg) {
			c.Next()
			return
		}

		// For safe methods (GET, HEAD, OPTIONS), issue a new CSRF token.
		if isSafeMethod(c.Request.Method) {
			token, err := generateCSRFToken(cfg.TokenLength)
			if err != nil {
				slog.Error("failed to generate CSRF token", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "内部错误",
				})
				c.Abort()
				return
			}

			csrfStore.store(token, time.Now().Add(cfg.TokenExpiry))
			setCSRFCookie(c, cfg, token)
			c.Header("X-CSRF-Token", token)
			c.Next()
			return
		}

		// For unsafe methods (POST, PUT, DELETE, PATCH), validate the CSRF token.
		clientToken := extractCSRFToken(c, cfg)
		if clientToken == "" {
			slog.Warn("CSRF token missing", "method", c.Request.Method, "path", c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "CSRF token 缺失",
			})
			c.Abort()
			return
		}

		if !csrfStore.validate(clientToken) {
			slog.Warn("CSRF token invalid or expired", "method", c.Request.Method, "path", c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "CSRF token 无效或已过期",
			})
			c.Abort()
			return
		}

		// Also validate that the cookie token matches the header token.
		cookieToken, err := c.Cookie(cfg.CookieName)
		if err != nil {
			slog.Warn("CSRF cookie missing")
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "CSRF cookie 缺失",
			})
			c.Abort()
			return
		}

		if subtle.ConstantTimeCompare([]byte(clientToken), []byte(cookieToken)) != 1 {
			slog.Warn("CSRF token mismatch between header and cookie")
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "CSRF 验证失败",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isCSRFExempt checks if a path is exempt from CSRF checks.
func isCSRFExempt(path string, cfg CSRFConfig) bool {
	for _, prefix := range cfg.ExemptPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	for _, exempt := range cfg.ExemptPaths {
		if path == exempt {
			return true
		}
	}
	return false
}

// isSafeMethod returns true for HTTP methods that don't modify state.
func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

// extractCSRFToken extracts the token from the configured location.
func extractCSRFToken(c *gin.Context, cfg CSRFConfig) string {
	parts := strings.SplitN(cfg.TokenLookup, ":", 2)
	if len(parts) != 2 {
		return ""
	}
	switch parts[0] {
	case "header":
		return c.GetHeader(parts[1])
	case "form":
		return c.PostForm(parts[1])
	case "query":
		return c.Query(parts[1])
	default:
		return ""
	}
}

// setCSRFCookie sets the CSRF token as a cookie with appropriate SameSite attributes.
func setCSRFCookie(c *gin.Context, cfg CSRFConfig, token string) {
	sameSite := parseSameSite(cfg.CookieSameSite)

	c.SetSameSite(sameSite)
	c.SetCookie(
		cfg.CookieName,
		token,
		int(cfg.TokenExpiry.Seconds()),
		cfg.CookiePath,
		cfg.CookieDomain,
		cfg.CookieSecure,
		cfg.CookieHTTPOnly,
	)
}

// parseSameSite converts a string to http.SameSite.
func parseSameSite(s string) http.SameSite {
	switch strings.ToLower(s) {
	case "strict":
		return http.SameSiteStrictMode
	case "lax":
		return http.SameSiteLaxMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

// store and validate manage the token store.
func (s *csrfTokenStore) store(token string, expiry time.Time) {
	s.mu.Lock()
	s.tokens[token] = expiry
	s.mu.Unlock()
}

func (s *csrfTokenStore) validate(token string) bool {
	s.mu.RLock()
	expiry, exists := s.tokens[token]
	s.mu.RUnlock()

	if !exists {
		return false
	}
	return time.Now().Before(expiry)
}

func (s *csrfTokenStore) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for token, expiry := range s.tokens {
			if now.After(expiry) {
				delete(s.tokens, token)
			}
		}
		s.mu.Unlock()
	}
}
