package dto

import "github.com/google/uuid"

// "контейнер" для переноски данных
type GameFieldDTO struct {
	Field [][]int `json:"field" db:"field"`
}

type CurrentGameDTO struct {
	UUID  uuid.UUID     `json:"uuid" db:"uuid"`
	Field *GameFieldDTO `json:"field" db:"field"`
}

type UserDTO struct {
	UUID     uuid.UUID `db:"uuid"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}
