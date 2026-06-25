//go:build !windows && !linux

package sysmon

import (
	"context"
	"fmt"
)

// DefaultSystemMonitor implements SystemMonitor for unsupported platforms.
type DefaultSystemMonitor struct{}

// NewSystemMonitor creates a new DefaultSystemMonitor.
func NewSystemMonitor() (SystemMonitor, error) {
	return &DefaultSystemMonitor{}, nil
}

// GetCPUInfo retrieves static CPU information for unsupported platforms.
func (d *DefaultSystemMonitor) GetCPUInfo(ctx context.Context) (*CPUInfo, error) {
	return &CPUInfo{Name: "N/A", Cores: 0, Threads: 0, Frequency: 0.0}, nil
}

// GetCPUUsage retrieves dynamic CPU usage for unsupported platforms.
func (d *DefaultSystemMonitor) GetCPUUsage(ctx context.Context) (*CPUUsage, error) {
	return &CPUUsage{Percent: 0.0}, nil
}

// GetMemoryInfo retrieves memory usage for unsupported platforms.
func (d *DefaultSystemMonitor) GetMemoryInfo(ctx context.Context) (*MemoryInfo, error) {
	return &MemoryInfo{TotalGB: 0.0, UsedGB: 0.0, Percent: 0.0}, nil
}

// GetDiskInfo retrieves disk usage for unsupported platforms.
func (d *DefaultSystemMonitor) GetDiskInfo(ctx context.Context) (*DiskInfo, error) {
	return &DiskInfo{TotalGB: 0.0, UsedGB: 0.0, Percent: 0.0}, nil
}

// GetNetworkInfo retrieves network I/O for unsupported platforms.
func (d *DefaultSystemMonitor) GetNetworkInfo(ctx context.Context) (*NetworkInfo, error) {
	return &NetworkInfo{BytesSentSec: 0, BytesRecvSec: 0}, nil
}

// GetOSInfo retrieves operating system information for unsupported platforms.
func (d *DefaultSystemMonitor) GetOSInfo(ctx context.Context) (*OSInfo, error) {
	return &OSInfo{Platform: "N/A", Family: "N/A", Version: "N/A"}, nil
}

// GetGPUInfo retrieves GPU information for unsupported platforms.
func (d *DefaultSystemMonitor) GetGPUInfo(ctx context.Context) (*GPUInfo, error) {
	return &GPUInfo{Name: "N/A", TemperatureC: 0.0, UtilizationPercent: 0.0}, nil
}
