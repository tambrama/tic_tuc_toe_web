package models

import "github.com/google/uuid"

type GameField struct {
	Field [][]int
}

type CurrentGame struct {
	UUID  uuid.UUID
	Field *GameField
}
