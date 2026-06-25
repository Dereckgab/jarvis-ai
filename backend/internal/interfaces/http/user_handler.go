package http

import (
	"context"
	"strings"
	"time"

	"jarvis/config"
	"jarvis/internal/application/service"
	"jarvis/internal/interfaces/http/dto"
	appErrors "jarvis/pkg/errors"
	"jarvis/pkg/security"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserHandler handles HTTP requests related to user management and authentication.
type UserHandler struct {
	userService service.UserService
	cfg         *config.Config
	validator   *validator.Validate
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(us service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: us,
		cfg:         cfg,
		validator:   validator.New(),
	}
}

// RegisterUser handles user registration.
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterUserRequest true "User registration details"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	var req dto.RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	user, err := h.userService.RegisterUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{Message: appErrors.ErrConflict.Error(), Details: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
	}

	accessToken, err := security.GenerateAccessToken(&h.cfg.Security, user.ID, user.Username, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: "failed to generate access token"})
	}

	refreshToken, err := security.GenerateRefreshToken(&h.cfg.Security, user.ID, user.Username, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: "failed to generate refresh token"})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

// LoginUser handles user login.
// @Summary Log in a user
// @Description Authenticate user and return JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.LoginUserRequest true "User login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	var req dto.LoginUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	user, err := h.userService.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Message: appErrors.ErrUnauthorized.Error(), Details: err.Error()})
	}

	accessToken, err := security.GenerateAccessToken(&h.cfg.Security, user.ID, user.Username, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: "failed to generate access token"})
	}

	refreshToken, err := security.GenerateRefreshToken(&h.cfg.Security, user.ID, user.Username, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: "failed to generate refresh token"})
	}

	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

// RefreshToken handles refresh token exchange and returns a new access token.
func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	claims, err := security.ValidateToken(&h.cfg.Security, req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Message: appErrors.ErrUnauthorized.Error(), Details: err.Error()})
	}

	accessToken, err := security.GenerateAccessToken(&h.cfg.Security, claims.UserID, claims.Username, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: "failed to generate access token"})
	}

	refreshToken, err := security.GenerateRefreshToken(&h.cfg.Security, claims.UserID, claims.Username, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: "failed to generate refresh token"})
	}

	return c.Status(fiber.StatusOK).JSON(dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// GetUserProfile retrieves the authenticated user's profile.
// @Summary Get user profile
// @Description Retrieve details of the authenticated user
// @Tags User
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/profile [get]
func (h *UserHandler) GetUserProfile(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	userID := c.Locals("userID").(uuid.UUID)

	user, err := h.userService.GetUserProfile(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: appErrors.ErrNotFound.Error(), Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}

// UpdateUserProfile updates the authenticated user's profile.
// @Summary Update user profile
// @Description Update details of the authenticated user
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body dto.UpdateUserProfileRequest true "User profile update details"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/profile [put]
func (h *UserHandler) UpdateUserProfile(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	userID := c.Locals("userID").(uuid.UUID)

	var req dto.UpdateUserProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	user, err := h.userService.UpdateUserProfile(ctx, userID, req.Username, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: appErrors.ErrInternalServerError.Error(), Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}

// ChangePassword handles changing the authenticated user's password.
// @Summary Change user password
// @Description Change the password for the authenticated user
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param password body dto.ChangePasswordRequest true "Old and new passwords"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/password [put]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.ReadTimeout)
	defer cancel()

	userID := c.Locals("userID").(uuid.UUID)

	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrBadRequest.Error(), Details: err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: appErrors.ErrValidation.Error(), Details: err.Error()})
	}

	if err := h.userService.ChangeUserPassword(ctx, userID, req.OldPassword, req.NewPassword); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Message: appErrors.ErrUnauthorized.Error(), Details: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
