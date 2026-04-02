package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/cinaclaw"
	"github.com/cinagroup/cinaseek/backend/internal/model"
	"github.com/cinagroup/cinaseek/backend/internal/repository"
	"gorm.io/gorm"
)

// VMService handles virtual machine business logic.
type VMService struct {
	vmRepo      *repository.VMRepo
	snapshotRepo *repository.SnapshotRepo
	vmLogRepo   *repository.VMLogRepo
	metricRepo  *repository.MetricRepo
	clientMgr   *cinaclaw.ClientManager
}

// NewVMService creates a new VMService.
func NewVMService(clientMgr *cinaclaw.ClientManager) *VMService {
	return &VMService{
		vmRepo:      repository.NewVMRepo(),
		snapshotRepo: repository.NewSnapshotRepo(),
		vmLogRepo:   repository.NewVMLogRepo(),
		metricRepo:  repository.NewMetricRepo(),
		clientMgr:   clientMgr,
	}
}

// ListVMsRequest holds the input for listing VMs.
type ListVMsRequest struct {
	Name     string `form:"name"`
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// CreateVMRequest holds the input for creating a VM.
type CreateVMRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Image       string `json:"image" binding:"required"`
	CPU         int    `json:"cpu" binding:"min=1,max=8"`
	Memory      int    `json:"memory" binding:"min=1,max=16"`
	Disk        int    `json:"disk" binding:"min=10,max=500"`
	NetworkType string `json:"network_type"`
	SSHKey      string `json:"ssh_key"`
	InitScript  string `json:"init_script"`
}

// OperateVMRequest holds the input for VM operations.
type OperateVMRequest struct {
	Operation string `json:"operation" binding:"required"` // start, stop, restart, pause, resume, delete
}

// UpdateVMConfigRequest holds the input for updating VM configuration.
type UpdateVMConfigRequest struct {
	CPU    int `json:"cpu" binding:"omitempty,min=1,max=8"`
	Memory int `json:"memory" binding:"omitempty,min=1,max=16"`
	Disk   int `json:"disk" binding:"omitempty,min=10,max=500"`
}

// CreateSnapshotRequest holds the input for creating a snapshot.
type CreateSnapshotRequest struct {
	Name string `json:"name" binding:"required,max=100"`
}

// RestoreSnapshotRequest holds the input for restoring a snapshot.
type RestoreSnapshotRequest struct {
	SnapshotID uint `json:"snapshot_id" binding:"required"`
}

var (
	ErrVMNotFound       = errors.New("虚拟机不存在")
	ErrVMInvalidConfig  = errors.New("虚拟机配置无效")
	ErrVMOperationFail  = errors.New("虚拟机操作失败")
	ErrSnapshotNotFound = errors.New("快照不存在")
)

// ListVMs returns a paginated list of VMs for a user.
func (s *VMService) ListVMs(userID uint, req *ListVMsRequest) ([]model.VM, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	return s.vmRepo.FindByUser(userID, req.Name, req.Status, req.Page, req.PageSize)
}

// GetVM returns a single VM by ID, scoped to the user.
func (s *VMService) GetVM(vmID, userID uint) (*model.VM, error) {
	vm, err := s.vmRepo.FindByID(vmID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVMNotFound
		}
		return nil, fmt.Errorf("find vm: %w", err)
	}
	return vm, nil
}

