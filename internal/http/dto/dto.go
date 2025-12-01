package dto

import "github.com/google/uuid"

type GameFieldResponse struct {
	Field [][]int `json:"field"`
}

type GameResponse struct {
	UUID    uuid.UUID          `json:"uuid"`
	Field   *GameFieldResponse `json:"field"`
	Status  string             `json:"status"`
	Message string             `json:"message,omitempty"`
}

type UserResponse struct {
	UUID     uuid.UUID `db:"uuid"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}
