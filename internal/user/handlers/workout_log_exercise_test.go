package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/user/models"
)

func TestListWorkoutLogExercises(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup: exercise -> scheme -> log -> section
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "bench-press", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED", "sets": 3, "reps": 10,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED", "sets": 5, "reps": 5,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Session", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-log-exercises", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.WorkoutLogExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 2, "position": 1,
	})

	t.Run("list all with sets preloaded", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-log-exercises", nil)
		var result []models.WorkoutLogExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Fatalf("expected 2, got %d", len(result))
		}
		if result[0].Position != 0 || result[1].Position != 1 {
			t.Error("exercises not ordered by position")
		}
		// First exercise from scheme with sets=3
		if len(result[0].Sets) != 3 {
			t.Errorf("expected 3 sets for exercise 1, got %d", len(result[0].Sets))
		}
		// Second exercise from scheme with sets=5
		if len(result[1].Sets) != 5 {
			t.Errorf("expected 5 sets for exercise 2, got %d", len(result[1].Sets))
		}
	})

	t.Run("filter by workoutLogSectionId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-log-exercises?workoutLogSectionId=1", nil)
		var result []models.WorkoutLogExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by nonexistent section", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-log-exercises?workoutLogSectionId=999", nil)
		var result []models.WorkoutLogExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/user/workout-log-exercises", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateWorkoutLogExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 5, "reps": 5, "weight": 100.0, "restBetweenSets": 180,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Leg Day Log", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0, "restBetweenExercises": 90,
	})

	t.Run("success with auto-created sets", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
			"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutLogExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if result.TargetMeasurementType != "REP_BASED" {
			t.Errorf("expected REP_BASED, got %s", result.TargetMeasurementType)
		}
		// Exercise-level BreakAfterSeconds should be snapshotted from section's RestBetweenExercises
		if result.BreakAfterSeconds == nil || *result.BreakAfterSeconds != 90 {
			t.Errorf("exercise breakAfterSeconds: expected 90, got %v", result.BreakAfterSeconds)
		}
		// Should have 5 auto-created sets
		if len(result.Sets) != 5 {
			t.Fatalf("expected 5 sets, got %d", len(result.Sets))
		}
		for i, s := range result.Sets {
			if s.SetNumber != i+1 {
				t.Errorf("set %d: expected setNumber %d, got %d", i, i+1, s.SetNumber)
			}
			if s.Completed {
				t.Errorf("set %d: should not be completed", i)
			}
			if s.TargetReps == nil || *s.TargetReps != 5 {
				t.Errorf("set %d: target reps mismatch", i)
			}
			if s.TargetWeight == nil || *s.TargetWeight != 100.0 {
				t.Errorf("set %d: target weight mismatch", i)
			}
			// Sets 1..4 should have BreakAfterSeconds=180, set 5 should be nil
			if i < 4 {
				if s.BreakAfterSeconds == nil || *s.BreakAfterSeconds != 180 {
					t.Errorf("set %d: expected breakAfterSeconds 180, got %v", i+1, s.BreakAfterSeconds)
				}
			} else {
				if s.BreakAfterSeconds != nil {
					t.Errorf("set %d (last): expected nil breakAfterSeconds, got %v", i+1, *s.BreakAfterSeconds)
				}
			}
		}
	})

	t.Run("success with dto breakAfterSeconds override", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
			"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 1,
			"breakAfterSeconds": 120,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutLogExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		// Should use the DTO value (120) instead of the section fallback (90)
		if result.BreakAfterSeconds == nil || *result.BreakAfterSeconds != 120 {
			t.Errorf("exercise breakAfterSeconds: expected 120, got %v", result.BreakAfterSeconds)
		}
	})

	t.Run("section not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
			"workoutLogSectionId": 999, "sourceExerciseSchemeId": 1, "position": 0,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("scheme not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
			"workoutLogSectionId": 1, "sourceExerciseSchemeId": 999, "position": 0,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/user/workout-log-exercises", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestUpdateWorkoutLogExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 5, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Leg Day Log", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	t.Run("update preserves target fields and returns sets", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-log-exercises/1", map[string]any{
			"position":               2,
			"targetMeasurementType":  "CHANGED",
			"workoutLogSectionId":    999,
			"sourceExerciseSchemeId": 999,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutLogExercise
		json.Unmarshal(w.Body.Bytes(), &result)
		// Target fields should be preserved
		if result.TargetMeasurementType != "REP_BASED" {
			t.Errorf("target measurement type changed to %q", result.TargetMeasurementType)
		}
		// Immutable fields should be preserved
		if result.WorkoutLogSectionID != 1 {
			t.Error("section ID changed")
		}
		if result.SourceExerciseSchemeID != 1 {
			t.Error("scheme ID changed")
		}
		// Position should be updated
		if result.Position != 2 {
			t.Errorf("expected position 2, got %d", result.Position)
		}
		// Sets should be preloaded
		if len(result.Sets) != 5 {
			t.Errorf("expected 5 sets, got %d", len(result.Sets))
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-log-exercises/999", map[string]any{
			"position": 1,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "PUT", "/api/user/workout-log-exercises/1", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestDeleteExercisePropagatesCompletion(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 1, "reps": 5, "weight": 100.0,
	})
	// Second exercise scheme
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 1, "reps": 8, "weight": 60.0,
	})

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Test", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	// Exercise 1: completed
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})
	// Exercise 2: incomplete
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 2, "position": 1,
	})

	// Complete exercise 1's only set
	doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"completed": true, "actualReps": 5, "actualWeight": 100.0,
	})

	// Section should not be completed (exercise 2 incomplete)
	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Completed {
		t.Fatal("log should not be completed yet")
	}

	// Delete the incomplete exercise 2 — now the section only has completed exercise 1
	doJSON(r, "DELETE", "/api/user/workout-log-exercises/2", nil)

	w = doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	json.Unmarshal(w.Body.Bytes(), &log)

	if !log.Sections[0].Completed {
		t.Error("section should be completed after deleting incomplete exercise")
	}
	if !log.Completed {
		t.Error("log should be completed after deleting incomplete exercise")
	}
}

func TestDeleteWorkoutLogExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "deadlift", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED", "sets": 3, "reps": 8,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Pull Day Log", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/workout-log-exercises/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/user/workout-log-exercises/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
