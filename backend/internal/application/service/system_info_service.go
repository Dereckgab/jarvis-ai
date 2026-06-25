package service

import (
	"context"
	"fmt"

	"jarvis/internal/domain/entity"
	"jarvis/internal/domain/repository"
	"jarvis/pkg/logger"
	"jarvis/pkg/sysmon"

	"github.com/google/uuid"
)

// SystemInfoService defines the application service interface for system information management.
type SystemInfoService interface {
	CollectAndSaveSystemInfo(ctx context.Context, userID uuid.UUID) (*entity.SystemInfo, error)
	GetLatestSystemInfo(ctx context.Context, userID uuid.UUID) (*entity.SystemInfo, error)
	GetSystemInfoHistory(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.SystemInfo, error)
}

// systemInfoService implements SystemInfoService.
type systemInfoService struct {
	systemInfoRepo repository.SystemInfoRepository
	sysMonitor     sysmon.SystemMonitor
}

// NewSystemInfoService creates a new SystemInfoService.
func NewSystemInfoService(repo repository.SystemInfoRepository) SystemInfoService {
	monitor, err := sysmon.NewSystemMonitor()
	if err != nil {
		logger.Warn("Failed to initialize system monitor", "error", err)
	}
	return &systemInfoService{
		systemInfoRepo: repo,
		sysMonitor:     monitor,
	}
}

// CollectAndSaveSystemInfo collects current system metrics and saves them to the database.
func (s *systemInfoService) CollectAndSaveSystemInfo(ctx context.Context, userID uuid.UUID) (*entity.SystemInfo, error) {
	if s.sysMonitor == nil {
		return nil, fmt.Errorf("system monitor not initialized")
	}

	// Collect CPU Info
	cpuInfo, err := s.sysMonitor.GetCPUInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}
	cpuUsage, err := s.sysMonitor.GetCPUUsage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %w", err)
	}

	// Collect Memory Info
	memInfo, err := s.sysMonitor.GetMemoryInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}

	// Collect Disk Info
	diskInfo, err := s.sysMonitor.GetDiskInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk info: %w", err)
	}

	// Collect Network Info
	netInfo, err := s.sysMonitor.GetNetworkInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get network info: %w", err)
	}

	// Collect OS Info
	osInfo, err := s.sysMonitor.GetOSInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get OS info: %w", err)
	}

	// Collect GPU Info (placeholder for now)
	gpuInfo, err := s.sysMonitor.GetGPUInfo(ctx)
	if err != nil {
		logger.Warn("GPU info unavailable", "error", err)
		gpuInfo = &sysmon.GPUInfo{Name: "N/A", TemperatureC: 0.0, UtilizationPercent: 0.0}
	}

	// Create SystemInfo entity
	systemInfo := entity.NewSystemInfo(userID)
	systemInfo.CPUName = cpuInfo.Name
	systemInfo.CPUCores = cpuInfo.Cores
	systemInfo.CPUThreads = cpuInfo.Threads
	systemInfo.CPUFrequency = cpuInfo.Frequency
	systemInfo.CPUPercent = cpuUsage.Percent

	systemInfo.TotalMemoryGB = memInfo.TotalGB
	systemInfo.UsedMemoryGB = memInfo.UsedGB
	systemInfo.MemoryPercent = memInfo.Percent

	systemInfo.TotalDiskGB = diskInfo.TotalGB
	systemInfo.UsedDiskGB = diskInfo.UsedGB
	systemInfo.DiskPercent = diskInfo.Percent

	systemInfo.GPUName = gpuInfo.Name
	systemInfo.GPUTemperatureC = gpuInfo.TemperatureC
	systemInfo.GPUUtilizationPercent = gpuInfo.UtilizationPercent

	systemInfo.BytesSentSec = netInfo.BytesSentSec
	systemInfo.BytesRecvSec = netInfo.BytesRecvSec

	systemInfo.OSPlatform = osInfo.Platform
	systemInfo.OSFamily = osInfo.Family
	systemInfo.OSVersion = osInfo.Version

	// Save to repository
	if err := s.systemInfoRepo.CreateSystemInfo(ctx, systemInfo); err != nil {
		return nil, fmt.Errorf("failed to save system info: %w", err)
	}

	return systemInfo, nil
}

// GetLatestSystemInfo retrieves the most recent system information for a user.
func (s *systemInfoService) GetLatestSystemInfo(ctx context.Context, userID uuid.UUID) (*entity.SystemInfo, error) {
	info, err := s.systemInfoRepo.GetLatestSystemInfoByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest system info: %w", err)
	}
	return info, nil
}

// GetSystemInfoHistory retrieves a paginated history of system information for a user.
func (s *systemInfoService) GetSystemInfoHistory(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.SystemInfo, error) {
	infos, err := s.systemInfoRepo.GetSystemInfoByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get system info history: %w", err)
	}
	return infos, nil
}
