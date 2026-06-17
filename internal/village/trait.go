package village

import (
	"encoding/json"
	"fmt"
)

// Trait represents a character trait of a player.
type Trait int

const (
	Cunning Trait = iota
	Aggressive
	Brave
	Timid
	Sly
)

func (t Trait) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Trait) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "Cunning":
		*t = Cunning
	case "Aggressive":
		*t = Aggressive
	case "Brave":
		*t = Brave
	case "Timid":
		*t = Timid
	case "Sly":
		*t = Sly
	default:
		return fmt.Errorf("unknown trait: %s", s)
	}
	return nil
}

func (t Trait) String() string {
	switch t {
	case Cunning:
		return "Cunning"
	case Aggressive:
		return "Aggressive"
	case Brave:
		return "Brave"
	case Timid:
		return "Timid"
	case Sly:
		return "Sly"
	default:
		return "Unknown"
	}
}
