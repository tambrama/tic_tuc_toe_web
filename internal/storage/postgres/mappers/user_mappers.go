package mappers

import (
	model "tic-tac-toe/internal/domain/model/user"
	dto "tic-tac-toe/internal/storage/postgres/dto"
)

func UserFromDomainToDB(u model.User) dto.UserDTO {
	return dto.UserDTO{
		UUID:     u.UUID,
		Login:    u.Login,
		Password: u.Password,
	}
}

func UserFromDBToDomain(dbU dto.UserDTO) model.User {
	return model.User{
		UUID:     dbU.UUID,
		Login:    dbU.Login,
		Password: dbU.Password,
	}
}
