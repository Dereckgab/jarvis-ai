package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SystemInfo represents a snapshot of the system's hardware and performance.
type SystemInfo struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;index" json:"user_id"`
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`

	// CPU Information
	CPUName      string  `gorm:"type:varchar(255)" json:"cpu_name"`
	CPUCores     int32   `json:"cpu_cores"`
	CPUThreads   int32   `json:"cpu_threads"`
	CPUFrequency float64 `json:"cpu_frequency"`
	CPUPercent   float64 `json:"cpu_percent"`

	// Memory Information
	TotalMemoryGB float64 `json:"total_memory_gb"`
	UsedMemoryGB  float64 `json:"used_memory_gb"`
	MemoryPercent float64 `json:"memory_percent"`

	// Disk Information
	TotalDiskGB float64 `json:"total_disk_gb"`
	UsedDiskGB  float64 `json:"used_disk_gb"`
	DiskPercent float64 `json:"disk_percent"`

	// GPU Information (simplified for now, can be expanded)
	GPUName               string  `gorm:"type:varchar(255)" json:"gpu_name"`
	GPUTemperatureC       float64 `json:"gpu_temperature_c"`
	GPUUtilizationPercent float64 `json:"gpu_utilization_percent"`

	// Network Information (simplified)
	BytesSentSec uint64 `json:"bytes_sent_sec"`
	BytesRecvSec uint64 `json:"bytes_recv_sec"`

	// OS Information
	OSPlatform string `gorm:"type:varchar(100)" json:"os_platform"`
	OSFamily   string `gorm:"type:varchar(100)" json:"os_family"`
	OSVersion  string `gorm:"type:varchar(100)" json:"os_version"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// NewSystemInfo creates a new SystemInfo entity with a generated ID and timestamps.
func NewSystemInfo(userID uuid.UUID) *SystemInfo {
	return &SystemInfo{
		ID:        uuid.New(),
		UserID:    userID,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
