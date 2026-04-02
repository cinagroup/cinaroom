package handler

import (
	"github.com/cinagroup/cinaseek/backend/internal/quota"
	"github.com/cinagroup/cinaseek/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct{}

func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{}
}

// GetPlans 获取所有套餐
func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
	plans := map[string]quota.PlanConfig{
		"free":       quota.GetPlan("free"),
		"pro":        quota.GetPlan("pro"),
		"enterprise": quota.GetPlan("enterprise"),
	}
	response.Success(c, plans)
}

// GetCurrentSubscription 获取当前用户订阅
func (h *SubscriptionHandler) GetCurrentSubscription(c *gin.Context) {
	userID, _ := c.Get("userID")
	// TODO: 从数据库查询
	response.Success(c, gin.H{
		"user_id": userID,
		"plan":    "free",
		"status":  "active",
	})
}
