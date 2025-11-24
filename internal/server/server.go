package server

import (
	"context"
	"log"
	"net/http"
	"tic-tac-toe/api"
	"tic-tac-toe/internal/config"
	"tic-tac-toe/internal/web/middleware"
)

type Server struct {
	config  *config.Config
	gameAPI *api.GameAPI
}

func NewServer(conf *config.Config, api *api.GameAPI) *Server {
	return &Server{
		config:  conf,
		gameAPI: api,
	}
}

func (s *Server) Start() error {
	gameMoveHandler := middleware.Chain(
		s.gameAPI.HandlerMakeMove,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
	)
	gameNewHandler := middleware.Chain(
		s.gameAPI.HandlerNewGame,
		middleware.EnableCORS,
		middleware.ContentTypeJSON,
	)

	http.HandleFunc("/game/new", gameNewHandler)
	http.HandleFunc("/game/", gameMoveHandler)

	log.Printf("Server starting on port %s", s.config.ServerPort)
	return http.ListenAndServe(":"+s.config.ServerPort, nil)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Завершение работы сервера...")
	return nil
}
