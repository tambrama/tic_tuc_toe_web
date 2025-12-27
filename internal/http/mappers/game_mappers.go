package mappers

import (
	"tic-tac-toe/internal/domain/models"
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
		Status: dbModel.StatusGame,
		PlayerX: dbModel.PlayerX,
		PlayerO: dbModel.PlayerO,
		CurrentTurn: dbModel.UUID,
		Symbols: dbModel.Symbols,
	}
}
// Domain -> Web
func CurrentGameFromDomainToWeb(model *models.CurrentGame, status models.GameStatus) *dto.GameResponse {
	if model == nil {
		return nil
	}

	return &dto.GameResponse{
		UUID:  model.UUID,
		Field: &dto.GameFieldResponse{Field: model.Field.Field},
		StatusGame: model.Status,
		PlayerX: model.PlayerX,
		PlayerO: model.PlayerO,
		CurrentTurn: model.CurrentTurn,
		Symbols: model.Symbols,
		Status: stringStatus(status),
	}
}

func stringStatus(status models.GameStatus) string {
	switch status {
	case models.Waiting:
		return "waiting"
	case models.Playing:
		return "in_progress"	
	case models.WonX:
		return "won_X"
	case models.WonO:
		return "won_X"
	case models.Draw:
		return "draw"
	default:
		return "unknown"
	}
}
