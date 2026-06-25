package service

import (
	"context"
	"fmt"

	"jarvis/internal/domain/entity"
	"jarvis/internal/domain/repository"

	"github.com/google/uuid"
)

// GameService defines the application service interface for game-related operations.
type GameService interface {
	CreateGame(ctx context.Context, name, description, developer, publisher, genre, imageURL string) (*entity.Game, error)
	GetGame(ctx context.Context, gameID uuid.UUID) (*entity.Game, error)
	SearchGames(ctx context.Context, query string, limit, offset int) ([]*entity.Game, error)
	UpdateGame(ctx context.Context, gameID uuid.UUID, name, description, developer, publisher, genre, imageURL string) (*entity.Game, error)
	DeleteGame(ctx context.Context, gameID uuid.UUID) error

	AddGameRequirement(ctx context.Context, gameID uuid.UUID, reqType, cpuModel, cpuSpeed, gpuModel, gpuMemory, os string, ramGB, storageGB float64) (*entity.GameRequirement, error)
	GetGameRequirements(ctx context.Context, gameID uuid.UUID) ([]*entity.GameRequirement, error)
	UpdateGameRequirement(ctx context.Context, reqID uuid.UUID, reqType, cpuModel, cpuSpeed, gpuModel, gpuMemory, os string, ramGB, storageGB float64) (*entity.GameRequirement, error)
	DeleteGameRequirement(ctx context.Context, reqID uuid.UUID) error
}

// gameService implements GameService.
type gameService struct {
	gameRepo repository.GameRepository
}

// NewGameService creates a new GameService.
func NewGameService(gameRepo repository.GameRepository) GameService {
	return &gameService{gameRepo: gameRepo}
}

// CreateGame handles the creation of a new game.
func (s *gameService) CreateGame(ctx context.Context, name, description, developer, publisher, genre, imageURL string) (*entity.Game, error) {
	game := entity.NewGame(name)
	game.Description = description
	game.Developer = developer
	game.Publisher = publisher
	game.Genre = genre
	game.ImageURL = imageURL

	if err := s.gameRepo.CreateGame(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}
	return game, nil
}

// GetGame retrieves a game by its ID.
func (s *gameService) GetGame(ctx context.Context, gameID uuid.UUID) (*entity.Game, error) {
	game, err := s.gameRepo.GetGameByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("game not found: %w", err)
	}
	return game, nil
}

// SearchGames searches for games based on a query string.
func (s *gameService) SearchGames(ctx context.Context, query string, limit, offset int) ([]*entity.Game, error) {
	games, err := s.gameRepo.SearchGames(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search games: %w", err)
	}
	return games, nil
}

// UpdateGame updates an existing game.
func (s *gameService) UpdateGame(ctx context.Context, gameID uuid.UUID, name, description, developer, publisher, genre, imageURL string) (*entity.Game, error) {
	game, err := s.gameRepo.GetGameByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("game not found: %w", err)
	}

	game.Name = name
	game.Description = description
	game.Developer = developer
	game.Publisher = publisher
	game.Genre = genre
	game.ImageURL = imageURL

	if err := s.gameRepo.UpdateGame(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}
	return game, nil
}

// DeleteGame deletes a game by its ID.
func (s *gameService) DeleteGame(ctx context.Context, gameID uuid.UUID) error {
	if err := s.gameRepo.DeleteGame(ctx, gameID); err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}
	return nil
}

// AddGameRequirement adds a new requirement for a game.
func (s *gameService) AddGameRequirement(ctx context.Context, gameID uuid.UUID, reqType, cpuModel, cpuSpeed, gpuModel, gpuMemory, os string, ramGB, storageGB float64) (*entity.GameRequirement, error) {
	// Check if game exists
	_, err := s.gameRepo.GetGameByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("game not found: %w", err)
	}

	req := entity.NewGameRequirement(gameID, reqType)
	req.CPUModel = cpuModel
	req.CPUSpeed = cpuSpeed
	req.GPUModel = gpuModel
	req.GPUMemory = gpuMemory
	req.OS = os
	req.RAMGB = ramGB
	req.StorageGB = storageGB

	if err := s.gameRepo.CreateGameRequirement(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to add game requirement: %w", err)
	}
	return req, nil
}

// GetGameRequirements retrieves all requirements for a specific game.
func (s *gameService) GetGameRequirements(ctx context.Context, gameID uuid.UUID) ([]*entity.GameRequirement, error) {
	reqs, err := s.gameRepo.GetGameRequirementsByGameID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game requirements: %w", err)
	}
	return reqs, nil
}

// UpdateGameRequirement updates an existing game requirement.
func (s *gameService) UpdateGameRequirement(ctx context.Context, reqID uuid.UUID, reqType, cpuModel, cpuSpeed, gpuModel, gpuMemory, os string, ramGB, storageGB float64) (*entity.GameRequirement, error) {
	req, err := s.gameRepo.GetGameRequirementByID(ctx, reqID) // Assuming a GetGameRequirementByID exists in repo
	if err != nil {
		return nil, fmt.Errorf("game requirement not found: %w", err)
	}

	req.Type = reqType
	req.CPUModel = cpuModel
	req.CPUSpeed = cpuSpeed
	req.GPUModel = gpuModel
	req.GPUMemory = gpuMemory
	req.OS = os
	req.RAMGB = ramGB
	req.StorageGB = storageGB

	if err := s.gameRepo.UpdateGameRequirement(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to update game requirement: %w", err)
	}
	return req, nil
}

// DeleteGameRequirement deletes a game requirement by its ID.
func (s *gameService) DeleteGameRequirement(ctx context.Context, reqID uuid.UUID) error {
	if err := s.gameRepo.DeleteGameRequirement(ctx, reqID); err != nil {
		return fmt.Errorf("failed to delete game requirement: %w", err)
	}
	return nil
}
