package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GameStatus int

const (
	Waiting GameStatus = iota //ожидание
	Playing                   //игра
	WonX                      //победа X
	WonO                      //победа O
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

type Game struct {
	UUID        uuid.UUID
	Field       *GameField
	Status      GameStatus
	PlayerX     uuid.UUID
	PlayerO     *uuid.UUID
	CurrentTurn uuid.UUID
	Symbols     map[uuid.UUID]Char
	DateCreate  time.Time
}

type UserLeaders struct {
	Login   string
	UserId  uuid.UUID
	WinRate string
}

type GameRepository interface {
	SaveGame(ctx context.Context, game Game) error
	GetCurrentGame(ctx context.Context, uuid uuid.UUID) (Game, error)
	GetAvailableGames(ctx context.Context) ([]Game, error)
	GetComplitedGames(ctx context.Context, userID uuid.UUID) ([]Game, error)
	GetLeaderBoard(ctx context.Context, count int) ([]UserLeaders, error)
}
