// Package cinaclaw provides a Go client for communicating with the CinaClaw VM engine
// via gRPC over Unix domain sockets.
package cinaclaw

import (
	"time"
)

// CreateVMRequest defines the parameters for creating a new virtual machine.
type CreateVMRequest struct {
	Name   string `json:"name"`   // Instance name
	CPU    int    `json:"cpu"`    // Number of CPU cores
	Memory string `json:"memory"` // Memory size, e.g. "2G"
	Disk   string `json:"disk"`   // Disk size, e.g. "20G"
	Image  string `json:"image"`  // Image identifier, e.g. "22.04"
}

// VMInfo holds information about a virtual machine instance.
type VMInfo struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"` // RUNNING, STOPPED, SUSPENDED, etc.
	CPU       int       `json:"cpu"`
	Memory    string    `json:"memory"`
	Disk      string    `json:"disk"`
	IP        string    `json:"ip"`
	OS        string    `json:"os"`
	Release   string    `json:"release"`
	CreatedAt time.Time `json:"created_at"`
}

// VMMetrics holds real-time resource usage metrics for a VM.
type VMMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}

// SnapshotInfo holds information about a VM snapshot.
type SnapshotInfo struct {
	Name      string    `json:"name"`
	Parent    string    `json:"parent,omitempty"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// SSHInfo holds SSH connection details for a VM.
type SSHInfo struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Key      string `json:"key"` // Base64-encoded private key
}

// MountTarget defines a mount target inside a VM.
type MountTarget struct {
	InstanceName string `json:"instance_name"`
	TargetPath   string `json:"target_path"`
}

// ImageInfo holds information about a available VM image.
type ImageInfo struct {
	OS        string   `json:"os"`
	Release   string   `json:"release"`
	Version   string   `json:"version"`
	Aliases   []string `json:"aliases"`
	Codename  string   `json:"codename,omitempty"`
	Remote    string   `json:"remote"`
}
