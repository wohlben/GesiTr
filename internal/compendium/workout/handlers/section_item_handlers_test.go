package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/compendium/workout/models"
)

func TestListWorkoutSectionItems(t *testing.T) {
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
	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Push Day",
	})
	doJSON(r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 0,
	})

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workout-section-items", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.WorkoutSectionItem
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	doJSON(r, "POST", "/api/workout-section-items", map[string]any{
		"workoutSectionId": 1, "type": "exercise", "exerciseSchemeId": 1, "position": 0,
	})
	doJSON(r, "POST", "/api/workout-section-items", map[string]any{
		"workoutSectionId": 1, "type": "exercise", "exerciseSchemeId": 2, "position": 1,
	})

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workout-section-items", nil)
		var result []models.WorkoutSectionItem
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
		// Should be ordered by position
		if result[0].Position != 0 || result[1].Position != 1 {
			t.Error("items not ordered by position")
		}
	})

	t.Run("filter by workoutSectionId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workout-section-items?workoutSectionId=1", nil)
		var result []models.WorkoutSectionItem
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by nonexistent section", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workout-section-items?workoutSectionId=999", nil)
		var result []models.WorkoutSectionItem
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/workout-section-items", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateWorkoutSectionItem(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup
	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Squat", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED", "sets": 5, "reps": 5,
	})
	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Leg Day",
	})
	doJSON(r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 0,
	})

	t.Run("success exercise type", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/workout-section-items", map[string]any{
			"workoutSectionId": 1, "type": "exercise", "exerciseSchemeId": 1, "position": 0,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutSectionItem
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.WorkoutSectionID != 1 || *result.ExerciseSchemeID != 1 {
			t.Error("create response mismatch")
		}
		if result.Type != models.WorkoutSectionItemTypeExercise {
			t.Errorf("expected type exercise, got %s", result.Type)
		}
	})

	t.Run("success exercise_group type", func(t *testing.T) {
		// Create an exercise group
		doJSON(r, "POST", "/api/exercise-groups", map[string]any{
			"name": "Push Group",
		})
		w := doJSON(r, "POST", "/api/workout-section-items", map[string]any{
			"workoutSectionId": 1, "type": "exercise_group", "exerciseGroupId": 1, "position": 1,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutSectionItem
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Type != models.WorkoutSectionItemTypeExerciseGroup {
			t.Errorf("expected type exercise_group, got %s", result.Type)
		}
		if *result.ExerciseGroupID != 1 {
			t.Error("exerciseGroupId mismatch")
		}
	})

	t.Run("exercise type missing schemeId", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/workout-section-items", map[string]any{
			"workoutSectionId": 1, "type": "exercise", "position": 0,
		})
		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", w.Code)
		}
	})

	t.Run("exercise_group type missing groupId", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/workout-section-items", map[string]any{
			"workoutSectionId": 1, "type": "exercise_group", "position": 0,
		})
		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", w.Code)
		}
	})

	t.Run("section not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/workout-section-items", map[string]any{
			"workoutSectionId": 999, "type": "exercise", "exerciseSchemeId": 1, "position": 0,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("scheme not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/workout-section-items", map[string]any{
			"workoutSectionId": 1, "type": "exercise", "exerciseSchemeId": 999, "position": 0,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/workout-section-items", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestDeleteWorkoutSectionItem(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup
	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Deadlift", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED", "sets": 3, "reps": 8,
	})
	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Pull Day",
	})
	doJSON(r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/workout-section-items", map[string]any{
		"workoutSectionId": 1, "type": "exercise", "exerciseSchemeId": 1, "position": 0,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/workout-section-items/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/workout-section-items/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
