package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"multipass-backend/internal/config"
	"multipass-backend/internal/model"
	"multipass-backend/internal/repository"
	"multipass-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type OpenClawHandler struct {
	cfg *config.Config
}

func NewOpenClawHandler(cfg *config.Config) *OpenClawHandler {
	return &OpenClawHandler{cfg: cfg}
}

// GetOpenClawStatus 获取 OpenClaw 状态
func (h *OpenClawHandler) GetOpenClawStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	vmIDStr := c.Query("vm_id")
	if vmIDStr == "" {
		response.BadRequest(c, "缺少 vm_id 参数")
		return
	}

	vmID, err := strconv.ParseUint(vmIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的虚拟机 ID")
		return
	}

	db := repository.GetDB()

	// 验证虚拟机归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 获取 OpenClaw 配置
	var ocConfig model.OpenClawConfig
	if err := db.Where("vm_id = ?", vmID).First(&ocConfig).Error; err != nil {
		response.NotFound(c, "OpenClaw 未配置")
		return
	}

	response.Success(c, gin.H{
		"status":      ocConfig.Status,
		"version":     ocConfig.Version,
		"running_time": ocConfig.RunningTime,
		"last_deployed_at": ocConfig.LastDeployedAt,
	})
}

// DeployOpenClaw 部署 OpenClaw
func (h *OpenClawHandler) DeployOpenClaw(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		VMID       uint   `json:"vm_id" binding:"required"`
		Version    string `json:"version"`
		APIKey     string `json:"api_key"`
		DefaultModel string `json:"default_model"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	db := repository.GetDB()

	// 验证虚拟机归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", req.VMID, userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 设置默认值
	if req.Version == "" {
		req.Version = "latest"
	}
	if req.DefaultModel == "" {
		req.DefaultModel = "qwencode/qwen3.5-plus"
	}

	// 创建或更新配置
	var ocConfig model.OpenClawConfig
	result := db.Where("vm_id = ?", req.VMID).First(&ocConfig)

	ocConfig.VMID = req.VMID
	ocConfig.Version = req.Version
	ocConfig.APIKey = req.APIKey
	ocConfig.DefaultModel = req.DefaultModel
	ocConfig.Status = "deploying"

	now := time.Now()
	ocConfig.LastDeployedAt = &now

	if result.Error != nil {
		db.Create(&ocConfig)
	} else {
		db.Save(&ocConfig)
	}

	// TODO: 实际部署 OpenClaw
	// 1. SSH 连接到虚拟机
	// 2. 安装依赖（Node.js, pnpm 等）
	// 3. 克隆或更新 OpenClaw
	// 4. 安装依赖
	// 5. 配置 openclaw.json
	// 6. 启动服务

	// 模拟部署过程
	go func() {
		// 模拟部署延迟
		time.Sleep(5 * time.Second)
		
		ocConfig.Status = "running"
		ocConfig.RunningTime = 0
		db.Save(&ocConfig)
	}()

	response.SuccessWithMessage(c, "OpenClaw 部署已开始", gin.H{
		"vm_id":     req.VMID,
		"version":   req.Version,
		"status":    "deploying",
	})
}

// OperateOpenClaw 操作 OpenClaw（启动/停止/重启/更新）
func (h *OpenClawHandler) OperateOpenClaw(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	vmIDStr := c.Param("id")
	vmID, err := strconv.ParseUint(vmIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的虚拟机 ID")
		return
	}

	var req struct {
		Operation string `json:"operation" binding:"required"` // start, stop, restart, update
		Version   string `json:"version"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	db := repository.GetDB()

	// 验证虚拟机归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 获取 OpenClaw 配置
	var ocConfig model.OpenClawConfig
	if err := db.Where("vm_id = ?", vmID).First(&ocConfig).Error; err != nil {
		response.NotFound(c, "OpenClaw 未配置")
		return
	}

	var message string

	switch req.Operation {
	case "start":
		// TODO: SSH 连接到虚拟机，启动 OpenClaw 服务
		ocConfig.Status = "running"
		message = "OpenClaw 已启动"

	case "stop":
		// TODO: SSH 连接到虚拟机，停止 OpenClaw 服务
		ocConfig.Status = "stopped"
		message = "OpenClaw 已停止"

	case "restart":
		// TODO: SSH 连接到虚拟机，重启 OpenClaw 服务
		ocConfig.Status = "running"
		message = "OpenClaw 已重启"

	case "update":
		if req.Version == "" {
			response.BadRequest(c, "更新操作需要指定版本号")
			return
		}
		// TODO: SSH 连接到虚拟机，更新 OpenClaw 到指定版本
		ocConfig.Version = req.Version
		ocConfig.Status = "updating"
		message = fmt.Sprintf("OpenClaw 正在更新到 %s", req.Version)

	default:
		response.BadRequest(c, "不支持的操作")
		return
	}

	db.Save(&ocConfig)

	response.SuccessWithMessage(c, message, ocConfig)
}

