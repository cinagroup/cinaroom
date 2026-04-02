package tests

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/middleware"
	"github.com/cinagroup/cinaseek/backend/internal/pool"

	"github.com/gin-gonic/gin"
)

// BenchmarkVMList benchmarks the JWT auth + list handler path.
func BenchmarkVMList(b *testing.B) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	jwtCfg := &config.JWTConfig{Secret: "bench-secret", ExpireTime: 24 * time.Hour}
	r.Use(middleware.JWTAuth(jwtCfg))
	r.GET("/api/v1/vm/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []struct{}{}})
	})

	token := generateTestToken("bench-secret", 1, "benchuser", time.Now().Add(24*time.Hour))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/v1/vm/list", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}
	})
}

// BenchmarkAuthLogin benchmarks JWT token generation.
func BenchmarkAuthLogin(b *testing.B) {
	cfg := &config.JWTConfig{Secret: "bench-secret", ExpireTime: 24 * time.Hour}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = middleware.GenerateToken(cfg, uint(i%1000), "benchuser")
	}
}

// BenchmarkTokenGeneration benchmarks token creation and validation separately.
func BenchmarkTokenGeneration(b *testing.B) {
	cfg := &config.JWTConfig{Secret: "bench-secret", ExpireTime: 24 * time.Hour}

	b.Run("generate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = middleware.GenerateToken(cfg, 1, "user1")
		}
	})

	validToken, _ := middleware.GenerateToken(cfg, 1, "user1")
	validateCfg := &config.JWTConfig{Secret: "bench-secret", ExpireTime: 24 * time.Hour}

	b.Run("validate", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r := newTestRouter()
			r.Use(middleware.JWTAuth(validateCfg))
			r.GET("/", okHandler)
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "Bearer "+validToken)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}
	})
}

// BenchmarkRateLimitCheck benchmarks the rate limiter middleware.
func BenchmarkRateLimitCheck(b *testing.B) {
	r := newTestRouter()
	r.Use(middleware.RateLimiter(middleware.RateLimiterConfig{
		Rate:  1000,
		Burst: 100,
	}))
	r.GET("/", okHandler)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}
	})
}

// BenchmarkConnectionPool benchmarks the GoroutinePool Submit/execute path.
func BenchmarkConnectionPool(b *testing.B) {
	p := pool.NewGoroutinePool(100)
	defer p.Stop()

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		_ = p.Submit(func() {
			wg.Done()
		})
	}
	wg.Wait()
}

// BenchmarkWSPool_RegisterUnregister benchmarks WS connection lifecycle.
func BenchmarkWSPool_RegisterUnregister(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wp := pool.NewWSPool()
		for j := 0; j < 100; j++ {
			userID := string(rune(j))
			wp.Register(userID, &mockWSConn{})
		}
		for j := 0; j < 50; j++ {
			wp.Unregister(string(rune(j)))
		}
		wp.CloseAll()
	}
}

// BenchmarkWSPool_Send benchmarks message sending to connected users.
func BenchmarkWSPool_Send(b *testing.B) {
	wp := pool.NewWSPool()
	defer wp.CloseAll()

	for i := 0; i < 100; i++ {
		userID := string(rune(i))
		wp.Register(userID, &mockWSConn{})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := string(rune(i % 100))
		_ = wp.Send(userID, []byte("benchmark message"))
	}
}

// BenchmarkGRPCPool_Creation benchmarks pool creation overhead.
func BenchmarkGRPCPool_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := pool.NewGRPCPool("/tmp/bench.sock", 10)
		p.Close()
	}
}
