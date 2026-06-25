package http

import (
	"context"
	"strconv"

	"jarvis/config"
	"jarvis/internal/application/service"
	"jarvis/internal/domain/entity"
	"jarvis/internal/interfaces/http/dto"
	appErrors "jarvis/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GameHandler handles HTTP requests related to game management.
type GameHandler struct {
	gameService service.GameService
	cfg         *config.Config
	validator   *validator.Validate
}

// NewGameHandler creates a new GameHandler.
func NewGameHandler(gs service.GameService, cfg *config.Config) *GameHandler {
	return &GameHandler{
		gameService: gs,
		cfg:         cfg,
		validator:   validator.New(),
	}
}

// CreateGame handles the creation of a new game.
// @Summary Create a new game
// @Description Creates a new game entry in the system
// @Tags Games
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param game body dto.CreateGameRequest true "Game creation details"
// @Success 201 {object} dto.GameResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /games [post]
func (h *GameHandler) CreateGame(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	var req dto.CreateGameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	game, err := h.gameService.CreateGame(ctx, req.Name, req.Description, req.Developer, req.Publisher, req.Genre, req.ImageURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(toGameResponse(game))
}

// GetGame handles retrieving a game by ID.
// @Summary Get game by ID
// @Description Retrieves details of a specific game by its ID
// @Tags Games
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} dto.GameResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /games/{id} [get]
func (h *GameHandler) GetGame(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	gameID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid game ID"})
	}

	game, err := h.gameService.GetGame(ctx, gameID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: appErrors.ErrNotFound.Error(), Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(toGameResponse(game))
}

// SearchGames handles searching for games.
// @Summary Search games
// @Description Searches for games based on a query string
// @Tags Games
// @Security ApiKeyAuth
// @Produce json
// @Param query query string false "Search query for game name"
// @Param limit query int false "Limit the number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} dto.GameResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /games/search [get]
func (h *GameHandler) SearchGames(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	query := c.Query("query", "")
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid limit parameter"})
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid offset parameter"})
	}

	games, err := h.gameService.SearchGames(ctx, query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
	}

	var responses []dto.GameResponse
	for _, game := range games {
		responses = append(responses, toGameResponse(game))
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

// UpdateGame handles updating an existing game.
// @Summary Update an existing game
// @Description Updates details of an existing game by its ID
// @Tags Games
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param game body dto.UpdateGameRequest true "Game update details"
// @Success 200 {object} dto.GameResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /games/{id} [put]
func (h *GameHandler) UpdateGame(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	gameID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid game ID"})
	}

	var req dto.UpdateGameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	game, err := h.gameService.UpdateGame(ctx, gameID, req.Name, req.Description, req.Developer, req.Publisher, req.Genre, req.ImageURL)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: appErrors.ErrNotFound.Error(), Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(toGameResponse(game))
}

// DeleteGame handles deleting a game by ID.
// @Summary Delete a game
// @Description Deletes a game entry by its ID
// @Tags Games
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Game ID"
// @Success 204 "No Content"
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /games/{id} [delete]
func (h *GameHandler) DeleteGame(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	gameID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid game ID"})
	}

	if err := h.gameService.DeleteGame(ctx, gameID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: appErrors.ErrNotFound.Error(), Details: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// AddGameRequirement handles adding a new requirement to a game.
// @Summary Add game requirement
// @Description Adds a new hardware requirement for a specific game
// @Tags Games
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param requirement body dto.AddGameRequirementRequest true "Game requirement details"
// @Success 201 {object} dto.GameRequirementDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /games/{id}/requirements [post]
func (h *GameHandler) AddGameRequirement(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	gameID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid game ID"})
	}

	var req dto.AddGameRequirementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	reqEntity, err := h.gameService.AddGameRequirement(ctx, gameID, req.Type, req.CPUModel, req.CPUSpeed, req.GPUModel, req.GPUMemory, req.OS, req.RAMGB, req.StorageGB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(toGameRequirementDTO(reqEntity))
}

// GetGameRequirements handles retrieving all requirements for a game.
// @Summary Get game requirements
// @Description Retrieves all hardware requirements for a specific game
// @Tags Games
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {array} dto.GameRequirementDTO
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /games/{id}/requirements [get]
func (h *GameHandler) GetGameRequirements(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	gameID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid game ID"})
	}

	reqs, err := h.gameService.GetGameRequirements(ctx, gameID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: appErrors.ErrNotFound.Error(), Details: err.Error()})
	}

	var responses []dto.GameRequirementDTO
	for _, req := range reqs {
		responses = append(responses, toGameRequirementDTO(req))
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

// UpdateGameRequirement handles updating an existing game requirement.
// @Summary Update game requirement
// @Description Updates an existing hardware requirement for a specific game
// @Tags Games
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param reqID path string true "Requirement ID"
// @Param requirement body dto.UpdateGameRequirementRequest true "Game requirement update details"
// @Success 200 {object} dto.GameRequirementDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /games/{id}/requirements/{reqID} [put]
func (h *GameHandler) UpdateGameRequirement(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	_, err := uuid.Parse(c.Params("id")) // Ensure game exists, though not directly used in service call
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid game ID"})
	}

	reqID, err := uuid.Parse(c.Params("reqID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid requirement ID"})
	}

	var req dto.UpdateGameRequirementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	reqEntity, err := h.gameService.UpdateGameRequirement(ctx, reqID, req.Type, req.CPUModel, req.CPUSpeed, req.GPUModel, req.GPUMemory, req.OS, req.RAMGB, req.StorageGB)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: appErrors.ErrNotFound.Error(), Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(toGameRequirementDTO(reqEntity))
}

// DeleteGameRequirement handles deleting a game requirement.
// @Summary Delete game requirement
// @Description Deletes a specific hardware requirement for a game
// @Tags Games
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Game ID"
// @Param reqID path string true "Requirement ID"
// @Success 204 "No Content"
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /games/{id}/requirements/{reqID} [delete]
func (h *GameHandler) DeleteGameRequirement(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	_, err := uuid.Parse(c.Params("id")) // Ensure game exists, though not directly used in service call
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid game ID"})
	}

	reqID, err := uuid.Parse(c.Params("reqID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid requirement ID"})
	}

	if err := h.gameService.DeleteGameRequirement(ctx, reqID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: appErrors.ErrNotFound.Error(), Details: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func toGameResponse(game *entity.Game) dto.GameResponse {
	var requirementsDTO []dto.GameRequirementDTO
	for _, req := range game.Requirements {
		reqData := toGameRequirementDTO(&req)
		requirementsDTO = append(requirementsDTO, reqData)
	}

	return dto.GameResponse{
		ID:           game.ID,
		Name:         game.Name,
		Description:  game.Description,
		ReleaseDate:  game.ReleaseDate,
		Developer:    game.Developer,
		Publisher:    game.Publisher,
		Genre:        game.Genre,
		ImageURL:     game.ImageURL,
		Requirements: requirementsDTO,
		CreatedAt:    game.CreatedAt,
		UpdatedAt:    game.UpdatedAt,
	}
}

func toGameRequirementDTO(req *entity.GameRequirement) dto.GameRequirementDTO {
	return dto.GameRequirementDTO{
		ID:        req.ID,
		Type:      req.Type,
		CPUModel:  req.CPUModel,
		CPUSpeed:  req.CPUSpeed,
		GPUModel:  req.GPUModel,
		GPUMemory: req.GPUMemory,
		RAMGB:     req.RAMGB,
		StorageGB: req.StorageGB,
		OS:        req.OS,
	}
}
