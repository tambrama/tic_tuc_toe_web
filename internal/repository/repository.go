package repository

import (
	"context"
	"tic-tac-toe/internal/domain/models"

	"github.com/google/uuid"
)
type GameRepository interface {
	SaveGame(ctx context.Context, game *models.CurrentGame) error
	GetCurrentGame(ctx context.Context, uuid uuid.UUID) (*models.CurrentGame, error)
	GetAvailableGames(ctx context.Context) ([]*models.CurrentGame, error)
}