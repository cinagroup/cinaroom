package tests

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/handler"
	"github.com/cinagroup/cinaseek/backend/internal/oauth"

	"crypto/rand"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestCinaTokenClient_GetAuthorizationURL 测试生成授权 URL
func TestCinaTokenClient_GetAuthorizationURL(t *testing.T) {
	cfg := &oauth.CinaTokenConfig{
		BaseURL:     "https://cinatoken.com",
		ClientID:    "test-client-id",
		RedirectURI: "http://localhost:3000/oauth/callback",
		Scopes:      "user:read user:email",
	}

	client := oauth.NewCinaTokenClient(cfg)
	authURL := client.GetAuthorizationURL("test-state-123")

	assert.Contains(t, authURL, "https://cinatoken.com/oauth/authorize")
	assert.Contains(t, authURL, "client_id=test-client-id")
	assert.Contains(t, authURL, "redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Foauth%2Fcallback")
	assert.Contains(t, authURL, "scope=user%3Aread+user%3Aemail")
	assert.Contains(t, authURL, "state=test-state-123")
}

// TestOAuthHandler_Providers 测试获取 OAuth 提供商列表
func TestOAuthHandler_Providers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Load()
	handler := handler.NewOAuthHandler(cfg)

	router := gin.Default()
	router.GET("/api/v1/oauth/providers", handler.GetOAuthProviders)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/oauth/providers", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Code int `json:"code"`
		Data struct {
			Providers []struct {
				Name    string `json:"name"`
				Display string `json:"display_name"`
				Enabled bool   `json:"enabled"`
			} `json:"providers"`
			Note string `json:"note"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Greater(t, len(response.Data.Providers), 0)
}

// TestGenerateState 测试随机 state 生成
func TestGenerateState(t *testing.T) {
	// Use crypto/rand to generate two states
	state1 := generateTestState()
	state2 := generateTestState()

	assert.NotEqual(t, state1, state2)
	assert.Len(t, state1, 64) // 32 bytes = 64 hex characters
	assert.Len(t, state2, 64)
}

// generateTestState generates a random state using crypto/rand
func generateTestState() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		t := testing.T{}
		t.Fatalf("failed to generate random state: %v", err)
	}
	return hex.EncodeToString(bytes)
}

// TestOAuthConfig_EnvironmentVariables 测试环境变量加载
func TestOAuthConfig_EnvironmentVariables(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("CINATOKEN_BASE_URL", "https://test.cinatoken.com")
	os.Setenv("CINATOKEN_CLIENT_ID", "test-client-id")
	os.Setenv("CINATOKEN_CLIENT_SECRET", "test-secret")
	os.Setenv("CINATOKEN_REDIRECT_URI", "http://test.com/callback")
	os.Setenv("CINATOKEN_SCOPES", "user:read")

	defer func() {
		os.Unsetenv("CINATOKEN_BASE_URL")
		os.Unsetenv("CINATOKEN_CLIENT_ID")
		os.Unsetenv("CINATOKEN_CLIENT_SECRET")
		os.Unsetenv("CINATOKEN_REDIRECT_URI")
		os.Unsetenv("CINATOKEN_SCOPES")
	}()

	cfg := config.Load()

	assert.Equal(t, "https://test.cinatoken.com", cfg.CinaToken.BaseURL)
	assert.Equal(t, "test-client-id", cfg.CinaToken.ClientID)
	assert.Equal(t, "test-secret", cfg.CinaToken.ClientSecret)
	assert.Equal(t, "http://test.com/callback", cfg.CinaToken.RedirectURI)
	assert.Equal(t, "user:read", cfg.CinaToken.Scopes)
}
