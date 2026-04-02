package repository

import (
	"fmt"
	"log/slog"

	"github.com/cinagroup/cinaseek/backend/internal/model"
	"gorm.io/gorm"
)

// SnapshotRepo provides CRUD operations for the vm_snapshots table.
type SnapshotRepo struct {
	db *gorm.DB
}

// NewSnapshotRepo creates a new SnapshotRepo.
func NewSnapshotRepo() *SnapshotRepo {
	return &SnapshotRepo{db: GetDB()}
}

// Create inserts a new snapshot record.
func (r *SnapshotRepo) Create(snapshot *model.VMSnapshot) error {
	if err := r.db.Create(snapshot).Error; err != nil {
		slog.Error("snapshot create failed", "error", err, "vm_id", snapshot.VMID, "name", snapshot.Name)
		return fmt.Errorf("create snapshot: %w", err)
	}
	slog.Info("snapshot created", "snapshot_id", snapshot.ID, "vm_id", snapshot.VMID, "name", snapshot.Name)
	return nil
}

// FindByID retrieves a snapshot by primary key, scoped to a VM.
func (r *SnapshotRepo) FindByID(id, vmID uint) (*model.VMSnapshot, error) {
	var snapshot model.VMSnapshot
	if err := r.db.Where("id = ? AND vm_id = ?", id, vmID).First(&snapshot).Error; err != nil {
		return nil, fmt.Errorf("find snapshot %d for vm %d: %w", id, vmID, err)
	}
	return &snapshot, nil
}

// ListByVM returns all snapshots for a given VM.
func (r *SnapshotRepo) ListByVM(vmID uint) ([]model.VMSnapshot, error) {
	var snapshots []model.VMSnapshot
	if err := r.db.Where("vm_id = ?", vmID).Order("created_at DESC").Find(&snapshots).Error; err != nil {
		return nil, fmt.Errorf("list snapshots for vm %d: %w", vmID, err)
	}
	return snapshots, nil
}

// Delete removes a snapshot record.
func (r *SnapshotRepo) Delete(id, vmID uint) error {
	result := r.db.Where("id = ? AND vm_id = ?", id, vmID).Delete(&model.VMSnapshot{})
	if result.Error != nil {
		return fmt.Errorf("delete snapshot: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("snapshot %d not found for vm %d", id, vmID)
	}
	slog.Info("snapshot deleted", "snapshot_id", id, "vm_id", vmID)
	return nil
}

// UpdateSize updates the size field of a snapshot.
func (r *SnapshotRepo) UpdateSize(id uint, size int64) error {
	result := r.db.Model(&model.VMSnapshot{}).Where("id = ?", id).Update("size", size)
	if result.Error != nil {
		return fmt.Errorf("update snapshot size: %w", result.Error)
	}
	return nil
}
