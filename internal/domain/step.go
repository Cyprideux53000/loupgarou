package domain

import (
	"fmt"
	"loupgarou/internal/llm"
	"math/rand"
)

type TargetChooser interface {
	ChooseTarget(candidateNames []string, phase string) (string, error)
}

type StepResult struct {
	Victim   Player  `json:"victim"`
	Phase    string  `json:"phase"`
	Message  string  `json:"message"`
	NewMayor *Player `json:"new_mayor,omitempty"`
}

type StepResponse struct {
	Game *Game      `json:"game"`
	Step StepResult `json:"step"`
}

func (g *Game) Step() (StepResult, error) {
	if g.GetStatus().IsGameOver {
		return StepResult{}, fmt.Errorf("game is already over")
	}

	phase := g.CurrentStep

	var candidates []Player
	switch phase {
	case "wolfAttack":
		candidates = g.AliveVillagers()
	case "DayVote":
		candidates = g.AlivePlayers()
	default:
		return StepResult{}, fmt.Errorf("invalid step: %s", phase)
	}

	if len(candidates) == 0 {
		return StepResult{}, fmt.Errorf("no eligible target for %s", phase)
	}

	var target Player
	if g.Mode == "llm" {
		names := make([]string, len(candidates))
		for i, c := range candidates {
			names[i] = c.Name
		}
		chooser := llm.New("llama3.1")
		chosenName, err := chooser.ChooseTarget(names, phase)
		if err != nil {
			return StepResult{}, fmt.Errorf("choose target: %w", err)
		}
		found := false
		for _, c := range candidates {
			if c.Name == chosenName {
				target = c
				found = true
				break
			}
		}
		if !found {
			return StepResult{}, fmt.Errorf("chosen target %q not found in candidates", chosenName)
		}
	} else {
		target = candidates[rand.Intn(len(candidates))]
	}

	victim, err := g.KillPlayer(target.Id)
	if err != nil {
		return StepResult{}, err
	}

	switch phase {
	case "wolfAttack":
		g.CurrentStep = "DayVote"
		g.Night = false
	case "DayVote":
		g.CurrentStep = "wolfAttack"
		g.Night = true
	}

	result := StepResult{
		Victim:  victim,
		Phase:   phase,
		Message: fmt.Sprintf("%s was eliminated during the %s", victim.Name, phase),
	}

	if victim.Mayor {
		g.AssignNewMayor()
		for i := range g.Players {
			if g.Players[i].Mayor {
				result.NewMayor = &g.Players[i]
				break
			}
		}
	}

	return result, nil
}
