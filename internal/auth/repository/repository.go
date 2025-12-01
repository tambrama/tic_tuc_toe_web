package repository

import (
	"context"
	"tic-tac-toe/internal/auth/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}