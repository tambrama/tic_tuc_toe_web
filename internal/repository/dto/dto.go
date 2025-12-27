package dto

import (
	"tic-tac-toe/internal/domain/models"

	"github.com/google/uuid"
)

// "контейнер" для переноски данных
type GameFieldDTO struct {
	Field [][]int `json:"field" db:"field"`
}

type CurrentGameDTO struct {
	UUID  uuid.UUID     `json:"uuid" db:"uuid"`
	Field *GameFieldDTO `json:"field" db:"field"`
	Status        models.GameStatus  `json:"status" db:"status"`
	PlayerX       uuid.UUID  `json:"player_x" db:"player_x"`
	PlayerO       *uuid.UUID `json:"player_o" db:"player_o"`
	CurrentTurn uuid.UUID `json:"current_turn" db:"current_turn"`
	Symbols       map[uuid.UUID]models.Char `json:"symbols" db:"symbols"`
}

type UserDTO struct {
	UUID     uuid.UUID `db:"uuid"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}
