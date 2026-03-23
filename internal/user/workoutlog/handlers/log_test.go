package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"gesitr/internal/user/workoutlog/models"
)

func itoa(v uint) string { return fmt.Sprint(v) }

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
		"name": "Monday Session", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name": "Tuesday Session", "date": "2026-03-08T10:00:00Z",
	})

	t.Run("list all (scoped to auth user)", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-logs", nil)
		var result []models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
		for _, log := range result {
			if log.Owner != "alice" {
				t.Errorf("expected owner alice, got %s", log.Owner)
			}
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

func TestUpdateWorkoutLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Session", "date": "2026-03-07T10:00:00Z",
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "PATCH", "/api/user/workout-logs/1", map[string]any{
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
		w := doJSON(r, "PATCH", "/api/user/workout-logs/999", map[string]any{
			"owner": "x", "name": "X", "date": "2026-03-07T10:00:00Z",
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "PATCH", "/api/user/workout-logs/1", "{bad")
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

func TestStartAdhocWorkoutLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("creates adhoc log with section", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-logs/adhoc", nil)
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Status != models.WorkoutLogStatusAdhoc {
			t.Errorf("expected adhoc status, got %s", result.Status)
		}
		if result.Name != "Ad-hoc Workout" {
			t.Errorf("expected 'Ad-hoc Workout', got %q", result.Name)
		}
		if result.WorkoutID != nil {
			t.Error("expected nil workoutId")
		}
		if result.Date == nil {
			t.Error("expected date to be set")
		}
		if len(result.Sections) != 1 {
			t.Fatalf("expected 1 section, got %d", len(result.Sections))
		}
		sec := result.Sections[0]
		if sec.Status != models.WorkoutLogItemStatusInProgress {
			t.Errorf("expected section in_progress, got %s", sec.Status)
		}
		if sec.Label == nil || *sec.Label != "Adhoc" {
			t.Error("expected section label 'Adhoc'")
		}
	})
}

func TestFinishWorkoutLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create prerequisite: user exercise + scheme
	ueW := doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Test Exercise", "slug": "test-ex-1", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	if ueW.Code != http.StatusCreated {
		t.Fatalf("create exercise: status = %d, body = %s", ueW.Code, ueW.Body.String())
	}
	schemeW := doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED", "sets": 2, "reps": 10,
	})
	if schemeW.Code != http.StatusCreated {
		t.Fatalf("create scheme: status = %d, body = %s", schemeW.Code, schemeW.Body.String())
	}

	t.Run("finishes adhoc with exercises", func(t *testing.T) {
		// Start adhoc
		w := doJSON(r, "POST", "/api/user/workout-logs/adhoc", nil)
		var log models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &log)

		// Add exercise (should work in adhoc)
		w = doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
			"workoutLogSectionId": log.Sections[0].ID, "sourceExerciseSchemeId": 1, "position": 0,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("create exercise: status = %d, body = %s", w.Code, w.Body.String())
		}
		var ex models.WorkoutLogExercise
		json.Unmarshal(w.Body.Bytes(), &ex)

		// Verify sets are in_progress
		if len(ex.Sets) != 2 {
			t.Fatalf("expected 2 sets, got %d", len(ex.Sets))
		}
		for _, s := range ex.Sets {
			if s.Status != models.WorkoutLogItemStatusInProgress {
				t.Errorf("expected set in_progress, got %s", s.Status)
			}
		}

		// Complete first set
		w = doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/"+itoa(ex.Sets[0].ID), map[string]any{
			"status": "finished", "actualReps": 10,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("complete set: status = %d, body = %s", w.Code, w.Body.String())
		}

		// Verify propagation stopped at exercise (log should still be adhoc)
		w = doJSON(r, "GET", "/api/user/workout-logs/"+itoa(log.ID), nil)
		json.Unmarshal(w.Body.Bytes(), &log)
		if log.Status != models.WorkoutLogStatusAdhoc {
			t.Errorf("expected log still adhoc, got %s", log.Status)
		}

		// Finish workout
		w = doJSON(r, "POST", "/api/user/workout-logs/"+itoa(log.ID)+"/finish", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("finish: status = %d, body = %s", w.Code, w.Body.String())
		}
		json.Unmarshal(w.Body.Bytes(), &log)
		if log.Status == models.WorkoutLogStatusAdhoc {
			t.Error("expected log to no longer be adhoc after finish")
		}
		// Should be partially_finished (one set finished, one skipped)
		if log.Status != models.WorkoutLogStatusPartiallyFinished {
			t.Errorf("expected partially_finished, got %s", log.Status)
		}
	})

	t.Run("finish non-adhoc rejected", func(t *testing.T) {
		// Create a regular planning log
		w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
			"name": "Regular", "date": "2026-03-07T10:00:00Z",
		})
		var log models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &log)

		w = doJSON(r, "POST", "/api/user/workout-logs/"+itoa(log.ID)+"/finish", nil)
		if w.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-logs/999/finish", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}
