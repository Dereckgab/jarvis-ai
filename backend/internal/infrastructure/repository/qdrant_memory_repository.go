//go:build qdrant
// +build qdrant

package repository

import (
	"context"
	"fmt"
	"time"

	"jarvis/config"
	"jarvis/internal/domain/entity"
	"jarvis/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

// QdrantMemoryRepository implements the MemoryRepository interface using Qdrant.
type QdrantMemoryRepository struct {
	client *qdrant.QdrantClient
	cfg    *config.QdrantConfig
}

// NewQdrantMemoryRepository creates a new QdrantMemoryRepository.
func NewQdrantMemoryRepository(client *qdrant.QdrantClient, cfg *config.QdrantConfig) repository.MemoryRepository {
	return &QdrantMemoryRepository{
		client: client,
		cfg:    cfg,
	}
}

// CreateMemory creates a new memory entry in Qdrant.
func (r *QdrantMemoryRepository) CreateMemory(ctx context.Context, memory *entity.Memory) error {
	pointsClient := qdrant.NewPointsClient(r.client.Target())

	payload := map[string]*qdrant.Value{
		"user_id":    qdrant.NewValue(memory.UserID.String()),
		"type":       qdrant.NewValue(memory.Type),
		"content":    qdrant.NewValue(memory.Content),
		"created_at": qdrant.NewValue(memory.CreatedAt.Format(time.RFC3339)),
		"updated_at": qdrant.NewValue(memory.UpdatedAt.Format(time.RFC3339)),
	}

	_, err := pointsClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: r.cfg.CollectionName,
		Wait:           true,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewUuid(memory.ID.String()),
				Vectors: &qdrant.Vectors{Vectors: &qdrant.Vectors_Vector{Vector: memory.Embedding}},
				Payload: payload,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create memory in Qdrant: %w", err)
	}
	return nil
}

