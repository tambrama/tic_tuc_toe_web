package di

import (
	"tic-tac-toe/internal/app"
	"tic-tac-toe/internal/config"
	"tic-tac-toe/internal/repository"
	"tic-tac-toe/internal/domain/services"
	"tic-tac-toe/internal/server"
	"tic-tac-toe/internal/http/handler"
	"tic-tac-toe/internal/database"

	"go.uber.org/fx"
)

// FX сам определяет порядок создания объектов
// FX управляет жизненным циклом компонентов
// Проверяет что все зависимости могут быть созданы
var Module = fx.Module("tic-tac-toe",
	fx.Provide(
		config.NewConfig,
		config.NewConfigDB,
		database.NewDB,
		repository.NewGameRepository,
		func(repo repository.GameRepository) services.GameServices {
			return services.NewGameService(repo)
		},
		handler.NewGameAPI,
		server.NewServer,
	),
	//запуск
	fx.Invoke(app.NewApp),
)
