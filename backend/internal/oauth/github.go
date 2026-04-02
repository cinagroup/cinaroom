package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GitHubConfig GitHub OAuth 配置
type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// GitHubClient GitHub OAuth 客户端
type GitHubClient struct {
	config *GitHubConfig
	client *http.Client
}

// NewGitHubClient 创建 GitHub OAuth 客户端
func NewGitHubClient(cfg *GitHubConfig) *GitHubClient {
	return &GitHubClient{
		config: cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetAuthorizationURL 获取 GitHub 授权 URL
func (g *GitHubClient) GetAuthorizationURL(state string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s",
		g.config.ClientID,
		g.config.RedirectURI,
		state,
	)
}

// ExchangeCode 用授权码换取 access_token
func (g *GitHubClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, "POST",
		"https://github.com/login/oauth/access_token"+
			"?client_id="+g.config.ClientID+
			"&client_secret="+g.config.ClientSecret+
			"&code="+code+
			"&redirect_uri="+g.config.RedirectURI,
		nil)
	req.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("GitHub token exchange failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode GitHub response failed: %w", err)
	}
	if result.Error != "" {
		return "", fmt.Errorf("GitHub OAuth error: %s", result.Error)
	}
	return result.AccessToken, nil
}

// GetUserInfo 用 access_token 获取 GitHub 用户信息
func (g *GitHubClient) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GitHub userinfo failed: %w", err)
	}
	defer resp.Body.Close()

	var ghUser struct {
		ID        uint   `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, fmt.Errorf("decode GitHub user failed: %w", err)
	}

	if ghUser.Email == "" {
		ghUser.Email = g.getUserEmail(ctx, accessToken)
	}

	return &UserInfo{
		ID:        ghUser.ID,
		Username:  ghUser.Login,
		Email:     ghUser.Email,
		Nickname:  ghUser.Name,
		Avatar:    ghUser.AvatarURL,
		Provider:  "github",
		CreatedAt: time.Now(),
	}, nil
}

func (g *GitHubClient) getUserEmail(ctx context.Context, accessToken string) string {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return ""
	}
	for _, e := range emails {
		if e.Primary {
			return e.Email
		}
	}
	if len(emails) > 0 {
		return emails[0].Email
	}
	return ""
}
