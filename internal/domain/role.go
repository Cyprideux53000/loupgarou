package domain

import (
	"encoding/json"
	"fmt"
)

type Role int

const (
	Villager Role = iota
	Wolf
)

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Role) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "Villager":
		*r = Villager
	case "Wolf":
		*r = Wolf
	default:
		return fmt.Errorf("unknown role: %s", s)
	}
	return nil
}

func (r Role) String() string {
	switch r {
	case Villager:
		return "Villager"
	case Wolf:
		return "Wolf"
	default:
		return "Unknown"
	}
}
