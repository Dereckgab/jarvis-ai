package sysmon

import (
	"context"
)

// SystemMonitor defines the interface for collecting system information.
type SystemMonitor interface {
	GetCPUInfo(ctx context.Context) (*CPUInfo, error)
	GetCPUUsage(ctx context.Context) (*CPUUsage, error)
	GetMemoryInfo(ctx context.Context) (*MemoryInfo, error)
	GetDiskInfo(ctx context.Context) (*DiskInfo, error)
	GetNetworkInfo(ctx context.Context) (*NetworkInfo, error)
	GetOSInfo(ctx context.Context) (*OSInfo, error)
	GetGPUInfo(ctx context.Context) (*GPUInfo, error) // Placeholder for GPU info
}

// CPUInfo represents static CPU information.
type CPUInfo struct {
	Name      string
	Cores     int32
	Threads   int32
	Frequency float64 // in MHz
}

// CPUUsage represents dynamic CPU usage.
type CPUUsage struct {
	Percent float64 // 0.0 - 100.0
}

// MemoryInfo represents memory usage.
type MemoryInfo struct {
	TotalGB float64
	UsedGB  float64
	Percent float64 // 0.0 - 100.0
}

// DiskInfo represents disk usage.
type DiskInfo struct {
	TotalGB float64
	UsedGB  float64
	Percent float64 // 0.0 - 100.0
}

// NetworkInfo represents network I/O.
type NetworkInfo struct {
	BytesSentSec uint64
	BytesRecvSec uint64
}

// OSInfo represents operating system information.
type OSInfo struct {
	Platform string
	Family   string
	Version  string
}

// GPUInfo represents GPU information (placeholder).
type GPUInfo struct {
	Name               string
	TemperatureC       float64
	UtilizationPercent float64
}

// Platform-specific implementations provide `NewSystemMonitor`.
