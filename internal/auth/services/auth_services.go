package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
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
	return &authServices{
		userRepository: userRepository,
		validator:      validator.New(),
	}
}

func (a *authServices) Registration(ctx context.Context, account *models.SignUpRequest) (user *models.User, err error) {
	if err := a.validator.Struct(account); err != nil {
		return nil, fmt.Errorf("Ошибка валидации: %w", err)
	}
	existUser, _ := a.userRepository.GetUserByLogin(ctx, account.Login)
	// if err != nil {
	// 	 return nil, fmt.Errorf("Ошибка при проверке существования пользователя: %w", err)
	// }
	if existUser != nil {
		return nil, fmt.Errorf("Пользователь с таким логином существует: %w", existUser.Login)
	}
	userUUID := uuid.New()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Ошибка хэширования пароля: %w", err)
	}

	user = &models.User{
		UUID:     userUUID,
		Login:    account.Login,
		Password: string(hashPassword),
	}
	return user, a.userRepository.CreateUser(ctx, user)
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

func ParseCredintials(authHeader string) (login, password string, err error) {
	if authHeader == "" {
		return "", "", fmt.Errorf("empty auth header")
	}
	
	// Удаляем префикс "Basic " если он есть
	authHeader = strings.TrimSpace(authHeader)
	if strings.HasPrefix(authHeader, "Basic ") {
		authHeader = strings.TrimPrefix(authHeader, "Basic ")
		authHeader = strings.TrimSpace(authHeader)
	}
	
	creds, err := base64.StdEncoding.DecodeString(authHeader)
	if err != nil {
		return "", "", fmt.Errorf("invalid base64: %w", err)
	}
	parts := strings.SplitN(string(creds), ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format")
	}
	return parts[0], parts[1], nil
}
