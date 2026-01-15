package service

import (
	"context"
	"errors"
	"log"

	model "tic-tac-toe/internal/domain/model/user"
	dto "tic-tac-toe/internal/dto/web"

	// service "tic-tac-toe/internal/service/auth_service"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userServices struct {
	userRepository model.UserRepository
	validator      *validator.Validate
}

func NewUserServices(repo model.UserRepository) UserService {
	return &userServices{
		userRepository: repo,
		validator:      validator.New(),
	}
}

func (u *userServices) Register(ctx context.Context, account dto.SignUpRequest) (user model.User, err error) {
	if err := u.validator.Struct(account); err != nil {
		return user, ErrValidationFailed
	}
	existUser, err := u.userRepository.GetUserByLogin(ctx, account.Login)
	if existUser != nil {
		return user, ErrUserAlreadyExists
	}
	if !errors.Is(err, ErrUserNotFound) {
		return user, err
	}

	log.Printf("AUTH login='%s'", account.Login)

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, ErrPasswordHash
	}

	user = model.User{
		UUID:     uuid.New(),
		Login:    account.Login,
		Password: string(hashPassword),
	}
	return user, u.userRepository.CreateUser(ctx, user)
}

func (u *userServices) Authenticate(ctx context.Context, login, password string) (model.User, error) {
	// Логируем попытку аутентификации — до любого ответа
	log.Printf("Попытка аутентификации: login=%s", login)

	user, err := u.userRepository.GetUserByLogin(ctx, login)
	if err != nil {
		log.Printf("пользователь не найден login=%s", login)
		return model.User{}, ErrUserNotFound
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Аутентификация не удалась: неверный пароль для login=%s", login)
		return model.User{}, ErrInvalidCredentials
	}

	log.Printf("Успешная аутентификация: user_id=%s, login=%s", user.UUID, user.Login)
	return *user, nil
}

func (u *userServices) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	return u.userRepository.GetUserByID(ctx, id)
}
