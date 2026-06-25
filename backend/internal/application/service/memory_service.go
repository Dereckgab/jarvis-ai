package service

import (
	"context"
	"fmt"

	"jarvis/internal/domain/entity"
	"jarvis/internal/domain/repository"

	"github.com/google/uuid"
)

// MemoryService defines the application service interface for memory-related operations.
type MemoryService interface {
	SaveMemory(ctx context.Context, userID uuid.UUID, memoryType, content string) (*entity.Memory, error)
	SearchMemories(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*entity.Memory, error)
	GetMemories(ctx context.Context, userID uuid.UUID, memoryType string, limit, offset int) ([]*entity.Memory, error)
	DeleteMemory(ctx context.Context, memoryID uuid.UUID) error
}

// memoryService implements MemoryService.
type memoryService struct {
	memoryRepo repository.MemoryRepository
	aiService  AIService // To generate embeddings
}

// NewMemoryService creates a new MemoryService.
func NewMemoryService(memoryRepo repository.MemoryRepository, aiService AIService) MemoryService {
	return &memoryService{
		memoryRepo: memoryRepo,
		aiService:  aiService,
	}
}

// SaveMemory processes content, generates embeddings, and saves it as a memory.
func (s *memoryService) SaveMemory(ctx context.Context, userID uuid.UUID, memoryType, content string) (*entity.Memory, error) {
	// TODO: Generate embedding using aiService.GetEmbedding(ctx, content)
	// For now, using a placeholder empty slice.
	var embedding []float32 // Placeholder

	memory := entity.NewMemory(userID, memoryType, content, embedding)

	if err := s.memoryRepo.CreateMemory(ctx, memory); err != nil {
		return nil, fmt.Errorf("failed to save memory: %w", err)
	}

	return memory, nil
}

// SearchMemories performs a semantic search using AI embeddings.
func (s *memoryService) SearchMemories(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*entity.Memory, error) {
	// TODO: Generate embedding for the query using aiService.GetEmbedding(ctx, query)
	var queryEmbedding []float32 // Placeholder

	memories, err := s.memoryRepo.SearchMemories(ctx, userID, queryEmbedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories: %w", err)
	}

	return memories, nil
}

// GetMemories retrieves memories for a user with optional type filtering and pagination.
func (s *memoryService) GetMemories(ctx context.Context, userID uuid.UUID, memoryType string, limit, offset int) ([]*entity.Memory, error) {
	memories, err := s.memoryRepo.GetMemoriesByUserID(ctx, userID, memoryType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories by user ID: %w", err)
	}
	return memories, nil
}

// DeleteMemory deletes a memory by its ID.
func (s *memoryService) DeleteMemory(ctx context.Context, memoryID uuid.UUID) error {
	if err := s.memoryRepo.DeleteMemory(ctx, memoryID); err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}
	return nil
}
