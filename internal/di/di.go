package di

import (
	"tic-tac-toe/internal/app"
	authRepo "tic-tac-toe/internal/auth/repository"
	authServic "tic-tac-toe/internal/auth/services"
	"tic-tac-toe/internal/config"
	"tic-tac-toe/internal/database"
	"tic-tac-toe/internal/domain/services"
	"tic-tac-toe/internal/http/handler"
	"tic-tac-toe/internal/repository"
	"tic-tac-toe/internal/server"

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
		authRepo.NewUserRepository,
		func(repo repository.GameRepository) services.GameServices {
			return services.NewGameService(repo)
		},
		func(repo authRepo.UserRepository) authServic.UserService {
			return authServic.NewAuthServices(repo)
		},
		handler.NewGameAPI,
		handler.NewUserAPI,
		server.NewServer,
	),
	//запуск
	fx.Invoke(app.NewApp),
)
