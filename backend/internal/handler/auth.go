package handler

import (
	"net/http"
	"regexp"
	"strings"

	"multipass-backend/internal/config"
	"multipass-backend/internal/middleware"
	"multipass-backend/internal/model"
	"multipass-backend/internal/repository"
	"multipass-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username       string `json:"username" binding:"required,min=3,max=20"`
		Email          string `json:"email" binding:"required,email"`
		Password       string `json:"password" binding:"required,min=8,max=20"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 验证密码强度
	if !validatePassword(req.Password) {
		response.BadRequest(c, "密码必须包含大小写字母和数字")
		return
	}

	db := repository.GetDB()

	// 检查用户名是否已存在
	var existingUser model.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		response.BadRequest(c, "用户名或邮箱已被注册")
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		response.InternalError(c, "注册失败："+err.Error())
		return
	}

	// 生成 Token
	token, err := middleware.GenerateToken(&h.cfg.JWT, user.ID, user.Username)
	if err != nil {
		response.InternalError(c, "Token 生成失败")
		return
	}

	response.SuccessWithMessage(c, "注册成功", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Remember bool   `json:"remember"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	db := repository.GetDB()

	// 查找用户
	var user model.User
	query := db.Where("username = ? OR email = ?", req.Username, req.Username)
	if err := query.First(&user).Error; err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 更新登录时间
	now := user.CreatedAt // 使用 CreatedAt 避免 import time
	user.LastLoginAt = &now
	db.Save(&user)

	// 记录登录日志
	clientIP := c.ClientIP()
	loginLog := model.LoginLog{
		UserID: user.ID,
		IP:     clientIP,
		Device: c.Request.UserAgent(),
	}
	db.Create(&loginLog)

	// 生成 Token
	token, err := middleware.GenerateToken(&h.cfg.JWT, user.ID, user.Username)
	if err != nil {
		response.InternalError(c, "Token 生成失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		},
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 在实际应用中，这里可以将 Token 加入黑名单
	// 由于 JWT 是无状态的，我们只需要让客户端删除 Token 即可

	response.SuccessWithMessage(c, "登出成功", nil)
}

// ResetPassword 重置密码
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		NewPassword string `json:"new_password" binding:"required,min=8,max=20"`
		Code        string `json:"code" binding:"required"` // 验证码（实际应用中需要验证）
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 验证密码强度
	if !validatePassword(req.NewPassword) {
		response.BadRequest(c, "密码必须包含大小写字母和数字")
		return
	}

	db := repository.GetDB()

	// 查找用户
	var user model.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		response.BadRequest(c, "邮箱未注册")
		return
	}

	// TODO: 这里应该验证验证码（通过邮件或短信发送的）
	// 为简化实现，暂时跳过验证码验证

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}

	// 更新密码
	user.Password = string(hashedPassword)
	if err := db.Save(&user).Error; err != nil {
		response.InternalError(c, "密码更新失败")
		return
	}

	response.SuccessWithMessage(c, "密码重置成功", nil)
}

// validatePassword 验证密码强度
func validatePassword(password string) bool {
	// 至少 8 个字符
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	hasLower := regexp.MustCompile(`[a-z]`).MatchString
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString

	return hasUpper(password) && hasLower(password) && hasNumber(password)
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	db := repository.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, gin.H{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"nickname":      user.Nickname,
		"phone":         user.Phone,
		"avatar":        user.Avatar,
		"created_at":    user.CreatedAt,
		"last_login_at": user.LastLoginAt,
		"two_factor_enabled": user.TwoFactorEnabled,
	})
}

// UpdateUserInfo 更新用户信息
func (h *AuthHandler) UpdateUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		Email    string `json:"email" binding:"omitempty,email"`
		Nickname string `json:"nickname" binding:"max=50"`
		Phone    string `json:"phone" binding:"max=20"`
		Avatar   string `json:"avatar" binding:"max=255"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	db := repository.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 如果修改邮箱，检查是否已被使用
	if req.Email != "" && req.Email != user.Email {
		var existingUser model.User
		if err := db.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			response.BadRequest(c, "邮箱已被使用")
			return
		}
		user.Email = req.Email
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := db.Save(&user).Error; err != nil {
		response.InternalError(c, "更新失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"nickname": user.Nickname,
		"phone":    user.Phone,
		"avatar":   user.Avatar,
	})
}

// UpdatePassword 修改密码
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8,max=20"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 验证密码强度
	if !validatePassword(req.NewPassword) {
		response.BadRequest(c, "密码必须包含大小写字母和数字")
		return
	}

	db := repository.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 验证当前密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		response.BadRequest(c, "当前密码错误")
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}

	user.Password = string(hashedPassword)
	if err := db.Save(&user).Error; err != nil {
		response.InternalError(c, "密码更新失败")
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// GetLoginLogs 获取登录日志
func (h *AuthHandler) GetLoginLogs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	db := repository.GetDB()
	var logs []model.LoginLog
	if err := db.Where("user_id = ?", userID).Order("login_time DESC").Limit(50).Find(&logs).Error; err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, logs)
}

// GetSessions 获取活跃会话
func (h *AuthHandler) GetSessions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	db := repository.GetDB()
	var sessions []model.Session
	if err := db.Where("user_id = ? AND expired_at > ?", userID, db.NowFunc()).Find(&sessions).Error; err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, sessions)
}

// RevokeSession 撤销会话
func (h *AuthHandler) RevokeSession(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		SessionID uint `json:"session_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	db := repository.GetDB()
	
	// 更新会话过期时间为当前时间
	result := db.Model(&model.Session{}).
		Where("id = ? AND user_id = ?", req.SessionID, userID).
		Update("expired_at", db.NowFunc())

	if result.RowsAffected == 0 {
		response.NotFound(c, "会话不存在")
		return
	}

	response.SuccessWithMessage(c, "会话已撤销", nil)
}
