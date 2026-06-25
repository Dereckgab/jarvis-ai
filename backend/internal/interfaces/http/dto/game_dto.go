package dto

import (
	"time"

	"github.com/google/uuid"
)

// GameRequirementDTO represents a game requirement in DTO format.
type GameRequirementDTO struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	CPUModel  string    `json:"cpu_model,omitempty"`
	CPUSpeed  string    `json:"cpu_speed,omitempty"`
	GPUModel  string    `json:"gpu_model,omitempty"`
	GPUMemory string    `json:"gpu_memory,omitempty"`
	RAMGB     float64   `json:"ram_gb,omitempty"`
	StorageGB float64   `json:"storage_gb,omitempty"`
	OS        string    `json:"os,omitempty"`
}

// GameResponse represents the response body for game details.
type GameResponse struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	ReleaseDate *time.Time           `json:"release_date,omitempty"`
	Developer   string               `json:"developer,omitempty"`
	Publisher   string               `json:"publisher,omitempty"`
	Genre       string               `json:"genre,omitempty"`
	ImageURL    string               `json:"image_url,omitempty"`	
	Requirements []GameRequirementDTO `json:"requirements,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// CreateGameRequest represents the request body for creating a new game.
type CreateGameRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Developer   string `json:"developer,omitempty"`
	Publisher   string `json:"publisher,omitempty"`
	Genre       string `json:"genre,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

// UpdateGameRequest represents the request body for updating an existing game.
type UpdateGameRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Developer   string `json:"developer,omitempty"`
	Publisher   string `json:"publisher,omitempty"`
	Genre       string `json:"genre,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

// AddGameRequirementRequest represents the request body for adding a game requirement.
type AddGameRequirementRequest struct {
	Type      string  `json:"type" validate:"required,oneof=minimum recommended"`
	CPUModel  string  `json:"cpu_model,omitempty"`
	CPUSpeed  string  `json:"cpu_speed,omitempty"`
	GPUModel  string  `json:"gpu_model,omitempty"`
	GPUMemory string  `json:"gpu_memory,omitempty"`
	RAMGB     float64 `json:"ram_gb,omitempty"`
	StorageGB float64 `json:"storage_gb,omitempty"`
	OS        string  `json:"os,omitempty"`
}

// UpdateGameRequirementRequest represents the request body for updating a game requirement.
type UpdateGameRequirementRequest struct {
	Type      string  `json:"type" validate:"required,oneof=minimum recommended"`
	CPUModel  string  `json:"cpu_model,omitempty"`
	CPUSpeed  string  `json:"cpu_speed,omitempty"`
	GPUModel  string  `json:"gpu_model,omitempty"`	
	GPUMemory string  `json:"gpu_memory,omitempty"`
	RAMGB     float64 `json:"ram_gb,omitempty"`
	StorageGB float64 `json:"storage_gb,omitempty"`
	OS        string  `json:"os,omitempty"`
}
