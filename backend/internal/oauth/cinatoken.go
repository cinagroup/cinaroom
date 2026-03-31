package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// CinaTokenConfig OAuth 配置
type CinaTokenConfig struct {
	BaseURL      string `mapstructure:"base_url" json:"base_url"`
	ClientID     string `mapstructure:"client_id" json:"client_id"`
	ClientSecret string `mapstructure:"client_secret" json:"client_secret"`
	RedirectURI  string `mapstructure:"redirect_uri" json:"redirect_uri"`
	Scopes       string `mapstructure:"scopes" json:"scopes"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Phone     string `json:"phone"`
	Provider  string `json:"provider"`  // OAuth 提供商：github/google/microsoft 等
	CreatedAt time.Time `json:"created_at"`
}

// TokenResponse Token 响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// CinaTokenClient CinaToken OAuth 客户端
type CinaTokenClient struct {
	config *CinaTokenConfig
	client *http.Client
	cache  sync.Map // 缓存验证过的 Token
}

// TokenCacheItem 缓存项
type TokenCacheItem struct {
	UserInfo  *UserInfo
	ExpiresAt time.Time
}

// NewCinaTokenClient 创建客户端
func NewCinaTokenClient(cfg *CinaTokenConfig) *CinaTokenClient {
	return &CinaTokenClient{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetAuthorizationURL 获取授权 URL
func (c *CinaTokenClient) GetAuthorizationURL(state string) string {
	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("redirect_uri", c.config.RedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", c.config.Scopes)
	params.Set("state", state)

	return fmt.Sprintf("%s/oauth/authorize?%s", c.config.BaseURL, params.Encode())
}

// ExchangeCode 用授权码换取 Token
func (c *CinaTokenClient) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", c.config.RedirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/oauth/token", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败：%w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(nil)
	req.URL.RawQuery = data.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败：%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OAuth 服务器返回错误：%d - %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败：%w", err)
	}

	return &tokenResp, nil
}

// ValidateToken 验证 Token 并获取用户信息
func (c *CinaTokenClient) ValidateToken(ctx context.Context, accessToken string) (*UserInfo, error) {
	// 检查缓存
	if cached, ok := c.cache.Load(accessToken); ok {
		item := cached.(*TokenCacheItem)
		if time.Now().Before(item.ExpiresAt) {
			return item.UserInfo, nil
		}
		c.cache.Delete(accessToken)
	}

	// 调用 CinaToken userinfo 接口
	req, err := http.NewRequestWithContext(ctx, "GET", c.config.BaseURL+"/oauth/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败：%w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败：%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("用户信息请求失败：%d - %s", resp.StatusCode, string(body))
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("解析用户信息失败：%w", err)
	}

	// 缓存 5 分钟
	c.cache.Store(accessToken, &TokenCacheItem{
		UserInfo:  &userInfo,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	return &userInfo, nil
}

// RefreshToken 刷新 Token
func (c *CinaTokenClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/oauth/token", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败：%w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败：%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("刷新 Token 失败：%d - %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败：%w", err)
	}

	return &tokenResp, nil
}

// RevokeToken 撤销 Token
func (c *CinaTokenClient) RevokeToken(ctx context.Context, token string) error {
	data := url.Values{}
	data.Set("token", token)
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/oauth/revoke", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败：%w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败：%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("撤销 Token 失败：%d - %s", resp.StatusCode, string(body))
	}

	return nil
}
