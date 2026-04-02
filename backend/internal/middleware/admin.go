package middleware

import (
	"net/http"

	"github.com/cinagroup/cinaseek/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// AdminRequired returns a middleware that requires admin role (>=10).
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}
		// JWT claims numbers are parsed as float64
		userRole, ok := role.(float64)
		if !ok {
			// also support int (for non-JWT contexts)
			if roleInt, okInt := role.(int); okInt {
				userRole = float64(roleInt)
			} else {
				response.Unauthorized(c, "角色信息异常")
				c.Abort()
				return
			}
		}
		if int(userRole) < 10 {
			response.Error(c, http.StatusForbidden, "需要管理员权限", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

// RootRequired returns a middleware that requires root role (>=100).
func RootRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}
		userRole, ok := role.(float64)
		if !ok {
			if roleInt, okInt := role.(int); okInt {
				userRole = float64(roleInt)
			} else {
				response.Unauthorized(c, "角色信息异常")
				c.Abort()
				return
			}
		}
		if int(userRole) < 100 {
			response.Error(c, http.StatusForbidden, "需要 Root 权限", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
