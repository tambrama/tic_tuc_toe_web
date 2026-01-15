package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	model "tic-tac-toe/internal/domain/model/game"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type gameRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) model.GameRepository {
	return &gameRepositoryDB{
		pool: pool,
	}
}

func (r *gameRepositoryDB) SaveGame(ctx context.Context, game model.Game) error {
	query := `INSERT INTO games(uuid, field, status, player_x, player_o, current_turn, symbols, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	ON CONFLICT (uuid)
	DO UPDATE SET field = $2, status = $3,
			player_o = $5,
			current_turn = $6,
			symbols = $7, 
			created_at = $8,
			updated_at = NOW()`

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
		game.PlayerX, game.PlayerO, game.CurrentTurn, symbolJSON, game.DateCreate)
	if err != nil {
		return fmt.Errorf("ошибка сохранения игры: %w", err)
	}
	return nil
}

func (r *gameRepositoryDB) GetCurrentGame(ctx context.Context, id uuid.UUID) (model.Game, error) {
	query := `SELECT uuid, field, status, player_x, player_o, current_turn, symbols, created_at 
	FROM games 
	WHERE uuid = $1`

	var (
		gameUUID    uuid.UUID
		fieldJSON   []byte
		status      model.GameStatus
		playerX     uuid.UUID
		playerO     *uuid.UUID
		currentTurn uuid.UUID
		symbolJSON  []byte
		dateCreate  time.Time
	)

	err := r.pool.QueryRow(ctx, query, id).Scan(&gameUUID, &fieldJSON, &status, &playerX, &playerO, &currentTurn, &symbolJSON, &dateCreate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Game{}, fmt.Errorf("game not found: %w", err)
		}
		return model.Game{}, fmt.Errorf("ошибка получения игры: %w", err)
	}

	// Десериализуем JSON в поле
	var fieldData [][]int
	if err := json.Unmarshal(fieldJSON, &fieldData); err != nil {
		return model.Game{}, fmt.Errorf("ошибка десериализации поля: %w", err)
	}
	var symbolData map[uuid.UUID]model.Char
	if err := json.Unmarshal(symbolJSON, &symbolData); err != nil {
		return model.Game{}, fmt.Errorf("ошибка десериализации поля: %w", err)
	}

	return model.Game{
		UUID:        gameUUID,
		Field:       &model.GameField{Field: fieldData},
		Status:      status,
		PlayerX:     playerX,
		PlayerO:     playerO,
		CurrentTurn: currentTurn,
		Symbols:     symbolData,
		DateCreate:  dateCreate,
	}, nil
}

func (r *gameRepositoryDB) GetAvailableGames(ctx context.Context) ([]model.Game, error) {
	query := `SELECT uuid, field, status, player_x, created_at  
	FROM games 
	WHERE status = $1 and player_o IS NULL`

	rows, err := r.pool.Query(ctx, query, model.Waiting)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения игры: %w", err)
	}
	defer rows.Close()
	var games []model.Game

	for rows.Next() {
		var (
			gameUUID   uuid.UUID
			fieldJSON  []byte
			status     model.GameStatus
			playerX    uuid.UUID
			dateCreate time.Time
		)

		if err := rows.Scan(&gameUUID, &fieldJSON, &status, &playerX, &dateCreate); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		// Десериализуем JSON в поле
		var fieldData [][]int
		if err := json.Unmarshal(fieldJSON, &fieldData); err != nil {
			return nil, fmt.Errorf("ошибка десериализации поля: %w", err)
		}

		games = append(games, model.Game{
			UUID:       gameUUID,
			Field:      &model.GameField{Field: fieldData},
			Status:     status,
			PlayerX:    playerX,
			DateCreate: dateCreate,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return games, nil
}

func (r *gameRepositoryDB) GetComplitedGames(ctx context.Context, userID uuid.UUID) ([]model.Game, error) {
	query := `SELECT uuid, field, status, player_x, player_o, symbols, created_at  
	FROM games 
	WHERE 
	(status = 2 AND player_x = $1)
	OR
	(status = 3 AND player_o = $1)
	OR
	(status = 4 AND (player_x = $1 OR player_o = $1))
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения игры: %w", err)
	}
	defer rows.Close()
	var games []model.Game

	for rows.Next() {
		var (
			gameUUID   uuid.UUID
			fieldJSON  []byte
			status     model.GameStatus
			playerX    uuid.UUID
			playerO    *uuid.UUID
			symbolJson []byte
			dateCreate time.Time
		)

		if err := rows.Scan(&gameUUID, &fieldJSON, &status, &playerX, &playerO, &symbolJson, &dateCreate); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		// Десериализуем JSON в поле
		var fieldData [][]int
		if err := json.Unmarshal(fieldJSON, &fieldData); err != nil {
			return nil, fmt.Errorf("ошибка десериализации поля: %w", err)
		}
		var symbolData map[uuid.UUID]model.Char
		if err := json.Unmarshal(symbolJson, &symbolData); err != nil {
			return nil, fmt.Errorf("ошибка десериализации поля: %w", err)
		}

		games = append(games, model.Game{
			UUID:       gameUUID,
			Field:      &model.GameField{Field: fieldData},
			Status:     status,
			PlayerX:    playerX,
			PlayerO:    playerO,
			Symbols:    symbolData,
			DateCreate: dateCreate,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return games, nil
}

func (r *gameRepositoryDB) GetLeaderBoard(ctx context.Context, count int) ([]model.UserLeaders, error) {
	query := `SELECT 
	u.login,
    u.uuid,
    ROUND(
        COUNT (*) FIlTER (
            WHERE g.player_x = u.uuid
                AND status = 2
                OR g.player_o = u.uuid
                AND status = 3
       		 ) * 100 / COUNT(*), 2
    	) AS win_rate
	FROM users u
	JOIN games g ON g.player_x = u.uuid OR g.player_o = u.uuid 
	GROUP BY u.uuid
	ORDER BY win_rate DESC
	LIMIT $1;`

	rows, err := r.pool.Query(ctx, query, count)

	if err != nil {
		return []model.UserLeaders{}, fmt.Errorf("ошибка получения таблицы: %w", err)
	}

	defer rows.Close()
	var leaders []model.UserLeaders

	for rows.Next() {
		var userID uuid.UUID
		var login, winRate string

		if err := rows.Scan(&login, &userID, &winRate); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		leaders = append(leaders, model.UserLeaders{
			Login:   login,
			UserId:  userID,
			WinRate: winRate,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return leaders, nil
}
