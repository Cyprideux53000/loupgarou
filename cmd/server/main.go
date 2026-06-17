package main

import (
	"encoding/json"
	"log"
	"loupgarou/internal/store"
	"loupgarou/internal/village"
	"net/http"
)

func gameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		game := village.NewGame(r)
		if err := store.Save(game); err != nil {
			log.Printf("[ERROR] POST /game - save failed: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("[GAME] New game created | id=%s players=%d wolves=%d", game.Id, len(game.Players), game.WolfNumber)
		json.NewEncoder(w).Encode(game)

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id query parameter", http.StatusBadRequest)
			return
		}
		game, err := store.Load(id)
		if err != nil {
			log.Printf("[ERROR] GET /game - game not found: %s", id)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("[GAME] Game loaded | id=%s", id)
		json.NewEncoder(w).Encode(game)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id query parameter", http.StatusBadRequest)
			return
		}
		game, err := store.Load(id)
		if err != nil {
			log.Printf("[ERROR] GET /status - game not found: %s", id)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		status := game.GetStatus()
		log.Printf("[STATUS] id=%s | wolves=%d villagers=%d next=%s game_over=%v",
			id, status.WolvesAlive, status.VillagersAlive, status.NextStep, status.IsGameOver)
		json.NewEncoder(w).Encode(status)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func stepHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id query parameter", http.StatusBadRequest)
			return
		}
		game, err := store.Load(id)
		if err != nil {
			log.Printf("[ERROR] POST /step - game not found: %s", id)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		stepResult, err := game.Step()
		if err != nil {
			log.Printf("[ERROR] POST /step - step failed: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := store.Save(game); err != nil {
			log.Printf("[ERROR] POST /step - save failed: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("[STEP] id=%s | phase=%s victim=%s next=%s",
			id, stepResult.Phase, stepResult.Victim.Name, game.CurrentStep)
		json.NewEncoder(w).Encode(village.StepResponse{
			Game: game,
			Step: stepResult,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/step", stepHandler)
	log.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
