package services

import (
	"tic-tac-toe/internal/auth/models"
)

type UserService interface {
	Registration(login *models.SignUpRequest) error
	Autorization(login string, password string) (uuid int)
}