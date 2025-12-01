package mappers

import (
	auth "tic-tac-toe/internal/auth/models"
	"tic-tac-toe/internal/repository/dto"
)

func UserFromDomainToDB(u *auth.User) *dto.UserDTO {
	if u == nil {
		return nil
	}
	return &dto.UserDTO{
		UUID:     u.UUID,
		Login:    u.Login,
		Password: u.Password,
	}
}

func UserFromDBToDomain(dbU *dto.UserDTO) *auth.User {
	if dbU == nil {
		return nil
	}
	return &auth.User{
		UUID:     dbU.UUID,
		Login:    dbU.Login,
		Password: dbU.Password,
	}
}
