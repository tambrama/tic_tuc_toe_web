package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"tic-tac-toe/internal/config"
	jwt "tic-tac-toe/internal/service/jwt_service"
	"tic-tac-toe/internal/web/handler"
	"tic-tac-toe/internal/web/middleware"
)

type Server struct {
	config  *config.Config
	gameAPI *handler.GameAPI
	userAPI *handler.AuthAPI
	jwt     jwt.JwtProvider
}

func NewServer(conf *config.Config, api *handler.GameAPI, user *handler.AuthAPI, jwt jwt.JwtProvider) *Server {
	return &Server{
		config:  conf,
		gameAPI: api,
		userAPI: user,
		jwt:     jwt,
	}
}

func (s *Server) Start() error {
	// middleware авторизации
	requireAuth := middleware.MiddlewareAuth(s.jwt)

	// без авторизации
	userRegistrationHandler := middleware.Chain(
		s.userAPI.HandlerRegistration,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
	)
	userAuthHandler := middleware.Chain(
		s.userAPI.HandlerAuthorization,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
	)
	refreshTokenHandler := middleware.Chain(
		s.userAPI.HandlerRefreshAccessToken,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
	)

	// с авторизацией
	gameMainHandler := middleware.Chain(
		s.mainHandler,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
		requireAuth,
	)
	gameNewHandler := middleware.Chain(
		s.gameAPI.HandlerNewGame,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
		requireAuth,
	)
	gamesListHandler := middleware.Chain(
		s.gameAPI.HandlerGetAvailableGames,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
		requireAuth,
	)
	userInfoHandler := middleware.Chain(
		s.userAPI.HandlerGetUserUUID,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
		requireAuth,
	)
	getUserHandler := middleware.Chain(
		s.userAPI.HandlerGetCurrentUser,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
		requireAuth,
	)
	getHistoryHandler := middleware.Chain(
		s.gameAPI.HandlerGetComplitedGames,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
		requireAuth,
	)

	getLeadersHandler := middleware.Chain(
		s.gameAPI.HandlerGetLeaderBoard,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
		requireAuth,
	)

	http.HandleFunc("/registration", userRegistrationHandler)
	http.HandleFunc("/auth", userAuthHandler)
	http.HandleFunc("/game/new", gameNewHandler)
	http.HandleFunc("/game/", gameMainHandler)
	http.HandleFunc("/game/list", gamesListHandler)
	http.HandleFunc("/user/", userInfoHandler)
	http.HandleFunc("/auth/refresh", refreshTokenHandler)
	http.HandleFunc("/auth/me", getUserHandler)
	http.HandleFunc("/game/history", getHistoryHandler)
	http.HandleFunc("/game/leaders", getLeadersHandler)

	log.Printf("Server starting on port %s", s.config.ServerPort)
	return http.ListenAndServe(":"+s.config.ServerPort, nil)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Завершение работы сервера...")
	return nil
}

func (s *Server) mainHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	//проверяем  путь
	if strings.HasSuffix(path, "/join") {
		s.gameAPI.HandlerJoinGame(w, r)
		return
	}

	if strings.HasSuffix(path, "/status") {
		s.gameAPI.HandlerGetCurrentGame(w, r)
		return
	}

	if path != "/game/new" && path != "/game/leaders" && path != "/game/history" && path != "/game/list" && !strings.Contains(path, "/join") && !strings.Contains(path, "/status") {
		s.gameAPI.HandlerMakeMove(w, r)
		return
	}

	http.NotFound(w, r)
}
