package domain

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type Game struct {
	Id          string   `json:"id"`
	Players     []Player `json:"players"`
	WolfNumber  int      `json:"wolf_number"`
	Night       bool     `json:"night"`
	CurrentStep string   `json:"current_step"`
}

type Player struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
	Trait Trait   `json:"trait"`
	Alive bool   `json:"alive"`
	Mayor bool   `json:"mayor"`
}

func NewGame(names []string, wolfCount int) Game {
	roles := make([]Role, len(names))
	for i := 0; i < wolfCount && i < len(roles); i++ {
		roles[i] = Wolf
	}
	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	mayorIndex := rand.Intn(len(names))

	players := make([]Player, len(names))
	for i, name := range names {
		players[i] = Player{
			Id:    uuid.New().String(),
			Name:  name,
			Role:  roles[i],
			Trait: Trait(rand.Intn(5)),
			Alive: true,
			Mayor: i == mayorIndex,
		}
	}

	return Game{
		Id:          uuid.New().String(),
		Players:     players,
		WolfNumber:  wolfCount,
		Night:       true,
		CurrentStep: "wolfAttack",
	}
}

func (g *Game) setMayor(idx int) {
	for i := range g.Players {
		g.Players[i].Mayor = false
	}
	g.Players[idx].Mayor = true
}

func (g *Game) AssignNewMayor() {
	var candidates []int
	for i, player := range g.Players {
		if player.Alive {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) == 0 {
		return
	}
	g.setMayor(candidates[rand.Intn(len(candidates))])
}

func (g *Game) mayorBreaksTie(tied []int) int {
	for _, p := range g.Players {
		if p.Mayor && p.Alive {
			return tied[rand.Intn(len(tied))]
		}
	}
	return tied[rand.Intn(len(tied))]
}

func (g *Game) KillRandomAliveVillager() (Player, error) {
	var villagers []int
	for i, p := range g.Players {
		if p.Alive && p.Role != Wolf {
			villagers = append(villagers, i)
		}
	}
	if len(villagers) == 0 {
		return Player{}, fmt.Errorf("no alive villager to vote")
	}

	votes := make(map[int]int)
	for i, p := range g.Players {
		if !p.Alive {
			continue
		}
		eligible := make([]int, 0)
		for _, idx := range villagers {
			if idx != i {
				eligible = append(eligible, idx)
			}
		}
		if len(eligible) == 0 {
			continue
		}
		votes[eligible[rand.Intn(len(eligible))]]++
	}

	maxVotes := -1
	var tied []int
	for idx, count := range votes {
		if count > maxVotes {
			maxVotes = count
			tied = []int{idx}
		} else if count == maxVotes {
			tied = append(tied, idx)
		}
	}

	var targetIdx int
	if len(tied) == 1 {
		targetIdx = tied[0]
	} else {
		targetIdx = g.mayorBreaksTie(tied)
	}

	g.Players[targetIdx].Alive = false
	return g.Players[targetIdx], nil
}

func (g *Game) KillRandomAlivePlayer() (Player, error) {
	var candidates []int
	for i, player := range g.Players {
		if player.Alive {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) == 0 {
		return Player{}, fmt.Errorf("no alive player to kill")
	}
	idx := candidates[rand.Intn(len(candidates))]
	g.Players[idx].Alive = false
	return g.Players[idx], nil
}

func (p Player) String() string {
	return fmt.Sprintf("Player: %s, Role: %s, Trait: %s", p.Name, p.Role, p.Trait)
}
