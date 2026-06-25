package repository

import (
	"context"
	"errors"
	"fmt"

	"jarvis/internal/domain/entity"
	"jarvis/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormUserRepository implements the UserRepository interface using GORM.
type GormUserRepository struct {
	DB *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository.
func NewGormUserRepository(db *gorm.DB) repository.UserRepository {
	return &GormUserRepository{DB: db}
}

// CreateUser creates a new user in the database.
func (r *GormUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	result := r.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

// GetUserByID retrieves a user by their ID.
func (r *GormUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	result := r.DB.WithContext(ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with ID %s: %w", id.String(), result.Error)
		}
		return nil, fmt.Errorf("failed to get user by ID %s: %w", id.String(), result.Error)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email address.
func (r *GormUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	result := r.DB.WithContext(ctx).First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with email %s: %w", email, result.Error)
		}
		return nil, fmt.Errorf("failed to get user by email %s: %w", email, result.Error)
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database.
func (r *GormUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	result := r.DB.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	return nil
}

// DeleteUser deletes a user by their ID.
func (r *GormUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&entity.User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user with ID %s: %w", id.String(), result.Error)
	}
	return nil
}
