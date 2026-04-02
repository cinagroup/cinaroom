package model

import (
	"net"
	"time"
)

// User represents a registered user.
type User struct {
	ID               uint       `gorm:"primaryKey" json:"id" form:"id"`
	CinatokenID      uint       `gorm:"uniqueIndex;not null" json:"cinatoken_id" form:"cinatoken_id"`
	Username         string     `gorm:"uniqueIndex;size:20;not null" json:"username" form:"username"`
	Email            string     `gorm:"uniqueIndex;size:100;not null" json:"email" form:"email" binding:"email"`
	Password         string     `gorm:"size:255" json:"-" form:"password"`
	Nickname         string     `gorm:"size:50" json:"nickname" form:"nickname"`
	Phone            string     `gorm:"size:20" json:"phone" form:"phone"`
	Avatar           string     `gorm:"size:255" json:"avatar" form:"avatar"`
	Provider         string     `gorm:"size:50" json:"provider" form:"provider"`
	Active           bool       `gorm:"default:true" json:"active" form:"active"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
	LastLoginAt      *time.Time `json:"last_login_at" form:"last_login_at"`
	TwoFactorEnabled bool       `gorm:"default:false" json:"two_factor_enabled" form:"two_factor_enabled"`
}

// TableName specifies the table name for User.
func (User) TableName() string { return "users" }

// VM represents a virtual machine instance.
type VM struct {
	ID          uint      `gorm:"primaryKey" json:"id" form:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id" form:"user_id"`
	Name        string    `gorm:"size:100;not null" json:"name" form:"name" binding:"max=100"`
	Status      string    `gorm:"size:20;default:stopped" json:"status" form:"status"`
	IP          string    `gorm:"size:50" json:"ip" form:"ip"`
	Image       string    `gorm:"size:50;not null" json:"image" form:"image" binding:"required"`
	CPU         int       `gorm:"default:1" json:"cpu" form:"cpu" binding:"min=1,max=8"`
	Memory      int       `gorm:"default:1" json:"memory" form:"memory" binding:"min=1,max=16"`
	Disk        int       `gorm:"default:10" json:"disk" form:"disk" binding:"min=10,max=500"`
	NetworkType string    `gorm:"size:20;default:nat" json:"network_type" form:"network_type"`
	SSHKey      string    `gorm:"size:500" json:"-" form:"ssh_key"`
	InitScript  string    `gorm:"type:text" json:"-" form:"init_script"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
}

// TableName specifies the table name for VM.
func (VM) TableName() string { return "vms" }

// VMSnapshot represents a point-in-time snapshot of a VM.
type VMSnapshot struct {
	ID        uint      `gorm:"primaryKey" json:"id" form:"id"`
	VMID      uint      `gorm:"index;not null" json:"vm_id" form:"vm_id"`
	Name      string    `gorm:"size:100;not null" json:"name" form:"name" binding:"max=100"`
	Size      int64     `json:"size" form:"size"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
}

// TableName specifies the table name for VMSnapshot.
func (VMSnapshot) TableName() string { return "vm_snapshots" }

// VMLog records an operation performed on a VM.
type VMLog struct {
	ID        uint      `gorm:"primaryKey" json:"id" form:"id"`
	VMID      uint      `gorm:"index;not null" json:"vm_id" form:"vm_id"`
	Operation string    `gorm:"size:50;not null" json:"operation" form:"operation"`
	Result    string    `gorm:"size:20;not null" json:"result" form:"result"`
	Message   string    `gorm:"type:text" json:"message" form:"message"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
}

// TableName specifies the table name for VMLog.
func (VMLog) TableName() string { return "vm_logs" }

// Mount represents a directory mount between host and VM.
type Mount struct {
	ID         uint      `gorm:"primaryKey" json:"id" form:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id" form:"user_id"`
	VMID       uint      `gorm:"index;not null" json:"vm_id" form:"vm_id"`
	Name       string    `gorm:"size:100;not null" json:"name" form:"name" binding:"max=100"`
	HostPath   string    `gorm:"size:500;not null" json:"host_path" form:"host_path" binding:"max=500"`
	VMPath     string    `gorm:"size:500;not null" json:"vm_path" form:"vm_path" binding:"max=500"`
	Status     string    `gorm:"size:20;default:unmounted" json:"status" form:"status"`
	Permission string    `gorm:"size:10;default:rw" json:"permission" form:"permission"`
	AutoMount  bool      `gorm:"default:false" json:"auto_mount" form:"auto_mount"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
}

// TableName specifies the table name for Mount.
func (Mount) TableName() string { return "mounts" }

// OpenClawConfig stores per-VM OpenClaw deployment settings.
type OpenClawConfig struct {
	ID               uint       `gorm:"primaryKey" json:"id" form:"id"`
	VMID             uint       `gorm:"index;not null" json:"vm_id" form:"vm_id"`
	Status           string     `gorm:"size:20;default:not_installed" json:"status" form:"status"`
	Version          string     `gorm:"size:20" json:"version" form:"version"`
	RunningTime      int64      `json:"running_time" form:"running_time"`
	DefaultModel     string     `gorm:"size:100" json:"default_model" form:"default_model"`
	APIKey           string     `gorm:"size:255" json:"-" form:"api_key"`
	EnabledTools     string     `gorm:"type:text" json:"enabled_tools" form:"enabled_tools"`
	EnabledSkills    string     `gorm:"type:text" json:"enabled_skills" form:"enabled_skills"`
	WorkspacePath    string     `gorm:"size:500" json:"workspace_path" form:"workspace_path"`
	SkillsPath       string     `gorm:"size:500" json:"skills_path" form:"skills_path"`
	SyncOpenClawJSON bool       `gorm:"default:true" json:"sync_openclaw_json" form:"sync_openclaw_json"`
	SyncToolConfigs  bool       `gorm:"default:true" json:"sync_tool_configs" form:"sync_tool_configs"`
	LastDeployedAt   *time.Time `json:"last_deployed_at" form:"last_deployed_at"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
}

