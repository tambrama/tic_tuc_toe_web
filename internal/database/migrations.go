package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	query_game := `
	CREATE TABLE IF NOT EXISTS games(
	uuid UUID PRIMARY KEY,
	field JSONB NOT NULL,
	status INTEGER NOT NULL DEFAULT 0,
	player_x UUID NOT NULL,
	player_o UUID,
	current_turn UUID NOT NULL,
	symbols JSONB NOT NULL DEFAULT '{}',
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW()
	); `

	_, err := pool.Exec(ctx, query_game)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	query_user := `
	CREATE TABLE IF NOT EXISTS users(
	uuid UUID PRIMARY KEY,
	login VARCHAR(32) UNIQUE NOT NULL,
	password TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW()
	); `

	_, er := pool.Exec(ctx, query_user)
	if er != nil {
		return fmt.Errorf("ошибка создания таблицы: %w", er)
	}

	return nil
}
