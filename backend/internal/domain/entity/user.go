package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the core user entity in the domain.
type User struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Username  string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-" validate:"required,min=8"` // Stored hashed
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// NewUser creates a new User entity with a generated ID and timestamps.
func NewUser(username, email, password string) (*User, error) {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  password, // Password should be hashed before saving
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
