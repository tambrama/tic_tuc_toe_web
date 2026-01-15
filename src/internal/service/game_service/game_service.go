package service

import (
	"context"
	"slices"
	model "tic-tac-toe/internal/domain/model/game"
	"time"

	"github.com/google/uuid"
)

const (
	SIZE_FIELD = 3
)

const (
	Empty = iota
	X
	O
)

type gameService struct {
	repo model.GameRepository
}

func NewGameService(repo model.GameRepository) GameServices {
	return &gameService{
		repo: repo,
	}
}

// создание новой игры
func (service *gameService) CreateNewGame(ctx context.Context, playerX uuid.UUID, withBot bool) (model.Game, error) {
	newUUID := uuid.New()

	newField := model.GameField{
		Field: [][]int{
			{Empty, Empty, Empty},
			{Empty, Empty, Empty},
			{Empty, Empty, Empty},
		},
	}
	var status model.GameStatus
	var playerO *uuid.UUID
	if withBot {
		playerO = nil
		status = model.Playing
	} else {
		playerO = nil
		status = model.Waiting
	}

	newGame := model.Game{
		UUID:        newUUID,
		Field:       &newField,
		Status:      status,
		PlayerX:     playerX,
		PlayerO:     playerO,
		CurrentTurn: playerX,
		Symbols: map[uuid.UUID]model.Char{
			playerX: model.CharX,
		},
		DateCreate: time.Now(),
	}

	return newGame, service.repo.SaveGame(ctx, newGame)
}

func (service *gameService) GetAvailableGames(ctx context.Context) ([]model.Game, error) {
	return service.repo.GetAvailableGames(ctx)
}

func (service *gameService) GetComplitedGames(ctx context.Context, userID uuid.UUID) ([]model.Game, error) {
	return service.repo.GetComplitedGames(ctx, userID)
}

func (service *gameService) JoinGame(ctx context.Context, gameID, playerO uuid.UUID) (model.Game, error) {
	gameCurrent, err := service.repo.GetCurrentGame(ctx, gameID)
	if err != nil {
		return model.Game{}, err
	}
	if gameCurrent.Status != model.Waiting {
		return model.Game{}, ErrGameNotWaiting
	}
	if gameCurrent.PlayerO != nil {
		return model.Game{}, ErrGameFull
	}
	if gameCurrent.PlayerX == playerO {
		return model.Game{}, ErrCannotJoinOwn
	}
	//обновляем игру
	gameCurrent.Symbols[playerO] = model.CharO
	gameCurrent.Status = model.Playing
	gameCurrent.PlayerO = &playerO
	gameCurrent.CurrentTurn = gameCurrent.PlayerX

	return gameCurrent, service.repo.SaveGame(ctx, gameCurrent)
}

func (service *gameService) GetCurrentGame(ctx context.Context, gameID uuid.UUID) (model.Game, error) {
	return service.repo.GetCurrentGame(ctx, gameID)
}

func (service *gameService) GetLeaderBoard(ctx context.Context, count int) ([]model.UserLeaders, error) {
	return service.repo.GetLeaderBoard(ctx, count)
}

// MakeMove обрабатывает ход игрока
// Принимает gameID, playerID и новое игровое поле после хода
// Возвращает обновленную игру
func (service *gameService) MakeMove(ctx context.Context, gameID, playerID uuid.UUID, newField *model.GameField) (model.Game, error) {
	// Загружаем текущую игру из БД
	gameCurrent, err := service.repo.GetCurrentGame(ctx, gameID)
	if err != nil {
		return model.Game{}, err
	}

	// Проверяем статус игры
	if gameCurrent.Status != model.Playing {
		return model.Game{}, ErrGameFinished
	}

	// Проверяем, что ходит правильный игрок
	if gameCurrent.CurrentTurn != playerID {
		return model.Game{}, ErrNotYourTurn
	}

	// Валидация хода
	if err := service.ValidationField(gameCurrent, model.Game{UUID: gameCurrent.UUID,
		Field: newField}); err != nil {
		return model.Game{}, err
	}

	// Применяем новое поле
	gameCurrent.Field = newField

	// Проверяем окончание игры после хода игрока
	if status := service.CheckEndGame(gameCurrent); status != model.Playing {
		gameCurrent.Status = status
		return gameCurrent, service.repo.SaveGame(ctx, gameCurrent)
	}

	///===== Игра с ботом ======
	if gameCurrent.PlayerO == nil {
		// Игра с ботом - делаем ход бота
		botGame, err := service.GetNextStep(gameCurrent)
		if err != nil {
			return model.Game{}, err
		}
		// Обновляем поле после хода бота
		gameCurrent.Field = botGame.Field

		// Проверяем окончание игры после хода бота
		if status := service.CheckEndGame(gameCurrent); status != model.Playing {
			gameCurrent.Status = status
			return gameCurrent, service.repo.SaveGame(ctx, gameCurrent)
		}

		// Возвращаем ход игроку X
		gameCurrent.CurrentTurn = gameCurrent.PlayerX
		return gameCurrent, service.repo.SaveGame(ctx, gameCurrent)
	}
	// Игра между двумя игроками - меняем текущего игрока
	if gameCurrent.CurrentTurn == gameCurrent.PlayerX {
		gameCurrent.CurrentTurn = *gameCurrent.PlayerO
	} else {
		gameCurrent.CurrentTurn = gameCurrent.PlayerX
	}

	return gameCurrent, service.repo.SaveGame(ctx, gameCurrent)
}

