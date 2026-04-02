package repository

import (
	"fmt"
	"log/slog"

	"github.com/cinagroup/cinaseek/backend/internal/model"
	"gorm.io/gorm"
)

// RemoteRepo provides CRUD operations for remote access tables.
type RemoteRepo struct {
	db *gorm.DB
}

// NewRemoteRepo creates a new RemoteRepo.
func NewRemoteRepo() *RemoteRepo {
	return &RemoteRepo{db: GetDB()}
}

// FindRemoteAccess retrieves the remote access config for a VM.
func (r *RemoteRepo) FindRemoteAccess(vmID uint) (*model.RemoteAccess, error) {
	var ra model.RemoteAccess
	if err := r.db.Where("vm_id = ?", vmID).First(&ra).Error; err != nil {
		return nil, fmt.Errorf("find remote access for vm %d: %w", vmID, err)
	}
	return &ra, nil
}

// SaveRemoteAccess creates or updates a remote access config.
func (r *RemoteRepo) SaveRemoteAccess(ra *model.RemoteAccess) error {
	if err := r.db.Save(ra).Error; err != nil {
		slog.Error("remote access save failed", "error", err, "vm_id", ra.VMID)
		return fmt.Errorf("save remote access: %w", err)
	}
	return nil
}

// CreateRemoteAccess creates a new remote access config.
func (r *RemoteRepo) CreateRemoteAccess(ra *model.RemoteAccess) error {
	if err := r.db.Create(ra).Error; err != nil {
		return fmt.Errorf("create remote access: %w", err)
	}
	return nil
}

// ListIPWhitelist returns all IP whitelist entries for a VM.
func (r *RemoteRepo) ListIPWhitelist(vmID uint) ([]model.IPWhitelist, error) {
	var list []model.IPWhitelist
	if err := r.db.Where("vm_id = ?", vmID).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list ip whitelist for vm %d: %w", vmID, err)
	}
	return list, nil
}

// CreateIPWhitelist inserts a new IP whitelist entry.
func (r *RemoteRepo) CreateIPWhitelist(entry *model.IPWhitelist) error {
	if err := r.db.Create(entry).Error; err != nil {
		return fmt.Errorf("create ip whitelist: %w", err)
	}
	return nil
}

// DeleteIPWhitelist removes an IP whitelist entry scoped to a VM.
func (r *RemoteRepo) DeleteIPWhitelist(id, vmID uint) (int64, error) {
	result := r.db.Where("id = ? AND vm_id = ?", id, vmID).Delete(&model.IPWhitelist{})
	if result.Error != nil {
		return 0, fmt.Errorf("delete ip whitelist: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// ListRemoteLogs returns paginated remote access logs for a VM.
func (r *RemoteRepo) ListRemoteLogs(vmID uint, ipFilter string, page, pageSize int) ([]model.RemoteLog, int64, error) {
	query := r.db.Model(&model.RemoteLog{}).Where("vm_id = ?", vmID)

	if ipFilter != "" {
		query = query.Where("access_ip = ?", ipFilter)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count remote logs: %w", err)
	}

	var logs []model.RemoteLog
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("access_time DESC").Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("list remote logs: %w", err)
	}

	return logs, total, nil
}
