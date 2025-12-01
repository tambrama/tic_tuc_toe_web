package mappers

import (
	"tic-tac-toe/internal/domain/models"
	"tic-tac-toe/internal/domain/services"
	"tic-tac-toe/internal/http/dto"
)
// Web -> Domain
func CurrentGameFromWebToDomain(dbModel *dto.GameResponse) *models.CurrentGame {
	if dbModel == nil {
		return nil
	}

	return &models.CurrentGame{
		UUID:  dbModel.UUID,
		Field: &models.GameField{Field: dbModel.Field.Field},
	}
}
// Domain -> Web
func CurrentGameFromDomainToWeb(model *models.CurrentGame, status services.GameStatus) *dto.GameResponse {
	if model == nil {
		return nil
	}

	return &dto.GameResponse{
		UUID:  model.UUID,
		Field: &dto.GameFieldResponse{Field: model.Field.Field},
		Status: stringStatus(status),
	}
}

func stringStatus(status services.GameStatus) string {
	switch status {
	case services.InProgress:
		return "in_progress"
	case services.BotWin:
		return "bot_win"
	case services.UserWin:
		return "user_win"
	case services.Draw:
		return "draw"
	default:
		return "unknown"
	}
}
