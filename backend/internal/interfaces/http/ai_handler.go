package http

import (
	"context"
	"fmt"

	"jarvis/config"
	"jarvis/internal/application/service"
	appErrors "jarvis/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AIHandler handles HTTP requests related to AI interactions.
type AIHandler struct {
	aiService         service.AIService
	systemInfoService service.SystemInfoService
	cfg               *config.Config
	validator         *validator.Validate
}

// NewAIHandler creates a new AIHandler.
func NewAIHandler(as service.AIService, sis service.SystemInfoService, cfg *config.Config) *AIHandler {
	return &AIHandler{
		aiService:         as,
		systemInfoService: sis,
		cfg:               cfg,
		validator:         validator.New(),
	}
}

// ChatRequest represents the request body for AI chat completion.
type ChatRequest struct {
	Prompt string `json:"prompt" validate:"required,min=1"`
}

// ChatResponse represents the response body for AI chat completion.
type ChatResponse struct {
	Response string `json:"response"`
}

// ChatCompletion handles AI chat completion requests.
func (h *AIHandler) ChatCompletion(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.WriteTimeout)
	defer cancel()

	userID := c.Locals("userID").(uuid.UUID)

	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	// Build enriched prompt with user's PC specs as context
	enrichedPrompt := h.buildPromptWithContext(ctx, userID, req.Prompt)

	response, err := h.aiService.ChatCompletion(ctx, userID.String(), enrichedPrompt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: fmt.Sprintf("AI service error: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(ChatResponse{Response: response})
}

// buildPromptWithContext prepends the user's PC specs to the prompt so JARVIS can give
// context-aware answers about game compatibility and performance.
func (h *AIHandler) buildPromptWithContext(ctx context.Context, userID uuid.UUID, userPrompt string) string {
	sysInfo, err := h.systemInfoService.GetLatestSystemInfo(ctx, userID)
	if err != nil || sysInfo == nil {
		// No system info available — send prompt as-is
		return userPrompt
	}

	freeMemGB := sysInfo.TotalMemoryGB - sysInfo.UsedMemoryGB

	systemContext := fmt.Sprintf(`[JARVIS SYSTEM CONTEXT]
The user's PC specifications are:
- CPU: %s (%d cores, %d threads, %.0f MHz)
- RAM: %.1f GB total, %.1f GB used, %.1f GB free
- Disk: %.1f GB total, %.1f GB used (%.1f%% full)
- GPU: %s
- OS: %s

Use these specs to give accurate, personalized answers about game compatibility,
performance expectations, and hardware requirements. Always reference the user's
actual hardware when relevant.
[END CONTEXT]

User question: %s`,
		sysInfo.CPUName, sysInfo.CPUCores, sysInfo.CPUThreads, sysInfo.CPUFrequency,
		sysInfo.TotalMemoryGB, sysInfo.UsedMemoryGB, freeMemGB,
		sysInfo.TotalDiskGB, sysInfo.UsedDiskGB, sysInfo.DiskPercent,
		sysInfo.GPUName,
		sysInfo.OSPlatform,
		userPrompt,
	)

	return systemContext
}
