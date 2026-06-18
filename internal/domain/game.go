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
	Mode        string   `json:"mode"`
}

type Player struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
	Trait Trait  `json:"trait"`
	Alive bool   `json:"alive"`
	Mayor bool   `json:"mayor"`
}

func NewGame(names []string, wolfCount int, mode string) Game {
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
		Mode:        mode,
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

func (g *Game) AliveVillagers() []Player {
	var result []Player
	for _, p := range g.Players {
		if p.Alive && p.Role != Wolf {
			result = append(result, p)
		}
	}
	return result
}

func (g *Game) AlivePlayers() []Player {
	var result []Player
	for _, p := range g.Players {
		if p.Alive {
			result = append(result, p)
		}
	}
	return result
}

func (g *Game) KillPlayer(id string) (Player, error) {
	for i, p := range g.Players {
		if p.Id == id {
			if !p.Alive {
				return Player{}, fmt.Errorf("player %s is already dead", p.Name)
			}
			g.Players[i].Alive = false
			return g.Players[i], nil
		}
	}
	return Player{}, fmt.Errorf("player not found: %s", id)
}

func (p Player) String() string {
	return fmt.Sprintf("Player: %s, Role: %s, Trait: %s", p.Name, p.Role, p.Trait)
}