// CreateVM creates a new VM and provisions it via CinaClaw.
func (s *VMService) CreateVM(userID uint, req *CreateVMRequest) (*model.VM, error) {
	// Set defaults
	if req.NetworkType == "" {
		req.NetworkType = "nat"
	}
	if req.CPU < 1 {
		req.CPU = 1
	}
	if req.Memory < 1 {
		req.Memory = 1
	}
	if req.Disk < 10 {
		req.Disk = 10
	}

	vm := &model.VM{
		UserID:      userID,
		Name:        req.Name,
		Status:      "creating",
		Image:       req.Image,
		CPU:         req.CPU,
		Memory:      req.Memory,
		Disk:        req.Disk,
		NetworkType: req.NetworkType,
		SSHKey:      req.SSHKey,
		InitScript:  req.InitScript,
	}

	// Persist to database first
	if err := s.vmRepo.Create(vm); err != nil {
		return nil, fmt.Errorf("create vm record: %w", err)
	}

	// Provision via CinaClaw in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		client, err := s.clientMgr.GetClient(strconv.FormatUint(uint64(userID), 10))
		if err != nil {
			slog.Error("cinaclaw client error", "error", err, "user_id", userID)
			_ = s.vmRepo.UpdateStatus(vm.ID, "error")
			s.logOperation(vm.ID, "create", "error", fmt.Sprintf("CinaClaw 连接失败: %v", err))
			return
		}

		info, err := client.CreateVM(ctx, &cinaclaw.CreateVMRequest{
			Name:   req.Name,
			Image:  req.Image,
			CPU:    req.CPU,
			Memory: fmt.Sprintf("%dG", req.Memory),
			Disk:   fmt.Sprintf("%dG", req.Disk),
		})
		if err != nil {
			slog.Error("cinaclaw create vm failed", "error", err, "vm_id", vm.ID)
			_ = s.vmRepo.UpdateStatus(vm.ID, "error")
			s.logOperation(vm.ID, "create", "error", fmt.Sprintf("创建失败: %v", err))
			return
		}

		// Update VM with actual info
		if info != nil {
			if info.IP != "" {
				_ = s.vmRepo.UpdateIP(vm.ID, info.IP)
			}
			_ = s.vmRepo.UpdateStatus(vm.ID, "stopped")
		} else {
			_ = s.vmRepo.UpdateStatus(vm.ID, "stopped")
		}

		s.logOperation(vm.ID, "create", "success", fmt.Sprintf("创建虚拟机: %s", req.Name))
		slog.Info("vm created via cinaclaw", "vm_id", vm.ID, "name", req.Name, "user_id", userID)
	}()

	return vm, nil
}

// OperateVM performs an operation on a VM (start/stop/restart/pause/resume/delete).
func (s *VMService) OperateVM(vmID, userID uint, req *OperateVMRequest) error {
	vm, err := s.vmRepo.FindByID(vmID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVMNotFound
		}
		return fmt.Errorf("find vm: %w", err)
	}

	// Validate operation
	validOps := map[string]bool{
		"start": true, "stop": true, "restart": true,
		"pause": true, "resume": true, "delete": true,
	}
	if !validOps[req.Operation] {
		return fmt.Errorf("不支持的操作: %s", req.Operation)
	}

	// Execute via CinaClaw
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := s.clientMgr.GetClient(strconv.FormatUint(uint64(userID), 10))
	if err != nil {
		slog.Error("cinaclaw client error", "error", err, "user_id", userID)
		return ErrVMOperationFail
	}

	var opErr error
	switch req.Operation {
	case "start":
		opErr = client.StartVM(ctx, vm.Name)
	case "stop":
		opErr = client.StopVM(ctx, vm.Name)
	case "restart":
		opErr = client.RestartVM(ctx, vm.Name)
	case "pause":
		opErr = client.SuspendVM(ctx, vm.Name)
	case "resume":
		opErr = client.StartVM(ctx, vm.Name) // resume = start from suspended
	case "delete":
		opErr = client.DeleteVM(ctx, vm.Name)
	}

	if opErr != nil {
		slog.Error("cinaclaw operation failed", "error", opErr, "vm_id", vm.ID, "op", req.Operation)
		s.logOperation(vm.ID, req.Operation, "error", opErr.Error())
		return ErrVMOperationFail
	}

	// Update status in database
	statusMap := map[string]string{
		"start": "running", "stop": "stopped", "restart": "running",
		"pause": "paused", "resume": "running",
	}

	if req.Operation == "delete" {
		if err := s.vmRepo.Delete(vmID, userID); err != nil {
			slog.Error("vm delete from db failed", "error", err, "vm_id", vmID)
		}
	} else {
		if newStatus, ok := statusMap[req.Operation]; ok {
			_ = s.vmRepo.UpdateStatus(vmID, newStatus)
		}
	}

	message := fmt.Sprintf("虚拟机操作: %s", req.Operation)
	s.logOperation(vm.ID, req.Operation, "success", message)
	slog.Info("vm operation completed", "vm_id", vm.ID, "operation", req.Operation)

	return nil
}

