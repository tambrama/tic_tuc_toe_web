package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	dsDto "tic-tac-toe/internal/repository/dto"
	"tic-tac-toe/internal/repository/mappers"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	// webDto "tic-tac-toe/internal/web/dto"
	"tic-tac-toe/internal/domain/models"
	"tic-tac-toe/internal/domain/services"
	"tic-tac-toe/internal/http/dto"
	webMappers "tic-tac-toe/internal/http/mappers"
	"tic-tac-toe/internal/http/middleware"
	"tic-tac-toe/internal/repository"
)

type GameAPI struct {
	gameServis services.GameServices
	gameRepo   repository.GameRepository
}

func NewGameAPI(servis services.GameServices, repo repository.GameRepository) *GameAPI {
	return &GameAPI{
		gameServis: servis,
		gameRepo:   repo,
	}
}
func (api *GameAPI) gameUUIDFromPath(path string) (uuid.UUID, error) {
	// Парсим путь: /game/{uuid}
	cleanPath := strings.TrimPrefix(path, "/game/")
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

func (api *GameAPI) validRequest(w http.ResponseWriter, r *http.Request) (uuid.UUID, *dsDto.CurrentGameDTO, error) {
	gameUUID, err := api.gameUUIDFromPath(r.URL.Path)
	if err != nil {
		log.Printf("Error parsing UUID from path '%s': %v", r.URL.Path, err)
		return uuid.Nil, nil, fmt.Errorf("Invalid format UUID")
	}

	var gameDTO dsDto.CurrentGameDTO

	if err := json.NewDecoder(r.Body).Decode(&gameDTO); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return uuid.Nil, nil, fmt.Errorf("Invalid JSON")
	}

	if gameDTO.UUID != gameUUID {
		log.Printf("UUID mismatch: path=%s, body=%s", r.URL.Path, gameDTO.UUID)
		return uuid.Nil, nil, fmt.Errorf("UUID not equal")
	}
	return gameUUID, &gameDTO, nil
}