// получение следующего хода
func (service *gameService) GetNextStep(game model.Game) (model.Game, error) {
	status := service.CheckEndGame(game)

	if status != model.Playing {
		return model.Game{}, ErrGameFinished
	}

	bestX, bestY := service.bestStep(game)

	if bestX == -1 || bestY == -1 {
		return model.Game{}, ErrGameFinished
	}

	copyField := service.copyField(game.Field)
	copyField.Field[bestX][bestY] = O

	return model.Game{
		UUID:  game.UUID,
		Field: &copyField,
	}, nil
}

// валидация игрового поля
func (service *gameService) ValidationField(myField, botField model.Game) error {

	for i := range myField.Field.Field {
		for j := range myField.Field.Field[i] {
			if myField.Field.Field[i][j] != Empty &&
				myField.Field.Field[i][j] != botField.Field.Field[i][j] {
				return ErrInvalidMove
			}
		}
	}
	return nil
}

func (service *gameService) bestStep(m model.Game) (x, y int) {
	if service.fullField(m.Field) {
		return -1, -1
	}
	bestVal := -1000
	x, y = -1, -1
	field := m.Field
	for i := range field.Field {
		for j := range field.Field[i] {
			if field.Field[i][j] == Empty {

				field.Field[i][j] = O

				score := service.minimax(m.Field, 0, false)

				field.Field[i][j] = Empty

				if score > bestVal {
					bestVal = score
					x = i
					y = j

				}
			}
		}
	}
	return x, y
}

func (service *gameService) CheckEndGame(f model.Game) model.GameStatus {
	field := f.Field.Field
	for i := 0; i < SIZE_FIELD; i++ {
		if field[i][0] != 0 && field[i][0] == field[i][1] && field[i][1] == field[i][2] {
			if field[i][0] == X {
				return model.WonX
			} else if field[i][0] == O {
				return model.WonO
			}
		}
	}
	for j := 0; j < SIZE_FIELD; j++ {
		if field[0][j] != 0 && field[0][j] == field[1][j] && field[1][j] == field[2][j] {
			if field[0][j] == X {
				return model.WonX
			} else if field[0][j] == O {
				return model.WonO
			}
		}
	}
	if field[0][0] != 0 && field[0][0] == field[1][1] && field[1][1] == field[2][2] {
		if field[0][0] == X {
			return model.WonX
		} else if field[0][0] == O {
			return model.WonO
		}
	}

	if field[0][2] != 0 && field[0][2] == field[1][1] && field[1][1] == field[2][0] {
		if field[0][2] == X {
			return model.WonX
		} else if field[0][2] == O {
			return model.WonO
		}
	}

	for i := 0; i < SIZE_FIELD; i++ {
		for j := 0; j < SIZE_FIELD; j++ {
			if field[i][j] == Empty {
				return model.Playing
			}
		}
	}

	return model.Draw
}

func (service *gameService) minimax(field *model.GameField, depth int, isMax bool) int {

	tempGame := model.Game{Field: field}
	status := service.CheckEndGame(tempGame)
	score := service.score(status)

	if status != 0 {
		return score
	}

	if isMax { //ход бота 0
		bestScore := -1000
		for i := range field.Field {
			for j := range field.Field[i] {
				if service.emptyField(i, j, field) {
					field.Field[i][j] = O
					best := service.minimax(field, depth+1, false)
					field.Field[i][j] = Empty

					if best > bestScore {
						bestScore = best
					}
				}
			}
		}
		return bestScore
	}
	bestScore := 1000
	for i := range field.Field {
		for j := range field.Field[i] {
			if service.emptyField(i, j, field) {
				field.Field[i][j] = X
				best := service.minimax(field, depth+1, true)
				field.Field[i][j] = Empty

				if best < bestScore {
					bestScore = best
				}
			}
		}
	}
	return bestScore
}

func (service *gameService) score(status model.GameStatus) int {
	switch status {
	case model.WonX:
		return -10
	case model.WonO:
		return 10
	default:
		return 0
	}

}

func (service *gameService) emptyField(row, col int, field *model.GameField) bool {
	return field.Field[row][col] == Empty
}

func (service *gameService) fullField(field *model.GameField) bool {
	for i := range field.Field {
		if slices.Contains(field.Field[i], Empty) {
			return false
		}
	}
	return true
}

func (service *gameService) copyField(field *model.GameField) model.GameField {
	copyField := model.GameField{
		Field: make([][]int, len(field.Field)),
	}
	for i := range field.Field {
		copyField.Field[i] = make([]int, len(field.Field[i]))
		for j := range field.Field[i] {
			copyField.Field[i][j] = field.Field[i][j]
		}
	}
	return copyField
}
