package repository

import (
	"fmt"

	"github.com/cinagroup/cinaseek/backend/internal/model"
	"gorm.io/gorm"
)

// SessionRepo provides operations for the sessions table.
type SessionRepo struct {
	db *gorm.DB
}

// NewSessionRepo creates a new SessionRepo.
func NewSessionRepo() *SessionRepo {
	return &SessionRepo{db: GetDB()}
}

// Create inserts a new session record.
func (r *SessionRepo) Create(session *model.Session) error {
	if err := r.db.Create(session).Error; err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

// FindActiveByUser returns all active (non-expired) sessions for a user.
func (r *SessionRepo) FindActiveByUser(userID uint) ([]model.Session, error) {
	var sessions []model.Session
	if err := r.db.Where("user_id = ? AND expired_at > ?", userID, gorm.Expr("NOW()")).Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("find active sessions: %w", err)
	}
	return sessions, nil
}

// RevokeByID expires a session, scoped to the user.
func (r *SessionRepo) RevokeByID(sessionID, userID uint) (int64, error) {
	result := r.db.Model(&model.Session{}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		Update("expired_at", gorm.Expr("NOW()"))
	if result.Error != nil {
		return 0, fmt.Errorf("revoke session: %w", result.Error)
	}
	return result.RowsAffected, nil
}
