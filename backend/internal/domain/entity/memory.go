package entity

import (
	"time"

	"github.com/google/uuid"
)

// Memory represents a piece of information stored in the system's memory, potentially for AI context.
type Memory struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"` // e.g., "short-term", "long-term", "semantic", "contextual"
	Content   string    `json:"content"`
	Embedding []float32 `json:"embedding,omitempty"` // Vector embedding for semantic search
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewMemory creates a new Memory entity.
func NewMemory(userID uuid.UUID, memoryType, content string, embedding []float32) *Memory {
	return &Memory{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      memoryType,
		Content:   content,
		Embedding: embedding,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
