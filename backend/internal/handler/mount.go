package handler

import (
	"strconv"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/service"
	"github.com/cinagroup/cinaseek/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type MountHandler struct {
	cfg          *config.Config
	mountService *service.MountService
}

func NewMountHandler(cfg *config.Config, mountService *service.MountService) *MountHandler {
	return &MountHandler{cfg: cfg, mountService: mountService}
}

// ListMounts 获取挂载列表
func (h *MountHandler) ListMounts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var vmID uint
	if vmIDStr := c.Query("vm_id"); vmIDStr != "" {
		if v, err := strconv.ParseUint(vmIDStr, 10, 32); err == nil {
			vmID = uint(v)
		}
	}

	mounts, err := h.mountService.ListMounts(userID.(uint), vmID)
	if err != nil {
		response.InternalError(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, mounts)
}

// AddMount 添加挂载
func (h *MountHandler) AddMount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.AddMountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	mount, err := h.mountService.AddMount(userID.(uint), &req)
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "创建失败："+err.Error())
		return
	}

	response.SuccessWithMessage(c, "挂载添加成功", mount)
}

// OperateMount 操作挂载（挂载/卸载/编辑/删除）
func (h *MountHandler) OperateMount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	mountID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的挂载 ID")
		return
	}

	var req service.OperateMountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.mountService.OperateMount(uint(mountID), userID.(uint), &req); err != nil {
		if err == service.ErrMountNotFound {
			response.NotFound(c, "挂载不存在")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	var message string
	switch req.Operation {
	case "mount":
		message = "挂载成功"
	case "unmount":
		message = "卸载成功"
	case "edit":
		message = "编辑成功"
	case "delete":
		message = "删除成功"
	default:
		message = "操作成功"
	}

	response.SuccessWithMessage(c, message, nil)
}

// GetOpenClawConfig 获取 OpenClaw 专属配置
func (h *MountHandler) GetOpenClawConfig(c *gin.Context) {
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

	result, err := h.mountService.GetOpenClawConfig(uint(vmID), userID.(uint))
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		if err == service.ErrOpenClawNotConfigured {
			response.NotFound(c, "OpenClaw 配置不存在")
			return
		}
		response.InternalError(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, result)
}

// ConfigureOpenClawMount 配置 OpenClaw 挂载
func (h *MountHandler) ConfigureOpenClawMount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.ConfigureOpenClawMountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	cfg, err := h.mountService.ConfigureOpenClawMount(userID.(uint), &req)
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "配置失败："+err.Error())
		return
	}

	response.SuccessWithMessage(c, "OpenClaw 挂载配置成功", cfg)
}
