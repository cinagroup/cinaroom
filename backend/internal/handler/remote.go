package handler

import (
	"net/http"
	"strconv"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/model"
	"github.com/cinagroup/cinaseek/backend/internal/repository"
	"github.com/cinagroup/cinaseek/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type RemoteHandler struct {
	cfg *config.Config
}

func NewRemoteHandler(cfg *config.Config) *RemoteHandler {
	return &RemoteHandler{cfg: cfg}
}

// GetRemoteStatus 获取远程访问状态
func (h *RemoteHandler) GetRemoteStatus(c *gin.Context) {
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

	// 获取远程访问配置
	var remoteAccess model.RemoteAccess
	if err := db.Where("vm_id = ?", vmID).First(&remoteAccess).Error; err != nil {
		// 不存在则返回默认状态
		response.Success(c, gin.H{
			"enabled": false,
			"status":  "disabled",
			"access_address": "",
			"qr_code": "",
		})
		return
	}

	status := "disabled"
	if remoteAccess.Enabled {
		status = "enabled"
	}

	response.Success(c, gin.H{
		"enabled":        remoteAccess.Enabled,
		"status":         status,
		"access_address": remoteAccess.AccessAddress,
		"qr_code":        remoteAccess.QRCode,
	})
}

// SwitchRemoteAccess 切换远程访问开关
func (h *RemoteHandler) SwitchRemoteAccess(c *gin.Context) {
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
		Enabled bool `json:"enabled" binding:"required"`
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

	// 获取或创建远程访问配置
	var remoteAccess model.RemoteAccess
	result := db.Where("vm_id = ?", vmID).First(&remoteAccess)

	remoteAccess.VMID = uint(vmID)
	remoteAccess.Enabled = req.Enabled

	if req.Enabled {
		// TODO: 生成访问地址和二维码
		// 例如：https://<vm-ip>:8080
		remoteAccess.AccessAddress = "https://" + vm.IP + ":8080"
		// TODO: 生成二维码（可以使用 qrcode 库）
		remoteAccess.QRCode = "data:image/png;base64,..."
	} else {
		remoteAccess.AccessAddress = ""
		remoteAccess.QRCode = ""
	}

	if result.Error != nil {
		db.Create(&remoteAccess)
	} else {
		db.Save(&remoteAccess)
	}

	message := "远程访问已启用"
	if !req.Enabled {
		message = "远程访问已禁用"
	}

	response.SuccessWithMessage(c, message, remoteAccess)
}

// GetIPWhitelist 获取 IP 白名单
func (h *RemoteHandler) GetIPWhitelist(c *gin.Context) {
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

	// 获取白名单列表
	var whitelist []model.IPWhitelist
	if err := db.Where("vm_id = ?", vmID).Order("created_at DESC").Find(&whitelist).Error; err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, whitelist)
}

// AddIPWhitelist 添加 IP 白名单
func (h *RemoteHandler) AddIPWhitelist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		VMID uint   `json:"vm_id" binding:"required"`
		IP   string `json:"ip" binding:"required"`
		Note string `json:"note"`
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

	// 创建白名单
	whitelist := model.IPWhitelist{
		VMID: req.VMID,
		IP:   req.IP,
		Note: req.Note,
	}

	if err := db.Create(&whitelist).Error; err != nil {
		response.InternalError(c, "添加失败："+err.Error())
		return
	}

	// TODO: 更新防火墙规则，允许该 IP 访问

	response.SuccessWithMessage(c, "IP 白名单添加成功", whitelist)
}

// RemoveIPWhitelist 删除 IP 白名单
func (h *RemoteHandler) RemoveIPWhitelist(c *gin.Context) {
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

	whitelistIDStr := c.Param("whitelist_id")
	whitelistID, err := strconv.ParseUint(whitelistIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的白名单 ID")
		return
	}

	db := repository.GetDB()

	// 验证虚拟机归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 删除白名单
	result := db.Where("id = ? AND vm_id = ?", uint(whitelistID), vmID).Delete(&model.IPWhitelist{})
	if result.RowsAffected == 0 {
		response.NotFound(c, "白名单不存在")
		return
	}

	// TODO: 更新防火墙规则，移除该 IP

	response.SuccessWithMessage(c, "IP 白名单删除成功", nil)
}

// GetRemoteLog 获取远程访问日志
func (h *RemoteHandler) GetRemoteLog(c *gin.Context) {
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

	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "50")
	ipFilter := c.Query("ip")
	statusFilter := c.Query("status")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	db := repository.GetDB()

	// 验证虚拟机归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 构建查询
	query := db.Model(&model.RemoteLog{}).Where("vm_id = ?", vmID)

	if ipFilter != "" {
		query = query.Where("access_ip = ?", ipFilter)
	}
	if statusFilter != "" {
		query = query.Where("response_code = ?", statusFilter)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var logs []model.RemoteLog
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("access_time DESC").Find(&logs)

	response.SuccessWithPage(c, logs, total, page, pageSize)
}