// обработчик запроса
func (api *GameAPI) HandlerMakeMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	gameUUID, gameDTO, err := api.validRequest(w, r)
	if err != nil {
		log.Printf("Error in validRequest: %v", err)
		switch {
		case strings.Contains(err.Error(), "method not allowed"):
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		case strings.Contains(err.Error(), "Invalid UUID"):
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
		case strings.Contains(err.Error(), "Invalid format UUID"):
			http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		case strings.Contains(err.Error(), "Invalid JSON"):
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
		case strings.Contains(err.Error(), "UUID not equal"):
			http.Error(w, "UUID mismatch", http.StatusBadRequest)
		default:
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	log.Printf("Processing move for game UUID: %s", gameUUID)
	//получаем uuid игрока
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	game := mappers.CurrentGameFromDBToDomain(gameDTO)

	//загрузка игры из бд
	currentGame, err := api.gameRepo.GetCurrentGame(ctx, gameUUID)
	if err != nil {
		log.Printf("Error getting game: %v", err)
		http.Error(w, "Error: game not found "+err.Error(), http.StatusInternalServerError)
		return
	}
	//проверяем, что ходит текущий игрок
	if currentGame.CurrentTurn != userID {
		http.Error(w, "Not your turn", http.StatusForbidden)
		return
	}
	//обновляем поле
	currentGame.Field.Field = gameDTO.Field.Field

	if err := api.gameServis.ValidationField(currentGame, game); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, "Невозможность сделать ход, проверьте вводимые данные: "+err.Error(), http.StatusBadRequest)
		return
	}

	currentStatus := api.gameServis.CheckEndGame(currentGame)
	if currentStatus != models.Playing {
		api.sendFinalResponse(w, currentGame, currentStatus)
		return
	}

	///===== Игра с ботом ======
	isBotGame := currentGame.PlayerO == nil
	if isBotGame {
		updateGame, err := api.gameServis.GetNextStep(currentGame)
		if err != nil {
			currentStatus = api.gameServis.CheckEndGame(currentGame)
			api.sendFinalResponse(w, currentGame, currentStatus)
			return
		}
		currentGame = updateGame
		currentStatus = api.gameServis.CheckEndGame(currentGame)
		if currentStatus != models.Playing {
			api.sendFinalResponse(w, currentGame, currentStatus)
			return
		}
		currentGame.CurrentTurn = currentGame.PlayerX
	} else {
		if currentGame.CurrentTurn == currentGame.PlayerX {
			currentGame.CurrentTurn = *currentGame.PlayerO
		} else {
			currentGame.CurrentTurn = currentGame.PlayerX
		}
	}

	if err := api.gameRepo.SaveGame(ctx, currentGame); err != nil {
		log.Printf("Error saving game: %v", err)
		http.Error(w, "Ошибка сохранения игры: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := webMappers.CurrentGameFromDomainToWeb(currentGame, models.Playing)

	if isBotGame {
		response.Message = "Your move + bot responded"
	} else {
		response.Message = "Move completed"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

}

func (api *GameAPI) sendFinalResponse(w http.ResponseWriter, game *models.CurrentGame, status models.GameStatus) {
	response := webMappers.CurrentGameFromDomainToWeb(game, status)
	switch status {
	case models.Draw:
		response.Message = "Draw!"
	case models.WonO:
		response.Message = "Game over"
	default:
		response.Message = "Game ended"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// новая игра
func (api *GameAPI) HandlerNewGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	newUUID := uuid.New()

	newField := &models.GameField{
		Field: [][]int{
			{services.Empty, services.Empty, services.Empty},
			{services.Empty, services.Empty, services.Empty},
			{services.Empty, services.Empty, services.Empty},
		},
	}
	playerX, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Ошибка получения uuid игрока X", http.StatusUnauthorized)
		return
	}

	var req dto.NewGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var status models.GameStatus
	var playerO *uuid.UUID
	if req.WithBot {
		playerO = nil
		status = models.Playing
	} else {
		playerO = nil
		status = models.Waiting
	}
	newGame := &models.CurrentGame{
		UUID:        newUUID,
		Field:       newField,
		Status:      status,
		PlayerX:     playerX,
		PlayerO:     playerO,
		CurrentTurn: playerX,
		Symbols: map[uuid.UUID]models.Char{
			playerX: models.CharX,
			// *playerO: models.CharO,
		},
	}

	log.Printf("Creating new game with UUID: %s", newUUID)

	if err := api.gameRepo.SaveGame(ctx, newGame); err != nil {
		log.Printf("Error saving game: %v", err)
		http.Error(w, "Ошибка создания игры: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := webMappers.CurrentGameFromDomainToWeb(newGame, models.Playing)
	response.Message = "Игра создана"

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
	games, err := api.gameRepo.GetAvailableGames(ctx)
	if err != nil {
		log.Printf("Error fetching available games: %v", err)
		http.Error(w, "Ошибка возарщения списка игр", http.StatusBadRequest)
		return
	}
	var response []dto.GameResponse
	for _, game := range games {
		gameDTO := webMappers.CurrentGameFromDomainToWeb(game, game.Status)
		response = append(response, *gameDTO)
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
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	playerO, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	gameCurrent, err := api.gameRepo.GetCurrentGame(ctx, gameUUID)
	if err != nil {
		log.Printf("Error find game: %v", err)
		http.Error(w, "Игра не найдена", http.StatusBadRequest)
		return
	}

	//валидация
	if gameCurrent.Status != models.Waiting {
		log.Printf("Game status is not Waiting: %d", gameCurrent.Status)
		http.Error(w, "Игра не в статусе ожидания", http.StatusBadRequest)
		return
	}

	if gameCurrent.PlayerO != nil {
		log.Printf("Error player: %v", err)
		http.Error(w, "Second player already joined", http.StatusBadRequest)
		return
	}

	if gameCurrent.PlayerX == playerO {
		log.Printf("Error player X: %v", err)
		http.Error(w, "Игрок X = игроку O", http.StatusBadRequest)
		return
	}
	//обновляем игру
	gameCurrent.Symbols[playerO] = models.CharO
	gameCurrent.Status = models.Playing
	gameCurrent.PlayerO = &playerO
	gameCurrent.CurrentTurn = gameCurrent.PlayerX

	if err := api.gameRepo.SaveGame(ctx, gameCurrent); err != nil {
		log.Printf("Error saving game: %v", err)
		http.Error(w, "Ошибка создания игры: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := webMappers.CurrentGameFromDomainToWeb(gameCurrent, models.Playing)
	response.Message = "Игра начата"

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
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	game, err := api.gameRepo.GetCurrentGame(ctx, gameUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || game == nil {
			http.Error(w, "Game not found", http.StatusNotFound)
		} else {
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	if game == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	response := webMappers.CurrentGameFromDomainToWeb(game, game.Status)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
