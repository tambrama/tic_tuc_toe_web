package auth

import (
	"context"
	"fmt"
	"tic-tac-toe/internal/auth/models"
	"tic-tac-toe/internal/auth/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authServices struct {
	userRepository repository.UserRepository
	validator      *validator.Validate
}

func NewAuthServices(userRepository repository.UserRepository) UserService {
	return &authServices{userRepository: userRepository}
}

func (a *authServices) Registration(ctx context.Context, account *models.SignUpRequest) error {
	if err := a.validator.Struct(account); err != nil {
		return fmt.Errorf("Ошибка валидации: %w", err)
	}
	existUser, _ := a.userRepository.GetUserByLogin(ctx, account.Login)
	if existUser != nil {
		return fmt.Errorf("Пользователь с таким логином существует: %w", existUser.Login)
	}
	userUUID := uuid.New()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Ошибка хэширования пароля: %w", err)
	}

	user := &models.User{
		UUID:     userUUID,
		Login:    account.Login,
		Password: string(hashPassword),
	}
	return a.userRepository.CreateUser(ctx, user)
}

func (a *authServices) Authenticate(ctx context.Context, login string, password string) (id uuid.UUID, err error) {
	existUser, err := a.userRepository.GetUserByLogin(ctx, login)
	if err != nil || existUser == nil {
		return uuid.Nil, fmt.Errorf("Не верные данные")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(password))
	if err != nil {
		return uuid.Nil, fmt.Errorf("Не верные данные")
	}

	return existUser.UUID, nil
}
