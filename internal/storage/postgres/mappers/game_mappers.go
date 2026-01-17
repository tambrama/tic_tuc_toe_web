package mappers

import (
	model "tic-tac-toe/internal/domain/model/game"
	dto "tic-tac-toe/internal/storage/postgres/dto"
)

func CurrentGameFromDBToDomain(dbModel dto.CurrentGameDTO) model.Game {
	return model.Game{
		UUID:        dbModel.UUID,
		Field:       &model.GameField{Field: dbModel.Field.Field},
		Status:      dbModel.Status,
		PlayerX:     dbModel.PlayerX,
		PlayerO:     dbModel.PlayerO,
		CurrentTurn: dbModel.CurrentTurn,
		Symbols:     dbModel.Symbols,
	}
}

func CurrentGameFromDomainToDB(model model.Game) dto.CurrentGameDTO {
	return dto.CurrentGameDTO{
		UUID:        model.UUID,
		Field:       &dto.GameFieldDTO{Field: model.Field.Field},
		Status:      model.Status,
		PlayerX:     model.PlayerX,
		PlayerO:     model.PlayerO,
		CurrentTurn: model.CurrentTurn,
		Symbols:     model.Symbols,
	}
}
