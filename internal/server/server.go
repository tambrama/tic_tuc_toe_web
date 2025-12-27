package server

import (
	"context"
	"log"
	"net/http"
	"strings"
	auth "tic-tac-toe/internal/auth/services"
	"tic-tac-toe/internal/config"
	"tic-tac-toe/internal/http/handler"
	"tic-tac-toe/internal/http/middleware"
)

type Server struct {
	config      *config.Config
	gameAPI     *handler.GameAPI
	userAPI     *handler.UserAPI
	userService auth.UserService
}

func NewServer(conf *config.Config, api *handler.GameAPI, user *handler.UserAPI, userService auth.UserService) *Server {
	return &Server{
		config:      conf,
		gameAPI:     api,
		userAPI:     user,
		userService: userService,
	}
}

func (s *Server) Start() error {
	// middleware авторизации
	requireAuth := middleware.MiddlewareAuth(s.userService)

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

	http.HandleFunc("/registr", userRegistrationHandler)
	http.HandleFunc("/auth", userAuthHandler)
	http.HandleFunc("/game/new", gameNewHandler)
	http.HandleFunc("/game/", gameMainHandler)
	http.HandleFunc("/game/available", gamesListHandler)
	http.HandleFunc("/user/", userInfoHandler)

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

	if path != "/game/new" && path != "/game/available" && !strings.Contains(path, "/join") && !strings.Contains(path, "/status") {
		s.gameAPI.HandlerMakeMove(w, r)
		return
	}

	http.NotFound(w, r)
}
