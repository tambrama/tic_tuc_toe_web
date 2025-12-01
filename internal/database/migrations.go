package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS games(
	uuid UUID PRIMARY KEY,
	field JSONB NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW()
	); `

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %w", err)
	}
	return nil
}