package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"tic-tac-toe/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type gameRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) GameRepository {
	return &gameRepositoryDB{
		pool: pool,
	}
}

func (r *gameRepositoryDB) SaveGame(ctx context.Context, game *models.CurrentGame) error {
	query := `INSERT INTO games(uuid, field, status, player_x, player_o, current_turn, symbols) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	ON CONFLICT (uuid)
	DO UPDATE SET field = $2, status = $3,
			player_o = $5,
			current_turn = $6,
			symbols = $7, updated_at = NOW()`

	// Сериализуем поле в JSON
	fieldJSON, err := json.Marshal(game.Field.Field)
	if err != nil {
		return fmt.Errorf("ошибка сериализации поля: %w", err)
	}

	symbolJSON, err := json.Marshal(game.Symbols)
	if err != nil {
		return fmt.Errorf("ошибка сериализации поля: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, game.UUID, fieldJSON, game.Status,
		game.PlayerX, game.PlayerO, game.CurrentTurn, symbolJSON)
	if err != nil {
		return fmt.Errorf("ошибка сохранения игры: %w", err)
	}
	return nil
}

func (r *gameRepositoryDB) GetCurrentGame(ctx context.Context, id uuid.UUID) (*models.CurrentGame, error) {
	query := `SELECT uuid, field, status, player_x, player_o, current_turn, symbols 
	FROM games 
	WHERE uuid = $1`

	var (
		gameUUID    uuid.UUID
		fieldJSON   []byte
		status      models.GameStatus
		playerX     uuid.UUID
		playerO     *uuid.UUID
		currentTurn uuid.UUID
		symbolJSON  []byte
	)

	err := r.pool.QueryRow(ctx, query, id).Scan(&gameUUID, &fieldJSON, &status, &playerX, &playerO, &currentTurn, &symbolJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // или errors.New("game not found")
		}
		return nil, fmt.Errorf("ошибка получения игры: %w", err)
	}

	// Десериализуем JSON в поле
	var fieldData [][]int
	if err := json.Unmarshal(fieldJSON, &fieldData); err != nil {
		return nil, fmt.Errorf("ошибка десериализации поля: %w", err)
	}
	var symbolData map[uuid.UUID]models.Char
	if err := json.Unmarshal(symbolJSON, &symbolData); err != nil {
		return nil, fmt.Errorf("ошибка десериализации поля: %w", err)
	}

	return &models.CurrentGame{
		UUID:        gameUUID,
		Field:       &models.GameField{Field: fieldData},
		Status:      status,
		PlayerX:     playerX,
		PlayerO:     playerO,
		CurrentTurn: currentTurn,
		Symbols:     symbolData,
	}, nil
}

func (r *gameRepositoryDB) GetAvailableGames(ctx context.Context) ([]*models.CurrentGame, error) {
	query := `SELECT uuid, field, status, player_x, player_o, current_turn, symbols 
	FROM games 
	WHERE status = $1 and player_o IS NULL`

	rows, err := r.pool.Query(ctx, query, models.Waiting)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения игры: %w", err)
	}
	defer rows.Close()
	var games []*models.CurrentGame

	for rows.Next() {
		var (
			gameUUID    uuid.UUID
			fieldJSON   []byte
			status      models.GameStatus
			playerX     uuid.UUID
			playerO     *uuid.UUID
			currentTurn uuid.UUID
			symbolJson  []byte
		)

		if err := rows.Scan(&gameUUID, &fieldJSON, &status, &playerX, &playerO, &currentTurn, &symbolJson); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		// Десериализуем JSON в поле
		var fieldData [][]int
		if err := json.Unmarshal(fieldJSON, &fieldData); err != nil {
			return nil, fmt.Errorf("ошибка десериализации поля: %w", err)
		}
		var symbolData map[uuid.UUID]models.Char
		if err := json.Unmarshal(symbolJson, &symbolData); err != nil {
			return nil, fmt.Errorf("ошибка десериализации поля: %w", err)
		}

		games = append(games, &models.CurrentGame{
			UUID:        gameUUID,
			Field:       &models.GameField{Field: fieldData},
			Status:      status,
			PlayerX:     playerX,
			PlayerO:     playerO,
			CurrentTurn: currentTurn,
			Symbols:     symbolData,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return games, nil
}
