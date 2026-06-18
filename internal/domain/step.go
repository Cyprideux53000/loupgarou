package domain

import (
	"fmt"
	"loupgarou/internal/llm"
	"math/rand"
)

type TargetChooser interface {
	ChooseTarget(candidates []llm.CandidateInfo, phase string) (string, error)
}

type PlayerVote struct {
	Voter  string `json:"voter"`
	Target string `json:"target"`
}

type StepResult struct {
	Victim     Player       `json:"victim"`
	Phase      string       `json:"phase"`
	Message    string       `json:"message"`
	NewMayor   *Player      `json:"new_mayor,omitempty"`
	Discussion []string     `json:"discussion,omitempty"`
	Votes      []PlayerVote `json:"votes,omitempty"`
}

type StepResponse struct {
	Game *Game      `json:"game"`
	Step StepResult `json:"step"`
}

func (g *Game) Step(userDiscussion []string) (StepResult, error) {
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
	var chooseResult llm.ChooseResult
	if g.Mode == "llm" {
		infos := make([]llm.CandidateInfo, len(candidates))
		for i, c := range candidates {
			infos[i] = llm.CandidateInfo{
				Name:  c.Name,
				Trait: c.Trait.String(),
				Role:  c.Role.String(),
				Mayor: c.Mayor,
			}
		}
		var err error
		chooseResult, err = llm.New("llama3.2:3b").ChooseTarget(infos, phase, userDiscussion)
		if err != nil {
			return StepResult{}, fmt.Errorf("choose target: %w", err)
		}
		found := false
		for _, c := range candidates {
			if c.Name == chooseResult.TargetName {
				target = c
				found = true
				break
			}
		}
		if !found {
			return StepResult{}, fmt.Errorf("chosen target %q not found in candidates", chooseResult.TargetName)
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

	var discussion []string
	var votes []PlayerVote
	if g.Mode == "llm" {
		discussion = chooseResult.Discussion
		for _, v := range chooseResult.Votes {
			votes = append(votes, PlayerVote{Voter: v.Voter, Target: v.Target})
		}
	}

	result := StepResult{
		Victim:     victim,
		Phase:      phase,
		Message:    fmt.Sprintf("%s was eliminated during the %s", victim.Name, phase),
		Discussion: discussion,
		Votes:      votes,
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
