package repository

import (
	"context"

	"jarvis/internal/domain/entity"

	"github.com/google/uuid"
)

// SystemInfoRepository defines the interface for managing SystemInfo entities.
type SystemInfoRepository interface {
	CreateSystemInfo(ctx context.Context, info *entity.SystemInfo) error
	GetSystemInfoByID(ctx context.Context, id uuid.UUID) (*entity.SystemInfo, error)
	GetSystemInfoByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.SystemInfo, error)
	GetLatestSystemInfoByUserID(ctx context.Context, userID uuid.UUID) (*entity.SystemInfo, error)
	DeleteSystemInfo(ctx context.Context, id uuid.UUID) error
}
