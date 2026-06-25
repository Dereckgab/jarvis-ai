package dto

import "github.com/google/uuid"

// RegisterUserRequest represents the request body for user registration.
type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginUserRequest represents the request body for user login.
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents the response body for user details.
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// AuthResponse represents the response body for authentication (login/register).
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// UpdateUserProfileRequest represents the request body for updating user profile.
type UpdateUserProfileRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

// RefreshTokenRequest represents the request body for refreshing an access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents the response body for refreshing JWT tokens.
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ChangePasswordRequest represents the request body for changing user password.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
