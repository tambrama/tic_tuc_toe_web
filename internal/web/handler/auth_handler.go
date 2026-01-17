package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	auth "tic-tac-toe/internal/service/auth_service"
	jwt "tic-tac-toe/internal/service/jwt_service"
	user "tic-tac-toe/internal/service/user_service"
	dto "tic-tac-toe/internal/web/dto"
	mappers "tic-tac-toe/internal/web/mappers"
	"tic-tac-toe/internal/web/middleware"

	uuid "github.com/google/uuid"
)

type AuthAPI struct {
	authServis auth.AuthService
	userServis user.UserService
	jwt        jwt.JwtProvider
}

func NewAuthAPI(servis auth.AuthService, user user.UserService, jwt jwt.JwtProvider) *AuthAPI {
	return &AuthAPI{
		authServis: servis,
		userServis: user,
		jwt:        jwt,
	}
}

func (api *AuthAPI) HandlerRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var req dto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	user, err := api.authServis.Registration(ctx, req)
	if err != nil {
		log.Printf("Registration error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := mappers.UserFromDomainToWeb(user)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

}

func (api *AuthAPI) HandlerAuthorization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.JwtRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	user, err := api.authServis.Login(r.Context(), req)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// получение информации о игроке по uuid
func (api *AuthAPI) HandlerGetUserUUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	userUUID, err := api.userUUIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := api.userServis.GetByID(ctx, userUUID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := mappers.UserFromDomainToWeb(user)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// получение пользователя по токену
func (api *AuthAPI) HandlerGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	userUUID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := api.userServis.GetByID(ctx, userUUID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := mappers.UserFromDomainToWeb(user)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// обновление аксес токена
func (api *AuthAPI) HandlerRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req dto.RefreshJwtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.RefreshToken == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	token, err := api.authServis.RefreshAccess(ctx, req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(token); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// /////////////////////////////////////////////////////////////////////
func (api *AuthAPI) userUUIDFromPath(path string) (uuid.UUID, error) {
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

	return uuid.Parse(cleanPath)
}
