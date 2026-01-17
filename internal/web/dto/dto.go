package dto

import (
	model "tic-tac-toe/internal/domain/model/game"

	"github.com/google/uuid"
)

type GameFieldResponse struct {
	Field [][]int `json:"field"`
}

type GameResponse struct {
	UUID        uuid.UUID                `json:"uuid"`
	Field       *GameFieldResponse       `json:"field"`
	StatusGame  model.GameStatus         `json:"status"`
	PlayerX     uuid.UUID                `json:"player_x"`
	PlayerO     *uuid.UUID               `json:"player_o"`
	CurrentTurn uuid.UUID                `json:"current_turn"`
	Symbols     map[uuid.UUID]model.Char `json:"symbols"`
	Status      string                   `json:"status_game"`
	Message     string                   `json:"message,omitempty"`
}

type SignUpRequest struct {
	Login    string `json:"login" validate:"required,min=5,max=32"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	UUID  uuid.UUID `json:"uuid"`
	Login string    `json:"login"`
}

type NewGameRequest struct {
	WithBot bool `json:"withBot"`
}

type CountLeaderRequest struct {
	Count int `json:"count"`
}

type LeaderResponse struct {
	Login   string    `json:"login"`
	UserId  uuid.UUID `json:"uuid"`
	WinRate string    `json:"win_rate"`
}


type JwtRequest struct {
	Login    string `json:"login" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type JwtResponse struct {
	Type         string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refresh_token"`
}
