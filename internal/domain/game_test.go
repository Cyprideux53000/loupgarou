package domain

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	names := []string{"Alice", "Bob", "Charlie", "Diana"}
	game := NewGame(names, 1)

	if game.Id == "" {
		t.Error("game should have an id")
	}
	if len(game.Players) != 4 {
		t.Errorf("expected 4 players, got %d", len(game.Players))
	}
	if game.WolfNumber != 1 {
		t.Errorf("expected 1 wolf, got %d", game.WolfNumber)
	}
	if !game.Night {
		t.Error("game should start at night")
	}
	if game.CurrentStep != "wolfAttack" {
		t.Errorf("expected wolfAttack, got %s", game.CurrentStep)
	}

	wolfCount := 0
	mayorCount := 0
	for _, p := range game.Players {
		if !p.Alive {
			t.Errorf("player %s should be alive", p.Name)
		}
		if p.Id == "" {
			t.Errorf("player %s should have an id", p.Name)
		}
		if p.Role == Wolf {
			wolfCount++
		}
		if p.Mayor {
			mayorCount++
		}
	}

	if wolfCount != 1 {
		t.Errorf("expected 1 wolf, got %d", wolfCount)
	}
	if mayorCount != 1 {
		t.Errorf("expected 1 mayor, got %d", mayorCount)
	}
}

func TestNewGameMultipleWolves(t *testing.T) {
	game := NewGame([]string{"A", "B", "C", "D", "E"}, 2)

	wolfCount := 0
	for _, p := range game.Players {
		if p.Role == Wolf {
			wolfCount++
		}
	}
	if wolfCount != 2 {
		t.Errorf("expected 2 wolves, got %d", wolfCount)
	}
}

func TestNewGamePlayerNames(t *testing.T) {
	names := []string{"Alice", "Bob"}
	game := NewGame(names, 1)

	found := map[string]bool{}
	for _, p := range game.Players {
		found[p.Name] = true
	}
	for _, name := range names {
		if !found[name] {
			t.Errorf("player %s not found in game", name)
		}
	}
}

func newTestGame() *Game {
	return &Game{
		Id: "test-id",
		Players: []Player{
			{Id: "1", Name: "Alice", Role: Villager, Trait: Brave, Alive: true, Mayor: true},
			{Id: "2", Name: "Bob", Role: Villager, Trait: Cunning, Alive: true, Mayor: false},
			{Id: "3", Name: "Charlie", Role: Wolf, Trait: Aggressive, Alive: true, Mayor: false},
		},
		WolfNumber:  1,
		Night:       true,
		CurrentStep: "wolfAttack",
	}
}

func TestKillRandomAliveVillager(t *testing.T) {
	game := newTestGame()

	victim, err := game.KillRandomAliveVillager()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if victim.Role == Wolf {
		t.Error("victim should not be a wolf")
	}
	if victim.Alive {
		t.Error("victim should be dead")
	}

	deadCount := 0
	for _, p := range game.Players {
		if !p.Alive {
			deadCount++
		}
	}
	if deadCount != 1 {
		t.Errorf("expected 1 dead player, got %d", deadCount)
	}
}

func TestKillRandomAliveVillagerNoVillagers(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "Wolf1", Role: Wolf, Alive: true},
			{Id: "2", Name: "Wolf2", Role: Wolf, Alive: true},
		},
		WolfNumber: 2,
	}

	_, err := game.KillRandomAliveVillager()
	if err == nil {
		t.Error("expected error when no villagers alive")
	}
}

func TestKillRandomAlivePlayer(t *testing.T) {
	game := newTestGame()

	victim, err := game.KillRandomAlivePlayer()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if victim.Alive {
		t.Error("victim should be dead")
	}
}

func TestKillRandomAlivePlayerNobodyAlive(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "A", Alive: false},
			{Id: "2", Name: "B", Alive: false},
		},
	}

	_, err := game.KillRandomAlivePlayer()
	if err == nil {
		t.Error("expected error when no players alive")
	}
}

func TestAssignNewMayor(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "A", Alive: false, Mayor: true},
			{Id: "2", Name: "B", Alive: true, Mayor: false},
			{Id: "3", Name: "C", Alive: true, Mayor: false},
		},
	}

	game.AssignNewMayor()

	mayorCount := 0
	for _, p := range game.Players {
		if p.Mayor {
			mayorCount++
			if !p.Alive {
				t.Error("mayor should be alive")
			}
		}
	}
	if mayorCount != 1 {
		t.Errorf("expected 1 mayor, got %d", mayorCount)
	}
}

func TestAssignNewMayorNobodyAlive(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "A", Alive: false, Mayor: false},
		},
	}

	game.AssignNewMayor()

	for _, p := range game.Players {
		if p.Mayor {
			t.Error("no mayor should be assigned when nobody is alive")
		}
	}
}

func TestPlayerString(t *testing.T) {
	p := Player{Name: "Alice", Role: Villager, Trait: Brave}
	expected := "Player: Alice, Role: Villager, Trait: Brave"
	if p.String() != expected {
		t.Errorf("expected %q, got %q", expected, p.String())
	}
}
