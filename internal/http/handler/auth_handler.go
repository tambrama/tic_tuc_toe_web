package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"tic-tac-toe/internal/auth/models"
	"tic-tac-toe/internal/auth/repository"
	auth "tic-tac-toe/internal/auth/services"
	"tic-tac-toe/internal/http/mappers"
)

type UserAPI struct {
	userServis auth.UserService
	userRepo   repository.UserRepository
}

func NewUserAPI(servis auth.UserService, repo repository.UserRepository) *UserAPI{
	return &UserAPI{
		userServis: servis,
		userRepo: repo,
	}
}

func (api *UserAPI) HandlerRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// authHeader := r.Header.Get("Registration")
	// if authHeader == "" {
	// 	http.Error(w, "Требуется заголовок регистрации", http.StatusUnauthorized)
	// }
	ctx := r.Context()
	var req models.SignUpRequest
	err := api.userServis.Registration(ctx, &req)
	if err != nil {
		api.handlerRegistrationError(w, err)
		return
	}
	response := mappers.UserFromDomainToWeb(user)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

	
}

func (api *UserAPI) handlerRegistrationError(w http.ResponseWriter, err error){
	switch err.Error() {
	case "Пользователь с таким логином существует":
		http.Error(w, "Пользователь с таким логином существует", http.StatusBadRequest)
	case "Ошибка валидации":
		http.Error(w, "Ошибка валидации", http.StatusBadRequest)
	case "Ошибка хэширования пароля":
		http.Error(w, "Ошибка хэширования пароля", http.StatusInternalServerError)
	default:
		http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
	}
	// log.Printf("Ошибка регистрации")
}

func (api *UserAPI) HandlerAuthorization(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
}