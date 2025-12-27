package dto

import (
	"tic-tac-toe/internal/domain/models"

	"github.com/google/uuid"
)

type GameFieldResponse struct {
	Field [][]int `json:"field"`
}

type GameResponse struct {
	UUID        uuid.UUID                 `json:"uuid"`
	Field       *GameFieldResponse        `json:"field"`
	StatusGame  models.GameStatus         `json:"status"`
	PlayerX     uuid.UUID                 `json:"player_x"`
	PlayerO     *uuid.UUID                `json:"player_o"`
	CurrentTurn uuid.UUID                 `json:"current_turn"`
	Symbols     map[uuid.UUID]models.Char `json:"symbols"`
	Status      string                    `json:"status_game"`
	Message     string                    `json:"message,omitempty"`
}

type UserResponse struct {
	UUID  uuid.UUID `json:"uuid"`
	Login string    `json:"login"`
}

type NewGameRequest struct{
	WithBot bool `json:"with_bot"`
}