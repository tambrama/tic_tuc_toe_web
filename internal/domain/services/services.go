package services

import (
	"tic-tac-toe/internal/domain/models"
)

type GameServices interface {
	GetNextStep(g *models.CurrentGame) (*models.CurrentGame, error) //минмакс
	CheckEndGame(g *models.CurrentGame) GameStatus
	ValidationField(myField, botField *models.CurrentGame) error
}

