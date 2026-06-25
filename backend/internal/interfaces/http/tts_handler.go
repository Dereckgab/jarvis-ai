package http

import (
	"context"
	"fmt"

	"jarvis/config"
	"jarvis/internal/application/service"
	appErrors "jarvis/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// TTSHandler handles HTTP requests related to Text-to-Speech.
type TTSHandler struct {
	ttsService service.TTSService
	cfg        *config.Config
	validator  *validator.Validate
}

// NewTTSHandler creates a new TTSHandler.
func NewTTSHandler(ts service.TTSService, cfg *config.Config) *TTSHandler {
	return &TTSHandler{
		ttsService: ts,
		cfg:        cfg,
		validator:  validator.New(),
	}
}

// GenerateSpeechRequest represents the request body for speech generation.
type GenerateSpeechRequest struct {
	Text string `json:"text" validate:"required,min=1"`
}

// GenerateSpeech handles requests to generate speech from text.
// @Summary Generate speech from text
// @Description Converts the provided text into speech audio
// @Tags TTS
// @Security ApiKeyAuth
// @Accept json
// @Produce audio/mpeg
// @Param request body GenerateSpeechRequest true "Text to convert to speech"
// @Success 200 {file} byte "Audio stream"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tts/generate [post]
func (h *TTSHandler) GenerateSpeech(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	var req GenerateSpeechRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	audioBytes, err := h.ttsService.GenerateSpeech(ctx, req.Text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: fmt.Sprintf("TTS service error: %v", err)})
	}

	c.Set("Content-Type", "audio/mpeg") // Or appropriate audio type
	return c.Send(audioBytes)
}
