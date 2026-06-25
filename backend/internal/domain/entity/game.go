package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Game represents a game entity with basic information.
type Game struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"name" validate:"required"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	ReleaseDate *time.Time     `json:"release_date,omitempty"`
	Developer   string         `gorm:"type:varchar(255)" json:"developer,omitempty"`
	Publisher   string         `gorm:"type:varchar(255)" json:"publisher,omitempty"`
	Genre       string         `gorm:"type:varchar(255)" json:"genre,omitempty"`
	ImageURL    string         `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	ExternalID  string         `gorm:"type:varchar(255);uniqueIndex" json:"external_id,omitempty"` // e.g., Steam App ID, RAWG ID
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Requirements []GameRequirement `gorm:"foreignKey:GameID" json:"requirements,omitempty"`
}

// NewGame creates a new Game entity.
func NewGame(name string) *Game {
	return &Game{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
