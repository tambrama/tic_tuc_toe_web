package auth

import (
	"context"
	"tic-tac-toe/internal/auth/models"

	"github.com/google/uuid"
)

type UserService interface {
	Registration(ctx context.Context, account *models.SignUpRequest) error
	Authenticate(ctx context.Context, login string, password string) (uuid uuid.UUID, err error)
}
