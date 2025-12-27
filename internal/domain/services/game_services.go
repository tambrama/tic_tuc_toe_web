package services

import (
	"errors"
	"slices"
	"tic-tac-toe/internal/domain/models"
	"tic-tac-toe/internal/repository"
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
	repo repository.GameRepository
}

func NewGameService(repo repository.GameRepository) GameServices {
	return &gameService{
		repo: repo,
	}
}

// получение следующего хода
func (service *gameService) GetNextStep(game *models.CurrentGame) (*models.CurrentGame, error) {
	status := service.CheckEndGame(game)

	if status != models.Playing {
		return nil, errors.New("Progress end")
	}

	bestX, bestY := service.bestStep(game)

	if bestX == -1 || bestY == -1 {
		if service.FullField(game.Field) {
			return nil, errors.New("Game end: draw")
		}
		return nil, errors.New("No moves available")
	}

	copyField := service.CopyField(game.Field)
	copyField.Field[bestX][bestY] = O

	return &models.CurrentGame{
		UUID:  game.UUID,
		Field: &copyField,
	}, nil
}

// валидация игрового поля
func (service *gameService) ValidationField(myField, botField *models.CurrentGame) error {

	for i := range myField.Field.Field {
		for j := range myField.Field.Field[i] {
			if myField.Field.Field[i][j] != Empty &&
				myField.Field.Field[i][j] != botField.Field.Field[i][j] {
				return errors.New("Move error, this cell is not empty")
			}
		}
	}
	return nil
}

func (service *gameService) bestStep(m *models.CurrentGame) (x, y int) {
	if service.FullField(m.Field) {
		return -1, -1
	}
	bestVal := -1000
	field := m.Field
	for i := range field.Field {
		for j := range field.Field[i] {
			if field.Field[i][j] == Empty {

				field.Field[i][j] = O

				score := service.minimax(m.Field, 0, true)

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

func (service *gameService) CheckEndGame(f *models.CurrentGame) models.GameStatus {
	field := f.Field.Field
	for i := 0; i < SIZE_FIELD; i++ {
		if field[i][0] != 0 && field[i][0] == field[i][1] && field[i][1] == field[i][2] {
			if field[i][0] == X {
				return models.WonX
			} else if field[i][0] == O {
				return models.WonO
			}
		}
	}
	for j := 0; j < SIZE_FIELD; j++ {
		if field[0][j] != 0 && field[0][j] == field[1][j] && field[1][j] == field[2][j] {
			if field[0][j] == X {
				return models.WonX
			} else if field[0][j] == O {
				return models.WonO
			}
		}
	}
	if field[0][0] != 0 && field[0][0] == field[1][1] && field[1][1] == field[2][2] {
		if field[0][0] == X {
			return models.WonX
		} else if field[0][0] == O {
			return models.WonO
		}
	}

	if field[0][2] != 0 && field[0][2] == field[1][1] && field[1][1] == field[2][0] {
		if field[0][2] == X {
			return models.WonX
		} else if field[0][2] == O {
			return models.WonO
		}
	}

	for i := 0; i < SIZE_FIELD; i++ {
		for j := 0; j < SIZE_FIELD; j++ {
			if field[i][j] == Empty {
				return models.Playing
			}
		}
	}

	return models.Draw
}

func (service *gameService) minimax(field *models.GameField, depth int, isMax bool) int {

	tempGame := &models.CurrentGame{Field: field}
	status := service.CheckEndGame(tempGame)
	score := service.score(status)

	if status != 0 {
		return score
	}

	if isMax {
		bestScore := -1000
		for i := range field.Field {
			for j := range field.Field[i] {
				if service.EmptyField(i, j, field) {
					field.Field[i][j] = X
					best := service.minimax(field, depth+1, !isMax)
					field.Field[i][j] = Empty

					if best > bestScore {
						bestScore = best
					}
				}
			}
		}
		return bestScore
	} else {
		bestScore := 1000
		for i := range field.Field {
			for j := range field.Field[i] {
				if service.EmptyField(i, j, field) {
					field.Field[i][j] = O
					best := service.minimax(field, depth+1, !isMax)
					field.Field[i][j] = Empty

					if best < bestScore {
						bestScore = best
					}
				}
			}
		}
		return bestScore
	}
}

func (service *gameService) score(status models.GameStatus) int {
	switch status {
	case models.WonX:
		return -10
	case models.WonO:
		return 10
	default:
		return 0
	}

}

func (service *gameService) EmptyField(row, col int, field *models.GameField) bool {
	return field.Field[row][col] == Empty
}

func (service *gameService) FullField(field *models.GameField) bool {
	for i := range field.Field {
		if slices.Contains(field.Field[i], Empty) {
			return false
		}
	}
	return true
}

func (service *gameService) CopyField(field *models.GameField) models.GameField {
	copyField := models.GameField{
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
