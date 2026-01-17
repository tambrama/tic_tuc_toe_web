package service

import (
	"context"
	"errors"
	"log"
	authModel "tic-tac-toe/internal/domain/model/auth"
	model "tic-tac-toe/internal/domain/model/user"
	jwtService "tic-tac-toe/internal/service/jwt_service"
	userService "tic-tac-toe/internal/service/user_service"
	db "tic-tac-toe/internal/storage/postgres/dto"
	"tic-tac-toe/internal/utils"
	dto "tic-tac-toe/internal/web/dto"

	"time"
)

type authServices struct {
	user   userService.UserService
	jwt    jwtService.JwtProvider
	tokens authModel.TokenRepository
}

func NewAuthServices(user userService.UserService, jwt jwtService.JwtProvider, tokens authModel.TokenRepository) AuthService {
	return &authServices{
		user:   user,
		jwt:    jwt,
		tokens: tokens,
	}
}

func (a *authServices) Registration(ctx context.Context, account dto.SignUpRequest) (user model.User, err error) {
	return a.user.Register(ctx, account)
}

func (a *authServices) Login(ctx context.Context, req dto.JwtRequest) (res dto.JwtResponse, err error) {

	user, err := a.user.Authenticate(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrInvalidCredentials) {
			return res, ErrInvalidCredentials
		}
		return res, err
	}

	accessToken, err := a.jwt.GenerateAccessToken(user)
	if err != nil {
		log.Printf("Ошибка генерации access token для user_id=%s: %v", user.UUID, err)
		return res, ErrTokenGeneration
	}

	refreshToken, err := a.jwt.GenerateRefreshToken(user)
	if err != nil {
		log.Printf("Ошибка генерации refresh token для user_id=%s: %v", user.UUID, err)
		return res, ErrTokenGeneration
	}

	token := db.RefreshToken{
		TokenHash: utils.HashToken(refreshToken),
		UserID:    user.UUID,
		ExpiresAt: time.Now().Add(jwtService.RefreshTokenTTL),
	}

	if err := a.tokens.Save(ctx, token); err != nil {
		log.Printf("Ошибка сохранения refresh token для user_id=%s: %v", user.UUID, err)
		return res, ErrTokenSaveFailed
	}

	log.Printf("Успешный вход: user_id=%s, access_token выдан", user.UUID)

	return dto.JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authServices) RefreshAccess(ctx context.Context, refreshToken string) (dto.JwtResponse, error) {
	claims, err := a.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return dto.JwtResponse{}, ErrTokenInvalid
	}
	hashToken := utils.HashToken(refreshToken)
	refreshTokenDB, err := a.tokens.FindByHash(ctx, hashToken)
	if err != nil {
		return dto.JwtResponse{}, ErrTokenNotFound
	}
	if time.Now().After(refreshTokenDB.ExpiresAt) {
		_ = a.tokens.DeleteByHash(ctx, hashToken)
		return dto.JwtResponse{}, ErrTokenExpired
	}
	user, err := a.user.GetByID(ctx, claims.UserID)
	if err != nil {
		return dto.JwtResponse{}, ErrUserNotFound
	}
	accessToken, err := a.jwt.GenerateAccessToken(user)
	if err != nil {
		return dto.JwtResponse{}, ErrTokenGeneration
	}
	return dto.JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authServices) RotateRefreshToken(ctx context.Context, refreshToken string) (dto.JwtResponse, error) {
	//валидируем
	claims, err := a.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return dto.JwtResponse{}, ErrTokenNotFound
	}
	//хэшируем
	hashToken := utils.HashToken(refreshToken)
	// //находим в базе данных
	refreshTokenDB, err := a.tokens.FindByHash(ctx, hashToken)
	if err != nil {
		return dto.JwtResponse{}, ErrTokenNotFound
	}
	//"вытаскиваем" пользователя
	user, err := a.user.GetByID(ctx, claims.UserID)
	if err != nil {
		return dto.JwtResponse{}, ErrTokenInvalid
	}
	//удаляем старый токен
	if time.Now().After(refreshTokenDB.ExpiresAt) {
		_ = a.tokens.DeleteByHash(ctx, hashToken)
		return dto.JwtResponse{}, ErrTokenExpired
	}
	//генерируем новый рефреш токен
	newRefreshToken, err := a.jwt.GenerateRefreshToken(user)
	if err != nil {
		return dto.JwtResponse{}, ErrTokenGeneration
	}
	//сохраняем новый токен в бд
	er := a.tokens.Save(ctx, db.RefreshToken{TokenHash: utils.HashToken(newRefreshToken), UserID: user.UUID})
	if er != nil {
		return dto.JwtResponse{}, ErrTokenSaveFailed
	}
	//генерируем новый аксес токен
	newAccessToken, err := a.jwt.GenerateAccessToken(user)
	if err != nil {
		return dto.JwtResponse{}, ErrTokenGeneration
	}

	return dto.JwtResponse{
		Type:         "Bearer",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
