package dto

type JwtRequest struct {
	Login    string `json:"login" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type JwtResponse struct {
	Type         string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refresh_token"`
}
