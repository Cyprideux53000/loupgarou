package village

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

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
	Trait Trait  `json:"trait"`
	Alive bool   `json:"alive"`
	Mayor bool   `json:"mayor"`
}

type PlayerRequest struct {
	Names     []string `json:"names"`
	WolfCount int      `json:"wolf_count"`
}

func NewGameWithPlayers(players []Player, wolfNumber int) Game {
	return Game{
		Id:          uuid.New().String(),
		Players:     players,
		WolfNumber:  wolfNumber,
		Night:       true,
		CurrentStep: "wolfAttack",
	}
}

func NewPlayer(name string, role Role, trait Trait, mayor bool) Player {
	return Player{
		Id:    uuid.New().String(),
		Name:  name,
		Role:  role,
		Trait: trait,
		Alive: true,
		Mayor: mayor,
	}
}

func NewGame(r *http.Request) Game {
	var req PlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return Game{}
	}

	names := req.Names
	if len(names) == 0 {
		return Game{}
	}

	roles := make([]Role, len(names))
	for i := 0; i < req.WolfCount && i < len(roles); i++ {
		roles[i] = Wolf
	}
	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	mayorIndex := rand.Intn(len(names))

	players := make([]Player, len(names))
	for i, name := range names {
		players[i] = NewPlayer(name, roles[i], Trait(rand.Intn(5)), i == mayorIndex)
	}

	newGame := NewGameWithPlayers(players, req.WolfCount)
	log.Printf("[GAME] Initialized | id=%s players=%d wolves=%d", newGame.Id, len(newGame.Players), newGame.WolfNumber)
	for _, player := range newGame.Players {
		log.Printf("[GAME] Player | name=%s role=%s trait=%s", player.Name, player.Role, player.Trait)
	}

	return newGame
}

func (g *Game) setMayor(idx int) {
	for i := range g.Players {
		g.Players[i].Mayor = false
	}
	g.Players[idx].Mayor = true
	log.Printf("[GAME] Mayor set | id=%s name=%s", g.Id, g.Players[idx].Name)
}

func (g *Game) assignNewMayor() {
	var candidates []int
	for i, player := range g.Players {
		if player.Alive {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) == 0 {
		log.Printf("[GAME] No alive players to assign as mayor | id=%s", g.Id)
		return
	}
	g.setMayor(candidates[rand.Intn(len(candidates))])
}

func (g *Game) killRandomAliveVillager() (Player, error) {
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
		log.Printf("[VOTE] %s votes against %s", p.Name, g.Players[eligible[rand.Intn(len(eligible))]].Name)
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
	log.Printf("[VOTE] Result | eliminated=%s votes=%d", g.Players[targetIdx].Name, maxVotes)
	return g.Players[targetIdx], nil
}

func (g *Game) mayorBreaksTie(tied []int) int {
	for _, p := range g.Players {
		if p.Mayor && p.Alive {
			choice := tied[rand.Intn(len(tied))]
			log.Printf("[VOTE] Tie broken by mayor | mayor=%s eliminates=%s", p.Name, g.Players[choice].Name)
			return choice
		}
	}
	return tied[rand.Intn(len(tied))]
}

func (g *Game) killRandomAlivePlayer() (Player, error) {
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
