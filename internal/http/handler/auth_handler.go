package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"tic-tac-toe/internal/auth/models"
	"tic-tac-toe/internal/auth/repository"
	auth "tic-tac-toe/internal/auth/services"
	"tic-tac-toe/internal/http/mappers"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserAPI struct {
	userServis auth.UserService
	userRepo   repository.UserRepository
}

func NewUserAPI(servis auth.UserService, repo repository.UserRepository) *UserAPI {
	return &UserAPI{
		userServis: servis,
		userRepo:   repo,
	}
}

func (api *UserAPI) HandlerRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// authHeader := r.Header.Get("Registration")
	// if authHeader == "" {
	// 	http.Error(w, "Требуется заголовок регистрации", http.StatusUnauthorized)
	// }
	ctx := r.Context()
	var req models.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}
	user, err := api.userServis.Registration(ctx, &req)
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

func (api *UserAPI) handlerRegistrationError(w http.ResponseWriter, err error) {
	errMsg := err.Error()
	log.Printf("Ошибка регистрации: %v", err)

	switch {
	case strings.Contains(errMsg, "Пользователь с таким логином существует"):
		http.Error(w, "Пользователь с таким логином существует", http.StatusBadRequest)
	case strings.Contains(errMsg, "Ошибка валидации"):
		http.Error(w, "Ошибка валидации", http.StatusBadRequest)
	case strings.Contains(errMsg, "Ошибка хэширования пароля"):
		http.Error(w, "Ошибка хэширования пароля", http.StatusInternalServerError)
	default:
		http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
	}
}

func (api *UserAPI) HandlerAuthorization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	login, password, err := auth.ParseCredintials(authHeader)
	if err != nil {
		http.Error(w, "Invalid credits", http.StatusUnauthorized)
		return
	}

	userID, err := api.userServis.Authenticate(r.Context(), login, password)
	if err != nil {
		http.Error(w, "Invalid credits", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		UUID string `json:"uuid"`
	}{UUID: userID.String()})
}

// ///////////////
func (api *UserAPI) userUUIDFromPath(path string) (uuid.UUID, error) {
	// Парсим путь: /user/{uuid}
	cleanPath := strings.TrimPrefix(path, "/user/")
	if cleanPath == path {
		// Если префикс не найден, пробуем другой способ
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) < 2 {
			log.Printf("Invalid path format: %s", path)
			return uuid.Nil, fmt.Errorf("Invalid UUID")
		}
		cleanPath = parts[1]
	}

	if idx := strings.Index(cleanPath, "/"); idx != -1 {
		cleanPath = cleanPath[:idx]
	}

	parseUUID, err := uuid.Parse(cleanPath)
	if err != nil {
		log.Printf("Error parsing UUID from path %w", err)
		return uuid.Nil, fmt.Errorf("Invalid format UUID")
	}

	return parseUUID, nil
}

// получение информации о игроке
func (api *UserAPI) HandlerGetUserUUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	userUUID, err := api.userUUIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	user, err := api.userRepo.GetUserByID(ctx, userUUID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || user == nil {
			http.Error(w, "Game not found", http.StatusNotFound)
		} else {
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := mappers.UserFromDomainToWeb(user)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
