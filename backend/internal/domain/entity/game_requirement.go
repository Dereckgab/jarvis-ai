package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GameRequirement represents the hardware requirements for a specific game.
type GameRequirement struct {
	ID     uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	GameID uuid.UUID `gorm:"type:char(36);not null;index" json:"game_id"`
	Type   string    `gorm:"type:varchar(50);not null" json:"type"` // e.g., "minimum", "recommended"

	// CPU Requirements
	CPUModel string `gorm:"type:varchar(255)" json:"cpu_model,omitempty"`
	CPUSpeed string `gorm:"type:varchar(100)" json:"cpu_speed,omitempty"` // e.g., "2.5 GHz"

	// GPU Requirements
	GPUModel  string `gorm:"type:varchar(255)" json:"gpu_model,omitempty"`
	GPUMemory string `gorm:"type:varchar(100)" json:"gpu_memory,omitempty"` // e.g., "4 GB"

	// RAM Requirements
	RAMGB float64 `json:"ram_gb,omitempty"`

	// Storage Requirements
	StorageGB float64 `json:"storage_gb,omitempty"`

	// OS Requirements
	OS string `gorm:"type:varchar(255)" json:"os,omitempty"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// NewGameRequirement creates a new GameRequirement entity.
func NewGameRequirement(gameID uuid.UUID, reqType string) *GameRequirement {
	return &GameRequirement{
		ID:        uuid.New(),
		GameID:    gameID,
		Type:      reqType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
