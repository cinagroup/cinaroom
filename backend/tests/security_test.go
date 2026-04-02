package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ─── helpers ────────────────────────────────────────────────────────────────

// newTestRouter creates a fresh gin.Engine in test mode.
func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// okHandler is a simple handler that returns 200 OK.
func okHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// generateToken creates a JWT signed with the given secret, expiring at exp.
func generateTestToken(secret string, userID uint, username string, exp time.Time) string {
	claims := middleware.Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cinaseek-backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(secret))
	return s
}

// jwtConfig returns a fixed JWT config for tests.
func jwtConfig() *config.JWTConfig {
	return &config.JWTConfig{
		Secret:     "test-secret-key",
		ExpireTime: 24 * time.Hour,
	}
}

// ─── Auth tests ─────────────────────────────────────────────────────────────

// TestInvalidToken_Returns401 verifies that a tampered token results in 401.
func TestInvalidToken_Returns401(t *testing.T) {
	r := newTestRouter()
	r.Use(middleware.JWTAuth(jwtConfig()))
	r.GET("/secure", okHandler)

	req, _ := http.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer this.is.not.valid")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// TestExpiredToken_Returns401 verifies that an expired token results in 401.
func TestExpiredToken_Returns401(t *testing.T) {
	cfg := jwtConfig()
	expiredToken := generateTestToken(cfg.Secret, 1, "tester", time.Now().Add(-1*time.Hour))

	r := newTestRouter()
	r.Use(middleware.JWTAuth(cfg))
	r.GET("/secure", okHandler)

	req, _ := http.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// TestNoToken_Returns401 verifies that missing Authorization header results in 401.
func TestNoToken_Returns401(t *testing.T) {
	r := newTestRouter()
	r.Use(middleware.JWTAuth(jwtConfig()))
	r.GET("/secure", okHandler)

	req, _ := http.NewRequest("GET", "/secure", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// ─── Rate limit test ────────────────────────────────────────────────────────

// TestRateLimit_Returns429 verifies that exceeding the rate limit returns 429.
func TestRateLimit_Returns429(t *testing.T) {
	r := newTestRouter()
	r.Use(middleware.RateLimiter(middleware.RateLimiterConfig{
		Rate:  100,    // very high to avoid flaky tests
		Burst: 2,      // only allow 2 burst
	}))
	r.GET("/", okHandler)

	// First two should succeed (burst = 2)
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// Third should be rate-limited
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 on 3rd burst request, got %d", w.Code)
	}
}

// ─── CSRF test ──────────────────────────────────────────────────────────────

// TestCSRFProtection verifies that a POST to a non-API path without a CSRF
// token returns 403.
func TestCSRFProtection(t *testing.T) {
	cfg := middleware.DefaultCSRFConfig()
	// Remove /api/ exemption for this test so the path is checked
	cfg.ExemptPrefixes = []string{"/health"}

	r := newTestRouter()
	r.Use(middleware.CSRF(cfg))
	r.POST("/form/submit", okHandler)

	req, _ := http.NewRequest("POST", "/form/submit", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for missing CSRF token, got %d", w.Code)
	}
}

// ─── IP filter test ─────────────────────────────────────────────────────────

// TestIPFilter_Blocked verifies that a blacklisted IP receives 403.
func TestIPFilter_Blocked(t *testing.T) {
	r := newTestRouter()
	r.Use(middleware.IPFilter(&middleware.IPFilterConfig{
		BlacklistCIDRs: []string{"192.168.1.100/32"},
		DefaultAllow:   true,
	}))
	r.GET("/", okHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.100:12345"

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for blacklisted IP, got %d", w.Code)
	}
}

// ─── Concurrency limit test ─────────────────────────────────────────────────

// TestConcurrencyLimit verifies that exceeding the concurrent connection limit
// returns 429.
func TestConcurrencyLimit(t *testing.T) {
	cfg := middleware.DefaultConcurrencyConfig()
	cfg.FreeLimit = 1 // Allow only 1 concurrent connection for free tier

	// Block handler — holds the connection until we release the channel.
	blockCh := make(chan struct{})

	r := newTestRouter()
	r.Use(middleware.ConcurrencyLimiter(cfg))
	r.GET("/", func(c *gin.Context) {
		<-blockCh // block until test releases
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	// First request — should block (holding the connection)
	req1, _ := http.NewRequest("GET", "/", nil)
	w1 := httptest.NewRecorder()

	done1 := make(chan struct{})
	go func() {
		r.ServeHTTP(w1, req1)
		close(done1)
	}()

	// Give the first request time to enter the handler
	time.Sleep(50 * time.Millisecond)

	// Second request — should be rejected (limit = 1)
	req2, _ := http.NewRequest("GET", "/", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 for concurrent connection overflow, got %d", w2.Code)
	}

	// Release the first request
	close(blockCh)
	<-done1
}
