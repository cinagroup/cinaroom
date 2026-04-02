package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/cinagroup/cinaseek/backend/internal/model"
	"github.com/cinagroup/cinaseek/backend/internal/repository"
	"gorm.io/gorm"
)

// RemoteService handles remote access management.
type RemoteService struct {
	remoteRepo *repository.RemoteRepo
	vmRepo     *repository.VMRepo
}

// NewRemoteService creates a new RemoteService.
func NewRemoteService() *RemoteService {
	return &RemoteService{
		remoteRepo: repository.NewRemoteRepo(),
		vmRepo:     repository.NewVMRepo(),
	}
}

// SwitchRemoteAccessRequest holds the input for toggling remote access.
type SwitchRemoteAccessRequest struct {
	Enabled bool `json:"enabled"`
}

// AddIPWhitelistRequest holds the input for adding an IP whitelist entry.
type AddIPWhitelistRequest struct {
	VMID uint   `json:"vm_id" binding:"required"`
	IP   string `json:"ip" binding:"required"`
	Note string `json:"note"`
}

var (
	ErrRemoteNotFound    = errors.New("远程访问配置不存在")
	ErrInvalidIP         = errors.New("IP 地址格式无效")
	ErrIPWhitelistExists = errors.New("IP 白名单已存在")
)

// GetStatus returns the remote access status for a VM.
func (s *RemoteService) GetStatus(vmID, userID uint) (map[string]interface{}, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	ra, err := s.remoteRepo.FindRemoteAccess(vmID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{
				"enabled":        false,
				"status":         "disabled",
				"access_address": "",
				"qr_code":        "",
			}, nil
		}
		return nil, fmt.Errorf("find remote access: %w", err)
	}

	status := "disabled"
	if ra.Enabled {
		status = "enabled"
	}

	return map[string]interface{}{
		"enabled":        ra.Enabled,
		"status":         status,
		"access_address": ra.AccessAddress,
		"qr_code":        ra.QRCode,
	}, nil
}

// SwitchRemoteAccess enables or disables remote access for a VM.
func (s *RemoteService) SwitchRemoteAccess(vmID, userID uint, req *SwitchRemoteAccessRequest) (*model.RemoteAccess, error) {
	vm, err := s.vmRepo.FindByID(vmID, userID)
	if err != nil {
		return nil, ErrVMNotFound
	}

	ra, err := s.remoteRepo.FindRemoteAccess(vmID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ra = &model.RemoteAccess{
				VMID:    vmID,
				Enabled: req.Enabled,
			}
		} else {
			return nil, fmt.Errorf("find remote access: %w", err)
		}
	}

	ra.Enabled = req.Enabled

	if req.Enabled {
		// Generate access address from VM IP
		if vm.IP != "" {
			ra.AccessAddress = "https://" + vm.IP + ":8080"
		} else {
			ra.AccessAddress = ""
		}
		// TODO: Generate QR code using a QR code library
		ra.QRCode = ""
	} else {
		ra.AccessAddress = ""
		ra.QRCode = ""
	}

	if ra.ID == 0 {
		if err := s.remoteRepo.CreateRemoteAccess(ra); err != nil {
			return nil, fmt.Errorf("create remote access: %w", err)
		}
	} else {
		if err := s.remoteRepo.SaveRemoteAccess(ra); err != nil {
			return nil, fmt.Errorf("save remote access: %w", err)
		}
	}

	action := "禁用"
	if req.Enabled {
		action = "启用"
	}
	slog.Info("remote access toggled", "vm_id", vmID, "user_id", userID, "action", action)

	return ra, nil
}

// GetIPWhitelist returns the IP whitelist for a VM.
func (s *RemoteService) GetIPWhitelist(vmID, userID uint) ([]model.IPWhitelist, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, ErrVMNotFound
	}
	return s.remoteRepo.ListIPWhitelist(vmID)
}

// AddIPWhitelist adds a new IP whitelist entry for a VM.
func (s *RemoteService) AddIPWhitelist(userID uint, req *AddIPWhitelistRequest) (*model.IPWhitelist, error) {
	if _, err := s.vmRepo.FindByID(req.VMID, userID); err != nil {
		return nil, ErrVMNotFound
	}

	// Validate IP or CIDR
	ip := net.ParseIP(req.IP)
	if ip == nil {
		// Try CIDR notation
		_, _, err := net.ParseCIDR(req.IP)
		if err != nil {
			return nil, ErrInvalidIP
		}
	}

	entry := &model.IPWhitelist{
		VMID: req.VMID,
		IP:   req.IP,
		Note: req.Note,
	}

	if err := s.remoteRepo.CreateIPWhitelist(entry); err != nil {
		return nil, fmt.Errorf("create ip whitelist: %w", err)
	}

	// TODO: Update firewall rules on the VM to allow this IP

	slog.Info("ip whitelist added", "vm_id", req.VMID, "ip", req.IP, "user_id", userID)
	return entry, nil
}

// RemoveIPWhitelist removes an IP whitelist entry.
func (s *RemoteService) RemoveIPWhitelist(entryID, vmID, userID uint) error {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return ErrVMNotFound
	}

	affected, err := s.remoteRepo.DeleteIPWhitelist(entryID, vmID)
	if err != nil {
		return fmt.Errorf("delete ip whitelist: %w", err)
	}
	if affected == 0 {
		return ErrIPWhitelistExists // reuse as "not found"
	}

	// TODO: Update firewall rules on the VM to remove this IP

	slog.Info("ip whitelist removed", "vm_id", vmID, "entry_id", entryID, "user_id", userID)
	return nil
}

// GetRemoteLogs returns paginated remote access logs for a VM.
func (s *RemoteService) GetRemoteLogs(vmID, userID uint, ipFilter string, page, pageSize int) ([]model.RemoteLog, int64, error) {
	if _, err := s.vmRepo.FindByID(vmID, userID); err != nil {
		return nil, 0, ErrVMNotFound
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}

	return s.remoteRepo.ListRemoteLogs(vmID, ipFilter, page, pageSize)
}

// VerifyVMOwnership is a helper that checks a VM belongs to a user and returns the VM ID from param.
func (s *RemoteService) VerifyVMOwnership(vmIDStr string, userID uint) (uint, error) {
	vmID, err := strconv.ParseUint(vmIDStr, 10, 32)
	if err != nil {
		return 0, errors.New("无效的虚拟机 ID")
	}
	if _, err := s.vmRepo.FindByID(uint(vmID), userID); err != nil {
		return 0, ErrVMNotFound
	}
	return uint(vmID), nil
}
