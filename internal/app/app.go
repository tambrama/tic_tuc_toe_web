package app

import (
	"context"
	"log"
	"tic-tac-toe/internal/database"
	"tic-tac-toe/internal/server"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewApp(lc fx.Lifecycle, server *server.Server, pool *pgxpool.Pool) {
	lc.Append(fx.Hook{
		OnStart: func (ctx context.Context) error {
			if err := database.CreateTable(ctx, pool); err != nil {
					log.Fatal("Ошибка создания таблицы:", err)
				}
			go func() {
				
				if err := server.Start(); err != nil{
					log.Fatal("Ошибка старта сервера:", err)
				}
			}()
			return nil
		},
		OnStop: func (ctx context.Context) error {
			pool.Close()
			return server.Shutdown(ctx)
		},
	})
}
