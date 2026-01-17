package dto

import (
	model "tic-tac-toe/internal/domain/model/game"
	"time"

	"github.com/google/uuid"
)

// "контейнер" для переноски данных
type GameFieldDTO struct {
	Field [][]int `db:"field"`
}

type CurrentGameDTO struct {
	UUID        uuid.UUID                `db:"uuid"`
	Field       *GameFieldDTO            `db:"field"`
	Status      model.GameStatus         `db:"status"`
	PlayerX     uuid.UUID                `db:"player_x"`
	PlayerO     *uuid.UUID               `db:"player_o"`
	CurrentTurn uuid.UUID                `db:"current_turn"`
	Symbols     map[uuid.UUID]model.Char `db:"symbols"`
	DateCreate  time.Time                `db:"created_at"`
}

type UserDTO struct {
	UUID     uuid.UUID `db:"uuid"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}

type RefreshToken struct {
	TokenHash string    `db:"token_hash"`
	UserID    uuid.UUID `db:"uuid"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
