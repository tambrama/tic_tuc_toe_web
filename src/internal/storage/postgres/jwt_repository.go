package postgres

import (
	"context"
	"errors"
	"fmt"
	model "tic-tac-toe/internal/domain/model/auth"
	dto "tic-tac-toe/internal/dto/database"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) model.TokenRepository {
	return &tokenRepository{
		pool: pool,
	}
}

func (t *tokenRepository) Save(ctx context.Context, token dto.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (token_hash, user_id, expires_at)
		VALUES ($1,$2,$3)`
	_, err := t.pool.Exec(ctx, query, token.TokenHash, token.UserID, token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("Ошибка добавления токена в бд: %w", err)
	}
	return nil
}

func (t *tokenRepository) FindByHash(ctx context.Context, hash string) (dto.RefreshToken, error) {
	query := `SELECT user_id, expires_at, created_at FROM refresh_tokens
		WHERE token_hash = $1`

	var (
		userId    uuid.UUID
		expiresAt time.Time
		createdAt time.Time
	)
	err := t.pool.QueryRow(ctx, query, hash).Scan(&userId, &expiresAt, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.RefreshToken{}, errors.New("refresh token not found")
		}
		return dto.RefreshToken{}, fmt.Errorf("ошибка поиска токена: %w", err)
	}
	return dto.RefreshToken{
		TokenHash: hash,
		UserID:    userId,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}, nil
}

func (t *tokenRepository) DeleteByHash(ctx context.Context, hash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`

	_, err := t.pool.Exec(ctx, query, hash)
	if err != nil {
		return fmt.Errorf("Ошибка удаления токена в бд: %w", err)
	}
	return nil
}

func (t *tokenRepository) DeleteAllByUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`

	_, err := t.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("Ошибка удаления токена в бд: %w", err)
	}
	return nil
}
