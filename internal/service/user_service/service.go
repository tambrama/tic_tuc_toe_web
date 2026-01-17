package service

import (
	"context"
	"errors"
	model "tic-tac-toe/internal/domain/model/user"
	dto "tic-tac-toe/internal/web/dto"

	"github.com/google/uuid"
)

var (
	ErrValidationFailed   = errors.New("validation failed")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordHash       = errors.New("password hash failed")
)

type UserService interface {
	Register(ctx context.Context, account dto.SignUpRequest) (model.User, error)
	Authenticate(ctx context.Context, login, password string) (model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
}
