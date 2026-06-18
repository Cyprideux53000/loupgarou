package service

import (
	"fmt"
	"loupgarou/internal/domain"
	"loupgarou/internal/storage"
)

type GameService interface {
	CreateGame(names []string, wolfCount int) (*domain.Game, error)
	GetGame(id string) (*domain.Game, error)
	GetStatus(id string) (*domain.Status, error)
	ExecuteStep(id string) (*domain.StepResponse, error)
}

type gameService struct {
	repo storage.GameRepository
}

func NewGameService(repo storage.GameRepository) GameService {
	return &gameService{repo: repo}
}

func (s *gameService) CreateGame(names []string, wolfCount int) (*domain.Game, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("at least one player is required")
	}
	if wolfCount <= 0 {
		return nil, fmt.Errorf("at least one wolf is required")
	}
	if wolfCount >= len(names) {
		return nil, fmt.Errorf("wolf count must be less than player count")
	}

	game := domain.NewGame(names, wolfCount)
	if err := s.repo.Save(&game); err != nil {
		return nil, fmt.Errorf("failed to save game: %w", err)
	}
	return &game, nil
}

func (s *gameService) GetGame(id string) (*domain.Game, error) {
	return s.repo.Load(id)
}

func (s *gameService) GetStatus(id string) (*domain.Status, error) {
	game, err := s.repo.Load(id)
	if err != nil {
		return nil, err
	}
	status := game.GetStatus()
	return &status, nil
}

func (s *gameService) ExecuteStep(id string) (*domain.StepResponse, error) {
	game, err := s.repo.Load(id)
	if err != nil {
		return nil, err
	}

	stepResult, err := game.Step()
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(game); err != nil {
		return nil, fmt.Errorf("failed to save game: %w", err)
	}

	return &domain.StepResponse{
		Game: game,
		Step: stepResult,
	}, nil
}
