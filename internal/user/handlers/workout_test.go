package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/user/models"
)

func TestListWorkouts(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workouts", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.Workout
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"owner": "alice", "name": "Push Day", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"owner": "bob", "name": "Pull Day", "date": "2026-03-08T10:00:00Z",
	})

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workouts", nil)
		var result []models.Workout
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by owner", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workouts?owner=alice", nil)
		var result []models.Workout
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Owner != "alice" {
			t.Errorf("owner filter: got %d results", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/user/workouts", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateWorkout(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("success", func(t *testing.T) {
		notes := "First workout"
		w := doJSON(r, "POST", "/api/user/workouts", map[string]any{
			"owner": "alice", "name": "Leg Day", "notes": notes, "date": "2026-03-07T10:00:00Z",
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Workout
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.Name != "Leg Day" || *result.Notes != notes {
			t.Error("create response mismatch")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/user/workouts", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/user/workouts", map[string]any{
			"owner": "x", "name": "X", "date": "2026-03-07T10:00:00Z",
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestGetWorkout(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"owner": "alice", "name": "Push Day", "date": "2026-03-07T10:00:00Z",
	})

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workouts/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.Workout
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Name != "Push Day" {
			t.Error("get response mismatch")
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workouts/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestGetWorkoutWithSectionsAndExercises(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create prerequisite chain: exercise -> scheme -> workout -> section -> section exercise
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "exerciseTemplateId": "bench-press", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED", "sets": 3, "reps": 10,
	})
	doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"owner": "alice", "name": "Full Workout", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-sections", map[string]any{
		"workoutId": 1, "type": "supplementary", "label": "Warmup", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 1,
	})
	doJSON(r, "POST", "/api/user/workout-section-exercises", map[string]any{
		"workoutSectionId": 2, "userExerciseSchemeId": 1, "position": 0,
	})

	w := doJSON(r, "GET", "/api/user/workouts/1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}

	var result models.Workout
	json.Unmarshal(w.Body.Bytes(), &result)

	if len(result.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(result.Sections))
	}
	// Sections should be ordered by position
	if result.Sections[0].Position != 0 || result.Sections[1].Position != 1 {
		t.Error("sections not ordered by position")
	}
	if result.Sections[0].Type != "supplementary" || result.Sections[1].Type != "main" {
		t.Error("section types mismatch")
	}
	if len(result.Sections[1].Exercises) != 1 {
		t.Fatalf("expected 1 exercise in main section, got %d", len(result.Sections[1].Exercises))
	}
	if result.Sections[1].Exercises[0].UserExerciseSchemeID != 1 {
		t.Error("exercise scheme ID mismatch")
	}
}

func TestUpdateWorkout(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"owner": "alice", "name": "Push Day", "date": "2026-03-07T10:00:00Z",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workouts/1", map[string]any{
			"owner": "alice", "name": "Heavy Push Day", "date": "2026-03-07T10:00:00Z",
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Workout
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Name != "Heavy Push Day" {
			t.Errorf("expected updated name, got %q", result.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workouts/999", map[string]any{
			"owner": "x", "name": "X", "date": "2026-03-07T10:00:00Z",
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "PUT", "/api/user/workouts/1", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestDeleteWorkout(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"owner": "alice", "name": "Push Day", "date": "2026-03-07T10:00:00Z",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/workouts/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/user/workouts/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