// UpdateVMConfig updates the configuration of a VM.
func (s *VMService) UpdateVMConfig(vmID, userID uint, req *UpdateVMConfigRequest) (*model.VM, error) {
	vm, err := s.vmRepo.FindByID(vmID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVMNotFound
		}
		return nil, fmt.Errorf("find vm: %w", err)
	}

	if err := s.vmRepo.UpdateConfig(vmID, req.CPU, req.Memory, req.Disk); err != nil {
		return nil, fmt.Errorf("update vm config: %w", err)
	}

	// Reload
	vm, _ = s.vmRepo.FindByID(vmID, userID)

	s.logOperation(vmID, "update_config", "success",
		fmt.Sprintf("更新配置: CPU=%d, Memory=%d, Disk=%d", vm.CPU, vm.Memory, vm.Disk))

	return vm, nil
}

// ListSnapshots returns all snapshots for a VM.
func (s *VMService) ListSnapshots(vmID, userID uint) ([]model.VMSnapshot, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}
	return s.snapshotRepo.ListByVM(vmID)
}

// CreateSnapshot creates a new snapshot for a VM.
func (s *VMService) CreateSnapshot(vmID, userID uint, req *CreateSnapshotRequest) (*model.VMSnapshot, error) {
	vm, err := s.vmRepo.FindByID(vmID, userID)
	if err != nil {
		return nil, ErrVMNotFound
	}

	snapshot := &model.VMSnapshot{
		VMID: vm.ID,
		Name: req.Name,
	}

	if err := s.snapshotRepo.Create(snapshot); err != nil {
		return nil, fmt.Errorf("create snapshot: %w", err)
	}

	// Create snapshot via CinaClaw
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		client, err := s.clientMgr.GetClient(strconv.FormatUint(uint64(userID), 10))
		if err != nil {
			slog.Error("cinaclaw client error for snapshot", "error", err, "vm_id", vmID)
			return
		}

		if err := client.SnapshotVM(ctx, vm.Name, req.Name); err != nil {
			slog.Error("cinaclaw snapshot failed", "error", err, "vm_id", vmID)
			s.logOperation(vmID, "create_snapshot", "error", err.Error())
			return
		}

		s.logOperation(vmID, "create_snapshot", "success", fmt.Sprintf("创建快照: %s", req.Name))
	}()

	return snapshot, nil
}

// RestoreSnapshot restores a VM from a snapshot.
func (s *VMService) RestoreSnapshot(vmID, userID uint, req *RestoreSnapshotRequest) error {
	vm, err := s.vmRepo.FindByID(vmID, userID)
	if err != nil {
		return ErrVMNotFound
	}

	snapshot, err := s.snapshotRepo.FindByID(req.SnapshotID, vmID)
	if err != nil {
		return ErrSnapshotNotFound
	}

	// Restore via CinaClaw
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		client, err := s.clientMgr.GetClient(strconv.FormatUint(uint64(userID), 10))
		if err != nil {
			slog.Error("cinaclaw client error for restore", "error", err, "vm_id", vmID)
			return
		}

		if err := client.RestoreVM(ctx, vm.Name, snapshot.Name); err != nil {
			slog.Error("cinaclaw restore failed", "error", err, "vm_id", vmID)
			s.logOperation(vmID, "restore_snapshot", "error", err.Error())
			return
		}

		s.logOperation(vmID, "restore_snapshot", "success", fmt.Sprintf("恢复快照: %s", snapshot.Name))
	}()

	return nil
}

