package service

import (
	"fmt"
	"loupgarou/internal/domain"
	"testing"
)

type mockRepo struct {
	games map[string]*domain.Game
}

func newMockRepo() *mockRepo {
	return &mockRepo{games: make(map[string]*domain.Game)}
}

func (m *mockRepo) Save(game *domain.Game) error {
	m.games[game.Id] = game
	return nil
}

func (m *mockRepo) Load(id string) (*domain.Game, error) {
	game, ok := m.games[id]
	if !ok {
		return nil, fmt.Errorf("game not found: %s", id)
	}
	return game, nil
}

func TestCreateGame(t *testing.T) {
	svc := NewGameService(newMockRepo())

	game, err := svc.CreateGame([]string{"Alice", "Bob", "Charlie"}, 1, "random")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if game.Id == "" {
		t.Error("game should have an id")
	}
	if len(game.Players) != 3 {
		t.Errorf("expected 3 players, got %d", len(game.Players))
	}
	if game.WolfNumber != 1 {
		t.Errorf("expected 1 wolf, got %d", game.WolfNumber)
	}
}

func TestCreateGameNoPlayers(t *testing.T) {
	svc := NewGameService(newMockRepo())

	_, err := svc.CreateGame([]string{}, 1, "random")
	if err == nil {
		t.Error("expected error with no players")
	}
}

func TestCreateGameZeroWolves(t *testing.T) {
	svc := NewGameService(newMockRepo())

	_, err := svc.CreateGame([]string{"Alice", "Bob"}, 0, "random")
	if err == nil {
		t.Error("expected error with zero wolves")
	}
}

func TestCreateGameTooManyWolves(t *testing.T) {
	svc := NewGameService(newMockRepo())

	_, err := svc.CreateGame([]string{"Alice", "Bob"}, 2, "random")
	if err == nil {
		t.Error("expected error when wolf count >= player count")
	}
}

func TestGetGame(t *testing.T) {
	repo := newMockRepo()
	svc := NewGameService(repo)

	created, _ := svc.CreateGame([]string{"Alice", "Bob", "Charlie"}, 1, "random")

	loaded, err := svc.GetGame(created.Id)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if loaded.Id != created.Id {
		t.Errorf("expected id %q, got %q", created.Id, loaded.Id)
	}
}

func TestGetGameNotFound(t *testing.T) {
	svc := NewGameService(newMockRepo())

	_, err := svc.GetGame("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent game")
	}
}

func TestGetStatus(t *testing.T) {
	repo := newMockRepo()
	svc := NewGameService(repo)

	game, _ := svc.CreateGame([]string{"Alice", "Bob", "Charlie"}, 1, "random")

	status, err := svc.GetStatus(game.Id)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if status.WolvesAlive != 1 {
		t.Errorf("expected 1 wolf alive, got %d", status.WolvesAlive)
	}
	if status.VillagersAlive != 2 {
		t.Errorf("expected 2 villagers alive, got %d", status.VillagersAlive)
	}
	if status.IsGameOver {
		t.Error("game should not be over")
	}
}

func TestGetStatusNotFound(t *testing.T) {
	svc := NewGameService(newMockRepo())

	_, err := svc.GetStatus("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent game")
	}
}

func TestExecuteStep(t *testing.T) {
	repo := newMockRepo()
	svc := NewGameService(repo)

	game, _ := svc.CreateGame([]string{"Alice", "Bob", "Charlie", "Diana"}, 1, "random")

	response, err := svc.ExecuteStep(game.Id)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if response.Step.Phase != "wolfAttack" {
		t.Errorf("expected wolfAttack phase, got %q", response.Step.Phase)
	}
	if response.Step.Victim.Alive {
		t.Error("victim should be dead")
	}
	if response.Game.CurrentStep != "DayVote" {
		t.Errorf("expected DayVote next, got %q", response.Game.CurrentStep)
	}
}

func TestExecuteStepNotFound(t *testing.T) {
	svc := NewGameService(newMockRepo())

	_, err := svc.ExecuteStep("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent game")
	}
}

func TestExecuteStepSavesState(t *testing.T) {
	repo := newMockRepo()
	svc := NewGameService(repo)

	game, _ := svc.CreateGame([]string{"Alice", "Bob", "Charlie", "Diana"}, 1, "random")

	if _, err := svc.ExecuteStep(game.Id); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	reloaded, _ := svc.GetGame(game.Id)
	if reloaded.CurrentStep != "DayVote" {
		t.Errorf("expected persisted state DayVote, got %q", reloaded.CurrentStep)
	}
}
