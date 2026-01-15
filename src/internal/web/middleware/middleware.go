package middleware

import (
	"context"
	"net/http"
	"strings"
	jwtService "tic-tac-toe/internal/service/jwt_service"

	"github.com/google/uuid"
)

// необходим, когда ваш фронтенд и бэкенд работают на разных портах/доменах
func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// автоматически устанавливает заголовок для всех ответов
func ContentTypeJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next(w, r)
	}
}

// позволяет комбинировать несколько middleware в цепочку,
// которая выполняется последовательно.
func Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

//Запрос → EnableCORS → ContentTypeJSON → Logging → Handler
//Ответ  ← EnableCORS ← ContentTypeJSON ← Logging ← Handler

// ключ для хранения UUID пользователя в контексте
type contextKey string

const UserIDKey contextKey = "userID"

// для проверки авторизации пользователя
// Ожидает заголовок Authorization с UUID пользователя
func MiddlewareAuth(jwt jwtService.JwtProvider) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodOptions {
				next(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			claims, err := jwt.ValidateAccessToken(token)
			if err != nil {
				http.Error(w, "Invalid base64", http.StatusUnauthorized)
				return
			}

			// Добавляем UUID пользователя в контекст
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			r = r.WithContext(ctx)

			next(w, r)
		}
	}
}

// для извлечения UUID пользователя из контекста
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}
