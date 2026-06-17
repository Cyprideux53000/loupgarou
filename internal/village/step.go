package village

import "fmt"

// StepResult describes what happened during a game step.
type StepResult struct {
	Victim  Player `json:"victim"`
	Phase   string `json:"phase"`
	Message string `json:"message"`
}

// StepResponse is the HTTP response after advancing a step.
type StepResponse struct {
	Game Game       `json:"game"`
	Step StepResult `json:"step"`
}

// Step advances the game by one phase.
// At night, wolves eliminate a random alive villager.
// During the day, the village votes to eliminate a random alive player.
// Returns an error if the game is already over or no valid target exists.
func (g *Game) Step() (StepResult, error) {
	if g.GetStatus().IsGameOver {
		return StepResult{}, fmt.Errorf("game is already over")
	}

	var victim Player
	var err error
	phase := g.CurrentStep

	if phase == "wolfAttack" {
		victim, err = g.killRandomAliveVillager()
		if err != nil {
			return StepResult{}, err
		}
		g.CurrentStep = "DayVote"
		g.Night = false
	} else if phase == "DayVote" {
		victim, err = g.killRandomAlivePlayer()
		if err != nil {
			return StepResult{}, err
		}
		g.CurrentStep = "wolfAttack"
		g.Night = true
	} else {
		return StepResult{}, fmt.Errorf("invalid step: %s", phase)
	}

	return StepResult{
		Victim:  victim,
		Phase:   phase,
		Message: fmt.Sprintf("%s was eliminated during the %s", victim.Name, phase),
	}, nil
}
