// Package openclaw provides monitoring and metrics collection for OpenClaw
// instances running inside CinaSeek-managed virtual machines.
package openclaw

import "time"

// Metrics holds a snapshot of OpenClaw runtime metrics collected from a VM.
type Metrics struct {
	// CPUUsage is the CPU utilization percentage (0-100).
	CPUUsage float64 `json:"cpu_usage"`

	// MemoryUsage is the memory utilization in MB.
	MemoryUsage float64 `json:"memory_usage"`

	// DiskUsage is the disk utilization percentage (0-100).
	DiskUsage float64 `json:"disk_usage"`

	// Uptime is how long the OpenClaw process has been running.
	Uptime time.Duration `json:"uptime"`

	// Version is the installed OpenClaw version string.
	Version string `json:"version"`

	// ActiveModels is the number of currently loaded AI models.
	ActiveModels int `json:"active_models"`

	// CollectedAt is the timestamp when these metrics were gathered.
	CollectedAt time.Time `json:"collected_at"`
}

// ResourceUsage holds system-level resource usage of the OpenClaw process.
type ResourceUsage struct {
	// CPUUsage is the CPU utilization percentage (0-100).
	CPUUsage float64 `json:"cpu_usage"`

	// MemoryUsageMB is the memory usage in megabytes.
	MemoryUsageMB float64 `json:"memory_usage_mb"`

	// DiskUsagePercent is the disk usage percentage for the workspace volume.
	DiskUsagePercent float64 `json:"disk_usage_percent"`

	// OpenFiles is the number of open file descriptors.
	OpenFiles int `json:"open_files"`

	// ProcessID is the OpenClaw main process PID.
	ProcessID int `json:"process_id"`
}

// Status represents the current state of an OpenClaw instance.
type Status struct {
	// Running indicates whether the OpenClaw process is active.
	Running bool `json:"running"`

	// Version is the installed version.
	Version string `json:"version"`

	// Uptime is the duration since the process started.
	Uptime time.Duration `json:"uptime"`

	// GatewayPort is the port the gateway is listening on.
	GatewayPort int `json:"gateway_port"`

	// ActiveConnections is the number of active client connections.
	ActiveConnections int `json:"active_connections"`
}

// Collector defines the interface for collecting OpenClaw metrics from a VM.
type Collector interface {
	// Collect retrieves current metrics from the specified VM.
	Collect(vmName string) (*Metrics, error)

	// GetStatus retrieves the current OpenClaw status for the specified VM.
	GetStatus(vmName string) (*Status, error)

	// GetResourceUsage retrieves resource usage for the specified VM.
	GetResourceUsage(vmName string) (*ResourceUsage, error)
}

// DefaultMetrics returns a zero-valued Metrics with CollectedAt set to now.
func DefaultMetrics() *Metrics {
	return &Metrics{
		CollectedAt: time.Now(),
	}
}
