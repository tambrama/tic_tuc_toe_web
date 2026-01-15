package di

import (
	"tic-tac-toe/internal/app"
	"tic-tac-toe/internal/config"
	"tic-tac-toe/internal/server"
	authService "tic-tac-toe/internal/service/auth_service"
	gameService "tic-tac-toe/internal/service/game_service"
	jwtService "tic-tac-toe/internal/service/jwt_service"
	userService "tic-tac-toe/internal/service/user_service"
	"tic-tac-toe/internal/storage/postgres"
	"tic-tac-toe/internal/web/handler"

	"go.uber.org/fx"
)

// FX сам определяет порядок создания объектов
// FX управляет жизненным циклом компонентов
// Проверяет что все зависимости могут быть созданы
var Module = fx.Module("tic-tac-toe",
	fx.Provide(
		config.NewConfig,
		func(cfg *config.Config) []byte {
			return cfg.JWT
		},
		postgres.NewDB,
		postgres.NewGameRepository,
		postgres.NewUserRepository,
		postgres.NewTokenRepository,
		jwtService.NewJwtProvider,
		gameService.NewGameService,
		userService.NewUserServices,
		authService.NewAuthServices,
		handler.NewGameAPI,
		handler.NewAuthAPI,
		server.NewServer,
	),
	//запуск
	fx.Invoke(app.NewApp),
)
