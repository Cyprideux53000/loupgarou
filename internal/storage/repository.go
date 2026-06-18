package storage

import "loupgarou/internal/domain"

type GameRepository interface {
	Save(game *domain.Game) error
	Load(id string) (*domain.Game, error)
}
