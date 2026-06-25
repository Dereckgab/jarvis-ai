package repository

import (
	"context"

	"jarvis/internal/domain/entity"

	"github.com/google/uuid"
)

// MemoryRepository defines the interface for managing Memory entities.
type MemoryRepository interface {
	CreateMemory(ctx context.Context, memory *entity.Memory) error
	GetMemoryByID(ctx context.Context, id uuid.UUID) (*entity.Memory, error)
	GetMemoriesByUserID(ctx context.Context, userID uuid.UUID, memoryType string, limit, offset int) ([]*entity.Memory, error)
	SearchMemories(ctx context.Context, userID uuid.UUID, queryEmbedding []float32, limit int) ([]*entity.Memory, error)
	DeleteMemory(ctx context.Context, id uuid.UUID) error
}
