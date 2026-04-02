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

// MountService handles directory mount management.
type MountService struct {
	mountRepo    *repository.MountRepo
	vmRepo       *repository.VMRepo
	openclawRepo *repository.OpenClawRepo
	clientMgr    *cinaclaw.ClientManager
}

// NewMountService creates a new MountService.
func NewMountService(clientMgr *cinaclaw.ClientManager) *MountService {
	return &MountService{
		mountRepo:    repository.NewMountRepo(),
		vmRepo:       repository.NewVMRepo(),
		openclawRepo: repository.NewOpenClawRepo(),
		clientMgr:    clientMgr,
	}
}

// AddMountRequest holds the input for adding a mount.
type AddMountRequest struct {
	VMID       uint   `json:"vm_id" binding:"required"`
	Name       string `json:"name" binding:"required,max=100"`
	HostPath   string `json:"host_path" binding:"required,max=500"`
	VMPath     string `json:"vm_path" binding:"required,max=500"`
	Permission string `json:"permission"` // ro, rw
	AutoMount  bool   `json:"auto_mount"`
}

// OperateMountRequest holds the input for mount operations.
type OperateMountRequest struct {
	Operation string `json:"operation" binding:"required"` // mount, unmount, edit, delete
	Name      string `json:"name"`
	VMPath    string `json:"vm_path"`
	Permission string `json:"permission"`
	AutoMount bool   `json:"auto_mount"`
}

// ConfigureOpenClawMountRequest holds the input for configuring OpenClaw mounts.
type ConfigureOpenClawMountRequest struct {
	VMID             uint   `json:"vm_id" binding:"required"`
	WorkspacePath    string `json:"workspace_path" binding:"required"`
	SkillsPath       string `json:"skills_path" binding:"required"`
	SyncOpenClawJSON bool   `json:"sync_openclaw_json"`
	SyncToolConfigs  bool   `json:"sync_tool_configs"`
}

var (
	ErrMountNotFound    = errors.New("挂载不存在")
	ErrMountInvalidPath = errors.New("挂载路径无效")
)

// ListMounts returns all mounts for a user, optionally filtered by VM.
func (s *MountService) ListMounts(userID uint, vmID uint) ([]model.Mount, error) {
	return s.mountRepo.ListByUser(userID, vmID)
}

