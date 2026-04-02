package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/workout/models"
)

func TestListWorkoutSections(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Push Day",
	})

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workout-sections", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.WorkoutSection
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	doJSON(r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": 1, "type": "supplementary", "label": "Warmup", "position": 0,
	})
	doJSON(r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 1,
	})

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workout-sections", nil)
		var result []models.WorkoutSection
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by workoutId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workout-sections?workoutId=1", nil)
		var result []models.WorkoutSection
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by nonexistent workoutId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workout-sections?workoutId=999", nil)
		var result []models.WorkoutSection
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/workout-sections", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateWorkoutSection(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Push Day",
	})

	t.Run("success main", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/workout-sections", map[string]any{
			"workoutId": 1, "type": "main", "position": 0,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutSection
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Type != "main" || result.WorkoutID != 1 {
			t.Error("create response mismatch")
		}
	})

	t.Run("success supplementary with label", func(t *testing.T) {
		label := "Warmup"
		w := doJSON(r, "POST", "/api/workout-sections", map[string]any{
			"workoutId": 1, "type": "supplementary", "label": label, "position": 0,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutSection
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Type != "supplementary" || *result.Label != label {
			t.Error("create supplementary response mismatch")
		}
	})

	t.Run("with restBetweenExercises", func(t *testing.T) {
		rest := 120
		w := doJSON(r, "POST", "/api/workout-sections", map[string]any{
			"workoutId": 1, "type": "main", "position": 2, "restBetweenExercises": rest,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutSection
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.RestBetweenExercises == nil || *result.RestBetweenExercises != rest {
			t.Error("restBetweenExercises mismatch")
		}
	})

	t.Run("workout not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/workout-sections", map[string]any{
			"workoutId": 999, "type": "main", "position": 0,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/workout-sections", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestGetWorkoutSection(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Push Day",
	})
	doJSON(r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 0,
	})

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workout-sections/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.WorkoutSection
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Type != "main" {
			t.Error("get response mismatch")
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workout-sections/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestDeleteWorkoutSection(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Push Day",
	})
	doJSON(r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 0,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/workout-sections/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/workout-sections/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
