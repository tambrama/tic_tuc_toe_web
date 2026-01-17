package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	model "tic-tac-toe/internal/domain/model/game"
	dsDto "tic-tac-toe/internal/storage/postgres/dto"
	dto "tic-tac-toe/internal/web/dto"
	webMappers "tic-tac-toe/internal/web/mappers"

	service "tic-tac-toe/internal/service/game_service"
	"tic-tac-toe/internal/web/middleware"

	"github.com/google/uuid"
)

type GameAPI struct {
	gameServis service.GameServices
}

func NewGameAPI(servis service.GameServices) *GameAPI {
	return &GameAPI{
		gameServis: servis,
	}
}

// обработчик запроса
func (api *GameAPI) HandlerMakeMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	gameUUID, gameDTO, err := api.validRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Processing move for game UUID: %s", gameUUID)

	// Получаем UUID игрока из контекста
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// // Создаем новое игровое поле из DTO
	newField := &model.GameField{
		Field: gameDTO.Field.Field,
	}

	// Вызываем сервис для обработки хода (вся бизнес-логика там)
	updatedGame, err := api.gameServis.MakeMove(ctx, gameUUID, userID, newField)
	if err != nil {
		switch err {
		case service.ErrNotYourTurn:
			http.Error(w, err.Error(), http.StatusForbidden)
		case service.ErrInvalidMove, service.ErrGameFinished:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error: ", http.StatusInternalServerError)
		}
		return
	}

	// Определяем статус для ответа
	gameStatus := updatedGame.Status
	if gameStatus == model.Playing {
		gameStatus = api.gameServis.CheckEndGame(updatedGame)
	}

	// Формируем ответ
	response := webMappers.CurrentGameFromDomainToWeb(updatedGame, gameStatus)

	// Определяем сообщение
	switch gameStatus {
	case model.Playing:
		if updatedGame.CurrentTurn == updatedGame.PlayerX {
			response.Message = "Player X's turn"
		} else {
			response.Message = "Player O's turn"
		}
	case model.Draw:
		response.Message = "Draw!"
	case model.WonX:
		response.Message = "Player X won!"
	case model.WonO:
		response.Message = "Player O won!"
	default:
		response.Message = "Game ended"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

}

// новая игра
func (api *GameAPI) HandlerNewGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	playerX, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.NewGameRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Вызываем сервис для создания игры (вся бизнес-логика там)
	newGame, err := api.gameServis.CreateNewGame(ctx, playerX, req.WithBot)
	if err != nil {
		http.Error(w, "Failed to create game: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := webMappers.CurrentGameFromDomainToWeb(newGame, model.Playing)
	response.Message = "Game created"

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// все игры в статусе ожидания
func (api *GameAPI) HandlerGetAvailableGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodOptions {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	games, err := api.gameServis.GetAvailableGames(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch available games", http.StatusInternalServerError)
		return
	}
	var response []dto.GameResponse
	for _, game := range games {
		response = append(response, webMappers.CurrentGameFromDomainToWeb(game, game.Status))
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// история игр
func (api *GameAPI) HandlerGetComplitedGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodOptions {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	games, err := api.gameServis.GetComplitedGames(ctx, userID)
	if err != nil {
		http.Error(w, "Failed to fetch completed games", http.StatusInternalServerError)
		return
	}
	var response []dto.GameResponse
	for _, game := range games {
		response = append(response, webMappers.CurrentGameFromDomainToWeb(game, game.Status))
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// присоединение к игре
func (api *GameAPI) HandlerJoinGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	gameUUID, err := api.gameUUIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	playerO, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	gameCurrent, err := api.gameServis.JoinGame(ctx, gameUUID, playerO)
	if err != nil {
		if errors.Is(err, service.ErrGameFull) || errors.Is(err, service.ErrCannotJoinOwn) || errors.Is(err, service.ErrGameNotWaiting) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to join game", http.StatusInternalServerError)
		}
		return
	}

	response := webMappers.CurrentGameFromDomainToWeb(gameCurrent, model.Playing)
	response.Message = "Game started"

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

}

// получение текущей игры
func (api *GameAPI) HandlerGetCurrentGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	gameUUID, err := api.gameUUIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game, err := api.gameServis.GetCurrentGame(ctx, gameUUID)
	if err != nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	response := webMappers.CurrentGameFromDomainToWeb(game, game.Status)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
func (api *GameAPI) HandlerGetLeaderBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CountLeaderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Вызываем сервис для создания игры (вся бизнес-логика там)
	leaderBoard, err := api.gameServis.GetLeaderBoard(ctx, req.Count)
	if err != nil {
		http.Error(w, "Failed to create game: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var leaders []model.UserLeaders
	for _, leader := range leaderBoard {
		leaders = append(leaders, leader)
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(leaders); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// ////////////////////////////////////////////////////////////////////////
func (api *GameAPI) gameUUIDFromPath(path string) (uuid.UUID, error) {
	// Парсим путь: /game/{uuid}
	cleanPath := strings.TrimPrefix(path, "/game/")
	if cleanPath == path {
		// Если префикс не найден, пробуем другой способ
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) < 2 {
			return uuid.Nil, fmt.Errorf("Invalid path format")
		}
		cleanPath = parts[1]
	}

	if idx := strings.Index(cleanPath, "/"); idx != -1 {
		cleanPath = cleanPath[:idx]
	}

	parseUUID, err := uuid.Parse(cleanPath)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid format UUID")
	}

	return parseUUID, nil
}

func (api *GameAPI) validRequest(r *http.Request) (uuid.UUID, dsDto.CurrentGameDTO, error) {
	gameUUID, err := api.gameUUIDFromPath(r.URL.Path)
	if err != nil {
		return uuid.Nil, dsDto.CurrentGameDTO{}, fmt.Errorf("Invalid format UUID")
	}

	var gameDTO dsDto.CurrentGameDTO

	if err := json.NewDecoder(r.Body).Decode(&gameDTO); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return uuid.Nil, gameDTO, fmt.Errorf("Invalid JSON")
	}

	if gameDTO.UUID != gameUUID {
		return uuid.Nil, gameDTO, fmt.Errorf("UUID not equal")
	}
	return gameUUID, gameDTO, nil
}
