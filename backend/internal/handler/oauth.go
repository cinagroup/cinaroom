package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/middleware"
	"github.com/cinagroup/cinaseek/backend/internal/model"
	"github.com/cinagroup/cinaseek/backend/internal/oauth"
	"github.com/cinagroup/cinaseek/backend/internal/repository"
	"github.com/cinagroup/cinaseek/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	cfg         *config.Config
	github      *oauth.GitHubClient
	cinatoken   *oauth.CinaTokenClient
}

func NewOAuthHandler(cfg *config.Config) *OAuthHandler {
	// GitHub OAuth（直接）
	githubClient := oauth.NewGitHubClient(&oauth.GitHubConfig{
		ClientID:     cfg.CinaToken.ClientID,     // 复用 GitHub ClientID
		ClientSecret: cfg.CinaToken.ClientSecret,  // 复用 GitHub ClientSecret
		RedirectURI:  "https://cinaseek.ai/oauth/callback",
	})

	// CinaToken OAuth（备用）
	cinatokenClient := oauth.NewCinaTokenClient(&oauth.CinaTokenConfig{
		BaseURL:      cfg.CinaToken.BaseURL,
		ClientID:     cfg.CinaToken.ClientID,
		ClientSecret: cfg.CinaToken.ClientSecret,
		RedirectURI:  cfg.CinaToken.RedirectURI,
		Scopes:       cfg.CinaToken.Scopes,
	})

	return &OAuthHandler{
		cfg:       cfg,
		github:    githubClient,
		cinatoken: cinatokenClient,
	}
}

// OAuthRedirect 重定向到 GitHub 授权页
func (h *OAuthHandler) OAuthRedirect(c *gin.Context) {
	state := generateState()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	authURL := h.github.GetAuthorizationURL(state)
	c.JSON(http.StatusOK, gin.H{
		"authorize_url": authURL,
	})
}

// OAuthCallback GitHub 授权回调
func (h *OAuthHandler) OAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		response.BadRequest(c, "Authorization code not obtained")
		return
	}

	// 验证 state
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		response.BadRequest(c, "State 验证失败")
		return
	}
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// 用授权码换取 GitHub access_token
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	accessToken, err := h.github.ExchangeCode(ctx, code)
	if err != nil {
		response.InternalError(c, "获取 Token 失败："+err.Error())
		return
	}

	// 获取 GitHub 用户信息
	userInfo, err := h.github.GetUserInfo(ctx, accessToken)
	if err != nil {
		response.InternalError(c, "获取用户信息失败："+err.Error())
		return
	}

	// 在本地数据库查找或创建用户
	db := repository.GetDB()
	var user model.User

	// 检查是否为第一个注册的用户（自动成为 root）
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)

	// 通过 GitHub ID 查找（provider + username 组合）
	if err := db.Where("provider = ? AND username = ?", "github", userInfo.Username).First(&user).Error; err != nil {
		// 尝试通过邮箱查找
		if userInfo.Email != "" {
			if err := db.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
				// 创建新用户
				newRole := 0
				if userCount == 0 {
					newRole = 100 // 第一个用户成为 root
				}
				user = model.User{
					Username: userInfo.Username,
					Email:    userInfo.Email,
					Nickname: userInfo.Nickname,
					Avatar:   userInfo.Avatar,
					Provider: "github",
					Active:   true,
					Role:     newRole,
				}
				if err := db.Create(&user).Error; err != nil {
					response.InternalError(c, "创建用户失败："+err.Error())
					return
				}
			}
		} else {
			// 无邮箱，直接创建
			newRole := 0
			if userCount == 0 {
				newRole = 100 // 第一个用户成为 root
			}
			user = model.User{
				Username: userInfo.Username,
				Email:    userInfo.Username + "@github",
				Nickname: userInfo.Nickname,
				Avatar:   userInfo.Avatar,
				Provider: "github",
				Active:   true,
				Role:     newRole,
			}
			if err := db.Create(&user).Error; err != nil {
				response.InternalError(c, "创建用户失败："+err.Error())
				return
			}
		}
	}

	// 更新用户信息
	user.Nickname = userInfo.Nickname
	user.Avatar = userInfo.Avatar
	now := time.Now()
	user.LastLoginAt = &now
	db.Save(&user)

	// 记录登录日志
	db.Create(&model.LoginLog{
		UserID: user.ID,
		IP:     c.ClientIP(),
		Device: c.Request.UserAgent(),
	})

	// 生成 CinaSeek JWT
	token, err := middleware.GenerateToken(&h.cfg.JWT, user.ID, user.Username, user.Role)
	if err != nil {
		response.InternalError(c, "Token 生成失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"provider": "github",
			"role":     user.Role,
		},
	})
}

// Logout 登出
func (h *OAuthHandler) Logout(c *gin.Context) {
	response.SuccessWithMessage(c, "登出成功", nil)
}

// GetOAuthProviders 获取支持的 OAuth 提供商
func (h *OAuthHandler) GetOAuthProviders(c *gin.Context) {
	providers := []gin.H{
		{"name": "github", "display_name": "GitHub", "enabled": true},
		{"name": "google", "display_name": "Google", "enabled": true},
		{"name": "microsoft", "display_name": "Microsoft", "enabled": true},
		{"name": "gitlab", "display_name": "GitLab", "enabled": true},
		{"name": "wechat", "display_name": "微信", "enabled": true},
		{"name": "feishu", "display_name": "飞书", "enabled": true},
		{"name": "dingtalk", "display_name": "钉钉", "enabled": true},
		{"name": "qq", "display_name": "QQ", "enabled": true},
		{"name": "weibo", "display_name": "微博", "enabled": true},
	}
	response.Success(c, gin.H{
		"providers": providers,
		"note":     "通过 CinaToken 统一认证，支持 9+ OAuth 提供商",
	})
}

func generateState() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
