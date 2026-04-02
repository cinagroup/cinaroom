package repository

import (
	"fmt"
	"log/slog"

	"github.com/cinagroup/cinaseek/backend/internal/model"
	"gorm.io/gorm"
)

// MountRepo provides CRUD operations for the mounts table.
type MountRepo struct {
	db *gorm.DB
}

// NewMountRepo creates a new MountRepo.
func NewMountRepo() *MountRepo {
	return &MountRepo{db: GetDB()}
}

// Create inserts a new mount record.
func (r *MountRepo) Create(mount *model.Mount) error {
	if err := r.db.Create(mount).Error; err != nil {
		slog.Error("mount create failed", "error", err, "name", mount.Name, "vm_id", mount.VMID)
		return fmt.Errorf("create mount: %w", err)
	}
	slog.Info("mount created", "mount_id", mount.ID, "name", mount.Name, "vm_id", mount.VMID)
	return nil
}

// FindByID retrieves a mount by primary key, scoped to the user.
func (r *MountRepo) FindByID(id, userID uint) (*model.Mount, error) {
	var mount model.Mount
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&mount).Error; err != nil {
		return nil, fmt.Errorf("find mount %d for user %d: %w", id, userID, err)
	}
	return &mount, nil
}

// ListByUser returns all mounts for a user, optionally filtered by vm_id.
func (r *MountRepo) ListByUser(userID uint, vmID uint) ([]model.Mount, error) {
	query := r.db.Where("user_id = ?", userID)
	if vmID > 0 {
		query = query.Where("vm_id = ?", vmID)
	}

	var mounts []model.Mount
	if err := query.Order("created_at DESC").Find(&mounts).Error; err != nil {
		return nil, fmt.Errorf("list mounts for user %d: %w", userID, err)
	}
	return mounts, nil
}

// ListByVM returns all mounts for a specific VM.
func (r *MountRepo) ListByVM(vmID uint) ([]model.Mount, error) {
	var mounts []model.Mount
	if err := r.db.Where("vm_id = ?", vmID).Order("created_at DESC").Find(&mounts).Error; err != nil {
		return nil, fmt.Errorf("list mounts for vm %d: %w", vmID, err)
	}
	return mounts, nil
}

// UpdateStatus updates only the status field.
func (r *MountRepo) UpdateStatus(id uint, status string) error {
	result := r.db.Model(&model.Mount{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("update mount status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("mount %d not found for status update", id)
	}
	return nil
}

// Save performs a full update of all fields on the mount.
func (r *MountRepo) Save(mount *model.Mount) error {
	if err := r.db.Save(mount).Error; err != nil {
		slog.Error("mount save failed", "error", err, "mount_id", mount.ID)
		return fmt.Errorf("save mount: %w", err)
	}
	return nil
}

// Delete removes a mount record.
func (r *MountRepo) Delete(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Mount{})
	if result.Error != nil {
		return fmt.Errorf("delete mount: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("mount %d not found for user %d", id, userID)
	}
	slog.Info("mount deleted", "mount_id", id, "user_id", userID)
	return nil
}

// FirstOrCreate finds or creates a mount matching the given conditions.
func (r *MountRepo) FirstOrCreate(mount *model.Mount, conditions model.Mount) error {
	if err := r.db.Where(conditions).FirstOrCreate(mount).Error; err != nil {
		return fmt.Errorf("first or create mount: %w", err)
	}
	return nil
}

// CountByUser returns the number of mounts for a user.
func (r *MountRepo) CountByUser(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Mount{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count mounts: %w", err)
	}
	return count, nil
}
