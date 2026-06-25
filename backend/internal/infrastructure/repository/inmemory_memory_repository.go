package repository

import (
	"context"
	"fmt"
	"sync"

	"jarvis/internal/domain/entity"

	"github.com/google/uuid"
)

// InMemoryMemoryRepository is a simple in-memory implementation of MemoryRepository used for local development and tests.
type InMemoryMemoryRepository struct {
	mu      sync.RWMutex
	storage map[uuid.UUID]*entity.Memory
}

// NewInMemoryMemoryRepository creates a new in-memory memory repository.
func NewInMemoryMemoryRepository() *InMemoryMemoryRepository {
	return &InMemoryMemoryRepository{
		storage: make(map[uuid.UUID]*entity.Memory),
	}
}

func (r *InMemoryMemoryRepository) CreateMemory(ctx context.Context, memory *entity.Memory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if memory == nil {
		return fmt.Errorf("memory is nil")
	}
	r.storage[memory.ID] = memory
	return nil
}

func (r *InMemoryMemoryRepository) GetMemoryByID(ctx context.Context, id uuid.UUID) (*entity.Memory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.storage[id]
	if !ok {
		return nil, fmt.Errorf("memory not found")
	}
	return m, nil
}

func (r *InMemoryMemoryRepository) GetMemoriesByUserID(ctx context.Context, userID uuid.UUID, memoryType string, limit, offset int) ([]*entity.Memory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*entity.Memory
	for _, m := range r.storage {
		if m.UserID == userID {
			if memoryType == "" || m.Type == memoryType {
				res = append(res, m)
			}
		}
	}
	// simple pagination
	if offset >= len(res) {
		return []*entity.Memory{}, nil
	}
	end := offset + limit
	if end > len(res) || limit <= 0 {
		end = len(res)
	}
	return res[offset:end], nil
}

func (r *InMemoryMemoryRepository) SearchMemories(ctx context.Context, userID uuid.UUID, queryEmbedding []float32, limit int) ([]*entity.Memory, error) {
	// Very simple implementation: return user's memories up to limit.
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*entity.Memory
	for _, m := range r.storage {
		if m.UserID == userID {
			res = append(res, m)
			if limit > 0 && len(res) >= limit {
				break
			}
		}
	}
	return res, nil
}

func (r *InMemoryMemoryRepository) DeleteMemory(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.storage, id)
	return nil
}
