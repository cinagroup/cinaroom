package handler

import (
	"net/http"
	"strconv"

	"multipass-backend/internal/config"
	"multipass-backend/internal/model"
	"multipass-backend/internal/repository"
	"multipass-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	cfg *config.Config
}

func NewSystemHandler(cfg *config.Config) *SystemHandler {
	return &SystemHandler{cfg: cfg}
}

// GetSystemSetting 获取系统设置
func (h *SystemHandler) GetSystemSetting(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	key := c.Query("key")
	if key == "" {
		response.BadRequest(c, "缺少 key 参数")
		return
	}

	db := repository.GetDB()
	var setting model.SystemSetting
	if err := db.Where("key = ?", key).First(&setting).Error; err != nil {
		response.NotFound(c, "设置不存在")
		return
	}

	response.Success(c, gin.H{
		"key":   setting.Key,
		"value": setting.Value,
	})
}

// UpdateSystemSetting 更新系统设置
func (h *SystemHandler) UpdateSystemSetting(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		Key   string `json:"key" binding:"required,max=100"`
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	db := repository.GetDB()

	// 获取或创建设置
	var setting model.SystemSetting
	result := db.Where("key = ?", req.Key).First(&setting)

	setting.Key = req.Key
	setting.Value = req.Value

	if result.Error != nil {
		db.Create(&setting)
	} else {
		db.Save(&setting)
	}

	response.SuccessWithMessage(c, "设置更新成功", setting)
}

// GetSystemVersion 获取系统版本信息
func (h *SystemHandler) GetSystemVersion(c *gin.Context) {
	response.Success(c, gin.H{
		"version":   "1.0.0",
		"build":     "20260401",
		"api_version": "v1",
		"go_version": "1.21.5",
		"features": []string{
			"user_management",
			"vm_management",
			"web_shell",
			"directory_mount",
			"openclaw_integration",
			"remote_access",
		},
	})
}

// GetAllSettings 获取所有系统设置
func (h *SystemHandler) GetAllSettings(c *gin.Context) {
	db := repository.GetDB()
	var settings []model.SystemSetting
	if err := db.Find(&settings).Error; err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	settingsMap := make(map[string]string)
	for _, setting := range settings {
		settingsMap[setting.Key] = setting.Value
	}

	response.Success(c, settingsMap)
}

// GetDashboard 获取仪表盘数据
func (h *SystemHandler) GetDashboard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	db := repository.GetDB()

	// 统计虚拟机数量
	var vmTotal int64
	db.Model(&model.VM{}).Where("user_id = ?", userID).Count(&vmTotal)

	var vmRunning int64
	db.Model(&model.VM{}).Where("user_id = ? AND status = ?", userID, "running").Count(&vmRunning)

	// 统计挂载数量
	var mountTotal int64
	db.Model(&model.Mount{}).Where("user_id = ?", userID).Count(&mountTotal)

	// 统计 OpenClaw 部署数量
	var openclawTotal int64
	db.Model(&model.OpenClawConfig{}).
		Joins("JOIN vms ON vms.id = openclaw_configs.vm_id").
		Where("vms.user_id = ?", userID).
		Count(&openclawTotal)

	response.Success(c, gin.H{
		"vm_total":      vmTotal,
		"vm_running":    vmRunning,
		"vm_stopped":    vmTotal - vmRunning,
		"mount_total":   mountTotal,
		"openclaw_total": openclawTotal,
		"recent_vms":    []model.VM{}, // 最近创建的虚拟机
		"recent_logs":   []model.VMLog{}, // 最近的操作日志
	})
}

// HealthCheck 健康检查
func (h *SystemHandler) HealthCheck(c *gin.Context) {
	db := repository.GetDB()

	// 检查数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		response.InternalError(c, "数据库连接失败")
		return
	}

	if err := sqlDB.Ping(); err != nil {
		response.InternalError(c, "数据库连接失败："+err.Error())
		return
	}

	response.Success(c, gin.H{
		"status":    "healthy",
		"database":  "connected",
		"timestamp": db.NowFunc(),
	})
}

// GetStatistics 获取统计数据
func (h *SystemHandler) GetStatistics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	db := repository.GetDB()

	// 虚拟机状态统计
	type VMStatusCount struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}

	var vmStatusCounts []VMStatusCount
	db.Model(&model.VM{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&vmStatusCounts)

	// 最近 7 天的操作统计
	type DailyOperation struct {
		Date  string `gorm:"column:date"`
		Count int64  `gorm:"column:count"`
	}

	var dailyOps []DailyOperation
	db.Model(&model.VMLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Joins("JOIN vms ON vms.id = vm_logs.vm_id").
		Where("vms.user_id = ? AND created_at >= ?", userID, db.NowFunc().AddDate(0, 0, -7)).
		Group("DATE(created_at)").
		Order("date DESC").
		Scan(&dailyOps)

	response.Success(c, gin.H{
		"vm_status_stats": vmStatusCounts,
		"daily_operations": dailyOps,
	})
}

// SearchVMs 搜索虚拟机
func (h *SystemHandler) SearchVMs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "缺少 keyword 参数")
		return
	}

	db := repository.GetDB()
	var vms []model.VM
	db.Where("user_id = ? AND (name LIKE ? OR ip LIKE ?)", userID, "%"+keyword+"%", "%"+keyword+"%").
		Limit(10).Find(&vms)

	response.Success(c, vms)
}

// BatchOperateVMs 批量操作虚拟机
func (h *SystemHandler) BatchOperateVMs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		VMIDs     []uint  `json:"vm_ids" binding:"required"`
		Operation string  `json:"operation" binding:"required"` // start, stop, restart, delete
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	db := repository.GetDB()

	// 验证所有虚拟机归属
	var vms []model.VM
	db.Where("id IN ? AND user_id = ?", req.VMIDs, userID).Find(&vms)

	if len(vms) != len(req.VMIDs) {
		response.BadRequest(c, "部分虚拟机不存在或无权限")
		return
	}

	successCount := 0
	failCount := 0

	for _, vm := range vms {
		var err error
		switch req.Operation {
		case "start":
			vm.Status = "running"
			err = db.Save(&vm).Error
		case "stop":
			vm.Status = "stopped"
			err = db.Save(&vm).Error
		case "restart":
			vm.Status = "running"
			err = db.Save(&vm).Error
		case "delete":
			err = db.Delete(&vm).Error
		}

		if err != nil {
			failCount++
		} else {
			successCount++
			// 记录日志
			vmLog := model.VMLog{
				VMID:      vm.ID,
				Operation: req.Operation,
				Result:    "success",
				Message:   "批量操作",
			}
			db.Create(&vmLog)
		}
	}

	response.Success(c, gin.H{
		"success_count": successCount,
		"fail_count":    failCount,
		"total":         len(req.VMIDs),
	})
}
