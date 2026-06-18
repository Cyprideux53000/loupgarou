package domain

import "testing"

func TestGetStatusGameInProgress(t *testing.T) {
	game := newTestGame()

	status := game.GetStatus()

	if status.WolvesAlive != 1 {
		t.Errorf("expected 1 wolf alive, got %d", status.WolvesAlive)
	}
	if status.VillagersAlive != 2 {
		t.Errorf("expected 2 villagers alive, got %d", status.VillagersAlive)
	}
	if status.IsGameOver {
		t.Error("game should not be over")
	}
	if status.Winner != "" {
		t.Errorf("expected no winner, got %q", status.Winner)
	}
	if status.NextStep != "wolves_attack" {
		t.Errorf("expected wolves_attack, got %q", status.NextStep)
	}
}

func TestGetStatusDayPhase(t *testing.T) {
	game := newTestGame()
	game.Night = false

	status := game.GetStatus()

	if status.NextStep != "village_vote" {
		t.Errorf("expected village_vote, got %q", status.NextStep)
	}
}

func TestGetStatusVillagersWin(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "A", Role: Villager, Alive: true},
			{Id: "2", Name: "B", Role: Villager, Alive: true},
			{Id: "3", Name: "C", Role: Wolf, Alive: false},
		},
		WolfNumber: 1,
	}

	status := game.GetStatus()

	if !status.IsGameOver {
		t.Error("game should be over")
	}
	if status.Winner != "Villagers" {
		t.Errorf("expected Villagers winner, got %q", status.Winner)
	}
	if status.NextStep != "" {
		t.Errorf("expected empty next step, got %q", status.NextStep)
	}
}

func TestGetStatusWolvesWin(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "A", Role: Villager, Alive: true},
			{Id: "2", Name: "B", Role: Wolf, Alive: true},
			{Id: "3", Name: "C", Role: Villager, Alive: false},
		},
		WolfNumber: 1,
	}

	status := game.GetStatus()

	if !status.IsGameOver {
		t.Error("game should be over (wolves >= villagers)")
	}
	if status.Winner != "Wolves" {
		t.Errorf("expected Wolves winner, got %q", status.Winner)
	}
}

func TestGetStatusWolvesWinMajority(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "A", Role: Villager, Alive: true},
			{Id: "2", Name: "B", Role: Wolf, Alive: true},
			{Id: "3", Name: "C", Role: Wolf, Alive: true},
		},
		WolfNumber: 2,
	}

	status := game.GetStatus()

	if !status.IsGameOver {
		t.Error("game should be over (wolves > villagers)")
	}
	if status.Winner != "Wolves" {
		t.Errorf("expected Wolves winner, got %q", status.Winner)
	}
}