// DeleteSnapshot deletes a snapshot.
func (s *VMService) DeleteSnapshot(vmID, snapshotID, userID uint) error {
	vm, err := s.vmRepo.FindByID(vmID, userID)
	if err != nil {
		return ErrVMNotFound
	}

	snapshot, err := s.snapshotRepo.FindByID(snapshotID, vmID)
	if err != nil {
		return ErrSnapshotNotFound
	}

	if err := s.snapshotRepo.Delete(snapshotID, vmID); err != nil {
		return fmt.Errorf("delete snapshot: %w", err)
	}

	// Delete via CinaClaw
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		client, err := s.clientMgr.GetClient(strconv.FormatUint(uint64(userID), 10))
		if err != nil {
			slog.Error("cinaclaw client error for snapshot delete", "error", err)
			return
		}

		if err := client.DeleteSnapshot(ctx, vm.Name, snapshot.Name); err != nil {
			slog.Error("cinaclaw delete snapshot failed", "error", err, "vm_id", vmID)
		}
	}()

	s.logOperation(vmID, "delete_snapshot", "success", fmt.Sprintf("删除快照: %s", snapshot.Name))
	return nil
}

// GetVMLogs returns operation logs for a VM.
func (s *VMService) GetVMLogs(vmID, userID uint, limit int) ([]model.VMLog, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}
	if limit <= 0 {
		limit = 100
	}
	return s.vmLogRepo.ListByVM(vmID, limit)
}

// GetVMMetrics returns metrics for a VM.
func (s *VMService) GetVMMetrics(vmID, userID uint, limit int) ([]model.VMMetric, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}
	if limit <= 0 {
		limit = 100
	}
	return s.metricRepo.ListByVM(vmID, limit)
}

// SearchVMs searches VMs by keyword.
func (s *VMService) SearchVMs(userID uint, keyword string, limit int) ([]model.VM, error) {
	if keyword == "" {
		return nil, errors.New("搜索关键词不能为空")
	}
	if limit <= 0 {
		limit = 10
	}
	return s.vmRepo.Search(userID, keyword, limit)
}

// BatchOperateVMs performs an operation on multiple VMs.
func (s *VMService) BatchOperateVMs(ids []uint, userID uint, operation string) (successCount, failCount int, err error) {
	vms, err := s.vmRepo.FindByIDsAndUser(ids, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("find vms: %w", err)
	}
	if len(vms) != len(ids) {
		return 0, 0, errors.New("部分虚拟机不存在或无权限")
	}

	for _, vm := range vms {
		switch operation {
		case "start":
			_ = s.vmRepo.UpdateStatus(vm.ID, "running")
		case "stop":
			_ = s.vmRepo.UpdateStatus(vm.ID, "stopped")
		case "restart":
			_ = s.vmRepo.UpdateStatus(vm.ID, "running")
		case "delete":
			_ = s.vmRepo.Delete(vm.ID, userID)
		default:
			failCount++
			continue
		}
		s.logOperation(vm.ID, operation, "success", "批量操作")
		successCount++
	}

	return successCount, failCount, nil
}

// GetVMDashboard returns dashboard statistics for a user.
func (s *VMService) GetVMDashboard(userID uint) (map[string]interface{}, error) {
	vmTotal, err := s.vmRepo.CountByUser(userID, "")
	if err != nil {
		return nil, fmt.Errorf("count vms: %w", err)
	}
	vmRunning, err := s.vmRepo.CountByUser(userID, "running")
	if err != nil {
		return nil, fmt.Errorf("count running vms: %w", err)
	}

	return map[string]interface{}{
		"vm_total":   vmTotal,
		"vm_running": vmRunning,
		"vm_stopped": vmTotal - vmRunning,
	}, nil
}

// logOperation is a helper to record a VM operation log.
func (s *VMService) logOperation(vmID uint, operation, result, message string) {
	if err := s.vmLogRepo.Create(&model.VMLog{
		VMID:      vmID,
		Operation: operation,
		Result:    result,
		Message:   message,
	}); err != nil {
		slog.Error("failed to write vm log", "error", err, "vm_id", vmID)
	}
}
