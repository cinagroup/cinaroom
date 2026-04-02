package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/service"
	"github.com/cinagroup/cinaseek/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type OpenClawHandler struct {
	cfg             *config.Config
	openclawService *service.OpenClawService
}

func NewOpenClawHandler(cfg *config.Config, openclawService *service.OpenClawService) *OpenClawHandler {
	return &OpenClawHandler{cfg: cfg, openclawService: openclawService}
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

	result, err := h.openclawService.GetStatus(uint(vmID), userID.(uint))
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, result)
}

// DeployOpenClaw 部署 OpenClaw
func (h *OpenClawHandler) DeployOpenClaw(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.DeployOpenClawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	result, err := h.openclawService.Deploy(userID.(uint), &req)
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "部署失败："+err.Error())
		return
	}

	response.SuccessWithMessage(c, "OpenClaw 部署已开始", result)
}

// OperateOpenClaw 操作 OpenClaw（启动/停止/重启/更新）
func (h *OpenClawHandler) OperateOpenClaw(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	vmID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的虚拟机 ID")
		return
	}

	var req service.OperateOpenClawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	cfg, err := h.openclawService.Operate(uint(vmID), userID.(uint), &req)
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		if err == service.ErrOpenClawNotConfigured {
			response.NotFound(c, "OpenClaw 未配置")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "操作成功", cfg)
}

// GetOpenClawLog 获取 OpenClaw 日志
func (h *OpenClawHandler) GetOpenClawLog(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	vmID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的虚拟机 ID")
		return
	}

	logs, err := h.openclawService.GetLogs(uint(vmID), userID.(uint))
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "查询失败："+err.Error())
		return
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

	vmID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的虚拟机 ID")
		return
	}

	var req service.UpdateOpenClawConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	cfg, err := h.openclawService.UpdateConfig(uint(vmID), userID.(uint), &req)
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		if err == service.ErrOpenClawNotConfigured {
			response.NotFound(c, "OpenClaw 未配置")
			return
		}
		response.InternalError(c, "更新失败："+err.Error())
		return
	}

	response.SuccessWithMessage(c, "配置更新成功", cfg)
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

	data, err := h.openclawService.GetMonitorData(uint(vmID), userID.(uint))
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		if err == service.ErrOpenClawNotConfigured {
			response.NotFound(c, "OpenClaw 未配置")
			return
		}
		response.InternalError(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, data)
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

	data, err := h.openclawService.GetWorkspaceList(uint(vmID), userID.(uint))
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, data)
}

// keep time import used for future actual log fetching
var _ = time.Now
var _ = fmt.Sprintf
