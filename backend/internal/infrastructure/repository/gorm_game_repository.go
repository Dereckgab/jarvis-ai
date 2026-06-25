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

// GormGameRepository implements the GameRepository interface using GORM.
type GormGameRepository struct {
	DB *gorm.DB
}

// NewGormGameRepository creates a new GormGameRepository.
func NewGormGameRepository(db *gorm.DB) repository.GameRepository {
	return &GormGameRepository{DB: db}
}

// CreateGame creates a new game in the database.
func (r *GormGameRepository) CreateGame(ctx context.Context, game *entity.Game) error {
	result := r.DB.WithContext(ctx).Create(game)
	if result.Error != nil {
		return fmt.Errorf("failed to create game: %w", result.Error)
	}
	return nil
}

// GetGameByID retrieves a game by its ID.
func (r *GormGameRepository) GetGameByID(ctx context.Context, id uuid.UUID) (*entity.Game, error) {
	var game entity.Game
	result := r.DB.WithContext(ctx).Preload("Requirements").First(&game, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("game not found with ID %s: %w", id.String(), result.Error)
		}
		return nil, fmt.Errorf("failed to get game by ID %s: %w", id.String(), result.Error)
	}
	return &game, nil
}

// GetGameByName retrieves a game by its name.
func (r *GormGameRepository) GetGameByName(ctx context.Context, name string) (*entity.Game, error) {
	var game entity.Game
	result := r.DB.WithContext(ctx).Preload("Requirements").First(&game, "name = ?", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("game not found with name %s: %w", name, result.Error)
		}
		return nil, fmt.Errorf("failed to get game by name %s: %w", name, result.Error)
	}
	return &game, nil
}

// SearchGames searches for games by a query string.
func (r *GormGameRepository) SearchGames(ctx context.Context, query string, limit, offset int) ([]*entity.Game, error) {
	var games []*entity.Game
	result := r.DB.WithContext(ctx).Preload("Requirements").Where("name LIKE ?", "%"+query+"%").Limit(limit).Offset(offset).Find(&games)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to search games: %w", result.Error)
	}
	return games, nil
}

// UpdateGame updates an existing game in the database.
func (r *GormGameRepository) UpdateGame(ctx context.Context, game *entity.Game) error {
	result := r.DB.WithContext(ctx).Save(game)
	if result.Error != nil {
		return fmt.Errorf("failed to update game: %w", result.Error)
	}
	return nil
}

// DeleteGame deletes a game by its ID.
func (r *GormGameRepository) DeleteGame(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&entity.Game{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete game with ID %s: %w", id.String(), result.Error)
	}
	return nil
}

// CreateGameRequirement creates a new game requirement in the database.
func (r *GormGameRepository) CreateGameRequirement(ctx context.Context, req *entity.GameRequirement) error {
	result := r.DB.WithContext(ctx).Create(req)
	if result.Error != nil {
		return fmt.Errorf("failed to create game requirement: %w", result.Error)
	}
	return nil
}

// GetGameRequirementsByGameID retrieves game requirements for a given game ID.
func (r *GormGameRepository) GetGameRequirementsByGameID(ctx context.Context, gameID uuid.UUID) ([]*entity.GameRequirement, error) {
	var reqs []*entity.GameRequirement
	result := r.DB.WithContext(ctx).Where("game_id = ?", gameID).Find(&reqs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get game requirements by game ID %s: %w", gameID.String(), result.Error)
	}
	return reqs, nil
}

// GetGameRequirementByID retrieves a game requirement by its ID.
func (r *GormGameRepository) GetGameRequirementByID(ctx context.Context, id uuid.UUID) (*entity.GameRequirement, error) {
	var req entity.GameRequirement
	result := r.DB.WithContext(ctx).First(&req, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("game requirement not found with ID %s: %w", id.String(), result.Error)
		}
		return nil, fmt.Errorf("failed to get game requirement by ID %s: %w", id.String(), result.Error)
	}
	return &req, nil
}

// UpdateGameRequirement updates an existing game requirement in the database.
func (r *GormGameRepository) UpdateGameRequirement(ctx context.Context, req *entity.GameRequirement) error {
	result := r.DB.WithContext(ctx).Save(req)
	if result.Error != nil {
		return fmt.Errorf("failed to update game requirement: %w", result.Error)
	}
	return nil
}

// DeleteGameRequirement deletes a game requirement by its ID.
func (r *GormGameRepository) DeleteGameRequirement(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&entity.GameRequirement{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete game requirement with ID %s: %w", id.String(), result.Error)
	}
	return nil
}
