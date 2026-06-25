package middleware

import (
	"jarvis/config"
	appErrors "jarvis/pkg/errors"
	"jarvis/pkg/security"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthRequired is a Fiber middleware to check for a valid JWT token.
func AuthRequired(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": appErrors.ErrUnauthorized.Error(),
				"details": "Authorization header is missing",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": appErrors.ErrUnauthorized.Error(),
				"details": "Invalid Authorization header format",
			})
		}

		tokenString := parts[1]
		claims, err := security.ValidateToken(&cfg.Security, tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": appErrors.ErrUnauthorized.Error(),
				"details": err.Error(),
			})
		}

		// Store user ID in context for downstream handlers
		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
