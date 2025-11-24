```
internal/
├── auth/                      # ⬅️ ОТДЕЛЬНЫЙ СЛОЙ АВТОРИЗАЦИИ
│   ├── domain/
│   │   ├── models/
│   │   │   └── user.go        # User модель
│   │   └── services/
│   │       ├── user_service.go    # UserService
│   │       └── auth_service.go    # AuthService
│   │
│   └── datasource/
│       └── repository/
│           └── user_repository.go # UserRepository
│
├── game/                      # ⬅️ ОТДЕЛЬНЫЙ СЛОЙ ИГР
│   ├── domain/
│   │   ├── models/
│   │   │   └── game.go        # Game модель
│   │   └── services/
│   │       └── game_service.go    # GameService
│   │
│   └── datasource/
│       └── repository/
│           └── game_repository.go # GameRepository
│
├── database/                  # ⬅️ ОТДЕЛЬНЫЙ СЛОЙ БАЗЫ ДАННЫХ
│   ├── connection.go
│   └── config.go
│
└── web/                      # ⬅️ WEB СЛОЙ
    ├── handlers/
    │   ├── auth_handler.go
    │   └── game_handler.go
    ├── middleware/
    │   └── auth_middleware.go
    └── dto/
        ├── auth_dto.go
        └── game_dto.go
```        