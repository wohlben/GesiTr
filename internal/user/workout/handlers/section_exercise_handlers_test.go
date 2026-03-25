package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/user/workout/models"
)

func TestListWorkoutSectionExercises(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup: exercise -> scheme -> workout -> section
	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Bench Press", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED", "sets": 3, "reps": 10,
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED", "sets": 5, "reps": 5,
	})
	doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"name": "Push Day",
	})
	doJSON(r, "POST", "/api/user/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 0,
	})

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-section-exercises", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.WorkoutSectionExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	doJSON(r, "POST", "/api/user/workout-section-exercises", map[string]any{
		"workoutSectionId": 1, "exerciseSchemeId": 1, "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-section-exercises", map[string]any{
		"workoutSectionId": 1, "exerciseSchemeId": 2, "position": 1,
	})

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-section-exercises", nil)
		var result []models.WorkoutSectionExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
		// Should be ordered by position
		if result[0].Position != 0 || result[1].Position != 1 {
			t.Error("exercises not ordered by position")
		}
	})

	t.Run("filter by workoutSectionId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-section-exercises?workoutSectionId=1", nil)
		var result []models.WorkoutSectionExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by nonexistent section", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-section-exercises?workoutSectionId=999", nil)
		var result []models.WorkoutSectionExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/user/workout-section-exercises", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateWorkoutSectionExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup
	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Squat", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED", "sets": 5, "reps": 5,
	})
	doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"name": "Leg Day",
	})
	doJSON(r, "POST", "/api/user/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 0,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-section-exercises", map[string]any{
			"workoutSectionId": 1, "exerciseSchemeId": 1, "position": 0,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutSectionExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.WorkoutSectionID != 1 || result.ExerciseSchemeID != 1 {
			t.Error("create response mismatch")
		}
	})

	t.Run("section not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-section-exercises", map[string]any{
			"workoutSectionId": 999, "exerciseSchemeId": 1, "position": 0,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("scheme not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-section-exercises", map[string]any{
			"workoutSectionId": 1, "exerciseSchemeId": 999, "position": 0,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/user/workout-section-exercises", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestDeleteWorkoutSectionExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup
	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Deadlift", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED", "sets": 3, "reps": 8,
	})
	doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"name": "Pull Day",
	})
	doJSON(r, "POST", "/api/user/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-section-exercises", map[string]any{
		"workoutSectionId": 1, "exerciseSchemeId": 1, "position": 0,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/workout-section-exercises/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/user/workout-section-exercises/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
