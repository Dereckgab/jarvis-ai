//go:build windows

package sysmon

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/host"
)

// WindowsSystemMonitor implements SystemMonitor for Windows.
type WindowsSystemMonitor struct{}

// NewSystemMonitor creates a new WindowsSystemMonitor.
func NewSystemMonitor() (SystemMonitor, error) {
	return &WindowsSystemMonitor{}, nil
}

// GetCPUInfo retrieves static CPU information for Windows.
func (w *WindowsSystemMonitor) GetCPUInfo(ctx context.Context) (*CPUInfo, error) {
	cpuInfos, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}
	if len(cpuInfos) == 0 {
		return nil, fmt.Errorf("no CPU info found")
	}

	// Assuming all CPUs are similar, take the first one
	info := cpuInfos[0]

	return &CPUInfo{
		Name:      info.ModelName,
		Cores:     info.Cores,
		Threads:   info.Cores * 2, // Heuristic for threads on Windows
		Frequency: info.Mhz,
	}, nil
}

// GetCPUUsage retrieves dynamic CPU usage for Windows.
func (w *WindowsSystemMonitor) GetCPUUsage(ctx context.Context) (*CPUUsage, error) {
	// Get CPU percentages for all CPUs, then average
	percents, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %w", err)
	}
	if len(percents) == 0 {
		return nil, fmt.Errorf("no CPU usage data found")
	}

	return &CPUUsage{
		Percent: percents[0],
	}, nil
}

// GetMemoryInfo retrieves memory usage for Windows.
func (w *WindowsSystemMonitor) GetMemoryInfo(ctx context.Context) (*MemoryInfo, error) {
	vmem, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual memory info: %w", err)
	}

	return &MemoryInfo{
		TotalGB: float64(vmem.Total) / (1024 * 1024 * 1024),
		UsedGB:  float64(vmem.Used) / (1024 * 1024 * 1024),
		Percent: vmem.UsedPercent,
	}, nil
}

// GetDiskInfo retrieves disk usage for Windows.
func (w *WindowsSystemMonitor) GetDiskInfo(ctx context.Context) (*DiskInfo, error) {
	// Get disk usage for the root partition (C: on Windows)
	usage, err := disk.UsageWithContext(ctx, "C:\\")
	if err != nil {
		return nil, fmt.Errorf("failed to get disk usage for C:\\: %w", err)
	}

	return &DiskInfo{
		TotalGB: float64(usage.Total) / (1024 * 1024 * 1024),
		UsedGB:  float64(usage.Used) / (1024 * 1024 * 1024),
		Percent: usage.UsedPercent,
	}, nil
}

// GetNetworkInfo retrieves network I/O for Windows.
func (w *WindowsSystemMonitor) GetNetworkInfo(ctx context.Context) (*NetworkInfo, error) {
	// Get network counters, then calculate delta over a short period
	netIOs, err := net.IOCountersWithContext(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get network IO counters: %w", err)
	}
	if len(netIOs) == 0 {
		return nil, fmt.Errorf("no network IO data found")
	}

	// Take a snapshot, wait, then take another to calculate per-second rates
	initialBytesSent := netIOs[0].BytesSent
	initialBytesRecv := netIOs[0].BytesRecv

	time.Sleep(time.Second)

	netIOs, err = net.IOCountersWithContext(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get network IO counters after delay: %w", err)
	}

	currentBytesSent := netIOs[0].BytesSent
	currentBytesRecv := netIOs[0].BytesRecv

	return &NetworkInfo{
		BytesSentSec: currentBytesSent - initialBytesSent,
		BytesRecvSec: currentBytesRecv - initialBytesRecv,
	}, nil
}

// GetOSInfo retrieves operating system information for Windows.
func (w *WindowsSystemMonitor) GetOSInfo(ctx context.Context) (*OSInfo, error) {
	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %w", err)
	}

	return &OSInfo{
		Platform: hostInfo.Platform,
		Family:   hostInfo.PlatformFamily,
		Version:  hostInfo.PlatformVersion,
	}, nil
}

// GetGPUInfo retrieves GPU information for Windows (placeholder).
func (w *WindowsSystemMonitor) GetGPUInfo(ctx context.Context) (*GPUInfo, error) {
	// gopsutil does not directly provide GPU info. This would require external libraries or WMI queries.
	// For now, return dummy data.
	return &GPUInfo{
		Name:             "N/A (Windows GPU)",
		TemperatureC:     0.0,
		UtilizationPercent: 0.0,
	}, nil
}
