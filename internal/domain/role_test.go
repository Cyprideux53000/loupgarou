package domain

import (
	"encoding/json"
	"testing"
)

func TestRoleString(t *testing.T) {
	tests := []struct {
		role     Role
		expected string
	}{
		{Villager, "Villager"},
		{Wolf, "Wolf"},
		{Role(99), "Unknown"},
	}
	for _, tt := range tests {
		if tt.role.String() != tt.expected {
			t.Errorf("Role(%d).String() = %q, want %q", tt.role, tt.role.String(), tt.expected)
		}
	}
}

func TestRoleMarshalJSON(t *testing.T) {
	tests := []struct {
		role     Role
		expected string
	}{
		{Villager, `"Villager"`},
		{Wolf, `"Wolf"`},
	}
	for _, tt := range tests {
		data, err := json.Marshal(tt.role)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if string(data) != tt.expected {
			t.Errorf("Marshal(%v) = %s, want %s", tt.role, data, tt.expected)
		}
	}
}

func TestRoleUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected Role
	}{
		{`"Villager"`, Villager},
		{`"Wolf"`, Wolf},
	}
	for _, tt := range tests {
		var r Role
		if err := json.Unmarshal([]byte(tt.input), &r); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if r != tt.expected {
			t.Errorf("Unmarshal(%s) = %v, want %v", tt.input, r, tt.expected)
		}
	}
}

func TestRoleUnmarshalJSONError(t *testing.T) {
	var r Role
	if err := json.Unmarshal([]byte(`"InvalidRole"`), &r); err == nil {
		t.Error("expected error for unknown role")
	}
	if err := json.Unmarshal([]byte(`123`), &r); err == nil {
		t.Error("expected error for non-string input")
	}
}
