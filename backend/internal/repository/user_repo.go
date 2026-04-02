package repository

import (
	"fmt"
	"log/slog"

	"github.com/cinagroup/cinaseek/backend/internal/model"
	"gorm.io/gorm"
)

// UserRepo provides CRUD operations for the users table.
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo() *UserRepo {
	return &UserRepo{db: GetDB()}
}

// Create inserts a new user record.
func (r *UserRepo) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		slog.Error("user create failed", "error", err, "username", user.Username)
		return fmt.Errorf("create user: %w", err)
	}
	slog.Info("user created", "user_id", user.ID, "username", user.Username)
	return nil
}

// FindByID retrieves a user by primary key.
func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("find user by id %d: %w", id, err)
	}
	return &user, nil
}

// FindByUsername retrieves a user by username.
func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by username %q: %w", username, err)
	}
	return &user, nil
}

// FindByEmail retrieves a user by email.
func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by email %q: %w", email, err)
	}
	return &user, nil
}

// FindByUsernameOrEmail retrieves a user by username or email.
func (r *UserRepo) FindByUsernameOrEmail(value string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ? OR email = ?", value, value).First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by username or email %q: %w", value, err)
	}
	return &user, nil
}

// FindByCinaTokenID retrieves a user by CinaToken ID.
func (r *UserRepo) FindByCinaTokenID(cinatokenID uint) (*model.User, error) {
	var user model.User
	if err := r.db.Where("cinatoken_id = ?", cinatokenID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by cinatoken_id %d: %w", cinatokenID, err)
	}
	return &user, nil
}

// Update updates an existing user record.
func (r *UserRepo) Update(user *model.User) error {
	if err := r.db.Save(user).Error; err != nil {
		slog.Error("user update failed", "error", err, "user_id", user.ID)
		return fmt.Errorf("update user: %w", err)
	}
	slog.Debug("user updated", "user_id", user.ID)
	return nil
}

// UpdatePassword updates only the password field.
func (r *UserRepo) UpdatePassword(userID uint, hashedPassword string) error {
	result := r.db.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword)
	if result.Error != nil {
		return fmt.Errorf("update user password: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user %d not found for password update", userID)
	}
	return nil
}

// UpdateLastLogin sets the last_login_at timestamp.
func (r *UserRepo) UpdateLastLogin(userID uint) error {
	result := r.db.Model(&model.User{}).Where("id = ?", userID).Update("last_login_at", gorm.Expr("NOW()"))
	if result.Error != nil {
		return fmt.Errorf("update last login: %w", result.Error)
	}
	return nil
}

// ExistsByUsernameOrEmail checks if a user with the given username or email exists.
func (r *UserRepo) ExistsByUsernameOrEmail(username, email string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("username = ? OR email = ?", username, email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check user existence: %w", err)
	}
	return count > 0, nil
}