// GetOpenClawLog 获取 OpenClaw 日志
func (h *OpenClawHandler) GetOpenClawLog(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	vmIDStr := c.Param("id")
	vmID, err := strconv.ParseUint(vmIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的虚拟机 ID")
		return
	}

	db := repository.GetDB()

	// 验证虚拟机归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// TODO: SSH 连接到虚拟机，获取 OpenClaw 日志文件内容
	// 例如：tail -n 100 /root/.openclaw/logs/openclaw.log

	// 模拟日志数据
	logs := []string{
		"[INFO] OpenClaw started successfully",
		"[INFO] Loaded 15 skills",
		"[INFO] Gateway is running on port 3000",
		"[INFO] Model: qwencode/qwen3.5-plus",
	}

	response.Success(c, gin.H{
		"vm_id": vmID,
		"logs":  logs,
	})
}

// UpdateOpenClawConfig 更新 OpenClaw 配置
func (h *OpenClawHandler) UpdateOpenClawConfig(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	vmIDStr := c.Param("id")
	vmID, err := strconv.ParseUint(vmIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的虚拟机 ID")
		return
	}

	var req struct {
		DefaultModel   string   `json:"default_model"`
		APIKey         string   `json:"api_key"`
		EnabledTools   []string `json:"enabled_tools"`
		EnabledSkills  []string `json:"enabled_skills"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	db := repository.GetDB()

	// 验证虚拟机归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 获取 OpenClaw 配置
	var ocConfig model.OpenClawConfig
	if err := db.Where("vm_id = ?", vmID).First(&ocConfig).Error; err != nil {
		response.NotFound(c, "OpenClaw 未配置")
		return
	}

	// 更新配置
	if req.DefaultModel != "" {
		ocConfig.DefaultModel = req.DefaultModel
	}
	if req.APIKey != "" {
		ocConfig.APIKey = req.APIKey
	}
	if req.EnabledTools != nil {
		// 转换为 JSON 字符串存储
		toolsJSON := "["
		for i, tool := range req.EnabledTools {
			if i > 0 {
				toolsJSON += ","
			}
			toolsJSON += fmt.Sprintf(`"%s"`, tool)
		}
		toolsJSON += "]"
		ocConfig.EnabledTools = toolsJSON
	}
	if req.EnabledSkills != nil {
		skillsJSON := "["
		for i, skill := range req.EnabledSkills {
			if i > 0 {
				skillsJSON += ","
			}
			skillsJSON += fmt.Sprintf(`"%s"`, skill)
		}
		skillsJSON += "]"
		ocConfig.EnabledSkills = skillsJSON
	}

	db.Save(&ocConfig)

	// TODO: 更新虚拟机中的 openclaw.json 配置文件

	response.SuccessWithMessage(c, "配置更新成功", ocConfig)
}

// GetOpenClawMonitor 获取 OpenClaw 监控数据
func (h *OpenClawHandler) GetOpenClawMonitor(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	vmIDStr := c.Query("vm_id")
	if vmIDStr == "" {
		response.BadRequest(c, "缺少 vm_id 参数")
		return
	}

	vmID, err := strconv.ParseUint(vmIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的虚拟机 ID")
		return
	}

	db := repository.GetDB()

	// 验证虚拟机归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 获取 OpenClaw 配置
	var ocConfig model.OpenClawConfig
	if err := db.Where("vm_id = ?", vmID).First(&ocConfig).Error; err != nil {
		response.NotFound(c, "OpenClaw 未配置")
		return
	}

	// TODO: 从虚拟机获取实际监控数据
	// 这里返回模拟数据
	response.Success(c, gin.H{
		"cpu_usage":     25.5,
		"memory_usage":  512.3,
		"disk_usage":    15.2,
		"today_requests": 1250,
		"avg_response_time": 1.2,
		"error_rate":    0.5,
		"active_sessions": 5,
		"enabled_tools_count": 10,
		"enabled_skills_count": 15,
	})
}

// GetWorkspaceList 获取工作空间列表
func (h *OpenClawHandler) GetWorkspaceList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	vmIDStr := c.Query("vm_id")
	if vmIDStr == "" {
		response.BadRequest(c, "缺少 vm_id 参数")
		return
	}

	vmID, err := strconv.ParseUint(vmIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的虚拟机 ID")
		return
	}

	db := repository.GetDB()

	// 验证虚拟机归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// TODO: SSH 连接到虚拟机，扫描工作空间目录
	// 返回工作空间列表和统计信息

	// 模拟数据
	workspaces := []gin.H{
		{
			"name":   "default",
			"path":   "/root/.openclaw/workspace",
			"size":   1024 * 1024 * 500, // 500MB
			"file_count": 1250,
			"last_modified": time.Now(),
		},
	}

	response.Success(c, gin.H{
		"workspaces": workspaces,
		"total_size": 1024 * 1024 * 500,
		"total_files": 1250,
	})
}
