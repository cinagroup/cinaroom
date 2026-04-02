package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/model"
	"github.com/cinagroup/cinaseek/backend/internal/repository"
	"gorm.io/gorm"
)

// SystemService handles system monitoring and settings.
type SystemService struct {
	systemRepo  *repository.SystemRepo
	vmRepo      *repository.VMRepo
	mountRepo   *repository.MountRepo
	openclawRepo *repository.OpenClawRepo
	metricRepo  *repository.MetricRepo
}

// NewSystemService creates a new SystemService.
func NewSystemService() *SystemService {
	return &SystemService{
		systemRepo:  repository.NewSystemRepo(),
		vmRepo:      repository.NewVMRepo(),
		mountRepo:   repository.NewMountRepo(),
		openclawRepo: repository.NewOpenClawRepo(),
		metricRepo:  repository.NewMetricRepo(),
	}
}

var (
	ErrSettingNotFound = errors.New("系统设置不存在")
)

// GetSetting retrieves a system setting by key.
func (s *SystemService) GetSetting(key string) (map[string]interface{}, error) {
	setting, err := s.systemRepo.GetSetting(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSettingNotFound
		}
		return nil, fmt.Errorf("get setting: %w", err)
	}

	return map[string]interface{}{
		"key":   setting.Key,
		"value": setting.Value,
	}, nil
}

// UpdateSetting creates or updates a system setting.
func (s *SystemService) UpdateSetting(key, value string) error {
	if err := s.systemRepo.SetSetting(key, value); err != nil {
		return fmt.Errorf("update setting: %w", err)
	}

	slog.Info("system setting updated", "key", key)
	return nil
}

// GetAllSettings returns all system settings.
func (s *SystemService) GetAllSettings() (map[string]string, error) {
	return s.systemRepo.GetAllSettings()
}

// GetVersion returns system version information.
func (s *SystemService) GetVersion() map[string]interface{} {
	return map[string]interface{}{
		"version":     "1.0.0",
		"build":       "20260401",
		"api_version": "v1",
		"go_version":  "1.21.5",
		"features": []string{
			"user_management",
			"vm_management",
			"web_shell",
			"directory_mount",
			"openclaw_integration",
			"remote_access",
		},
	}
}

// GetDashboard returns dashboard statistics for a user.
func (s *SystemService) GetDashboard(userID uint) (map[string]interface{}, error) {
	vmTotal, err := s.vmRepo.CountByUser(userID, "")
	if err != nil {
		return nil, fmt.Errorf("count vms: %w", err)
	}
	vmRunning, err := s.vmRepo.CountByUser(userID, "running")
	if err != nil {
		return nil, fmt.Errorf("count running vms: %w", err)
	}

	mountTotal, err := s.mountRepo.CountByUser(userID)
	if err != nil {
		slog.Warn("failed to count mounts", "error", err)
		mountTotal = 0
	}

	openclawTotal, err := s.openclawRepo.CountByUser(userID)
	if err != nil {
		slog.Warn("failed to count openclaw configs", "error", err)
		openclawTotal = 0
	}

	return map[string]interface{}{
		"vm_total":       vmTotal,
		"vm_running":     vmRunning,
		"vm_stopped":     vmTotal - vmRunning,
		"mount_total":    mountTotal,
		"openclaw_total": openclawTotal,
	}, nil
}

// GetStatistics returns detailed statistics for a user.
func (s *SystemService) GetStatistics(userID uint) (map[string]interface{}, error) {
	// VM status breakdown
	vmTotal, err := s.vmRepo.CountByUser(userID, "")
	if err != nil {
		return nil, fmt.Errorf("count vms: %w", err)
	}
	vmRunning, err := s.vmRepo.CountByUser(userID, "running")
	if err != nil {
		return nil, fmt.Errorf("count running vms: %w", err)
	}
	vmStopped, err := s.vmRepo.CountByUser(userID, "stopped")
	if err != nil {
		slog.Warn("failed to count stopped vms", "error", err)
		vmStopped = 0
	}

	return map[string]interface{}{
		"vm_total":   vmTotal,
		"vm_running": vmRunning,
		"vm_stopped": vmStopped,
		"vm_paused":  vmTotal - vmRunning - vmStopped,
	}, nil
}

// HealthCheck performs a system health check.
func (s *SystemService) HealthCheck() error {
	return repository.HealthCheck()
}

// GetHealthStatus returns a full health status map.
func (s *SystemService) GetHealthStatus() map[string]interface{} {
	dbErr := repository.HealthCheck()
	dbStatus := "connected"
	if dbErr != nil {
		dbStatus = "disconnected"
	}

	return map[string]interface{}{
		"status":    "healthy",
		"database":  dbStatus,
		"timestamp": time.Now(),
	}
}

// SearchVMs searches VMs by keyword for a user.
func (s *SystemService) SearchVMs(userID uint, keyword string) ([]model.VM, error) {
	return s.vmRepo.Search(userID, keyword, 10)
}

// BatchOperateVMs performs an operation on multiple VMs.
func (s *SystemService) BatchOperateVMs(ids []uint, userID uint, operation string) (int, int, error) {
	vms, err := s.vmRepo.FindByIDsAndUser(ids, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("find vms: %w", err)
	}
	if len(vms) != len(ids) {
		return 0, 0, errors.New("部分虚拟机不存在或无权限")
	}

	successCount := 0
	failCount := 0

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
		successCount++
	}

	return successCount, failCount, nil
}
