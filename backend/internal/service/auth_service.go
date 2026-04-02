package service

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/middleware"
	"github.com/cinagroup/cinaseek/backend/internal/model"
	"github.com/cinagroup/cinaseek/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo    *repository.UserRepo
	loginLogRepo *repository.LoginLogRepo
	sessionRepo *repository.SessionRepo
	jwtSecret   string
	jwtExpire   time.Duration
}

// NewAuthService creates a new AuthService.
func NewAuthService(jwtSecret string, jwtExpire time.Duration) *AuthService {
	return &AuthService{
		userRepo:    repository.NewUserRepo(),
		loginLogRepo: repository.NewLoginLogRepo(),
		sessionRepo: repository.NewSessionRepo(),
		jwtSecret:   jwtSecret,
		jwtExpire:   jwtExpire,
	}
}

// RegisterRequest holds the input for user registration.
type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// LoginRequest holds the input for user login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

// ResetPasswordRequest holds the input for password reset.
type ResetPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
	Code        string `json:"code"`
}

// UpdatePasswordRequest holds the input for changing password.
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// UpdateUserRequest holds the input for updating user profile.
type UpdateUserRequest struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

// AuthResult holds the result of a successful authentication.
type AuthResult struct {
	Token    string      `json:"token"`
	UserInfo interface{} `json:"user"`
}

var (
	ErrUserExists         = errors.New("用户名或邮箱已被注册")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrPasswordWeak       = errors.New("密码必须包含大小写字母和数字")
	ErrPasswordMismatch   = errors.New("密码不一致")
	ErrPasswordWrong      = errors.New("当前密码错误")
	ErrEmailUsed          = errors.New("邮箱已被使用")
	ErrEmailNotRegistered = errors.New("邮箱未注册")
	ErrInvalidToken       = errors.New("无效的认证令牌")
)

// Register creates a new user account.
func (s *AuthService) Register(req *RegisterRequest, ip string) (*AuthResult, error) {
	// Validate password
	if req.Password != req.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}
	if !isPasswordStrong(req.Password) {
		return nil, ErrPasswordWeak
	}

	// Check if username or email already exists
	exists, err := s.userRepo.ExistsByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		return nil, fmt.Errorf("check user existence: %w", err)
	}
	if exists {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Active:   true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	slog.Info("user registered", "user_id", user.ID, "username", user.Username, "ip", ip)

	return &AuthResult{
		Token: token,
		UserInfo: map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	}, nil
}

// Login authenticates a user with username/email and password.
func (s *AuthService) Login(req *LoginRequest, ip, userAgent string) (*AuthResult, error) {
	user, err := s.userRepo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login time
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		slog.Warn("failed to update last login time", "error", err, "user_id", user.ID)
	}

	// Record login log
	_ = s.loginLogRepo.Create(&model.LoginLog{
		UserID:    user.ID,
		LoginTime: time.Now(),
		IP:        ip,
		Device:    userAgent,
	})

	// Generate JWT token
	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	slog.Info("user logged in", "user_id", user.ID, "username", user.Username, "ip", ip)

	return &AuthResult{
		Token: token,
		UserInfo: map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		},
	}, nil
}

// ResetPassword resets a user's password using a verification code.
func (s *AuthService) ResetPassword(req *ResetPasswordRequest) error {
	if !isPasswordStrong(req.NewPassword) {
		return ErrPasswordWeak
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEmailNotRegistered
		}
		return fmt.Errorf("find user by email: %w", err)
	}

	// TODO: Verify verification code (Phase 2 – CinaToken email verification)
	// For now, we just check the code is non-empty
	if req.Code == "" {
		return errors.New("验证码不能为空")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(user.ID, string(hashedPassword)); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	slog.Info("password reset", "user_id", user.ID)
	return nil
}

// UpdatePassword changes the password for an authenticated user.
func (s *AuthService) UpdatePassword(userID uint, req *UpdatePasswordRequest) error {
	if !isPasswordStrong(req.NewPassword) {
		return ErrPasswordWeak
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return ErrPasswordWrong
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(user.ID, string(hashedPassword)); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	slog.Info("password changed", "user_id", userID)
	return nil
}

// GetUserInfo returns the profile of the authenticated user.
func (s *AuthService) GetUserInfo(userID uint) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return map[string]interface{}{
		"id":                 user.ID,
		"username":           user.Username,
		"email":              user.Email,
		"nickname":           user.Nickname,
		"phone":              user.Phone,
		"avatar":             user.Avatar,
		"created_at":         user.CreatedAt,
		"last_login_at":      user.LastLoginAt,
		"two_factor_enabled": user.TwoFactorEnabled,
	}, nil
}

// UpdateUser updates the profile of the authenticated user.
func (s *AuthService) UpdateUser(userID uint, req *UpdateUserRequest) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check email uniqueness if changing
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepo.ExistsByUsernameOrEmail("", req.Email)
		if err != nil {
			return nil, fmt.Errorf("check email existence: %w", err)
		}
		if exists {
			return nil, ErrEmailUsed
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

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"nickname": user.Nickname,
		"phone":    user.Phone,
		"avatar":   user.Avatar,
	}, nil
}

// GetLoginLogs returns recent login logs for a user.
func (s *AuthService) GetLoginLogs(userID uint) ([]model.LoginLog, error) {
	return s.loginLogRepo.ListByUser(userID, 50)
}

// GetActiveSessions returns active sessions for a user.
func (s *AuthService) GetActiveSessions(userID uint) ([]model.Session, error) {
	return s.sessionRepo.FindActiveByUser(userID)
}

// RevokeSession expires a session.
func (s *AuthService) RevokeSession(sessionID, userID uint) error {
	affected, err := s.sessionRepo.RevokeByID(sessionID, userID)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	if affected == 0 {
		return errors.New("会话不存在")
	}
	return nil
}

// FindOrCreateByOAuth handles OAuth login (find existing user or create new).
func (s *AuthService) FindOrCreateByOAuth(info *OAuthUserInfo) (*model.User, error) {
	// Try to find by CinaToken ID
	user, err := s.userRepo.FindByCinaTokenID(info.ID)
	if err == nil {
		// Update user info from OAuth
		user.Username = info.Username
		user.Nickname = info.Nickname
		user.Avatar = info.Avatar
		user.Phone = info.Phone
		_ = s.userRepo.Update(user)
		return user, nil
	}

	// Try to find by email
	if info.Email != "" {
		user, err = s.userRepo.FindByEmail(info.Email)
		if err == nil {
			user.CinatokenID = info.ID
			user.Provider = info.Provider
			user.Avatar = info.Avatar
			if info.Nickname != "" {
				user.Nickname = info.Nickname
			}
			_ = s.userRepo.Update(user)
			return user, nil
		}
	}

	// Create new user
	user = &model.User{
		CinatokenID: info.ID,
		Username:    info.Username,
		Email:       info.Email,
		Nickname:    info.Nickname,
		Avatar:      info.Avatar,
		Phone:       info.Phone,
		Provider:    info.Provider,
		Active:      true,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("create OAuth user: %w", err)
	}
	return user, nil
}

// OAuthUserInfo holds user info from an OAuth provider.
type OAuthUserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Provider string `json:"provider"`
}

// GenerateToken generates a JWT token for the user.
func (s *AuthService) generateToken(userID uint, username string) (string, error) {
	claims := middleware.Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cinaseek-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// isPasswordStrong checks password strength (upper + lower + digit).
func isPasswordStrong(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	hasLower := regexp.MustCompile(`[a-z]`).MatchString
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString
	return hasUpper(password) && hasLower(password) && hasNumber(password)
}
