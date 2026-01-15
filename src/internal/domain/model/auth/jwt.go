package model

import (
	"context"
	dto "tic-tac-toe/internal/dto/database"

	"github.com/google/uuid"
)

type TokenRepository interface {
	Save(ctx context.Context, token dto.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (dto.RefreshToken, error)
	DeleteByHash(ctx context.Context, hash string) error
	DeleteAllByUser(ctx context.Context, userID uuid.UUID) error
}
