package database

import (
	"context"
	"log"
	"tic-tac-toe/internal/database/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(cfg *config.ConfigDB)(*pgxpool.Pool, error){
	pool, err := pgxpool.New(context.Background(), cfg.URL)
	if err != nil {
		log.Printf("Не удалось подключиться к бд: %v", err)
		return nil, err
	}
	if err := pool.Ping((context.Background())); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
