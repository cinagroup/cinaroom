package handler

import (
	"log/slog"
	"regexp"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/middleware"
	"github.com/cinagroup/cinaseek/backend/internal/model"
	"github.com/cinagroup/cinaseek/backend/internal/repository"
	"github.com/cinagroup/cinaseek/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication-related requests.
type AuthHandler struct {
	cfg *config.Config
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

// Register godoc
// @Summary      用户注册
// @Description  注册新用户账号
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "注册信息"  example({"username":"admin","email":"admin@example.com","password":"Admin123","confirm_password":"Admin123"})
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username        string `json:"username" binding:"required,min=3,max=20"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=8,max=20"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	if !validatePassword(req.Password) {
		response.BadRequest(c, "密码必须包含大小写字母和数字")
		return
	}

	db := repository.GetDB()

	var existingUser model.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		response.BadRequest(c, "用户名或邮箱已被注册")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("password hashing failed", "error", err)
		response.InternalError(c, "密码加密失败")
		return
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Active:   true,
		Role:     0, // 普通用户
	}

	if err := db.Create(&user).Error; err != nil {
		slog.Error("user creation failed", "error", err, "username", req.Username)
		response.InternalError(c, "注册失败")
		return
	}

	token, err := middleware.GenerateToken(&h.cfg.JWT, user.ID, user.Username, user.Role)
	if err != nil {
		slog.Error("token generation failed", "error", err)
		response.InternalError(c, "Token 生成失败")
		return
	}

	slog.Info("user registered", "user_id", user.ID, "username", user.Username)

	response.SuccessWithMessage(c, "注册成功", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login godoc
// @Summary      用户登录
// @Description  使用用户名/邮箱和密码登录
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "登录信息"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /auth/login [post]
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

	var user model.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	now := time.Now()
	user.LastLoginAt = &now
	db.Save(&user)

	loginLog := model.LoginLog{
		UserID:    user.ID,
		LoginTime: now,
		IP:        c.ClientIP(),
		Device:    c.Request.UserAgent(),
	}
	db.Create(&loginLog)

	token, err := middleware.GenerateToken(&h.cfg.JWT, user.ID, user.Username, user.Role)
	if err != nil {
		slog.Error("token generation failed", "error", err)
		response.InternalError(c, "Token 生成失败")
		return
	}

	slog.Info("user logged in", "user_id", user.ID, "ip", c.ClientIP())

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	})
}

// Logout godoc
// @Summary      用户登出
// @Description  登出当前用户（客户端删除 token）
// @Tags         auth
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	response.SuccessWithMessage(c, "登出成功", nil)
}

// ResetPassword godoc
// @Summary      重置密码
// @Description  通过邮箱和验证码重置密码
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "重置密码信息"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Router       /auth/reset-pwd [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		NewPassword string `json:"new_password" binding:"required,min=8,max=20"`
		Code        string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if !validatePassword(req.NewPassword) {
		response.BadRequest(c, "密码必须包含大小写字母和数字")
		return
	}

	db := repository.GetDB()

	var user model.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		response.BadRequest(c, "邮箱未注册")
		return
	}

	// TODO: verify verification code (Phase 2 – CinaToken email verification)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("password hashing failed", "error", err)
		response.InternalError(c, "密码加密失败")
		return
	}

	user.Password = string(hashedPassword)
	if err := db.Save(&user).Error; err != nil {
		slog.Error("password update failed", "error", err, "user_id", user.ID)
		response.InternalError(c, "密码更新失败")
		return
	}

	slog.Info("password reset", "user_id", user.ID)
	response.SuccessWithMessage(c, "密码重置成功", nil)
}

// validatePassword checks password strength (upper + lower + digit).
func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	hasLower := regexp.MustCompile(`[a-z]`).MatchString
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString
	return hasUpper(password) && hasLower(password) && hasNumber(password)
}

// GetUserInfo godoc
// @Summary      获取用户信息
// @Description  获取当前登录用户的详细信息
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /auth/user-info [get]
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
		"id":                 user.ID,
		"username":           user.Username,
		"email":              user.Email,
		"nickname":           user.Nickname,
		"phone":              user.Phone,
		"avatar":             user.Avatar,
		"created_at":         user.CreatedAt,
		"last_login_at":      user.LastLoginAt,
		"two_factor_enabled": user.TwoFactorEnabled,
		"role":               user.Role,
	})
}

// UpdateUserInfo godoc
// @Summary      更新用户信息
// @Description  更新当前用户的昵称、邮箱、手机号等
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  object  true  "用户信息"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /auth/user-info [put]
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
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	db := repository.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	if req.Email != "" && req.Email != user.Email {
		var existing model.User
		if err := db.Where("email = ? AND id != ?", req.Email, userID).First(&existing).Error; err == nil {
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
		slog.Error("user update failed", "error", err, "user_id", userID)
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

// UpdatePassword godoc
// @Summary      修改密码
// @Description  修改当前用户密码
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  object  true  "密码信息"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /auth/user-pwd [put]
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		response.BadRequest(c, "当前密码错误")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("password hashing failed", "error", err)
		response.InternalError(c, "密码加密失败")
		return
	}

	user.Password = string(hashedPassword)
	if err := db.Save(&user).Error; err != nil {
		slog.Error("password save failed", "error", err, "user_id", userID)
		response.InternalError(c, "密码更新失败")
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// GetLoginLogs godoc
// @Summary      获取登录日志
// @Description  获取当前用户的登录日志
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /auth/login-logs [get]
func (h *AuthHandler) GetLoginLogs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	db := repository.GetDB()
	var logs []model.LoginLog
	if err := db.Where("user_id = ?", userID).Order("login_time DESC").Limit(50).Find(&logs).Error; err != nil {
		slog.Error("login logs query failed", "error", err, "user_id", userID)
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, logs)
}

// GetSessions godoc
// @Summary      获取活跃会话
// @Description  获取当前用户的活跃会话列表
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /auth/sessions [get]
func (h *AuthHandler) GetSessions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	db := repository.GetDB()
	var sessions []model.Session
	if err := db.Where("user_id = ? AND expired_at > ?", userID, db.NowFunc()).Find(&sessions).Error; err != nil {
		slog.Error("sessions query failed", "error", err, "user_id", userID)
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, sessions)
}

// RevokeSession godoc
// @Summary      撤销会话
// @Description  撤销指定的用户会话
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  object  true  "会话信息"
// @Success      200  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /auth/sessions/revoke [post]
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

	result := db.Model(&model.Session{}).
		Where("id = ? AND user_id = ?", req.SessionID, userID).
		Update("expired_at", db.NowFunc())

	if result.RowsAffected == 0 {
		response.NotFound(c, "会话不存在")
		return
	}

	response.SuccessWithMessage(c, "会话已撤销", nil)
}
