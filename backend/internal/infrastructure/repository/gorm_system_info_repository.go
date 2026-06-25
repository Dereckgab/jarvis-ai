package repository

import (
	"context"
	"errors"
	"fmt"

	"jarvis/internal/domain/entity"
	"jarvis/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormSystemInfoRepository implements the SystemInfoRepository interface using GORM.
type GormSystemInfoRepository struct {
	DB *gorm.DB
}

// NewGormSystemInfoRepository creates a new GormSystemInfoRepository.
func NewGormSystemInfoRepository(db *gorm.DB) repository.SystemInfoRepository {
	return &GormSystemInfoRepository{DB: db}
}

// CreateSystemInfo creates a new system information record in the database.
func (r *GormSystemInfoRepository) CreateSystemInfo(ctx context.Context, info *entity.SystemInfo) error {
	result := r.DB.WithContext(ctx).Create(info)
	if result.Error != nil {
		return fmt.Errorf("failed to create system info: %w", result.Error)
	}
	return nil
}

// GetSystemInfoByID retrieves a system information record by its ID.
func (r *GormSystemInfoRepository) GetSystemInfoByID(ctx context.Context, id uuid.UUID) (*entity.SystemInfo, error) {
	var info entity.SystemInfo
	result := r.DB.WithContext(ctx).First(&info, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("system info not found with ID %s: %w", id.String(), result.Error)
		}
		return nil, fmt.Errorf("failed to get system info by ID %s: %w", id.String(), result.Error)
	}
	return &info, nil
}

// GetSystemInfoByUserID retrieves system information records for a given user ID with pagination.
func (r *GormSystemInfoRepository) GetSystemInfoByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.SystemInfo, error) {
	var infos []*entity.SystemInfo
	result := r.DB.WithContext(ctx).Where("user_id = ?", userID).Order("timestamp DESC").Limit(limit).Offset(offset).Find(&infos)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get system info by user ID %s: %w", userID.String(), result.Error)
	}
	return infos, nil
}

// GetLatestSystemInfoByUserID retrieves the latest system information record for a given user ID.
func (r *GormSystemInfoRepository) GetLatestSystemInfoByUserID(ctx context.Context, userID uuid.UUID) (*entity.SystemInfo, error) {
	var info entity.SystemInfo
	result := r.DB.WithContext(ctx).Where("user_id = ?", userID).Order("timestamp DESC").First(&info)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no system info found for user ID %s: %w", userID.String(), result.Error)
		}
		return nil, fmt.Errorf("failed to get latest system info by user ID %s: %w", userID.String(), result.Error)
	}
	return &info, nil
}

// DeleteSystemInfo deletes a system information record by its ID.
func (r *GormSystemInfoRepository) DeleteSystemInfo(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&entity.SystemInfo{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete system info with ID %s: %w", id.String(), result.Error)
	}
	return nil
}
