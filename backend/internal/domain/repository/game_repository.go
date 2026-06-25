package repository

import (
	"context"

	"jarvis/internal/domain/entity"

	"github.com/google/uuid"
)

// GameRepository defines the interface for managing Game entities.
type GameRepository interface {
	CreateGame(ctx context.Context, game *entity.Game) error
	GetGameByID(ctx context.Context, id uuid.UUID) (*entity.Game, error)
	GetGameByName(ctx context.Context, name string) (*entity.Game, error)
	SearchGames(ctx context.Context, query string, limit, offset int) ([]*entity.Game, error)
	UpdateGame(ctx context.Context, game *entity.Game) error
	DeleteGame(ctx context.Context, id uuid.UUID) error

	CreateGameRequirement(ctx context.Context, req *entity.GameRequirement) error
	GetGameRequirementsByGameID(ctx context.Context, gameID uuid.UUID) ([]*entity.GameRequirement, error)
	UpdateGameRequirement(ctx context.Context, req *entity.GameRequirement) error
	DeleteGameRequirement(ctx context.Context, id uuid.UUID) error

	GetGameRequirementByID(ctx context.Context, id uuid.UUID) (*entity.GameRequirement, error)
}
