package village

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
)

// Game represents a werewolf game session.
type Game struct {
	Id          string   `json:"id"`
	Players     []Player `json:"players"`
	WolfNumber  int      `json:"wolf_number"`
	Night       bool     `json:"night"`
	CurrentStep string   `json:"current_step"`
}

// Player represents a player in the game.
type Player struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
	Trait Trait  `json:"trait"`
	Alive bool   `json:"alive"`
}

// PlayerRequest is the HTTP request body for creating a game.
type PlayerRequest struct {
	Names     []string `json:"names"`
	WolfCount int      `json:"wolf_count"`
}

// NewGameWithPlayers creates a new Game from a list of players.
func NewGameWithPlayers(players []Player, wolfNumber int) Game {
	return Game{
		Id:          uuid.New().String(),
		Players:     players,
		WolfNumber:  wolfNumber,
		Night:       true,
		CurrentStep: "wolfAttack",
	}
}

// NewPlayer creates a new Player with the given name, role and trait.
func NewPlayer(name string, role Role, trait Trait) Player {
	return Player{
		Id:    uuid.New().String(),
		Name:  name,
		Role:  role,
		Trait: trait,
		Alive: true,
	}
}

// NewGame parses the HTTP request and initializes a game.
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

	players := make([]Player, len(names))
	for i, name := range names {
		players[i] = NewPlayer(name, roles[i], Trait(rand.Intn(5)))
	}

	newGame := NewGameWithPlayers(players, req.WolfCount)
	log.Printf("[GAME] Initialized | id=%s players=%d wolves=%d", newGame.Id, len(newGame.Players), newGame.WolfNumber)
	for _, player := range newGame.Players {
		log.Printf("[GAME] Player | name=%s role=%s trait=%s", player.Name, player.Role, player.Trait)
	}

	return newGame
}

func (g *Game) killRandomAliveVillager() (Player, error) {
	var candidates []int
	for i, player := range g.Players {
		if player.Alive && player.Role != Wolf {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) == 0 {
		return Player{}, fmt.Errorf("no alive villager to kill")
	}
	idx := candidates[rand.Intn(len(candidates))]
	g.Players[idx].Alive = false
	return g.Players[idx], nil
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
