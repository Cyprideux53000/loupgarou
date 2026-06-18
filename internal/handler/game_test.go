package handler

import (
	"fmt"
	"loupgarou/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockService struct {
	game *domain.Game
}

func newMockService() *mockService {
	return &mockService{
		game: &domain.Game{
			Id: "mock-id",
			Players: []domain.Player{
				{Id: "1", Name: "Alice", Role: domain.Villager, Trait: domain.Brave, Alive: true, Mayor: true},
				{Id: "2", Name: "Bob", Role: domain.Wolf, Trait: domain.Cunning, Alive: true, Mayor: false},
				{Id: "3", Name: "Charlie", Role: domain.Villager, Trait: domain.Sly, Alive: true, Mayor: false},
			},
			WolfNumber:  1,
			Night:       true,
			CurrentStep: "wolfAttack",
		},
	}
}

func (m *mockService) CreateGame(names []string, wolfCount int, mode string) (*domain.Game, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("at least one player is required")
	}
	return m.game, nil
}

func (m *mockService) GetGame(id string) (*domain.Game, error) {
	if id != m.game.Id {
		return nil, fmt.Errorf("game not found: %s", id)
	}
	return m.game, nil
}

func (m *mockService) GetStatus(id string) (*domain.Status, error) {
	if id != m.game.Id {
		return nil, fmt.Errorf("game not found: %s", id)
	}
	status := m.game.GetStatus()
	return &status, nil
}

func (m *mockService) ExecuteStep(id string) (*domain.StepResponse, error) {
	if id != m.game.Id {
		return nil, fmt.Errorf("game not found: %s", id)
	}
	result, err := m.game.Step(nil)
	if err != nil {
		return nil, err
	}
	return &domain.StepResponse{Game: m.game, Step: result}, nil
}

// --- HandleGame tests ---

func TestHandleGamePost(t *testing.T) {
	h := NewGameHandler(newMockService())
	body := `{"names":["Alice","Bob","Charlie"],"wolf_count":1}`
	req := httptest.NewRequest(http.MethodPost, "/game", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.HandleGame(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
	if !strings.Contains(w.Body.String(), "mock-id") {
		t.Error("response should contain game id")
	}
}

func TestHandleGamePostInvalidBody(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodPost, "/game", strings.NewReader("invalid"))
	w := httptest.NewRecorder()

	h.HandleGame(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleGamePostValidationError(t *testing.T) {
	h := NewGameHandler(newMockService())
	body := `{"names":[],"wolf_count":1}`
	req := httptest.NewRequest(http.MethodPost, "/game", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.HandleGame(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleGameGet(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodGet, "/game?id=mock-id", nil)
	w := httptest.NewRecorder()

	h.HandleGame(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "mock-id") {
		t.Error("response should contain game id")
	}
}

func TestHandleGameGetMissingId(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodGet, "/game", nil)
	w := httptest.NewRecorder()

	h.HandleGame(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleGameGetNotFound(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodGet, "/game?id=unknown", nil)
	w := httptest.NewRecorder()

	h.HandleGame(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleGameMethodNotAllowed(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodDelete, "/game", nil)
	w := httptest.NewRecorder()

	h.HandleGame(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// --- HandleStatus tests ---

func TestHandleStatus(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodGet, "/status?id=mock-id", nil)
	w := httptest.NewRecorder()

	h.HandleStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "wolves_alive") {
		t.Error("response should contain wolves_alive")
	}
}

func TestHandleStatusMissingId(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()

	h.HandleStatus(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleStatusNotFound(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodGet, "/status?id=unknown", nil)
	w := httptest.NewRecorder()

	h.HandleStatus(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleStatusMethodNotAllowed(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodPost, "/status", nil)
	w := httptest.NewRecorder()

	h.HandleStatus(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// --- HandleStep tests ---

func TestHandleStep(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodPost, "/step?id=mock-id", nil)
	w := httptest.NewRecorder()

	h.HandleStep(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "victim") {
		t.Error("response should contain victim")
	}
}

func TestHandleStepMissingId(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodPost, "/step", nil)
	w := httptest.NewRecorder()

	h.HandleStep(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleStepNotFound(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodPost, "/step?id=unknown", nil)
	w := httptest.NewRecorder()

	h.HandleStep(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleStepMethodNotAllowed(t *testing.T) {
	h := NewGameHandler(newMockService())
	req := httptest.NewRequest(http.MethodGet, "/step", nil)
	w := httptest.NewRecorder()

	h.HandleStep(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}
