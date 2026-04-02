package repository

import (
	"fmt"
	"log/slog"

	"github.com/cinagroup/cinaseek/backend/internal/model"
	"gorm.io/gorm"
)

// VMRepo provides CRUD operations for the vms table.
type VMRepo struct {
	db *gorm.DB
}

// NewVMRepo creates a new VMRepo.
func NewVMRepo() *VMRepo {
	return &VMRepo{db: GetDB()}
}

// Create inserts a new VM record.
func (r *VMRepo) Create(vm *model.VM) error {
	if err := r.db.Create(vm).Error; err != nil {
		slog.Error("vm create failed", "error", err, "name", vm.Name, "user_id", vm.UserID)
		return fmt.Errorf("create vm: %w", err)
	}
	slog.Info("vm created", "vm_id", vm.ID, "name", vm.Name, "user_id", vm.UserID)
	return nil
}

// FindByID retrieves a VM by primary key, scoped to the given user.
func (r *VMRepo) FindByID(id, userID uint) (*model.VM, error) {
	var vm model.VM
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&vm).Error; err != nil {
		return nil, fmt.Errorf("find vm %d for user %d: %w", id, userID, err)
	}
	return &vm, nil
}

// FindByIDNoScope retrieves a VM by primary key without user scoping.
func (r *VMRepo) FindByIDNoScope(id uint) (*model.VM, error) {
	var vm model.VM
	if err := r.db.First(&vm, id).Error; err != nil {
		return nil, fmt.Errorf("find vm %d: %w", id, err)
	}
	return &vm, nil
}

// FindByUser returns all VMs belonging to a user, with optional filters and pagination.
func (r *VMRepo) FindByUser(userID uint, name, status string, page, pageSize int) ([]model.VM, int64, error) {
	query := r.db.Model(&model.VM{}).Where("user_id = ?", userID)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count vms: %w", err)
	}

	var vms []model.VM
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&vms).Error; err != nil {
		return nil, 0, fmt.Errorf("list vms: %w", err)
	}

	return vms, total, nil
}

// ListByUser returns all VMs for a user (no pagination).
func (r *VMRepo) ListByUser(userID uint) ([]model.VM, error) {
	var vms []model.VM
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&vms).Error; err != nil {
		return nil, fmt.Errorf("list vms for user %d: %w", userID, err)
	}
	return vms, nil
}

// UpdateStatus updates only the status field of a VM.
func (r *VMRepo) UpdateStatus(id uint, status string) error {
	result := r.db.Model(&model.VM{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("update vm status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("vm %d not found for status update", id)
	}
	slog.Info("vm status updated", "vm_id", id, "status", status)
	return nil
}

// UpdateConfig updates CPU, Memory, Disk fields.
func (r *VMRepo) UpdateConfig(id uint, cpu, memory, disk int) error {
	updates := map[string]interface{}{}
	if cpu > 0 {
		updates["cpu"] = cpu
	}
	if memory > 0 {
		updates["memory"] = memory
	}
	if disk > 0 {
		updates["disk"] = disk
	}

	if len(updates) == 0 {
		return nil
	}

	result := r.db.Model(&model.VM{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("update vm config: %w", result.Error)
	}
	slog.Info("vm config updated", "vm_id", id, "updates", updates)
	return nil
}

// UpdateIP updates the IP address of a VM.
func (r *VMRepo) UpdateIP(id uint, ip string) error {
	result := r.db.Model(&model.VM{}).Where("id = ?", id).Update("ip", ip)
	if result.Error != nil {
		return fmt.Errorf("update vm ip: %w", result.Error)
	}
	return nil
}

// Save performs a full update of all fields on the VM.
func (r *VMRepo) Save(vm *model.VM) error {
	if err := r.db.Save(vm).Error; err != nil {
		slog.Error("vm save failed", "error", err, "vm_id", vm.ID)
		return fmt.Errorf("save vm: %w", err)
	}
	return nil
}

// Delete removes a VM record (soft or hard).
func (r *VMRepo) Delete(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.VM{})
	if result.Error != nil {
		return fmt.Errorf("delete vm: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("vm %d not found for user %d", id, userID)
	}
	slog.Info("vm deleted", "vm_id", id, "user_id", userID)
	return nil
}

// BatchUpdateStatus updates status for multiple VMs owned by the user.
func (r *VMRepo) BatchUpdateStatus(ids []uint, userID uint, status string) (int64, error) {
	result := r.db.Model(&model.VM{}).
		Where("id IN ? AND user_id = ?", ids, userID).
		Update("status", status)
	if result.Error != nil {
		return 0, fmt.Errorf("batch update vm status: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// BatchDelete deletes multiple VMs owned by the user.
func (r *VMRepo) BatchDelete(ids []uint, userID uint) (int64, error) {
	result := r.db.Where("id IN ? AND user_id = ?", ids, userID).Delete(&model.VM{})
	if result.Error != nil {
		return 0, fmt.Errorf("batch delete vms: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// FindByIDsAndUser returns VMs matching the given IDs and userID.
func (r *VMRepo) FindByIDsAndUser(ids []uint, userID uint) ([]model.VM, error) {
	var vms []model.VM
	if err := r.db.Where("id IN ? AND user_id = ?", ids, userID).Find(&vms).Error; err != nil {
		return nil, fmt.Errorf("find vms by ids: %w", err)
	}
	return vms, nil
}

// CountByUser returns the number of VMs for a user, optionally filtered by status.
func (r *VMRepo) CountByUser(userID uint, status string) (int64, error) {
	query := r.db.Model(&model.VM{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count vms: %w", err)
	}
	return count, nil
}

// Search returns VMs matching a keyword in name or IP.
func (r *VMRepo) Search(userID uint, keyword string, limit int) ([]model.VM, error) {
	var vms []model.VM
	pattern := "%" + keyword + "%"
	if err := r.db.Where("user_id = ? AND (name LIKE ? OR ip LIKE ?)", userID, pattern, pattern).
		Limit(limit).Find(&vms).Error; err != nil {
		return nil, fmt.Errorf("search vms: %w", err)
	}
	return vms, nil
}
