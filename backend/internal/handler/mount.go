package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"cinaroom-backend/internal/config"
	"cinaroom-backend/internal/model"
	"cinaroom-backend/internal/repository"
	"cinaroom-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type MountHandler struct {
	cfg *config.Config
}

func NewMountHandler(cfg *config.Config) *MountHandler {
	return &MountHandler{cfg: cfg}
}

// ListMounts 获取挂载列表
func (h *MountHandler) ListMounts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	db := repository.GetDB()

	// 获取查询参数
	vmIDStr := c.Query("vm_id")
	
	query := db.Model(&model.Mount{}).Where("user_id = ?", userID)
	
	if vmIDStr != "" {
		vmID, err := strconv.ParseUint(vmIDStr, 10, 32)
		if err == nil {
			query = query.Where("vm_id = ?", uint(vmID))
		}
	}

	var mounts []model.Mount
	if err := query.Order("created_at DESC").Find(&mounts).Error; err != nil {
		response.InternalError(c, "查询失败")
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

	var req struct {
		VMID       uint   `json:"vm_id" binding:"required"`
		Name       string `json:"name" binding:"required,max=100"`
		HostPath   string `json:"host_path" binding:"required,max=500"`
		VMPath     string `json:"vm_path" binding:"required,max=500"`
		Permission string `json:"permission"` // ro, rw
		AutoMount  bool   `json:"auto_mount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 设置默认值
	if req.Permission == "" {
		req.Permission = "rw"
	}

	db := repository.GetDB()

	// 验证虚拟机归属
	var vm model.VM
	if err := db.Where("id = ? AND user_id = ?", req.VMID, userID).First(&vm).Error; err != nil {
		response.NotFound(c, "虚拟机不存在")
		return
	}

	// 创建挂载
	mount := model.Mount{
		UserID:     userID.(uint),
		VMID:       req.VMID,
		Name:       req.Name,
		HostPath:   req.HostPath,
		VMPath:     req.VMPath,
		Status:     "unmounted",
		Permission: req.Permission,
		AutoMount:  req.AutoMount,
	}

	if err := db.Create(&mount).Error; err != nil {
		response.InternalError(c, "创建失败："+err.Error())
		return
	}

	// TODO: 如果 AutoMount 为 true，调用 multipass mount <host_path> <vm_name>:<vm_path>

	response.SuccessWithMessage(c, "挂载添加成功", mount)
}

// OperateMount 操作挂载（挂载/卸载/编辑/删除）
func (h *MountHandler) OperateMount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	mountIDStr := c.Param("id")
	mountID, err := strconv.ParseUint(mountIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的挂载 ID")
		return
	}

	var req struct {
		Operation string `json:"operation" binding:"required"` // mount, unmount, edit, delete
		Name      string `json:"name"`
		VMPath    string `json:"vm_path"`
		Permission string `json:"permission"`
		AutoMount bool   `json:"auto_mount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	db := repository.GetDB()
	var mount model.Mount
	if err := db.Where("id = ? AND user_id = ?", uint(mountID), userID).First(&mount).Error; err != nil {
		response.NotFound(c, "挂载不存在")
		return
	}

	switch req.Operation {
	case "mount":
		// TODO: 调用 multipass mount <host_path> <vm_name>:<vm_path>
		mount.Status = "mounted"
		db.Save(&mount)
		response.SuccessWithMessage(c, "挂载成功", nil)

	case "unmount":
		// TODO: 调用 multipass umount <vm_name>:<vm_path>
		mount.Status = "unmounted"
		db.Save(&mount)
		response.SuccessWithMessage(c, "卸载成功", nil)

	case "edit":
		if req.Name != "" {
			mount.Name = req.Name
		}
		if req.VMPath != "" {
			mount.VMPath = req.VMPath
		}
		if req.Permission != "" {
			mount.Permission = req.Permission
		}
		mount.AutoMount = req.AutoMount
		db.Save(&mount)
		response.SuccessWithMessage(c, "编辑成功", mount)

	case "delete":
		// TODO: 如果已挂载，先调用 multipass umount
		db.Delete(&mount)
		response.SuccessWithMessage(c, "删除成功", nil)

	default:
		response.BadRequest(c, "不支持的操作")
	}
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
		response.NotFound(c, "OpenClaw 配置不存在")
		return
	}

	// 获取相关的挂载配置
	var mounts []model.Mount
	db.Where("vm_id = ? AND name LIKE ?", vmID, "%openclaw%").Find(&mounts)

	response.Success(c, gin.H{
		"openclaw_config": ocConfig,
		"mounts":          mounts,
	})
}

// ConfigureOpenClawMount 配置 OpenClaw 挂载
func (h *MountHandler) ConfigureOpenClawMount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		VMID              uint   `json:"vm_id" binding:"required"`
		WorkspacePath     string `json:"workspace_path" binding:"required"`
		SkillsPath        string `json:"skills_path" binding:"required"`
		SyncOpenClawJSON  bool   `json:"sync_openclaw_json"`
		SyncToolConfigs   bool   `json:"sync_tool_configs"`
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

	// 创建或更新 OpenClaw 配置
	var ocConfig model.OpenClawConfig
	result := db.Where("vm_id = ?", req.VMID).First(&ocConfig)

	ocConfig.VMID = req.VMID
	ocConfig.WorkspacePath = req.WorkspacePath
	ocConfig.SkillsPath = req.SkillsPath
	ocConfig.SyncOpenClawJSON = req.SyncOpenClawJSON
	ocConfig.SyncToolConfigs = req.SyncToolConfigs

	if result.Error != nil {
		db.Create(&ocConfig)
	} else {
		db.Save(&ocConfig)
	}

	// 创建工作空间挂载
	workspaceMount := model.Mount{
		UserID:     userID.(uint),
		VMID:       req.VMID,
		Name:       "openclaw-workspace",
		HostPath:   req.WorkspacePath,
		VMPath:     "/root/.openclaw/workspace",
		Status:     "unmounted",
		Permission: "rw",
		AutoMount:  true,
	}
	db.FirstOrCreate(&workspaceMount, model.Mount{UserID: userID.(uint), VMID: req.VMID, Name: "openclaw-workspace"})

	// 创建技能目录挂载
	skillsMount := model.Mount{
		UserID:     userID.(uint),
		VMID:       req.VMID,
		Name:       "openclaw-skills",
		HostPath:   req.SkillsPath,
		VMPath:     "/root/.openclaw/workspace/skills",
		Status:     "unmounted",
		Permission: "rw",
		AutoMount:  true,
	}
	db.FirstOrCreate(&skillsMount, model.Mount{UserID: userID.(uint), VMID: req.VMID, Name: "openclaw-skills"})

	// TODO: 自动挂载
	// multipass mount <workspace_path> <vm_name>:/root/.openclaw/workspace
	// multipass mount <skills_path> <vm_name>:/root/.openclaw/workspace/skills

	response.SuccessWithMessage(c, "OpenClaw 挂载配置成功", ocConfig)
}
