package handler

import (
	"encoding/json"
	"log"
	"loupgarou/internal/service"
	"net/http"
)

type GameHandler struct {
	service service.GameService
}

func NewGameHandler(service service.GameService) *GameHandler {
	return &GameHandler{service: service}
}

type createGameRequest struct {
	Names     []string `json:"names"`
	WolfCount int      `json:"wolf_count"`
	Mode      string   `json:"mode"`
}

func (h *GameHandler) HandleGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var req createGameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		game, err := h.service.CreateGame(req.Names, req.WolfCount, req.Mode)
		if err != nil {
			log.Printf("[ERROR] POST /game - %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("[GAME] New game created | id=%s players=%d wolves=%d", game.Id, len(game.Players), game.WolfNumber)
		if err := json.NewEncoder(w).Encode(game); err != nil {
			log.Printf("[ERROR] POST /game - encode failed: %s", err)
		}

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id query parameter", http.StatusBadRequest)
			return
		}

		game, err := h.service.GetGame(id)
		if err != nil {
			log.Printf("[ERROR] GET /game - game not found: %s", id)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Printf("[GAME] Game loaded | id=%s", id)
		if err := json.NewEncoder(w).Encode(game); err != nil {
			log.Printf("[ERROR] GET /game - encode failed: %s", err)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GameHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id query parameter", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetStatus(id)
	if err != nil {
		log.Printf("[ERROR] GET /status - game not found: %s", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("[STATUS] id=%s | wolves=%d villagers=%d next=%s game_over=%v",
		id, status.WolvesAlive, status.VillagersAlive, status.NextStep, status.IsGameOver)
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("[ERROR] GET /status - encode failed: %s", err)
	}
}

func (h *GameHandler) HandleStep(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id query parameter", http.StatusBadRequest)
		return
	}

	response, err := h.service.ExecuteStep(id)
	if err != nil {
		log.Printf("[ERROR] POST /step - %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[STEP] id=%s | phase=%s victim=%s", id, response.Step.Phase, response.Step.Victim.Name)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[ERROR] POST /step - encode failed: %s", err)
	}
}
