package village

// Status represents the current state of a game.
type Status struct {
	WolvesAlive    int    `json:"wolves_alive"`
	VillagersAlive int    `json:"villagers_alive"`
	NextStep       string `json:"next_step"`
	IsGameOver     bool   `json:"is_game_over"`
	Winner         string `json:"winner"`
}

// GetStatus computes the current status of the game.
func (g Game) GetStatus() Status {
	wolvesAlive := 0
	villagersAlive := 0

	for _, player := range g.Players {
		if player.Alive {
			if player.Role == Wolf {
				wolvesAlive++
			} else {
				villagersAlive++
			}
		}
	}

	isGameOver := wolvesAlive == 0 || wolvesAlive >= villagersAlive
	winner := ""
	nextStep := ""

	if isGameOver {
		if wolvesAlive == 0 {
			winner = "Villagers"
		} else {
			winner = "Wolves"
		}
	} else if g.Night {
		nextStep = "wolves_attack"
	} else {
		nextStep = "village_vote"
	}

	return Status{
		WolvesAlive:    wolvesAlive,
		VillagersAlive: villagersAlive,
		NextStep:       nextStep,
		IsGameOver:     isGameOver,
		Winner:         winner,
	}
}
