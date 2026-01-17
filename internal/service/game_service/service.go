package service

import (
	"context"
	"errors"
	model "tic-tac-toe/internal/domain/model/game"

	"github.com/google/uuid"
)

var (
	ErrNotYourTurn    = errors.New("not your turn")
	ErrInvalidMove    = errors.New("invalid move")
	ErrGameNotFound   = errors.New("game not found")
	ErrGameFinished   = errors.New("game finished")
	ErrGameNotWaiting = errors.New("game is not waiting")
	ErrGameFull       = errors.New("game is already full")
	ErrCannotJoinOwn  = errors.New("cannot join your own game")
)

type GameServices interface {
	GetNextStep(g model.Game) (model.Game, error) //минмакс
	CheckEndGame(g model.Game) model.GameStatus
	ValidationField(myField, botField model.Game) error

	CreateNewGame(ctx context.Context, playerX uuid.UUID, withBot bool) (model.Game, error)
	GetAvailableGames(ctx context.Context) ([]model.Game, error)
	GetComplitedGames(ctx context.Context, userID uuid.UUID) ([]model.Game, error)
	JoinGame(ctx context.Context, gameID, playerO uuid.UUID) (model.Game, error)
	MakeMove(ctx context.Context, gameID, player uuid.UUID, newField *model.GameField) (model.Game, error)
	GetCurrentGame(ctx context.Context, gameID uuid.UUID) (model.Game, error)

	GetLeaderBoard(ctx context.Context, count int) ([]model.UserLeaders, error)
}
