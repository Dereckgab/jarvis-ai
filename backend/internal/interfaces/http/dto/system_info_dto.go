package dto

import (
	"time"

	"github.com/google/uuid"
)

// SystemInfoResponse represents the response body for system information.
type SystemInfoResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`

	CPUName      string  `json:"cpu_name"`
	CPUCores     int32   `json:"cpu_cores"`
	CPUThreads   int32   `json:"cpu_threads"`
	CPUFrequency float64 `json:"cpu_frequency"`
	CPUPercent   float64 `json:"cpu_percent"`

	TotalMemoryGB float64 `json:"total_memory_gb"`
	UsedMemoryGB  float64 `json:"used_memory_gb"`
	MemoryPercent float64 `json:"memory_percent"`

	TotalDiskGB float64 `json:"total_disk_gb"`
	UsedDiskGB  float64 `json:"used_disk_gb"`
	DiskPercent float64 `json:"disk_percent"`

	GPUName             string  `json:"gpu_name"`
	GPUTemperatureC     float64 `json:"gpu_temperature_c"`
	GPUUtilizationPercent float64 `json:"gpu_utilization_percent"`

	BytesSentSec uint64 `json:"bytes_sent_sec"`
	BytesRecvSec uint64 `json:"bytes_recv_sec"`

	OSPlatform string `json:"os_platform"`
	OSFamily   string `json:"os_family"`
	OSVersion  string `json:"os_version"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
