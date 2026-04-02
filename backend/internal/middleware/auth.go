package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT payload.
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

// RefreshWindow is how long before expiry we allow a proactive refresh.
const RefreshWindow = 30 * time.Minute

// JWTAuth returns a middleware that validates Bearer tokens with auto-refresh.
func JWTAuth(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "未提供认证信息", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, "认证格式错误，应为 Bearer <token>", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Secret), nil
		})

		if err != nil || !token.Valid {
			slog.Warn("JWT validation failed", "error", err)
			response.Error(c, http.StatusUnauthorized, "认证失败或 Token 已过期", nil)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		// Proactive refresh: if the token is within RefreshWindow of expiry,
		// generate a fresh token and return it in the response header.
		if shouldRefresh(claims) {
			newToken, refreshErr := GenerateToken(cfg, claims.UserID, claims.Username, claims.Role)
			if refreshErr == nil {
				c.Header("X-New-Token", newToken)
				c.Header("X-Token-Refreshed", "true")
				slog.Debug("JWT token proactively refreshed",
					"user_id", claims.UserID,
					"username", claims.Username,
				)
			} else {
				slog.Warn("failed to refresh token", "error", refreshErr)
			}
		}

		c.Next()
	}
}

// shouldRefresh checks if the token is close to expiry and should be refreshed.
func shouldRefresh(claims *Claims) bool {
	if claims.ExpiresAt == nil {
		return false
	}
	timeUntilExpiry := time.Until(claims.ExpiresAt.Time)
	return timeUntilExpiry > 0 && timeUntilExpiry <= RefreshWindow
}

// GenerateToken creates a signed JWT for the given user.
func GenerateToken(cfg *config.JWTConfig, userID uint, username string, role int) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cinaseek-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}
