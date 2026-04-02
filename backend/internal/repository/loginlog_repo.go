package repository

import (
	"fmt"

	"github.com/cinagroup/cinaseek/backend/internal/model"
)

// LoginLogRepo provides operations for the login_logs table.
type LoginLogRepo struct {
}

// NewLoginLogRepo creates a new LoginLogRepo.
func NewLoginLogRepo() *LoginLogRepo {
	return &LoginLogRepo{}
}

// Create inserts a new login log record.
func (r *LoginLogRepo) Create(log *model.LoginLog) error {
	db := GetDB()
	if err := db.Create(log).Error; err != nil {
		return fmt.Errorf("create login log: %w", err)
	}
	return nil
}

// ListByUser returns login logs for a user, limited to n entries.
func (r *LoginLogRepo) ListByUser(userID uint, limit int) ([]model.LoginLog, error) {
	db := GetDB()
	var logs []model.LoginLog
	q := db.Where("user_id = ?", userID).Order("login_time DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("list login logs: %w", err)
	}
	return logs, nil
}
