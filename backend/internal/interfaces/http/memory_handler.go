package http

import (
	"context"
	"strconv"
	"time"

	"jarvis/config"
	"jarvis/internal/application/service"
	"jarvis/internal/domain/entity"
	appErrors "jarvis/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// MemoryHandler handles HTTP requests related to memory management.
type MemoryHandler struct {
	memoryService service.MemoryService
	cfg           *config.Config
	validator     *validator.Validate
}

// NewMemoryHandler creates a new MemoryHandler.
func NewMemoryHandler(ms service.MemoryService, cfg *config.Config) *MemoryHandler {
	return &MemoryHandler{
		memoryService: ms,
		cfg:           cfg,
		validator:     validator.New(),
	}
}

// SaveMemoryRequest represents the request body for saving a memory.
type SaveMemoryRequest struct {
	Type    string `json:"type" validate:"required,oneof=short-term long-term semantic contextual"`
	Content string `json:"content" validate:"required,min=1"`
}

// MemoryResponse represents the response body for a memory entry.
type MemoryResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// SaveMemory handles saving a new memory entry.
// @Summary Save a new memory
// @Description Saves a new memory entry for the authenticated user
// @Tags Memory
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param memory body SaveMemoryRequest true "Memory details"
// @Success 201 {object} MemoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /memory [post]
func (h *MemoryHandler) SaveMemory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	userID := c.Locals("userID").(uuid.UUID)

	var req SaveMemoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	memory, err := h.memoryService.SaveMemory(ctx, userID, req.Type, req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(toMemoryResponse(memory))
}

// GetMemories handles retrieving memories for the authenticated user.
// @Summary Get user memories
// @Description Retrieves a paginated list of memories for the authenticated user, optionally filtered by type
// @Tags Memory
// @Security ApiKeyAuth
// @Produce json
// @Param type query string false "Memory type (short-term, long-term, semantic, contextual)"
// @Param limit query int false "Limit the number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} MemoryResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /memory [get]
func (h *MemoryHandler) GetMemories(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	userID := c.Locals("userID").(uuid.UUID)

	memoryType := c.Query("type", "")
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid limit parameter"})
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid offset parameter"})
	}

	memories, err := h.memoryService.GetMemories(ctx, userID, memoryType, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
	}

	var responses []MemoryResponse
	for _, memory := range memories {
		responses = append(responses, toMemoryResponse(memory))
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

// SearchMemories handles semantic search for memories.
// @Summary Search memories
// @Description Performs a semantic search across user memories
// @Tags Memory
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param query body SearchMemoryRequest true "Search query"
// @Success 200 {array} MemoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /memory/search [post]
func (h *MemoryHandler) SearchMemories(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	userID := c.Locals("userID").(uuid.UUID)

	var req SearchMemoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	memories, err := h.memoryService.SearchMemories(ctx, userID, req.Query, req.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
	}

	var responses []MemoryResponse
	for _, memory := range memories {
		responses = append(responses, toMemoryResponse(memory))
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

// SearchMemoryRequest represents the request body for searching memories.
type SearchMemoryRequest struct {
	Query string `json:"query" validate:"required,min=1"`
	Limit int    `json:"limit" validate:"required,gt=0"`
}

// DeleteMemory handles deleting a memory entry.
// @Summary Delete a memory
// @Description Deletes a specific memory entry by its ID
// @Tags Memory
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Memory ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /memory/{id} [delete]
func (h *MemoryHandler) DeleteMemory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	memoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: "invalid memory ID"})
	}

	if err := h.memoryService.DeleteMemory(ctx, memoryID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: appErrors.ErrNotFound.Error(), Details: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func toMemoryResponse(memory *entity.Memory) MemoryResponse {
	return MemoryResponse{
		ID:        memory.ID,
		UserID:    memory.UserID,
		Type:      memory.Type,
		Content:   memory.Content,
		CreatedAt: memory.CreatedAt.Format(time.RFC3339),
		UpdatedAt: memory.UpdatedAt.Format(time.RFC3339),
	}
}
