package service

import (
	"context"
	"encoding/json"
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

// OpenClawService handles OpenClaw deployment management.
type OpenClawService struct {
	openclawRepo *repository.OpenClawRepo
	vmRepo       *repository.VMRepo
	vmLogRepo    *repository.VMLogRepo
	clientMgr    *cinaclaw.ClientManager
}

// NewOpenClawService creates a new OpenClawService.
func NewOpenClawService(clientMgr *cinaclaw.ClientManager) *OpenClawService {
	return &OpenClawService{
		openclawRepo: repository.NewOpenClawRepo(),
		vmRepo:       repository.NewVMRepo(),
		vmLogRepo:    repository.NewVMLogRepo(),
		clientMgr:    clientMgr,
	}
}

// DeployOpenClawRequest holds the input for deploying OpenClaw.
type DeployOpenClawRequest struct {
	VMID         uint   `json:"vm_id" binding:"required"`
	Version      string `json:"version"`
	APIKey       string `json:"api_key"`
	DefaultModel string `json:"default_model"`
}

// OperateOpenClawRequest holds the input for OpenClaw operations.
type OperateOpenClawRequest struct {
	Operation string `json:"operation" binding:"required"` // start, stop, restart, update
	Version   string `json:"version"`
}

// UpdateOpenClawConfigRequest holds the input for updating OpenClaw config.
type UpdateOpenClawConfigRequest struct {
	DefaultModel  string   `json:"default_model"`
	APIKey        string   `json:"api_key"`
	EnabledTools  []string `json:"enabled_tools"`
	EnabledSkills []string `json:"enabled_skills"`
}

var (
	ErrOpenClawNotConfigured = errors.New("OpenClaw 未配置")
	ErrOpenClawDeployFail    = errors.New("OpenClaw 部署失败")
)

// GetStatus returns the OpenClaw status for a VM.
func (s *OpenClawService) GetStatus(vmID, userID uint) (map[string]interface{}, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	cfg, err := s.openclawRepo.FindByVMID(vmID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{
				"status":  "not_installed",
				"version": "",
			}, nil
		}
		return nil, fmt.Errorf("find openclaw config: %w", err)
	}

	return map[string]interface{}{
		"status":           cfg.Status,
		"version":          cfg.Version,
		"running_time":     cfg.RunningTime,
		"last_deployed_at": cfg.LastDeployedAt,
		"default_model":    cfg.DefaultModel,
	}, nil
}

// Deploy deploys OpenClaw on a VM.
func (s *OpenClawService) Deploy(userID uint, req *DeployOpenClawRequest) (map[string]interface{}, error) {
	vm, err := s.vmRepo.FindByID(req.VMID, userID)
	if err != nil {
		return nil, ErrVMNotFound
	}

	// Set defaults
	if req.Version == "" {
		req.Version = "latest"
	}
	if req.DefaultModel == "" {
		req.DefaultModel = "qwencode/qwen3.5-plus"
	}

	// Create or update config
	cfg := &model.OpenClawConfig{
		VMID:         req.VMID,
		Version:      req.Version,
		APIKey:       req.APIKey,
		DefaultModel: req.DefaultModel,
		Status:       "deploying",
	}

	now := time.Now()
	cfg.LastDeployedAt = &now

	// Check if config already exists
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

	// Deploy via CinaClaw in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		client, err := s.clientMgr.GetClient(strconv.FormatUint(uint64(userID), 10))
		if err != nil {
			slog.Error("cinaclaw client error for openclaw deploy", "error", err, "vm_id", req.VMID)
			_ = s.openclawRepo.UpdateStatus(req.VMID, "error")
			return
		}

		// Check VM is running
		info, err := client.GetVMInfo(ctx, vm.Name)
		if err != nil {
			slog.Error("cinaclaw get vm info failed", "error", err, "vm_id", req.VMID)
			_ = s.openclawRepo.UpdateStatus(req.VMID, "error")
			return
		}

		if info.Status != "RUNNING" {
			slog.Error("vm not running for openclaw deploy", "vm_status", info.Status, "vm_id", req.VMID)
			_ = s.openclawRepo.UpdateStatus(req.VMID, "error")
			return
		}

		// Simulate deployment steps
		// TODO: Actual SSH-based deployment:
		// 1. SSH connect to VM
		// 2. Install Node.js, pnpm
		// 3. Clone/update OpenClaw
		// 4. Configure openclaw.json
		// 5. Start service

		slog.Info("openclaw deployment simulation completed", "vm_id", req.VMID)
		_ = s.openclawRepo.UpdateStatus(req.VMID, "running")

		s.logVMOperation(req.VMID, "deploy_openclaw", "success",
			fmt.Sprintf("部署 OpenClaw %s", req.Version))
	}()

	return map[string]interface{}{
		"vm_id":   req.VMID,
		"version": req.Version,
		"status":  "deploying",
	}, nil
}

