package service

import (
	"context"
	"fmt"

	"jarvis/internal/domain/entity"
	"jarvis/internal/domain/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the application service interface for user-related operations.
type UserService interface {
	RegisterUser(ctx context.Context, username, email, password string) (*entity.User, error)
	LoginUser(ctx context.Context, email, password string) (*entity.User, error)
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, username, email string) (*entity.User, error)
	ChangeUserPassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
}

// userService implements UserService.
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// RegisterUser handles user registration logic.
func (s *userService) RegisterUser(ctx context.Context, username, email, password string) (*entity.User, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := entity.NewUser(username, email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create new user entity: %w", err)
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

// LoginUser handles user login logic.
func (s *userService) LoginUser(ctx context.Context, email, password string) (*entity.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	return user, nil
}

// GetUserProfile retrieves a user's profile by ID.
func (s *userService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// UpdateUserProfile updates a user's profile information.
func (s *userService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, username, email string) (*entity.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	user.Username = username
	user.Email = email

	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return user, nil
}

// ChangeUserPassword changes a user's password.
func (s *userService) ChangeUserPassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("incorrect old password: %w", err)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	user.Password = string(hashedPassword)
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
