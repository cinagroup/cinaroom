package repository

import (
	"fmt"
	"log/slog"

	"github.com/cinagroup/cinaseek/backend/internal/model"
	"gorm.io/gorm"
)

// VMLogRepo provides operations for the vm_logs table.
type VMLogRepo struct {
	db *gorm.DB
}

// NewVMLogRepo creates a new VMLogRepo.
func NewVMLogRepo() *VMLogRepo {
	return &VMLogRepo{db: GetDB()}
}

// Create inserts a new VM log record.
func (r *VMLogRepo) Create(log *model.VMLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return fmt.Errorf("create vm log: %w", err)
	}
	slog.Debug("vm log created", "vm_id", log.VMID, "operation", log.Operation, "result", log.Result)
	return nil
}

// ListByVM returns logs for a VM, ordered by most recent first.
func (r *VMLogRepo) ListByVM(vmID uint, limit int) ([]model.VMLog, error) {
	var logs []model.VMLog
	q := r.db.Where("vm_id = ?", vmID).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("list vm logs: %w", err)
	}
	return logs, nil
}

// DailyOperationCounts returns operation counts per day for the last N days.
func (r *VMLogRepo) DailyOperationCounts(userID uint, days int) ([]struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}, error) {
	var results []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	err := r.db.Model(&model.VMLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Joins("JOIN vms ON vms.id = vm_logs.vm_id").
		Where("vms.user_id = ? AND vm_logs.created_at >= NOW() - INTERVAL '? days'", userID, days).
		Group("DATE(created_at)").
		Order("date DESC").
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("daily operation counts: %w", err)
	}
	return results, nil
}
