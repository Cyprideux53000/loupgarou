package domain

import "testing"

func TestStepWolfAttack(t *testing.T) {
	game := newTestGame()

	result, err := game.Step(nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if result.Phase != "wolfAttack" {
		t.Errorf("expected phase wolfAttack, got %q", result.Phase)
	}
	if result.Victim.Role == Wolf {
		t.Error("wolf attack should not kill a wolf")
	}
	if result.Victim.Alive {
		t.Error("victim should be dead")
	}
	if result.Message == "" {
		t.Error("message should not be empty")
	}
	if game.CurrentStep != "DayVote" {
		t.Errorf("expected next step DayVote, got %q", game.CurrentStep)
	}
	if game.Night {
		t.Error("should be day after wolf attack")
	}
}

func TestStepDayVote(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "Alice", Role: Villager, Alive: true, Mayor: true},
			{Id: "2", Name: "Bob", Role: Villager, Alive: true},
			{Id: "3", Name: "Charlie", Role: Wolf, Alive: true},
			{Id: "4", Name: "Diana", Role: Villager, Alive: true},
		},
		WolfNumber:  1,
		Night:       false,
		CurrentStep: "DayVote",
	}

	result, err := game.Step(nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if result.Phase != "DayVote" {
		t.Errorf("expected phase DayVote, got %q", result.Phase)
	}
	if result.Victim.Alive {
		t.Error("victim should be dead")
	}
	if game.CurrentStep != "wolfAttack" {
		t.Errorf("expected next step wolfAttack, got %q", game.CurrentStep)
	}
	if !game.Night {
		t.Error("should be night after day vote")
	}
}

func TestStepGameAlreadyOver(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "A", Role: Villager, Alive: true},
			{Id: "2", Name: "B", Role: Wolf, Alive: false},
		},
		WolfNumber:  1,
		CurrentStep: "wolfAttack",
		Night:       true,
	}

	_, err := game.Step(nil)
	if err == nil {
		t.Error("expected error when game is already over")
	}
}

func TestStepInvalidPhase(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "A", Role: Villager, Alive: true},
			{Id: "2", Name: "B", Role: Villager, Alive: true},
			{Id: "3", Name: "C", Role: Wolf, Alive: true},
		},
		WolfNumber:  1,
		CurrentStep: "invalidPhase",
	}

	_, err := game.Step(nil)
	if err == nil {
		t.Error("expected error for invalid step")
	}
}

func TestStepMayorKilledReassigned(t *testing.T) {
	game := &Game{
		Id: "test",
		Players: []Player{
			{Id: "1", Name: "Alice", Role: Villager, Alive: true, Mayor: true},
			{Id: "2", Name: "Bob", Role: Wolf, Alive: true},
			{Id: "3", Name: "Charlie", Role: Villager, Alive: true},
			{Id: "4", Name: "Diana", Role: Villager, Alive: true},
		},
		WolfNumber:  1,
		Night:       false,
		CurrentStep: "DayVote",
	}

	for i := 0; i < 50; i++ {
		testGame := &Game{
			Id:          game.Id,
			Players:     make([]Player, len(game.Players)),
			WolfNumber:  game.WolfNumber,
			Night:       game.Night,
			CurrentStep: game.CurrentStep,
		}
		copy(testGame.Players, game.Players)

		result, err := testGame.Step(nil)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if result.Victim.Mayor && result.NewMayor == nil {
			t.Error("new mayor should be assigned when mayor is killed")
		}

		if result.Victim.Mayor && result.NewMayor != nil {
			if !result.NewMayor.Alive {
				t.Error("new mayor should be alive")
			}
		}
	}
}
