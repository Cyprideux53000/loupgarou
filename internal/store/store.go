package store

import (
	"encoding/json"
	"fmt"
	"loupgarou/internal/village"
	"os"
)

const dir = "games"

func Save(game village.Game) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create games directory: %w", err)
	}
	data, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("failed to encode game: %w", err)
	}
	return os.WriteFile(dir+"/"+game.Id+".json", data, 0644)
}

func Load(id string) (village.Game, error) {
	data, err := os.ReadFile(dir + "/" + id + ".json")
	if err != nil {
		return village.Game{}, fmt.Errorf("game not found: %w", err)
	}
	var game village.Game
	if err := json.Unmarshal(data, &game); err != nil {
		return village.Game{}, fmt.Errorf("failed to decode game: %w", err)
	}
	return game, nil
}
