package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"tic-tac-toe/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

//хранилище данных
// type gameRepository struct {
// 	games map[int]*models.CurrentGame
// 	mu sync.RWMutex
// }

type gameRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) GameRepository {
	return &gameRepositoryDB{
		pool: pool,
	}
}

func (r *gameRepositoryDB) SaveGame(ctx context.Context, game *models.CurrentGame) error {
	query := `INSERT INTO games(uuid, field) 
	VALUES ($1, $2) 
	ON CONFLICT (uuid)
	DO UPDATE SET field = $2, updated_at = NOW()`

	// Сериализуем поле в JSON
	fieldJSON, err := json.Marshal(game.Field.Field)
	if err != nil {
		return fmt.Errorf("ошибка сериализации поля: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, game.UUID, fieldJSON)
	if err != nil {
		return fmt.Errorf("ошибка сохранения игры: %w", err)
	}
	return nil
}

func (r *gameRepositoryDB) GetCurrentGame(ctx context.Context, id uuid.UUID) (*models.CurrentGame, error) {
	query := `SELECT uuid, field 
	FROM games 
	WHERE uuid = $1`

	var gameUUID uuid.UUID
	var fieldJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&gameUUID, &fieldJSON)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения игры: %w", err)
	}

	// Десериализуем JSON в поле
	var fieldData [][]int
	if err := json.Unmarshal(fieldJSON, &fieldData); err != nil {
		return nil, fmt.Errorf("ошибка десериализации поля: %w", err)
	}

	return &models.CurrentGame{
		UUID:  gameUUID,
		Field: &models.GameField{Field: fieldData},
	}, nil
}