// AddMount creates a new mount and optionally auto-mounts it.
func (s *MountService) AddMount(userID uint, req *AddMountRequest) (*model.Mount, error) {
	// Verify VM ownership
	if _, err := s.vmRepo.FindByID(req.VMID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	// Set default permission
	if req.Permission == "" {
		req.Permission = "rw"
	}

	mount := &model.Mount{
		UserID:     userID,
		VMID:       req.VMID,
		Name:       req.Name,
		HostPath:   req.HostPath,
		VMPath:     req.VMPath,
		Status:     "unmounted",
		Permission: req.Permission,
		AutoMount:  req.AutoMount,
	}

	if err := s.mountRepo.Create(mount); err != nil {
		return nil, fmt.Errorf("create mount: %w", err)
	}

	// Auto-mount if requested
	if req.AutoMount {
		go s.performMount(mount, userID)
	}

	slog.Info("mount added", "mount_id", mount.ID, "vm_id", req.VMID, "user_id", userID)
	return mount, nil
}

// OperateMount performs an operation on a mount.
func (s *MountService) OperateMount(mountID, userID uint, req *OperateMountRequest) error {
	mount, err := s.mountRepo.FindByID(mountID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMountNotFound
		}
		return fmt.Errorf("find mount: %w", err)
	}

	switch req.Operation {
	case "mount":
		go s.performMount(mount, userID)
		return nil

	case "unmount":
		go s.performUnmount(mount, userID)
		return nil

	case "edit":
		if req.Name != "" {
			mount.Name = req.Name
		}
		if req.VMPath != "" {
			mount.VMPath = req.VMPath
		}
		if req.Permission != "" {
			mount.Permission = req.Permission
		}
		mount.AutoMount = req.AutoMount
		if err := s.mountRepo.Save(mount); err != nil {
			return fmt.Errorf("save mount: %w", err)
		}
		return nil

	case "delete":
		// Unmount first if mounted
		if mount.Status == "mounted" {
			s.performUnmount(mount, userID)
		}
		if err := s.mountRepo.Delete(mountID, userID); err != nil {
			return fmt.Errorf("delete mount: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("不支持的操作: %s", req.Operation)
	}
}

// GetOpenClawConfig returns the OpenClaw mount configuration for a VM.
func (s *MountService) GetOpenClawConfig(vmID, userID uint) (map[string]interface{}, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	cfg, err := s.openclawRepo.FindByVMID(vmID)
	if err != nil {
		return nil, ErrOpenClawNotConfigured
	}

	// Get related mounts
	mounts, _ := s.mountRepo.ListByVM(vmID)

	return map[string]interface{}{
		"openclaw_config": cfg,
		"mounts":          mounts,
	}, nil
}

// ConfigureOpenClawMount configures OpenClaw workspace and skills mounts.
func (s *MountService) ConfigureOpenClawMount(userID uint, req *ConfigureOpenClawMountRequest) (*model.OpenClawConfig, error) {
	if _, err := s.vmRepo.FindByID(req.VMID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	// Create or update OpenClaw config
	cfg := &model.OpenClawConfig{
		VMID:             req.VMID,
		WorkspacePath:    req.WorkspacePath,
		SkillsPath:       req.SkillsPath,
		SyncOpenClawJSON: req.SyncOpenClawJSON,
		SyncToolConfigs:  req.SyncToolConfigs,
	}

	existing, err := s.openclawRepo.FindByVMID(req.VMID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.openclawRepo.Create(cfg); err != nil {
				return nil, fmt.Errorf("create openclaw config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("find openclaw config: %w", err)
		}
	} else {
		cfg.ID = existing.ID
		if err := s.openclawRepo.Save(cfg); err != nil {
			return nil, fmt.Errorf("save openclaw config: %w", err)
		}
	}

	// Create workspace mount
	workspaceMount := &model.Mount{
		UserID:     userID,
		VMID:       req.VMID,
		Name:       "openclaw-workspace",
		HostPath:   req.WorkspacePath,
		VMPath:     "/root/.openclaw/workspace",
		Status:     "unmounted",
		Permission: "rw",
		AutoMount:  true,
	}
	_ = s.mountRepo.FirstOrCreate(workspaceMount, model.Mount{
		UserID: userID,
		VMID:   req.VMID,
		Name:   "openclaw-workspace",
	})

	// Create skills mount
	skillsMount := &model.Mount{
		UserID:     userID,
		VMID:       req.VMID,
		Name:       "openclaw-skills",
		HostPath:   req.SkillsPath,
		VMPath:     "/root/.openclaw/workspace/skills",
		Status:     "unmounted",
		Permission: "rw",
		AutoMount:  true,
	}
	_ = s.mountRepo.FirstOrCreate(skillsMount, model.Mount{
		UserID: userID,
		VMID:   req.VMID,
		Name:   "openclaw-skills",
	})

	slog.Info("openclaw mount configured", "vm_id", req.VMID, "user_id", userID)
	return cfg, nil
}

// performMount executes the mount operation via CinaClaw.
func (s *MountService) performMount(mount *model.Mount, userID uint) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := s.clientMgr.GetClient(strconv.FormatUint(uint64(userID), 10))
	if err != nil {
		slog.Error("cinaclaw client error for mount", "error", err, "mount_id", mount.ID)
		_ = s.mountRepo.UpdateStatus(mount.ID, "error")
		return
	}

	// Get VM info for instance name
	vm, err := s.vmRepo.FindByIDNoScope(mount.VMID)
	if err != nil {
		slog.Error("vm not found for mount", "error", err, "vm_id", mount.VMID)
		return
	}

	err = client.Mount(ctx, mount.HostPath, []cinaclaw.MountTarget{
		{
			InstanceName: vm.Name,
			TargetPath:   mount.VMPath,
		},
	})
	if err != nil {
		slog.Error("cinaclaw mount failed", "error", err, "mount_id", mount.ID)
		_ = s.mountRepo.UpdateStatus(mount.ID, "error")
		return
	}

	_ = s.mountRepo.UpdateStatus(mount.ID, "mounted")
	slog.Info("mount completed", "mount_id", mount.ID, "vm_id", mount.VMID)
}

// performUnmount executes the unmount operation via CinaClaw.
func (s *MountService) performUnmount(mount *model.Mount, userID uint) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := s.clientMgr.GetClient(strconv.FormatUint(uint64(userID), 10))
	if err != nil {
		slog.Error("cinaclaw client error for unmount", "error", err, "mount_id", mount.ID)
		return
	}

	vm, err := s.vmRepo.FindByIDNoScope(mount.VMID)
	if err != nil {
		slog.Error("vm not found for unmount", "error", err, "vm_id", mount.VMID)
		return
	}

	err = client.Unmount(ctx, []cinaclaw.MountTarget{
		{
			InstanceName: vm.Name,
			TargetPath:   mount.VMPath,
		},
	})
	if err != nil {
		slog.Error("cinaclaw unmount failed", "error", err, "mount_id", mount.ID)
		return
	}

	_ = s.mountRepo.UpdateStatus(mount.ID, "unmounted")
	slog.Info("unmount completed", "mount_id", mount.ID, "vm_id", mount.VMID)
}
