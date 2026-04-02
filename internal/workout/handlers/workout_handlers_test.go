package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/database"
	"gesitr/internal/workout/models"
)

type paginatedJSON struct {
	Items  json.RawMessage `json:"items"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

func TestListWorkouts(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workouts", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Workout
		json.Unmarshal(page.Items, &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Push Day",
	})
	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Pull Day",
	})
	database.DB.Create(&models.WorkoutEntity{Owner: "bob", Name: "Bob's Workout"})

	t.Run("lists own + public (bob's is private)", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workouts", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Workout
		json.Unmarshal(page.Items, &result)
		if len(result) != 2 {
			t.Errorf("expected 2 (own only, bob's is private), got %d", len(result))
		}
		for _, wo := range result {
			if wo.Owner != "alice" {
				t.Errorf("expected owner alice, got %q", wo.Owner)
			}
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/workouts", nil)
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
		w := doJSON(r, "POST", "/api/workouts", map[string]any{
			"name": "Leg Day", "notes": notes,
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
		w := doRaw(r, "POST", "/api/workouts", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/workouts", map[string]any{
			"name": "X",
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestGetWorkout(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Push Day",
	})
	database.DB.Create(&models.WorkoutEntity{Owner: "bob", Name: "Bob's Workout"})

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workouts/1", nil)
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
		w := doJSON(r, "GET", "/api/workouts/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("forbidden for other owner", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/workouts/2", nil)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})
}

func TestGetWorkoutWithSectionsAndExercises(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create prerequisite chain: exercise -> scheme -> workout -> section -> section exercise
	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Bench Press", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED", "sets": 3, "reps": 10,
	})
	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Full Workout",
	})
	doJSON(r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": 1, "type": "supplementary", "label": "Warmup", "position": 0,
	})
	doJSON(r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": 1, "type": "main", "position": 1,
	})
	doJSON(r, "POST", "/api/workout-section-items", map[string]any{
		"workoutSectionId": 2, "type": "exercise", "exerciseSchemeId": 1, "position": 0,
	})

	w := doJSON(r, "GET", "/api/workouts/1", nil)
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
	if len(result.Sections[1].Items) != 1 {
		t.Fatalf("expected 1 item in main section, got %d", len(result.Sections[1].Items))
	}
	if *result.Sections[1].Items[0].ExerciseSchemeID != 1 {
		t.Error("exercise scheme ID mismatch")
	}
}

func TestUpdateWorkout(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Push Day",
	})
	database.DB.Create(&models.WorkoutEntity{Owner: "bob", Name: "Bob's Workout"})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/workouts/1", map[string]any{
			"name": "Heavy Push Day",
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
		w := doJSON(r, "PUT", "/api/workouts/999", map[string]any{
			"name": "X",
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("forbidden for other owner", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/workouts/2", map[string]any{
			"name": "Hijacked",
		})
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "PUT", "/api/workouts/1", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestDeleteWorkout(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/workouts", map[string]any{
		"name": "Push Day",
	})
	database.DB.Create(&models.WorkoutEntity{Owner: "bob", Name: "Bob's Workout"})

	t.Run("forbidden for other owner", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/workouts/2", nil)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/workouts/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/workouts/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}
