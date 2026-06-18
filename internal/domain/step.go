package domain

import "fmt"

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

	var victim Player
	var err error
	phase := g.CurrentStep

	switch phase {
	case "wolfAttack":
		victim, err = g.KillRandomAliveVillager()
		if err != nil {
			return StepResult{}, err
		}
		g.CurrentStep = "DayVote"
		g.Night = false
	case "DayVote":
		victim, err = g.KillRandomAlivePlayer()
		if err != nil {
			return StepResult{}, err
		}
		g.CurrentStep = "wolfAttack"
		g.Night = true
	default:
		return StepResult{}, fmt.Errorf("invalid step: %s", phase)
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