// TableName specifies the table name for OpenClawConfig.
func (OpenClawConfig) TableName() string { return "openclaw_configs" }

// RemoteAccess stores remote-access settings for a VM.
type RemoteAccess struct {
	ID            uint      `gorm:"primaryKey" json:"id" form:"id"`
	VMID          uint      `gorm:"uniqueIndex;not null" json:"vm_id" form:"vm_id"`
	Enabled       bool      `gorm:"default:false" json:"enabled" form:"enabled"`
	AccessAddress string    `gorm:"size:255" json:"access_address" form:"access_address"`
	QRCode        string    `gorm:"size:500" json:"qr_code" form:"qr_code"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
}

// TableName specifies the table name for RemoteAccess.
func (RemoteAccess) TableName() string { return "remote_access" }

// IPWhitelist stores allowed IPs/CIDRs for remote access.
type IPWhitelist struct {
	ID        uint      `gorm:"primaryKey" json:"id" form:"id"`
	VMID      uint      `gorm:"index;not null" json:"vm_id" form:"vm_id"`
	IP        string    `gorm:"size:50;not null" json:"ip" form:"ip"`
	IsCIDR    bool      `gorm:"default:false" json:"is_cidr" form:"is_cidr"`
	Note      string    `gorm:"size:200" json:"note" form:"note"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
}

// TableName specifies the table name for IPWhitelist.
func (IPWhitelist) TableName() string { return "ip_whitelists" }

// IsCIDRValid checks whether the IP field is a valid CIDR or plain IP.
func (w *IPWhitelist) IsCIDRValid() bool {
	if _, _, err := net.ParseCIDR(w.IP); err == nil {
		return true
	}
	if net.ParseIP(w.IP) != nil {
		return true
	}
	return false
}

// RemoteLog records a remote access event.
type RemoteLog struct {
	ID           uint      `gorm:"primaryKey" json:"id" form:"id"`
	VMID         uint      `gorm:"index;not null" json:"vm_id" form:"vm_id"`
	AccessTime   time.Time `gorm:"index" json:"access_time" form:"access_time"`
	AccessIP     string    `gorm:"size:50;not null" json:"access_ip" form:"access_ip"`
	AccessPath   string    `gorm:"size:500" json:"access_path" form:"access_path"`
	UserAgent    string    `gorm:"size:500" json:"user_agent" form:"user_agent"`
	ResponseCode int       `json:"response_code" form:"response_code"`
}

// TableName specifies the table name for RemoteLog.
func (RemoteLog) TableName() string { return "remote_logs" }

// LoginLog records a user login event.
type LoginLog struct {
	ID        uint      `gorm:"primaryKey" json:"id" form:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id" form:"user_id"`
	LoginTime time.Time `gorm:"index" json:"login_time" form:"login_time"`
	IP        string    `gorm:"size:50;not null" json:"ip" form:"ip"`
	Location  string    `gorm:"size:200" json:"location" form:"location"`
	Device    string    `gorm:"size:200" json:"device" form:"device"`
}

// TableName specifies the table name for LoginLog.
func (LoginLog) TableName() string { return "login_logs" }

// Session represents an active user session.
type Session struct {
	ID           uint      `gorm:"primaryKey" json:"id" form:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id" form:"user_id"`
	Token        string    `gorm:"uniqueIndex;size:255;not null" json:"-" form:"token"`
	Device       string    `gorm:"size:200" json:"device" form:"device"`
	Location     string    `gorm:"size:200" json:"location" form:"location"`
	LoginTime    time.Time `json:"login_time" form:"login_time"`
	LastActiveAt time.Time `json:"last_active_at" form:"last_active_at"`
	ExpiredAt    time.Time `json:"expired_at" form:"expired_at"`
}

// TableName specifies the table name for Session.
func (Session) TableName() string { return "sessions" }

// SystemSetting stores a key-value system setting.
type SystemSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id" form:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key" form:"key" binding:"max=100"`
	Value     string    `gorm:"type:text" json:"value" form:"value"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`
}

// TableName specifies the table name for SystemSetting.
func (SystemSetting) TableName() string { return "system_settings" }

// Subscription 用户订阅
type Subscription struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Plan      string    `gorm:"size:20;not null;default:'free'" json:"plan"`      // free, pro, enterprise
	Status    string    `gorm:"size:20;not null;default:'active'" json:"status"`  // active, expired, cancelled
	StartedAt time.Time `json:"started_at"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for Subscription.
func (Subscription) TableName() string { return "subscriptions" }

// VMMetric stores a point-in-time resource usage sample for a VM.
type VMMetric struct {
	ID          uint      `gorm:"primaryKey" json:"id" form:"id"`
	VMID        uint      `gorm:"index;not null" json:"vm_id" form:"vm_id"`
	CPUUsage    float64   `json:"cpu_usage" form:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage" form:"memory_usage"`
	DiskIO      float64   `json:"disk_io" form:"disk_io"`
	NetworkRX   float64   `json:"network_rx" form:"network_rx"`
	NetworkTX   float64   `json:"network_tx" form:"network_tx"`
	Timestamp   time.Time `gorm:"index" json:"timestamp" form:"timestamp"`
}

// TableName specifies the table name for VMMetric.
func (VMMetric) TableName() string { return "vm_metrics" }
