package domain

import (
	"encoding/json"
	"testing"
)

func TestTraitString(t *testing.T) {
	tests := []struct {
		trait    Trait
		expected string
	}{
		{Cunning, "Cunning"},
		{Aggressive, "Aggressive"},
		{Brave, "Brave"},
		{Timid, "Timid"},
		{Sly, "Sly"},
		{Trait(99), "Unknown"},
	}
	for _, tt := range tests {
		if tt.trait.String() != tt.expected {
			t.Errorf("Trait(%d).String() = %q, want %q", tt.trait, tt.trait.String(), tt.expected)
		}
	}
}

func TestTraitMarshalJSON(t *testing.T) {
	tests := []struct {
		trait    Trait
		expected string
	}{
		{Cunning, `"Cunning"`},
		{Aggressive, `"Aggressive"`},
		{Brave, `"Brave"`},
		{Timid, `"Timid"`},
		{Sly, `"Sly"`},
	}
	for _, tt := range tests {
		data, err := json.Marshal(tt.trait)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if string(data) != tt.expected {
			t.Errorf("Marshal(%v) = %s, want %s", tt.trait, data, tt.expected)
		}
	}
}

func TestTraitUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected Trait
	}{
		{`"Cunning"`, Cunning},
		{`"Aggressive"`, Aggressive},
		{`"Brave"`, Brave},
		{`"Timid"`, Timid},
		{`"Sly"`, Sly},
	}
	for _, tt := range tests {
		var tr Trait
		if err := json.Unmarshal([]byte(tt.input), &tr); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if tr != tt.expected {
			t.Errorf("Unmarshal(%s) = %v, want %v", tt.input, tr, tt.expected)
		}
	}
}

func TestTraitUnmarshalJSONError(t *testing.T) {
	var tr Trait
	if err := json.Unmarshal([]byte(`"InvalidTrait"`), &tr); err == nil {
		t.Error("expected error for unknown trait")
	}
	if err := json.Unmarshal([]byte(`123`), &tr); err == nil {
		t.Error("expected error for non-string input")
	}
}
