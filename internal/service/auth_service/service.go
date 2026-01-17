package service

import (
	"context"
	"errors"
	model "tic-tac-toe/internal/domain/model/user"
	dto "tic-tac-toe/internal/web/dto"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenInvalid       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenNotFound      = errors.New("token not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenGeneration    = errors.New("token generation failed")
	ErrTokenSaveFailed    = errors.New("token save failed")
)

type AuthService interface {
	Registration(ctx context.Context, req dto.SignUpRequest) (user model.User, err error)
	Login(ctx context.Context, req dto.JwtRequest) (res dto.JwtResponse, err error)
	RefreshAccess(ctx context.Context, refreshToken string) (dto.JwtResponse, error)
	RotateRefreshToken(ctx context.Context, refreshToken string) (dto.JwtResponse, error)
}
