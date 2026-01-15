package mappers

import (
	model "tic-tac-toe/internal/domain/model/user"
	dto "tic-tac-toe/internal/dto/web"
)

// ответ
func UserFromDomainToWeb(u model.User) dto.UserResponse {
	return dto.UserResponse{
		UUID:  u.UUID,
		Login: u.Login,
	}
}

func UserFromWebToDomain(dbU dto.UserResponse) model.User {
	return model.User{
		UUID:  dbU.UUID,
		Login: dbU.Login,
		// Password не передается из DTO по соображениям безопасности
	}
}
