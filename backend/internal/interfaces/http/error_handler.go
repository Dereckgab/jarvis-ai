package http

import (
	"errors"

	appErrors "jarvis/pkg/errors"
	"jarvis/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles application-wide errors.
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError
	message := appErrors.ErrInternalServerError.Error()
	details := err.Error()

	// Retrieve the custom status code if it's a *fiber.Error
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
		details = fiberErr.Error()
	}

	// Log the error
	logger.Error("HTTP Request Error",
		"path", c.Path(),
		"method", c.Method(),
		"status", code,
		"message", message,
		"details", details,
	)

	// Send custom error response
	return c.Status(code).JSON(ErrorResponse{
		Message: message,
		Details: details,
	})
}
