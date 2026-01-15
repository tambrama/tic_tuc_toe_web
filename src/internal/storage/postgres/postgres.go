package postgres

import (
	"context"
	"log"
	"tic-tac-toe/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DB.URL)
	if err != nil {
		log.Fatal("Не удалось подключиться к бд", err)
		return nil, err
	}
	if err := pool.Ping((context.Background())); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
