package service

import (
	model "tic-tac-toe/internal/domain/model/user"
)

type JwtProvider interface {
	GenerateAccessToken(user model.User) (string, error)
	GenerateRefreshToken(user model.User) (string, error)
	ValidateAccessToken(tokenStr string) (*CustomClaims, error)
	ValidateRefreshToken(tokenStr string) (*CustomClaims, error)
}
