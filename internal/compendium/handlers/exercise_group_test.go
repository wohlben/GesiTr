package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/compendium/models"
)

func TestListExerciseGroups(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-groups", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
	})

	doJSON(r, "POST", "/api/exercise-groups", map[string]any{
		"templateId": "g1", "name": "Push Day", "createdBy": "user",
	})
	doJSON(r, "POST", "/api/exercise-groups", map[string]any{
		"templateId": "g2", "name": "Pull Day", "createdBy": "user",
	})

	t.Run("all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-groups", nil)
		var result []models.ExerciseGroup
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter q", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-groups?q=Push", nil)
		var result []models.ExerciseGroup
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Name != "Push Day" {
			t.Errorf("q filter: got %d results", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/exercise-groups", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateExerciseGroup(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/exercise-groups", map[string]any{
			"templateId": "g1", "name": "Legs", "createdBy": "user",
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.ExerciseGroup
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.Name != "Legs" {
			t.Error("create response mismatch")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/exercise-groups", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/exercise-groups", map[string]any{
			"templateId": "x", "name": "X", "createdBy": "s",
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestGetExerciseGroup(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercise-groups", map[string]any{
		"templateId": "g1", "name": "Arms", "createdBy": "user",
	})

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-groups/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.ExerciseGroup
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Name != "Arms" {
			t.Error("get response mismatch")
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-groups/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestDeleteExerciseGroup(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercise-groups", map[string]any{
		"templateId": "g1", "name": "Core", "createdBy": "user",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/exercise-groups/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/exercise-groups/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