// GetMemoryByID retrieves a memory entry by its ID from Qdrant.
func (r *QdrantMemoryRepository) GetMemoryByID(ctx context.Context, id uuid.UUID) (*entity.Memory, error) {
	pointsClient := qdrant.NewPointsClient(r.client.Target())

	retrieveResponse, err := pointsClient.Retrieve(ctx, &qdrant.RetrievePoints{
		CollectionName: r.cfg.CollectionName,
		Ids:            []*qdrant.PointId{qdrant.NewUuid(id.String())},
		WithVectors:    true,
		WithPayload:    qdrant.NewWithPayloadSelector(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve memory from Qdrant: %w", err)
	}

	if len(retrieveResponse.GetResult()) == 0 {
		return nil, fmt.Errorf("memory not found with ID %s", id.String())
	}

	point := retrieveResponse.GetResult()[0]
	return r.mapQdrantPointToMemory(point)
}

// GetMemoriesByUserID retrieves memory entries for a given user ID and type from Qdrant.
func (r *QdrantMemoryRepository) GetMemoriesByUserID(ctx context.Context, userID uuid.UUID, memoryType string, limit, offset int) ([]*entity.Memory, error) {
	pointsClient := qdrant.NewPointsClient(r.client.Target())

	filter := &qdrant.Filter{
		Must: []*qdrant.Condition{
			{
				FilterCondition: &qdrant.Condition_Field{Field: &qdrant.FieldCondition{
					Key:   "user_id",
					Range: &qdrant.Range{Eq: qdrant.NewValue(userID.String())},
				}},
			},
		},
	}

	if memoryType != "" {
		filter.Must = append(filter.Must, &qdrant.Condition{
			FilterCondition: &qdrant.Condition_Field{Field: &qdrant.FieldCondition{
				Key:   "type",
				Range: &qdrant.Range{Eq: qdrant.NewValue(memoryType)},
			}},
		})
	}

	scrollResponse, err := pointsClient.Scroll(ctx, &qdrant.ScrollPoints{
		CollectionName: r.cfg.CollectionName,
		Filter:         filter,
		Limit:          qdrant.NewUint64(uint64(limit)),
		Offset:         qdrant.NewUint64(uint64(offset)),
		WithVectors:    true,
		WithPayload:    qdrant.NewWithPayloadSelector(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scroll memories from Qdrant: %w", err)
	}

	var memories []*entity.Memory
	for _, point := range scrollResponse.GetResult() {
		memory, err := r.mapQdrantPointToMemory(point)
		if err != nil {
			return nil, fmt.Errorf("failed to map Qdrant point to memory: %w", err)
		}
		memories = append(memories, memory)
	}

	return memories, nil
}

// SearchMemories performs a semantic search for memories based on a query embedding.
func (r *QdrantMemoryRepository) SearchMemories(ctx context.Context, userID uuid.UUID, queryEmbedding []float32, limit int) ([]*entity.Memory, error) {
	pointsClient := qdrant.NewPointsClient(r.client.Target())

	filter := &qdrant.Filter{
		Must: []*qdrant.Condition{
			{
				FilterCondition: &qdrant.Condition_Field{Field: &qdrant.FieldCondition{
					Key:   "user_id",
					Range: &qdrant.Range{Eq: qdrant.NewValue(userID.String())},
				}},
			},
		},
	}

	searchResponse, err := pointsClient.Search(ctx, &qdrant.SearchPoints{
		CollectionName: r.cfg.CollectionName,
		Vector:         queryEmbedding,
		Filter:         filter,
		Limit:          uint64(limit),
		WithVectors:    true,
		WithPayload:    qdrant.NewWithPayloadSelector(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search memories in Qdrant: %w", err)
	}

	var memories []*entity.Memory
	for _, scoredPoint := range searchResponse.GetResult() {
		memory, err := r.mapQdrantPointToMemory(scoredPoint.GetPoint())
		if err != nil {
			return nil, fmt.Errorf("failed to map Qdrant scored point to memory: %w", err)
		}
		memories = append(memories, memory)
	}

	return memories, nil
}

// DeleteMemory deletes a memory entry by its ID from Qdrant.
func (r *QdrantMemoryRepository) DeleteMemory(ctx context.Context, id uuid.UUID) error {
	pointsClient := qdrant.NewPointsClient(r.client.Target())

	_, err := pointsClient.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: r.cfg.CollectionName,
		PointsSelector: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Points{
				Points: &qdrant.PointIdsList{
					Ids: []*qdrant.PointId{qdrant.NewUuid(id.String())},
				},
			},
		},
		Wait: true,
	})
	if err != nil {
		return fmt.Errorf("failed to delete memory from Qdrant: %w", err)
	}
	return nil
}

func (r *QdrantMemoryRepository) mapQdrantPointToMemory(point *qdrant.PointStruct) (*entity.Memory, error) {
	id, err := uuid.Parse(point.GetId().GetUuid())
	if err != nil {
		return nil, fmt.Errorf("invalid UUID in Qdrant point ID: %w", err)
	}

	payload := point.GetPayload()

	userIDStr, ok := payload["user_id"].GetStringValue()
	if !ok {
		return nil, fmt.Errorf("user_id not found or not a string in Qdrant payload")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id UUID in Qdrant payload: %w", err)
	}

	memoryType, ok := payload["type"].GetStringValue()
	if !ok {
		return nil, fmt.Errorf("type not found or not a string in Qdrant payload")
	}

	content, ok := payload["content"].GetStringValue()
	if !ok {
		return nil, fmt.Errorf("content not found or not a string in Qdrant payload")
	}

	createdAtStr, ok := payload["created_at"].GetStringValue()
	if !ok {
		return nil, fmt.Errorf("created_at not found or not a string in Qdrant payload")
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at format in Qdrant payload: %w", err)
	}

	updatedAtStr, ok := payload["updated_at"].GetStringValue()
	if !ok {
		return nil, fmt.Errorf("updated_at not found or not a string in Qdrant payload")
	}
	updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at format in Qdrant payload: %w", err)
	}

	var embedding []float32
	if vectors := point.GetVectors(); vectors != nil {
		if vec := vectors.GetVector(); vec != nil {
			embedding = vec
		}
	}

	return &entity.Memory{
		ID:        id,
		UserID:    userID,
		Type:      memoryType,
		Content:   content,
		Embedding: embedding,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
