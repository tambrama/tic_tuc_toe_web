package mappers

import (
	"tic-tac-toe/internal/auth/models"
	"tic-tac-toe/internal/http/dto"
)
//ответ 
func UserFromDomainToWeb(u *models.User) *dto.UserResponse {
	if u == nil {
		return nil
	}
	return &dto.UserResponse{
		UUID:  u.UUID,
		Login: u.Login,
	}
}

func UserFromWebToDomain(dbU *dto.UserResponse) *models.User {
	if dbU == nil {
		return nil
	}
	return &models.User{
		UUID:  dbU.UUID,
		Login: dbU.Login,
		// Password не передается из DTO по соображениям безопасности
	}
}
