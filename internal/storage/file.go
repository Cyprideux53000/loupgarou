package storage

import (
	"encoding/json"
	"fmt"
	"loupgarou/internal/domain"
	"os"
	"path/filepath"
)

type FileGameRepository struct {
	dir string
}

func NewFileGameRepository(dir string) *FileGameRepository {
	return &FileGameRepository{dir: dir}
}

func (r *FileGameRepository) Save(game *domain.Game) error {
	if err := os.MkdirAll(r.dir, 0755); err != nil {
		return fmt.Errorf("failed to create games directory: %w", err)
	}
	data, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("failed to encode game: %w", err)
	}
	return os.WriteFile(filepath.Join(r.dir, game.Id+".json"), data, 0644)
}

func (r *FileGameRepository) Load(id string) (*domain.Game, error) {
	data, err := os.ReadFile(filepath.Join(r.dir, id+".json"))
	if err != nil {
		return nil, fmt.Errorf("game not found: %w", err)
	}
	var game domain.Game
	if err := json.Unmarshal(data, &game); err != nil {
		return nil, fmt.Errorf("failed to decode game: %w", err)
	}
	return &game, nil
}
