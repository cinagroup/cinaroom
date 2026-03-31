package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"multipass-backend/internal/config"
	"multipass-backend/internal/middleware"
	"multipass-backend/internal/model"
	"multipass-backend/internal/oauth"
	"multipass-backend/internal/repository"
	"multipass-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	cfg        *config.Config
	oauthClient *oauth.CinaTokenClient
}

func NewOAuthHandler(cfg *config.Config) *OAuthHandler {
	oauthCfg := &oauth.CinaTokenConfig{
		BaseURL:      cfg.CinaToken.BaseURL,
		ClientID:     cfg.CinaToken.ClientID,
		ClientSecret: cfg.CinaToken.ClientSecret,
		RedirectURI:  cfg.CinaToken.RedirectURI,
		Scopes:       cfg.CinaToken.Scopes,
	}

	return &OAuthHandler{
		cfg:        cfg,
		oauthClient: oauth.NewCinaTokenClient(oauthCfg),
	}
}

// OAuthRedirect 重定向到 CinaToken 授权页
func (h *OAuthHandler) OAuthRedirect(c *gin.Context) {
	// 生成随机 state 防止 CSRF
	state := generateState()
	
	// 将 state 存入 session/cookie，用于回调验证
	c.SetCookie("oauth_state", state, 600, "/oauth", "", false, true)
	
	authURL := h.oauthClient.GetAuthorizationURL(state)
	c.JSON(http.StatusOK, gin.H{
		"authorize_url": authURL,
	})
}

// OAuthCallback CinaToken 授权回调
func (h *OAuthHandler) OAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	
	if code == "" {
		response.BadRequest(c, "授权码为空")
		return
	}
	
	// 验证 state
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		response.BadRequest(c, "State 验证失败，可能存在 CSRF 攻击")
		return
	}
	
	// 清除 state cookie
	c.SetCookie("oauth_state", "", -1, "/oauth", "", false, true)
	
	// 用授权码换取 Token
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	tokenResp, err := h.oauthClient.ExchangeCode(ctx, code)
	if err != nil {
		response.InternalError(c, "获取 Token 失败："+err.Error())
		return
	}
	
	// 验证 Token 并获取用户信息
	userInfo, err := h.oauthClient.ValidateToken(ctx, tokenResp.AccessToken)
	if err != nil {
		response.InternalError(c, "获取用户信息失败："+err.Error())
		return
	}
	
	// 在本地数据库查找或创建用户
	db := repository.GetDB()
	var user model.User
	
	// 优先通过 CinaToken 用户 ID 查找
	if err := db.Where("cinatoken_id = ?", userInfo.ID).First(&user).Error; err != nil {
		// 未找到，尝试通过邮箱查找
		if err := db.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
			// 邮箱也未注册，创建新用户
			user = model.User{
				CinatokenID: userInfo.ID,
				Username:    userInfo.Username,
				Email:       userInfo.Email,
				Nickname:    userInfo.Nickname,
				Avatar:      userInfo.Avatar,
				Phone:       userInfo.Phone,
				Provider:    userInfo.Provider,
				Active:      true,
			}
			
			if err := db.Create(&user).Error; err != nil {
				response.InternalError(c, "创建用户失败："+err.Error())
				return
			}
		} else {
			// 邮箱已注册，关联 CinaToken ID
			user.CinatokenID = userInfo.ID
			user.Provider = userInfo.Provider
			user.Avatar = userInfo.Avatar
			if userInfo.Nickname != "" {
				user.Nickname = userInfo.Nickname
			}
			db.Save(&user)
		}
	} else {
		// 已存在，更新用户信息
		user.Username = userInfo.Username
		user.Nickname = userInfo.Nickname
		user.Avatar = userInfo.Avatar
		user.Phone = userInfo.Phone
		db.Save(&user)
	}
	
	// 更新登录时间
	now := user.CreatedAt
	user.LastLoginAt = &now
	db.Save(&user)
	
	// 记录登录日志
	loginLog := model.LoginLog{
		UserID: user.ID,
		IP:     c.ClientIP(),
		Device: c.Request.UserAgent(),
	}
	db.Create(&loginLog)
	
	// 生成 CinaRoom 的 JWT Token
	token, err := middleware.GenerateToken(&h.cfg.JWT, user.ID, user.Username)
	if err != nil {
		response.InternalError(c, "Token 生成失败")
		return
	}
	
	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"nickname":  user.Nickname,
			"avatar":    user.Avatar,
			"provider":  user.Provider,
		},
		"oauth": gin.H{
			"provider":    userInfo.Provider,
			"expires_in":  tokenResp.ExpiresIn,
		},
	})
}

// Logout 登出（撤销 CinaToken Token）
func (h *OAuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	
	// 获取 CinaToken Token（从请求头或数据库）
	cinatokenToken := c.GetHeader("X-CinaToken-Token")
	if cinatokenToken != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		h.oauthClient.RevokeToken(ctx, cinatokenToken)
	}
	
	response.SuccessWithMessage(c, "登出成功", nil)
}

// GetOAuthProviders 获取支持的 OAuth 提供商列表
func (h *OAuthHandler) GetOAuthProviders(c *gin.Context) {
	// 这里可以动态获取 CinaToken 支持的 Provider
	// 暂时返回硬编码列表
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
		"note": "通过 CinaToken 统一认证，支持 9+ OAuth 提供商",
	})
}

// generateState 生成随机 state
func generateState() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
