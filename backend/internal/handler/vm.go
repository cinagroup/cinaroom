package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"multipass-backend/internal/config"
	"multipass-backend/internal/model"
	"multipass-backend/internal/repository"
	"multipass-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type VMHandler struct {
	cfg *config.Config
}

func NewVMHandler(cfg *config.Config) *VMHandler {
	return &VMHandler{cfg: cfg}
}

// ListVMs 获取虚拟机列表
func (h *VMHandler) ListVMs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	db := repository.GetDB()

	// 获取查询参数
	name := c.Query("name")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 构建查询
	query := db.Model(&model.VM{}).Where("user_id = ?", userID)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var vms []model.VM
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&vms)

	response.SuccessWithPage(c, vms, total, page, pageSize)
}

// GetVMDetail 获取虚拟机详情
func (h *VMHandler) GetVMDetail(c *gin.Context) {
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
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
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

	var req struct {
		Name        string `json:"name" binding:"required,max=100"`
		Image       string `json:"image" binding:"required"`
		CPU         int    `json:"cpu" binding:"min=1,max=8"`
		Memory      int    `json:"memory" binding:"min=1,max=16"`
		Disk        int    `json:"disk" binding:"min=10,max=500"`
		NetworkType string `json:"network_type"`
		SSHKey      string `json:"ssh_key"`
		InitScript  string `json:"init_script"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 设置默认值
	if req.NetworkType == "" {
		req.NetworkType = "nat"
	}
	if req.CPU < 1 {
		req.CPU = 1
	}
	if req.Memory < 1 {
		req.Memory = 1
	}
	if req.Disk < 10 {
		req.Disk = 10
	}

	db := repository.GetDB()

	// 创建虚拟机
	vm := model.VM{
		UserID:      userID.(uint),
		Name:        req.Name,
		Status:      "stopped",
		Image:       req.Image,
		CPU:         req.CPU,
		Memory:      req.Memory,
		Disk:        req.Disk,
		NetworkType: req.NetworkType,
		SSHKey:      req.SSHKey,
		InitScript:  req.InitScript,
	}

	if err := db.Create(&vm).Error; err != nil {
		response.InternalError(c, "创建失败："+err.Error())
		return
	}

	// TODO: 这里应该调用 Multipass 命令实际创建虚拟机
	// 例如：multipass launch --name <name> --cpus <cpu> --memory <memory>G --disk <disk>G <image>

	// 记录日志
	vmLog := model.VMLog{
		VMID:      vm.ID,
		Operation: "create",
		Result:    "success",
		Message:   fmt.Sprintf("创建虚拟机：%s", vm.Name),
	}
	db.Create(&vmLog)

	response.SuccessWithMessage(c, "虚拟机创建成功", vm)
}

// OperateVM 操作虚拟机（启动/停止/重启/暂停/删除）
func (h *VMHandler) OperateVM(c *gin.Context) {
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
		Operation string `json:"operation" binding:"required"` // start, stop, restart, pause, resume, delete
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	db := repository.GetDB()
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 执行操作
	var result string
	var message string

	switch req.Operation {
	case "start":
		// TODO: 调用 multipass start <name>
		vm.Status = "running"
		result = "success"
		message = "虚拟机已启动"
	case "stop":
		// TODO: 调用 multipass stop <name>
		vm.Status = "stopped"
		result = "success"
		message = "虚拟机已停止"
	case "restart":
		// TODO: 调用 multipass restart <name>
		vm.Status = "running"
		result = "success"
		message = "虚拟机已重启"
	case "pause":
		// TODO: 调用 multipass suspend <name>
		vm.Status = "paused"
		result = "success"
		message = "虚拟机已暂停"
	case "resume":
		// TODO: 调用 multipass recover <name>
		vm.Status = "running"
		result = "success"
		message = "虚拟机已恢复"
	case "delete":
		// TODO: 调用 multipass delete <name> 和 multipass purge <name>
		result = "success"
		message = "虚拟机已删除"
	default:
		response.BadRequest(c, "不支持的操作")
		return
	}

	// 更新数据库
	if req.Operation != "delete" {
		db.Save(&vm)
	} else {
		db.Delete(&vm)
	}

	// 记录日志
	vmLog := model.VMLog{
		VMID:      vm.ID,
		Operation: req.Operation,
		Result:    result,
		Message:   message,
	}
	db.Create(&vmLog)

	response.SuccessWithMessage(c, message, nil)
}

// UpdateVMConfig 更新虚拟机配置
func (h *VMHandler) UpdateVMConfig(c *gin.Context) {
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
		CPU    int `json:"cpu" binding:"omitempty,min=1,max=8"`
		Memory int `json:"memory" binding:"omitempty,min=1,max=16"`
		Disk   int `json:"disk" binding:"omitempty,min=10,max=500"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	db := repository.GetDB()
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 更新配置
	if req.CPU > 0 {
		vm.CPU = req.CPU
	}
	if req.Memory > 0 {
		vm.Memory = req.Memory
	}
	if req.Disk > 0 {
		vm.Disk = req.Disk
	}

	if err := db.Save(&vm).Error; err != nil {
		response.InternalError(c, "更新失败")
		return
	}

	// TODO: 这里应该调用 multipass set <name>.cpus/memory/disc 实际修改配置

	// 记录日志
	vmLog := model.VMLog{
		VMID:      vm.ID,
		Operation: "update_config",
		Result:    "success",
		Message:   fmt.Sprintf("更新配置：CPU=%d, Memory=%d, Disk=%d", vm.CPU, vm.Memory, vm.Disk),
	}
	db.Create(&vmLog)

	response.SuccessWithMessage(c, "配置更新成功", vm)
}

// CreateSnapshot 创建快照
func (h *VMHandler) CreateSnapshot(c *gin.Context) {
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
		Name string `json:"name" binding:"required,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	db := repository.GetDB()
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 创建快照
	snapshot := model.VMSnapshot{
		VMID: vm.ID,
		Name: req.Name,
	}

	if err := db.Create(&snapshot).Error; err != nil {
		response.InternalError(c, "创建快照失败")
		return
	}

	// TODO: 这里应该调用 multipass snapshot <name> 实际创建快照

	// 记录日志
	vmLog := model.VMLog{
		VMID:      vm.ID,
		Operation: "create_snapshot",
		Result:    "success",
		Message:   fmt.Sprintf("创建快照：%s", req.Name),
	}
	db.Create(&vmLog)

	response.SuccessWithMessage(c, "快照创建成功", snapshot)
}

// ListSnapshots 获取快照列表
func (h *VMHandler) ListSnapshots(c *gin.Context) {
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

	// 获取快照列表
	var snapshots []model.VMSnapshot
	if err := db.Where("vm_id = ?", vmID).Order("created_at DESC").Find(&snapshots).Error; err != nil {
		response.InternalError(c, "查询失败")
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

	vmIDStr := c.Param("id")
	vmID, err := strconv.ParseUint(vmIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的虚拟机 ID")
		return
	}

	var req struct {
		SnapshotID uint `json:"snapshot_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	db := repository.GetDB()
	
	// 验证虚拟机和快照归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	var snapshot model.VMSnapshot
	if err := db.Where("id = ? AND vm_id = ?", req.SnapshotID, vmID).First(&snapshot).Error; err != nil {
		response.NotFound(c, "快照不存在")
		return
	}

	// TODO: 这里应该调用 multipass restore <name> <snapshot> 实际恢复快照

	// 记录日志
	vmLog := model.VMLog{
		VMID:      vm.ID,
		Operation: "restore_snapshot",
		Result:    "success",
		Message:   fmt.Sprintf("恢复快照：%s", snapshot.Name),
	}
	db.Create(&vmLog)

	response.SuccessWithMessage(c, "快照恢复成功", nil)
}

// DeleteSnapshot 删除快照
func (h *VMHandler) DeleteSnapshot(c *gin.Context) {
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

	snapshotIDStr := c.Param("snapshot_id")
	snapshotID, err := strconv.ParseUint(snapshotIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的快照 ID")
		return
	}

	db := repository.GetDB()
	
	// 验证虚拟机和快照归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", uint(vmID), userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	var snapshot model.VMSnapshot
	if err := db.Where("id = ? AND vm_id = ?", uint(snapshotID), vmID).First(&snapshot).Error; err != nil {
		response.NotFound(c, "快照不存在")
		return
	}

	// TODO: 这里应该调用 multipass delete-snapshot <name> <snapshot> 实际删除快照

	// 删除数据库记录
	db.Delete(&snapshot)

	// 记录日志
	vmLog := model.VMLog{
		VMID:      vm.ID,
		Operation: "delete_snapshot",
		Result:    "success",
		Message:   fmt.Sprintf("删除快照：%s", snapshot.Name),
	}
	db.Create(&vmLog)

	response.SuccessWithMessage(c, "快照删除成功", nil)
}

// GetVMLogs 获取虚拟机操作日志
func (h *VMHandler) GetVMLogs(c *gin.Context) {
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

	// 获取日志列表
	var logs []model.VMLog
	if err := db.Where("vm_id = ?", vmID).Order("created_at DESC").Limit(100).Find(&logs).Error; err != nil {
		response.InternalError(c, "查询失败")
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

	// 获取最近的监控数据
	var metrics []model.VMMetric
	if err := db.Where("vm_id = ?", vmID).Order("timestamp DESC").Limit(100).Find(&metrics).Error; err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, metrics)
}
