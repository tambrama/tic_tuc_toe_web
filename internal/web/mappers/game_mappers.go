package mappers

import (
	model "tic-tac-toe/internal/domain/model/game"
	dto "tic-tac-toe/internal/web/dto"
)

// Web -> Domain
func CurrentGameFromWebToDomain(dbModel dto.GameResponse) model.Game {
	return model.Game{
		UUID:        dbModel.UUID,
		Field:       &model.GameField{Field: dbModel.Field.Field},
		Status:      dbModel.StatusGame,
		PlayerX:     dbModel.PlayerX,
		PlayerO:     dbModel.PlayerO,
		CurrentTurn: dbModel.UUID,
		Symbols:     dbModel.Symbols,
	}
}

// Domain -> Web
func CurrentGameFromDomainToWeb(model model.Game, status model.GameStatus) dto.GameResponse {
	return dto.GameResponse{
		UUID:        model.UUID,
		Field:       &dto.GameFieldResponse{Field: model.Field.Field},
		StatusGame:  model.Status,
		PlayerX:     model.PlayerX,
		PlayerO:     model.PlayerO,
		CurrentTurn: model.CurrentTurn,
		Symbols:     model.Symbols,
		Status:      stringStatus(status),
	}
}

// func CurrentGameFromDomainToWeb(model model.UserLeaders) dto.GameResponse {
// 	return dto.GameResponse{
// 		UUID:  model.UUID,
// 		Field: &dto.GameFieldResponse{Field: model.Field.Field},
// 		StatusGame: model.Status,
// 		PlayerX: model.PlayerX,
// 		PlayerO: model.PlayerO,
// 		CurrentTurn: model.CurrentTurn,
// 		Symbols: model.Symbols,
// 		Status: stringStatus(status),
// 	}
// }

func stringStatus(status model.GameStatus) string {
	switch status {
	case model.Waiting:
		return "waiting"
	case model.Playing:
		return "in_progress"
	case model.WonX:
		return "won_X"
	case model.WonO:
		return "won_X"
	case model.Draw:
		return "draw"
	default:
		return "unknown"
	}
}
