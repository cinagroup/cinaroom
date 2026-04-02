package middleware

import (
	"net/http"

	"github.com/cinagroup/cinaseek/backend/internal/quota"
	"github.com/cinagroup/cinaseek/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// QuotaCheck checks if the user's plan allows the requested resource.
func QuotaCheck(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		plan, _ := c.Get("user_plan")
		planStr := "free"
		if p, ok := plan.(string); ok && p != "" {
			planStr = p
		}

		current, _ := c.Get("quota_current")
		currentInt := 0
		if v, ok := current.(int); ok {
			currentInt = v
		}

		if err := quota.CheckQuota(planStr, resource, currentInt); err != nil {
			response.Error(c, http.StatusPaymentRequired, err.Error(), nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
