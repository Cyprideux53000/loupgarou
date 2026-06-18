package storage

import (
	"loupgarou/internal/domain"
	"os"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	repo := NewFileGameRepository(dir)

	game := &domain.Game{
		Id: "test-123",
		Players: []domain.Player{
			{Id: "p1", Name: "Alice", Role: domain.Villager, Trait: domain.Brave, Alive: true, Mayor: true},
			{Id: "p2", Name: "Bob", Role: domain.Wolf, Trait: domain.Cunning, Alive: true, Mayor: false},
		},
		WolfNumber:  1,
		Night:       true,
		CurrentStep: "wolfAttack",
	}

	if err := repo.Save(game); err != nil {
		t.Fatalf("Save failed: %s", err)
	}

	loaded, err := repo.Load("test-123")
	if err != nil {
		t.Fatalf("Load failed: %s", err)
	}

	if loaded.Id != game.Id {
		t.Errorf("expected id %q, got %q", game.Id, loaded.Id)
	}
	if len(loaded.Players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(loaded.Players))
	}
	if loaded.Players[0].Name != "Alice" {
		t.Errorf("expected Alice, got %q", loaded.Players[0].Name)
	}
	if loaded.Players[0].Role != domain.Villager {
		t.Errorf("expected Villager, got %v", loaded.Players[0].Role)
	}
	if loaded.Players[1].Role != domain.Wolf {
		t.Errorf("expected Wolf, got %v", loaded.Players[1].Role)
	}
	if loaded.WolfNumber != 1 {
		t.Errorf("expected wolf_number 1, got %d", loaded.WolfNumber)
	}
	if !loaded.Night {
		t.Error("expected night to be true")
	}
}

func TestLoadNotFound(t *testing.T) {
	dir := t.TempDir()
	repo := NewFileGameRepository(dir)

	_, err := repo.Load("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent game")
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	dir := t.TempDir() + "/subdir/games"
	repo := NewFileGameRepository(dir)

	game := &domain.Game{Id: "test-456", Players: []domain.Player{}}

	if err := repo.Save(game); err != nil {
		t.Fatalf("Save failed: %s", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("directory should have been created")
	}
}

func TestSaveOverwrite(t *testing.T) {
	dir := t.TempDir()
	repo := NewFileGameRepository(dir)

	game := &domain.Game{
		Id: "test-789",
		Players: []domain.Player{
			{Id: "p1", Name: "Alice", Alive: true},
		},
	}

	if err := repo.Save(game); err != nil {
		t.Fatalf("first Save failed: %s", err)
	}

	game.Players[0].Alive = false
	if err := repo.Save(game); err != nil {
		t.Fatalf("second Save failed: %s", err)
	}

	loaded, err := repo.Load("test-789")
	if err != nil {
		t.Fatalf("Load failed: %s", err)
	}
	if loaded.Players[0].Alive {
		t.Error("expected player to be dead after overwrite")
	}
}
