package repository

import (
	"fmt"
	"log/slog"

	"github.com/cinagroup/cinaseek/backend/internal/model"
	"gorm.io/gorm"
)

// OpenClawRepo provides CRUD operations for the openclaw_configs table.
type OpenClawRepo struct {
	db *gorm.DB
}

// NewOpenClawRepo creates a new OpenClawRepo.
func NewOpenClawRepo() *OpenClawRepo {
	return &OpenClawRepo{db: GetDB()}
}

// FindByVMID retrieves the OpenClaw config for a given VM.
func (r *OpenClawRepo) FindByVMID(vmID uint) (*model.OpenClawConfig, error) {
	var cfg model.OpenClawConfig
	if err := r.db.Where("vm_id = ?", vmID).First(&cfg).Error; err != nil {
		return nil, fmt.Errorf("find openclaw config for vm %d: %w", vmID, err)
	}
	return &cfg, nil
}

// Create inserts a new OpenClaw config.
func (r *OpenClawRepo) Create(cfg *model.OpenClawConfig) error {
	if err := r.db.Create(cfg).Error; err != nil {
		slog.Error("openclaw config create failed", "error", err, "vm_id", cfg.VMID)
		return fmt.Errorf("create openclaw config: %w", err)
	}
	slog.Info("openclaw config created", "config_id", cfg.ID, "vm_id", cfg.VMID)
	return nil
}

// Save creates or updates an OpenClaw config.
func (r *OpenClawRepo) Save(cfg *model.OpenClawConfig) error {
	if err := r.db.Save(cfg).Error; err != nil {
		slog.Error("openclaw config save failed", "error", err, "vm_id", cfg.VMID)
		return fmt.Errorf("save openclaw config: %w", err)
	}
	return nil
}

// UpdateStatus updates only the status field.
func (r *OpenClawRepo) UpdateStatus(vmID uint, status string) error {
	result := r.db.Model(&model.OpenClawConfig{}).Where("vm_id = ?", vmID).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("update openclaw status: %w", result.Error)
	}
	return nil
}

// CountByUser returns the number of OpenClaw deployments for a user.
func (r *OpenClawRepo) CountByUser(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.OpenClawConfig{}).
		Joins("JOIN vms ON vms.id = openclaw_configs.vm_id").
		Where("vms.user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count openclaw configs: %w", err)
	}
	return count, nil
}
