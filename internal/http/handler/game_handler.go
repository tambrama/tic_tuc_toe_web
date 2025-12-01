package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	dsDto "tic-tac-toe/internal/repository/dto"

	uuid "github.com/google/uuid"

	// webDto "tic-tac-toe/internal/web/dto"
	"tic-tac-toe/internal/domain/models"
	"tic-tac-toe/internal/domain/services"
	webMappers "tic-tac-toe/internal/http/mappers"
	"tic-tac-toe/internal/repository"
	dsMappers "tic-tac-toe/internal/repository/mappers"
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

func (api *GameAPI) validRequest(w http.ResponseWriter, r *http.Request) (uuid.UUID, *dsDto.CurrentGameDTO, error) {
	// Парсим путь: /game/{uuid}
	path := strings.TrimPrefix(r.URL.Path, "/game/")
	if path == r.URL.Path {
		// Если префикс не найден, пробуем другой способ
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 2 {
			log.Printf("Invalid path format: %s", r.URL.Path)
			return uuid.Nil, nil, fmt.Errorf("Invalid UUID")
		}
		path = parts[1]
	}

	parseUUID, err := uuid.Parse(path)
	if err != nil {
		log.Printf("Error parsing UUID from path '%s': %v", path, err)
		return uuid.Nil, nil, fmt.Errorf("Invalid format UUID")
	}

	var gameDTO dsDto.CurrentGameDTO

	if err := json.NewDecoder(r.Body).Decode(&gameDTO); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return uuid.Nil, nil, fmt.Errorf("Invalid JSON")
	}

	if gameDTO.UUID.String() != path {
		log.Printf("UUID mismatch: path=%s, body=%s", path, gameDTO.UUID)
		return uuid.Nil, nil, fmt.Errorf("UUID not equal")
	}
	return parseUUID, &gameDTO, nil
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
	currentGame := dsMappers.CurrentGameFromDBToDomain(gameDTO)

	game, err := api.gameRepo.GetCurrentGame(ctx, gameUUID)
	if err != nil {
		log.Printf("Error getting game: %v", err)
		http.Error(w, "Error: game not found "+err.Error(), http.StatusInternalServerError)
		return
	}
	if game != nil {
		if err := api.gameServis.ValidationField(game, currentGame); err != nil {
			log.Printf("Validation error: %v", err)
			http.Error(w, "Невозможность сделать ход, проверьте вводимые данные: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	currentStatus := api.gameServis.CheckEndGame(currentGame)
	if currentStatus != services.InProgress {
		response := webMappers.CurrentGameFromDomainToWeb(currentGame, currentStatus)
		response.Message = "Игра окончена"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	updateGame, err := api.gameServis.GetNextStep(currentGame)
	if err != nil {
		log.Printf("Error getting next step: %v", err)
		switch err.Error() {
		case "Progress end":
			response := webMappers.CurrentGameFromDomainToWeb(currentGame, currentStatus)
			response.Message = "Игра окончена"
			w.WriteHeader(http.StatusOK)
			if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
				log.Printf("Error encoding response: %v", encodeErr)
			}
			return
		case "Game end: draw":
		case "No moves available":
			response := webMappers.CurrentGameFromDomainToWeb(currentGame, currentStatus)
			response.Message = "Игра окончена в ничью"
			w.WriteHeader(http.StatusOK)
			if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
				log.Printf("Error encoding response: %v", encodeErr)
			}
			return
		default:
			http.Error(w, "Ошибка при выполнении хода: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := api.gameRepo.SaveGame(ctx, updateGame); err != nil {
		log.Printf("Error saving game: %v", err)
		http.Error(w, "Ошибка сохранения игры: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newStatus := api.gameServis.CheckEndGame(updateGame)
	response := webMappers.CurrentGameFromDomainToWeb(updateGame, newStatus)

	if newStatus != services.InProgress {
		switch newStatus {
		case services.BotWin:
			response.Message = "Bot - Win"
		case services.UserWin:
			response.Message = "User - Win"
		case services.Draw:
			response.Message = "Draw!"
		}
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

}

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

	newGame := &models.CurrentGame{
		UUID:  newUUID,
		Field: newField,
	}

	log.Printf("Creating new game with UUID: %s", newUUID)

	if err := api.gameRepo.SaveGame(ctx, newGame); err != nil {
		log.Printf("Error saving game: %v", err)
		http.Error(w, "Ошибка создания игры: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := webMappers.CurrentGameFromDomainToWeb(newGame, services.InProgress)
	response.Message = "Игра создана"

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