// Operate performs an operation on OpenClaw (start/stop/restart/update).
func (s *OpenClawService) Operate(vmID, userID uint, req *OperateOpenClawRequest) (*model.OpenClawConfig, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	cfg, err := s.openclawRepo.FindByVMID(vmID)
	if err != nil {
		return nil, ErrOpenClawNotConfigured
	}

	// Update status based on operation
	switch req.Operation {
	case "start":
		cfg.Status = "running"
	case "stop":
		cfg.Status = "stopped"
	case "restart":
		cfg.Status = "running"
	case "update":
		if req.Version == "" {
			return nil, errors.New("更新操作需要指定版本号")
		}
		cfg.Version = req.Version
		cfg.Status = "updating"
	default:
		return nil, fmt.Errorf("不支持的操作: %s", req.Operation)
	}

	if err := s.openclawRepo.Save(cfg); err != nil {
		return nil, fmt.Errorf("save openclaw config: %w", err)
	}

	// Execute via CinaClaw in background for update operation
	if req.Operation == "update" {
		go func() {
			// Simulate update delay
			time.Sleep(10 * time.Second)
			_ = s.openclawRepo.UpdateStatus(vmID, "running")
			s.logVMOperation(vmID, "update_openclaw", "success",
				fmt.Sprintf("更新到 %s", req.Version))
		}()
	}

	return cfg, nil
}

// GetLogs retrieves OpenClaw logs for a VM.
func (s *OpenClawService) GetLogs(vmID, userID uint) ([]string, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	// TODO: SSH into VM and retrieve actual log files
	// For now return placeholder
	return []string{
		"[INFO] OpenClaw started successfully",
		"[INFO] Loaded 15 skills",
		"[INFO] Gateway is running on port 3000",
	}, nil
}

// UpdateConfig updates the OpenClaw configuration for a VM.
func (s *OpenClawService) UpdateConfig(vmID, userID uint, req *UpdateOpenClawConfigRequest) (*model.OpenClawConfig, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	cfg, err := s.openclawRepo.FindByVMID(vmID)
	if err != nil {
		return nil, ErrOpenClawNotConfigured
	}

	if req.DefaultModel != "" {
		cfg.DefaultModel = req.DefaultModel
	}
	if req.APIKey != "" {
		cfg.APIKey = req.APIKey
	}
	if req.EnabledTools != nil {
		toolsJSON, _ := json.Marshal(req.EnabledTools)
		cfg.EnabledTools = string(toolsJSON)
	}
	if req.EnabledSkills != nil {
		skillsJSON, _ := json.Marshal(req.EnabledSkills)
		cfg.EnabledSkills = string(skillsJSON)
	}

	if err := s.openclawRepo.Save(cfg); err != nil {
		return nil, fmt.Errorf("save openclaw config: %w", err)
	}

	// TODO: Sync config to VM via SSH

	s.logVMOperation(vmID, "update_openclaw_config", "success", "更新 OpenClaw 配置")
	return cfg, nil
}

// GetMonitorData returns monitoring data for OpenClaw on a VM.
func (s *OpenClawService) GetMonitorData(vmID, userID uint) (map[string]interface{}, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	cfg, err := s.openclawRepo.FindByVMID(vmID)
	if err != nil {
		return nil, ErrOpenClawNotConfigured
	}

	// TODO: Retrieve actual metrics from VM
	_ = cfg
	return map[string]interface{}{
		"cpu_usage":            25.5,
		"memory_usage":         512.3,
		"disk_usage":           15.2,
		"today_requests":       1250,
		"avg_response_time":    1.2,
		"error_rate":           0.5,
		"active_sessions":      5,
		"enabled_tools_count":  10,
		"enabled_skills_count": 15,
	}, nil
}

// GetWorkspaceList returns workspace information for OpenClaw on a VM.
func (s *OpenClawService) GetWorkspaceList(vmID, userID uint) (map[string]interface{}, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	// TODO: SSH into VM and scan workspace directory
	return map[string]interface{}{
		"workspaces": []map[string]interface{}{
			{
				"name":          "default",
				"path":          "/root/.openclaw/workspace",
				"size":          1024 * 1024 * 500,
				"file_count":    1250,
				"last_modified": time.Now(),
			},
		},
		"total_size":  1024 * 1024 * 500,
		"total_files": 1250,
	}, nil
}

// logVMOperation is a helper to record a VM operation log.
func (s *OpenClawService) logVMOperation(vmID uint, operation, result, message string) {
	if err := s.vmLogRepo.Create(&model.VMLog{
		VMID:      vmID,
		Operation: operation,
		Result:    result,
		Message:   message,
	}); err != nil {
		slog.Error("failed to write vm log", "error", err, "vm_id", vmID)
	}
}
