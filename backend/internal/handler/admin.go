package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/model"
	"github.com/cinagroup/cinaseek/backend/internal/repository"
	"github.com/cinagroup/cinaseek/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin-only requests.
type AdminHandler struct {
	cfg *config.Config
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(cfg *config.Config) *AdminHandler {
	return &AdminHandler{cfg: cfg}
}

// ListUsers godoc
// @Summary      用户列表（管理员）
// @Description  获取所有用户列表，支持分页
// @Tags         admin
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int  false  "页码"     default(1)
// @Param        page_size  query  int  false  "每页数量" default(20)
// @Param        keyword    query  string  false  "搜索关键词（用户名/邮箱）"
// @Success      200  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Router       /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := repository.GetDB()
	query := db.Model(&model.User{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var users []model.User
	offset := (page - 1) * pageSize
	if err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		slog.Error("admin list users failed", "error", err)
		response.InternalError(c, "查询用户列表失败")
		return
	}

	response.SuccessWithPage(c, users, total, page, pageSize)
}

// GetUser godoc
// @Summary      用户详情（管理员）
// @Description  获取指定用户的详细信息
// @Tags         admin
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "用户 ID"
// @Success      200  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "缺少用户 ID")
		return
	}

	db := repository.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 获取用户关联的统计信息
	var vmCount int64
	db.Model(&model.VM{}).Where("user_id = ?", user.ID).Count(&vmCount)

	var loginCount int64
	db.Model(&model.LoginLog{}).Where("user_id = ?", user.ID).Count(&loginCount)

	response.Success(c, gin.H{
		"user":         user,
		"vm_count":     vmCount,
		"login_count":  loginCount,
	})
}

// UpdateUserRole godoc
// @Summary      修改用户角色（root）
// @Description  修改指定用户的角色（仅 root 可操作）
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  int  true  "用户 ID"
// @Param        body  body  object  true  "角色信息"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Router       /admin/users/{id}/role [put]
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "缺少用户 ID")
		return
	}

	var req struct {
		Role int `json:"role" binding:"required,oneof=0 10 100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误：role 必须为 0(用户)、10(管理员) 或 100(root)")
		return
	}

	db := repository.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	user.Role = req.Role
	if err := db.Save(&user).Error; err != nil {
		slog.Error("update user role failed", "error", err, "user_id", id)
		response.InternalError(c, "更新角色失败")
		return
	}

	slog.Info("user role updated", "user_id", user.ID, "new_role", req.Role, "operator", c.GetString("username"))

	response.SuccessWithMessage(c, "角色更新成功", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

// DeleteUser godoc
// @Summary      删除用户（root）
// @Description  删除指定用户及其关联数据（仅 root 可操作）
// @Tags         admin
// @Security     BearerAuth
// @Param        id  path  int  true  "用户 ID"
// @Success      200  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "缺少用户 ID")
		return
	}

	// 不允许删除自己
	currentUserID, _ := c.Get("userID")
	if idUint, err := strconv.ParseUint(id, 10, 64); err == nil && uint(idUint) == currentUserID.(uint) {
		response.Error(c, http.StatusForbidden, "不能删除自己", nil)
		return
	}

	db := repository.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 不允许删除 root 用户（保护机制）
	if user.Role >= 100 {
		response.Error(c, http.StatusForbidden, "不能删除 Root 用户", nil)
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		slog.Error("delete user failed", "error", err, "user_id", id)
		response.InternalError(c, "删除用户失败")
		return
	}

	slog.Info("user deleted", "user_id", user.ID, "username", user.Username, "operator", c.GetString("username"))

	response.SuccessWithMessage(c, "用户已删除", nil)
}

// GetSystemStats godoc
// @Summary      系统统计（管理员）
// @Description  获取系统全局统计数据
// @Tags         admin
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /admin/stats [get]
func (h *AdminHandler) GetSystemStats(c *gin.Context) {
	db := repository.GetDB()

	// 用户统计
	var userTotal int64
	var userActive int64
	db.Model(&model.User{}).Count(&userTotal)
	db.Model(&model.User{}).Where("active = ?", true).Count(&userActive)

	// 虚拟机统计
	var vmTotal int64
	var vmRunning int64
	var vmStopped int64
	db.Model(&model.VM{}).Count(&vmTotal)
	db.Model(&model.VM{}).Where("status = ?", "running").Count(&vmRunning)
	db.Model(&model.VM{}).Where("status = ?", "stopped").Count(&vmStopped)

	// 挂载统计
	var mountTotal int64
	db.Model(&model.Mount{}).Count(&mountTotal)

	// OpenClaw 部署统计
	var openclawTotal int64
	var openclawRunning int64
	db.Model(&model.OpenClawConfig{}).Count(&openclawTotal)
	db.Model(&model.OpenClawConfig{}).Where("status = ?", "running").Count(&openclawRunning)

	// 最近 7 天新用户
	var newUsers7d int64
	db.Model(&model.User{}).Where("created_at >= ?", db.NowFunc().AddDate(0, 0, -7)).Count(&newUsers7d)

	// 最近 7 天操作日志
	var ops7d int64
	db.Model(&model.VMLog{}).Where("created_at >= ?", db.NowFunc().AddDate(0, 0, -7)).Count(&ops7d)

	response.Success(c, gin.H{
		"users": gin.H{
			"total":  userTotal,
			"active": userActive,
			"new_7d": newUsers7d,
		},
		"vms": gin.H{
			"total":   vmTotal,
			"running": vmRunning,
			"stopped": vmStopped,
		},
		"mounts": gin.H{
			"total": mountTotal,
		},
		"openclaw": gin.H{
			"total":   openclawTotal,
			"running": openclawRunning,
		},
		"operations_7d": ops7d,
	})
}
