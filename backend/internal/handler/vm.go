package handler

import (
	"net/http"
	"strconv"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/quota"
	"github.com/cinagroup/cinaseek/backend/internal/service"
	"github.com/cinagroup/cinaseek/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type VMHandler struct {
	cfg       *config.Config
	vmService *service.VMService
}

func NewVMHandler(cfg *config.Config, vmService *service.VMService) *VMHandler {
	return &VMHandler{cfg: cfg, vmService: vmService}
}

// ListVMs 获取虚拟机列表
func (h *VMHandler) ListVMs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	req := &service.ListVMsRequest{
		Name:     c.Query("name"),
		Status:   c.Query("status"),
		Page:     0,
		PageSize: 0,
	}
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		req.Page, _ = strconv.Atoi(pageStr)
	}
	if pageSizeStr := c.DefaultQuery("page_size", "10"); pageSizeStr != "" {
		req.PageSize, _ = strconv.Atoi(pageSizeStr)
	}

	vms, total, err := h.vmService.ListVMs(userID.(uint), req)
	if err != nil {
		response.InternalError(c, "查询失败："+err.Error())
		return
	}

	response.SuccessWithPage(c, vms, total, req.Page, req.PageSize)
}

// GetVMDetail 获取虚拟机详情
func (h *VMHandler) GetVMDetail(c *gin.Context) {
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

	vm, err := h.vmService.GetVM(uint(vmID), userID.(uint))
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, vm)
}

// CreateVM 创建虚拟机
func (h *VMHandler) CreateVM(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.CreateVMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 配额检查
	plan := "free" // TODO: 从用户订阅获取
	currentVMCount, err := h.vmService.CountVMs(userID.(uint))
	if err != nil {
		response.InternalError(c, "查询配额失败："+err.Error())
		return
	}
	if err := quota.CheckQuota(plan, "vms", int(currentVMCount)); err != nil {
		response.Error(c, http.StatusPaymentRequired, err.Error(), nil)
		return
	}

	vm, err := h.vmService.CreateVM(userID.(uint), &req)
	if err != nil {
		response.InternalError(c, "创建失败："+err.Error())
		return
	}

	response.SuccessWithMessage(c, "虚拟机创建成功", vm)
}

// OperateVM 操作虚拟机（启动/停止/重启/暂停/删除）
func (h *VMHandler) OperateVM(c *gin.Context) {
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

	var req service.OperateVMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.vmService.OperateVM(uint(vmID), userID.(uint), &req); err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		if err == service.ErrVMOperationFail {
			response.InternalError(c, "虚拟机操作失败")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "操作成功", nil)
}

// UpdateVMConfig 更新虚拟机配置
func (h *VMHandler) UpdateVMConfig(c *gin.Context) {
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

	var req service.UpdateVMConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	vm, err := h.vmService.UpdateVMConfig(uint(vmID), userID.(uint), &req)
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "更新失败："+err.Error())
		return
	}

	response.SuccessWithMessage(c, "配置更新成功", vm)
}

// CreateSnapshot 创建快照
func (h *VMHandler) CreateSnapshot(c *gin.Context) {
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

	var req service.CreateSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	snapshot, err := h.vmService.CreateSnapshot(uint(vmID), userID.(uint), &req)
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "创建快照失败："+err.Error())
		return
	}

	response.SuccessWithMessage(c, "快照创建成功", snapshot)
}

// ListSnapshots 获取快照列表
func (h *VMHandler) ListSnapshots(c *gin.Context) {
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

	snapshots, err := h.vmService.ListSnapshots(uint(vmID), userID.(uint))
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, snapshots)
}

// RestoreSnapshot 恢复快照
func (h *VMHandler) RestoreSnapshot(c *gin.Context) {
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

	var req service.RestoreSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.vmService.RestoreSnapshot(uint(vmID), userID.(uint), &req); err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		if err == service.ErrSnapshotNotFound {
			response.NotFound(c, "快照不存在")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "快照恢复成功", nil)
}

// DeleteSnapshot 删除快照
func (h *VMHandler) DeleteSnapshot(c *gin.Context) {
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

	snapshotID, err := strconv.ParseUint(c.Param("snapshot_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的快照 ID")
		return
	}

	if err := h.vmService.DeleteSnapshot(uint(vmID), uint(snapshotID), userID.(uint)); err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		if err == service.ErrSnapshotNotFound {
			response.NotFound(c, "快照不存在")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "快照删除成功", nil)
}

// GetVMLogs 获取虚拟机操作日志
func (h *VMHandler) GetVMLogs(c *gin.Context) {
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

	limit := 100
	if l := c.Query("limit"); l != "" {
		if v, e := strconv.Atoi(l); e == nil && v > 0 {
			limit = v
		}
	}

	logs, err := h.vmService.GetVMLogs(uint(vmID), userID.(uint), limit)
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, logs)
}

// GetVMMetrics 获取虚拟机监控指标
func (h *VMHandler) GetVMMetrics(c *gin.Context) {
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

	limit := 100
	if l := c.Query("limit"); l != "" {
		if v, e := strconv.Atoi(l); e == nil && v > 0 {
			limit = v
		}
	}

	metrics, err := h.vmService.GetVMMetrics(uint(vmID), userID.(uint), limit)
	if err != nil {
		if err == service.ErrVMNotFound {
			response.NotFound(c, "虚拟机不存在")
			return
		}
		response.InternalError(c, "查询失败："+err.Error())
		return
	}

	response.Success(c, metrics)
}
