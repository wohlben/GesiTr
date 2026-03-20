package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/compendium/exercisegroup/models"
)

func TestListExerciseGroupMembers(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-group-members", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
	})

	doJSON(r, "POST", "/api/exercise-group-members", map[string]any{
		"groupTemplateId": "g1", "exerciseTemplateId": "ex1", "addedBy": "user",
	})
	doJSON(r, "POST", "/api/exercise-group-members", map[string]any{
		"groupTemplateId": "g1", "exerciseTemplateId": "ex2", "addedBy": "user",
	})
	doJSON(r, "POST", "/api/exercise-group-members", map[string]any{
		"groupTemplateId": "g2", "exerciseTemplateId": "ex1", "addedBy": "user",
	})

	t.Run("all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-group-members", nil)
		var result []models.ExerciseGroupMember
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 3 {
			t.Errorf("expected 3, got %d", len(result))
		}
	})

	t.Run("filter groupTemplateId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-group-members?groupTemplateId=g1", nil)
		var result []models.ExerciseGroupMember
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter exerciseTemplateId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-group-members?exerciseTemplateId=ex1", nil)
		var result []models.ExerciseGroupMember
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/exercise-group-members", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateExerciseGroupMember(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/exercise-group-members", map[string]any{
			"groupTemplateId": "g1", "exerciseTemplateId": "ex1", "addedBy": "user",
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.ExerciseGroupMember
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.GroupTemplateID != "g1" {
			t.Error("create response mismatch")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/exercise-group-members", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/exercise-group-members", map[string]any{
			"groupTemplateId": "x", "exerciseTemplateId": "y", "addedBy": "s",
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestDeleteExerciseGroupMember(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercise-group-members", map[string]any{
		"groupTemplateId": "g1", "exerciseTemplateId": "ex1", "addedBy": "user",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/exercise-group-members/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/exercise-group-members/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
