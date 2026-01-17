package model

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	UUID     uuid.UUID
	Login    string
	Password string
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByLogin(ctx context.Context, login string) (*User, error)
}
