package di

import (
	"tic-tac-toe/api"
	"tic-tac-toe/internal/app"
	configServer "tic-tac-toe/internal/config"
	configDB "tic-tac-toe/internal/database/config"
	"tic-tac-toe/internal/datasource/repository"
	"tic-tac-toe/internal/domain/services"
	"tic-tac-toe/internal/server"

	database "tic-tac-toe/internal/database/connection"

	"go.uber.org/fx"
)

//FX сам определяет порядок создания объектов
//FX управляет жизненным циклом компонентов
//Проверяет что все зависимости могут быть созданы
var Module = fx.Module("tic-tac-toe",
	fx.Provide(
		configServer.NewConfig,
		configDB.NewConfigDB,
		database.NewDB,
		repository.NewGameRepository,
		func(repo repository.GameRepository) services.GameServices {
			return services.NewGameService(repo)
		},
		api.NewGameAPI,
		server.NewServer,
	),
	//запуск
	fx.Invoke(app.NewApp),
)
