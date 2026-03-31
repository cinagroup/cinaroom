package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	CinatokenID     uint       `gorm:"uniqueIndex;not null" json:"cinatoken_id"` // CinaToken 用户 ID
	Username        string     `gorm:"uniqueIndex;size:20;not null" json:"username"`
	Email           string     `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password        string     `gorm:"size:255" json:"-"` // OAuth 用户为空
	Nickname        string     `gorm:"size:50" json:"nickname"`
	Phone           string     `gorm:"size:20" json:"phone"`
	Avatar          string     `gorm:"size:255" json:"avatar"`
	Provider        string     `gorm:"size:50" json:"provider"` // OAuth 提供商：github/google/microsoft 等
	Active          bool       `gorm:"default:true" json:"active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	TwoFactorEnabled bool      `gorm:"default:false" json:"two_factor_enabled"`
}

// VM 虚拟机模型
type VM struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Status       string    `gorm:"size:20;default:stopped" json:"status"` // running, stopped, paused
	IP           string    `gorm:"size:50" json:"ip"`
	Image        string    `gorm:"size:50;not null" json:"image"`
	CPU          int       `gorm:"default:1" json:"cpu"`
	Memory       int       `gorm:"default:1" json:"memory"` // GB
	Disk         int       `gorm:"default:10" json:"disk"`  // GB
	NetworkType  string    `gorm:"size:20;default:nat" json:"network_type"`
	SSHKey       string    `gorm:"size:500" json:"-"`
	InitScript   string    `gorm:"type:text" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// VMSnapshot 虚拟机快照
type VMSnapshot struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VMID      uint      `gorm:"index;not null" json:"vm_id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Size      int64     `json:"size"` // bytes
	CreatedAt time.Time `json:"created_at"`
}

// VMLog 虚拟机操作日志
type VMLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	VMID       uint      `gorm:"index;not null" json:"vm_id"`
	Operation  string    `gorm:"size:50;not null" json:"operation"`
	Result     string    `gorm:"size:20;not null" json:"result"` // success, failed
	Message    string    `gorm:"type:text" json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

// Mount 目录挂载
type Mount struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"index;not null" json:"user_id"`
	VMID           uint      `gorm:"index;not null" json:"vm_id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	HostPath       string    `gorm:"size:500;not null" json:"host_path"`
	VMPath         string    `gorm:"size:500;not null" json:"vm_path"`
	Status         string    `gorm:"size:20;default:unmounted" json:"status"` // mounted, unmounted
	Permission     string    `gorm:"size:10;default:rw" json:"permission"`    // ro, rw
	AutoMount      bool      `gorm:"default:false" json:"auto_mount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// OpenClawConfig OpenClaw 配置
type OpenClawConfig struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	VMID               uint      `gorm:"index;not null" json:"vm_id"`
	Status             string    `gorm:"size:20;default:not_installed" json:"status"`
	Version            string    `gorm:"size:20" json:"version"`
	RunningTime        int64     `json:"running_time"` // seconds
	DefaultModel       string    `gorm:"size:100" json:"default_model"`
	APIKey             string    `gorm:"size:255" json:"-"`
	EnabledTools       string    `gorm:"type:text" json:"enabled_tools"` // JSON array
	EnabledSkills      string    `gorm:"type:text" json:"enabled_skills"` // JSON array
	WorkspacePath      string    `gorm:"size:500" json:"workspace_path"`
	SkillsPath         string    `gorm:"size:500" json:"skills_path"`
	SyncOpenClawJSON   bool      `gorm:"default:true" json:"sync_openclaw_json"`
	SyncToolConfigs    bool      `gorm:"default:true" json:"sync_tool_configs"`
	LastDeployedAt     *time.Time `json:"last_deployed_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// RemoteAccess 远程访问配置
type RemoteAccess struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	VMID          uint      `gorm:"uniqueIndex;not null" json:"vm_id"`
	Enabled       bool      `gorm:"default:false" json:"enabled"`
	AccessAddress string    `gorm:"size:255" json:"access_address"`
	QRCode        string    `gorm:"size:500" json:"qr_code"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// IPWhitelist IP 白名单
type IPWhitelist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VMID      uint      `gorm:"index;not null" json:"vm_id"`
	IP        string    `gorm:"size:50;not null" json:"ip"`
	Note      string    `gorm:"size:200" json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// RemoteLog 远程访问日志
type RemoteLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	VMID         uint      `gorm:"index;not null" json:"vm_id"`
	AccessTime   time.Time `gorm:"index" json:"access_time"`
	AccessIP     string    `gorm:"size:50;not null" json:"access_ip"`
	AccessPath   string    `gorm:"size:500" json:"access_path"`
	UserAgent    string    `gorm:"size:500" json:"user_agent"`
	ResponseCode int       `json:"response_code"`
}

// LoginLog 用户登录日志
type LoginLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	LoginTime  time.Time `gorm:"index" json:"login_time"`
	IP         string    `gorm:"size:50;not null" json:"ip"`
	Location   string    `gorm:"size:200" json:"location"`
	Device     string    `gorm:"size:200" json:"device"`
}

// Session 用户会话
type Session struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Token        string    `gorm:"uniqueIndex;size:255;not null" json:"-"`
	Device       string    `gorm:"size:200" json:"device"`
	Location     string    `gorm:"size:200" json:"location"`
	LoginTime    time.Time `json:"login_time"`
	LastActiveAt time.Time `json:"last_active_at"`
	ExpiredAt    time.Time `json:"expired_at"`
}

// SystemSetting 系统设置
type SystemSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// VMMetric 虚拟机监控指标
type VMMetric struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	VMID       uint      `gorm:"index;not null" json:"vm_id"`
	CPUUsage   float64   `json:"cpu_usage"`
	MemoryUsage float64  `json:"memory_usage"`
	DiskIO     float64   `json:"disk_io"`
	NetworkRX  float64   `json:"network_rx"`
	NetworkTX  float64   `json:"network_tx"`
	Timestamp  time.Time `gorm:"index" json:"timestamp"`
}
