package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/cinaclaw"
	ocmetrics "github.com/cinagroup/cinaseek/backend/internal/openclaw"
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

// ─── Stage 2: OpenClaw Service Management API ─────────────────────────────────

// vmNameByID resolves a VM name from its database ID, scoped to the user.
func (s *OpenClawService) vmNameByID(vmID, userID uint) (string, error) {
	vm, err := s.vmRepo.FindByID(vmID, userID)
	if err != nil {
		return "", ErrVMNotFound
	}
	return vm.Name, nil
}

// InstallOpenClaw deploys OpenClaw into the specified VM by executing the
// install-openclaw.sh script via CinaClaw.
func (s *OpenClawService) InstallOpenClaw(vmName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	client, err := s.clientMgr.GetClient("system")
	if err != nil {
		return fmt.Errorf("cinaclaw client: %w", err)
	}

	// Verify VM exists and is running.
	info, err := client.GetVMInfo(ctx, vmName)
	if err != nil {
		return fmt.Errorf("get vm info: %w", err)
	}
	if info.Status != "RUNNING" {
		return fmt.Errorf("vm %q is not running (status: %s)", vmName, info.Status)
	}

	// TODO: SSH into VM and execute the install script.
	//   The actual implementation would use SSH to run:
	//   curl -fsSL <script-url> | bash
	//   or copy scripts/install-openclaw.sh and execute it.
	slog.Info("openclaw install initiated", "vm_name", vmName)
	return nil
}

// GetOpenClawStatus returns the runtime status of OpenClaw inside a VM.
func (s *OpenClawService) GetOpenClawStatus(vmName string) (*ocmetrics.Status, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := s.clientMgr.GetClient("system")
	if err != nil {
		return nil, fmt.Errorf("cinaclaw client: %w", err)
	}

	info, err := client.GetVMInfo(ctx, vmName)
	if err != nil {
		return nil, fmt.Errorf("get vm info: %w", err)
	}

	// TODO: SSH into VM and check openclaw service status.
	//   systemctl is-active openclaw
	//   openclaw --version
	//   Parse the output to populate the struct.
	_ = info
	return &ocmetrics.Status{
		Running:       false,
		Version:       "",
		Uptime:        0,
		GatewayPort:   3271,
		ActiveConnections: 0,
	}, nil
}

// StartOpenClaw starts the OpenClaw service inside the specified VM.
func (s *OpenClawService) StartOpenClaw(vmName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := s.clientMgr.GetClient("system")
	if err != nil {
		return fmt.Errorf("cinaclaw client: %w", err)
	}

	info, err := client.GetVMInfo(ctx, vmName)
	if err != nil {
		return fmt.Errorf("get vm info: %w", err)
	}
	if info.Status != "RUNNING" {
		return fmt.Errorf("vm %q is not running", vmName)
	}

	// TODO: SSH into VM and run: systemctl start openclaw
	slog.Info("openclaw start initiated", "vm_name", vmName)
	return nil
}

// StopOpenClaw stops the OpenClaw service inside the specified VM.
func (s *OpenClawService) StopOpenClaw(vmName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := s.clientMgr.GetClient("system")
	if err != nil {
		return fmt.Errorf("cinaclaw client: %w", err)
	}

	info, err := client.GetVMInfo(ctx, vmName)
	if err != nil {
		return fmt.Errorf("get vm info: %w", err)
	}
	if info.Status != "RUNNING" {
		return fmt.Errorf("vm %q is not running", vmName)
	}

	// TODO: SSH into VM and run: systemctl stop openclaw
	slog.Info("openclaw stop initiated", "vm_name", vmName)
	return nil
}

// RestartOpenClaw restarts the OpenClaw service inside the specified VM.
func (s *OpenClawService) RestartOpenClaw(vmName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := s.clientMgr.GetClient("system")
	if err != nil {
		return fmt.Errorf("cinaclaw client: %w", err)
	}

	info, err := client.GetVMInfo(ctx, vmName)
	if err != nil {
		return fmt.Errorf("get vm info: %w", err)
	}
	if info.Status != "RUNNING" {
		return fmt.Errorf("vm %q is not running", vmName)
	}

	// TODO: SSH into VM and run: systemctl restart openclaw
	slog.Info("openclaw restart initiated", "vm_name", vmName)
	return nil
}

// GetOpenClawLogs retrieves the last `lines` lines of OpenClaw journal logs from
// the specified VM.
func (s *OpenClawService) GetOpenClawLogs(vmName string, lines int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := s.clientMgr.GetClient("system")
	if err != nil {
		return nil, fmt.Errorf("cinaclaw client: %w", err)
	}

	info, err := client.GetVMInfo(ctx, vmName)
	if err != nil {
		return nil, fmt.Errorf("get vm info: %w", err)
	}
	_ = info

	if lines <= 0 {
		lines = 100
	}

	// TODO: SSH into VM and run:
	//   journalctl -u openclaw -n <lines> --no-pager
	//   Parse output lines.
	return []string{
		"[INFO] OpenClaw log placeholder — connect via SSH for real logs",
	}, nil
}

// GetOpenClawResourceUsage returns resource usage metrics for OpenClaw running
// inside the specified VM.
func (s *OpenClawService) GetOpenClawResourceUsage(vmName string) (*ocmetrics.ResourceUsage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := s.clientMgr.GetClient("system")
	if err != nil {
		return nil, fmt.Errorf("cinaclaw client: %w", err)
	}

	// Fetch VM-level metrics from CinaClaw.
	vmMetrics, err := client.GetMetrics(ctx, vmName)
	if err != nil {
		return nil, fmt.Errorf("get vm metrics: %w", err)
	}

	usage := &ocmetrics.ResourceUsage{
		CPUUsage:         vmMetrics.CPUUsage,
		MemoryUsageMB:    vmMetrics.MemoryUsage,
		DiskUsagePercent: vmMetrics.DiskUsage,
	}

	// TODO: SSH into VM and get process-level metrics:
	//   pid=$(pgrep -f "openclaw gateway")
	//   Read /proc/$pid/stat, /proc/$pid/status, /proc/$pid/fd
	//   Parse into ResourceUsage fields.

	return usage, nil
}

// keep strings import used.
var _ = strings.TrimSpace
