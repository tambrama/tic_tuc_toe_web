package models

import "github.com/google/uuid"

type SignUpRequest struct {
	Login    string `json:"login" validate:"required,min=5,max=32"`
	Password string `json:"password" validate:"required,min=6"`
}

type User struct {
	UUID     uuid.UUID `db:"uuid"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}
