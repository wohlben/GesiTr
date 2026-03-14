package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/user/models"
)

func TestListWorkoutLogs(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-logs", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Monday Session", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "bob", "name": "Tuesday Session", "date": "2026-03-08T10:00:00Z",
	})

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-logs", nil)
		var result []models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by owner", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-logs?owner=alice", nil)
		var result []models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Owner != "alice" {
			t.Errorf("owner filter: got %d results", len(result))
		}
	})

	t.Run("filter by workoutId", func(t *testing.T) {
		// Create a workout and a log referencing it
		doJSON(r, "POST", "/api/user/workouts", map[string]any{
			"owner": "alice", "name": "Template", "date": "2026-03-07T10:00:00Z",
		})
		doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
			"owner": "alice", "name": "From Template", "date": "2026-03-09T10:00:00Z", "workoutId": 1,
		})
		w := doJSON(r, "GET", "/api/user/workout-logs?workoutId=1", nil)
		var result []models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/user/workout-logs", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateWorkoutLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("ad-hoc (no workoutId)", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
			"owner": "alice", "name": "Ad-hoc Session", "date": "2026-03-07T10:00:00Z",
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.Name != "Ad-hoc Session" || result.WorkoutID != nil {
			t.Error("create response mismatch")
		}
	})

	t.Run("with workoutId", func(t *testing.T) {
		doJSON(r, "POST", "/api/user/workouts", map[string]any{
			"owner": "alice", "name": "Template", "date": "2026-03-07T10:00:00Z",
		})
		w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
			"owner": "alice", "name": "From Template", "date": "2026-03-07T10:00:00Z", "workoutId": 1,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.WorkoutID == nil || *result.WorkoutID != 1 {
			t.Error("workoutId mismatch")
		}
	})

	t.Run("workout not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
			"owner": "alice", "name": "Bad", "date": "2026-03-07T10:00:00Z", "workoutId": 999,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/user/workout-logs", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
			"owner": "x", "name": "X", "date": "2026-03-07T10:00:00Z",
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestGetWorkoutLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Session", "date": "2026-03-07T10:00:00Z",
	})

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Name != "Session" {
			t.Error("get response mismatch")
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-logs/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestGetWorkoutLogWithSectionsAndExercises(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup: exercise -> scheme -> workout log -> section -> exercise
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "bench-press", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED", "sets": 3, "reps": 10,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Full Session", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "supplementary", "label": "Warmup", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 1,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 2, "sourceExerciseSchemeId": 1, "position": 0,
	})

	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}

	var result models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &result)

	if len(result.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(result.Sections))
	}
	if result.Sections[0].Position != 0 || result.Sections[1].Position != 1 {
		t.Error("sections not ordered by position")
	}
	if len(result.Sections[1].Exercises) != 1 {
		t.Fatalf("expected 1 exercise in main section, got %d", len(result.Sections[1].Exercises))
	}
	if result.Sections[1].Exercises[0].SourceExerciseSchemeID != 1 {
		t.Error("exercise scheme ID mismatch")
	}
}

func TestUpdateWorkoutLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Session", "date": "2026-03-07T10:00:00Z",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-logs/1", map[string]any{
			"owner": "alice", "name": "Updated Session", "date": "2026-03-07T10:00:00Z",
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Name != "Updated Session" {
			t.Errorf("expected updated name, got %q", result.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-logs/999", map[string]any{
			"owner": "x", "name": "X", "date": "2026-03-07T10:00:00Z",
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "PUT", "/api/user/workout-logs/1", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestDeleteWorkoutLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Session", "date": "2026-03-07T10:00:00Z",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/workout-logs/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/user/workout-logs/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
