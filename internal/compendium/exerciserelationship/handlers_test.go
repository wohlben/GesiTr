package exerciserelationship

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListExerciseRelationships(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-relationships", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
	})

	doJSON(r, "POST", "/api/exercise-relationships", map[string]any{
		"relationshipType": "similar", "strength": 0.8, "createdBy": "system",
		"fromExerciseTemplateId": "ex1", "toExerciseTemplateId": "ex2",
	})
	doJSON(r, "POST", "/api/exercise-relationships", map[string]any{
		"relationshipType": "variation", "strength": 0.5, "createdBy": "system",
		"fromExerciseTemplateId": "ex3", "toExerciseTemplateId": "ex4",
	})

	t.Run("all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-relationships", nil)
		var result []ExerciseRelationship
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter fromExerciseTemplateId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-relationships?fromExerciseTemplateId=ex1", nil)
		var result []ExerciseRelationship
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("filter toExerciseTemplateId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-relationships?toExerciseTemplateId=ex4", nil)
		var result []ExerciseRelationship
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("filter relationshipType", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercise-relationships?relationshipType=similar", nil)
		var result []ExerciseRelationship
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/exercise-relationships", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateExerciseRelationship(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/exercise-relationships", map[string]any{
			"relationshipType": "similar", "strength": 0.9, "createdBy": "system",
			"fromExerciseTemplateId": "a", "toExerciseTemplateId": "b",
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d", w.Code)
		}
		var result ExerciseRelationship
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 {
			t.Error("expected non-zero ID")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/exercise-relationships", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/exercise-relationships", map[string]any{
			"relationshipType": "x", "strength": 0.1, "createdBy": "s",
			"fromExerciseTemplateId": "a", "toExerciseTemplateId": "b",
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestDeleteExerciseRelationship(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercise-relationships", map[string]any{
		"relationshipType": "similar", "strength": 0.8, "createdBy": "system",
		"fromExerciseTemplateId": "a", "toExerciseTemplateId": "b",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/exercise-relationships/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/exercise-relationships/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
