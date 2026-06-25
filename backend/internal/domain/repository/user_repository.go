package repository

import (
	"context"

	"jarvis/internal/domain/entity"

	"github.com/google/uuid"
)

// UserRepository defines the interface for managing User entities.
type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
