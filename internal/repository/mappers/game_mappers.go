package mappers

import (
	"tic-tac-toe/internal/domain/models"
	"tic-tac-toe/internal/repository/dto"
)

func CurrentGameFromDBToDomain(dbModel *dto.CurrentGameDTO) *models.CurrentGame {
	if dbModel == nil {
		return nil
	}

	return &models.CurrentGame{
		UUID:        dbModel.UUID,
		Field:       &models.GameField{Field: dbModel.Field.Field},
		Status:      dbModel.Status,
		PlayerX:     dbModel.PlayerX,
		PlayerO:     dbModel.PlayerO,
		CurrentTurn: dbModel.CurrentTurn,
		Symbols:     dbModel.Symbols,
	}
}

func CurrentGameFromDomainToDB(model *models.CurrentGame) *dto.CurrentGameDTO {
	if model == nil {
		return nil
	}

	return &dto.CurrentGameDTO{
		UUID:        model.UUID,
		Field:       &dto.GameFieldDTO{Field: model.Field.Field},
		Status:      model.Status,
		PlayerX:     model.PlayerX,
		PlayerO:     model.PlayerO,
		CurrentTurn: model.CurrentTurn,
		Symbols:     model.Symbols,
	}
}
