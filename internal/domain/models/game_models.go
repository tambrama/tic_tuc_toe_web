package models

import "github.com/google/uuid"

type GameStatus int

const (
	Waiting GameStatus = iota //ожидание
	Playing                   //игра
	WonX                       //победа X
	WonO                       //победа O
	Draw                      //ничья
)

type Char string

const (
	CharX Char = "X"
	CharO Char = "O"
)

type GameField struct {
	Field [][]int
}

type CurrentGame struct {
	UUID          uuid.UUID
	Field         *GameField
	Status        GameStatus
	PlayerX       uuid.UUID
	PlayerO       *uuid.UUID
	CurrentTurn uuid.UUID
	Symbols       map[uuid.UUID]Char
}
