package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/user/models"
)

func TestListWorkoutLogExerciseSets(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup exercise + scheme + log + section + exercise (auto-creates 3 sets)
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 3, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Test", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	t.Run("list all sets", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-log-exercise-sets?workoutLogExerciseId=1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.WorkoutLogExerciseSet
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 3 {
			t.Fatalf("expected 3 sets, got %d", len(result))
		}
		for i, s := range result {
			if s.SetNumber != i+1 {
				t.Errorf("set %d: expected setNumber %d, got %d", i, i+1, s.SetNumber)
			}
			if s.TargetReps == nil || *s.TargetReps != 5 {
				t.Errorf("set %d: target reps mismatch", i)
			}
			if s.TargetWeight == nil || *s.TargetWeight != 100.0 {
				t.Errorf("set %d: target weight mismatch", i)
			}
		}
	})

	t.Run("filter by nonexistent exercise", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-log-exercise-sets?workoutLogExerciseId=999", nil)
		var result []models.WorkoutLogExerciseSet
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/user/workout-log-exercise-sets", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateWorkoutLogExerciseSet(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 2, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Test", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	t.Run("create additional set", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-log-exercise-sets", map[string]any{
			"workoutLogExerciseId": 1, "setNumber": 3, "targetReps": 5, "targetWeight": 105.0,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutLogExerciseSet
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.SetNumber != 3 {
			t.Errorf("expected setNumber 3, got %d", result.SetNumber)
		}
		if result.TargetWeight == nil || *result.TargetWeight != 105.0 {
			t.Errorf("target weight mismatch")
		}
	})

	t.Run("exercise not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-log-exercise-sets", map[string]any{
			"workoutLogExerciseId": 999, "setNumber": 1,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/user/workout-log-exercise-sets", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestUpdateWorkoutLogExerciseSet(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 3, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Test", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	t.Run("update actual fields and completed", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/1", map[string]any{
			"completed": true, "actualReps": 5, "actualWeight": 100.0,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutLogExerciseSet
		json.Unmarshal(w.Body.Bytes(), &result)
		if !result.Completed {
			t.Error("expected completed to be true")
		}
		if result.ActualReps == nil || *result.ActualReps != 5 {
			t.Error("actual reps mismatch")
		}
		if result.ActualWeight == nil || *result.ActualWeight != 100.0 {
			t.Error("actual weight mismatch")
		}
	})

	t.Run("update preserves target fields when not provided", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/1", map[string]any{
			"completed":  true,
			"actualReps": 4,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.WorkoutLogExerciseSet
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.TargetReps == nil || *result.TargetReps != 5 {
			t.Errorf("target reps changed to %v", result.TargetReps)
		}
		if result.TargetWeight == nil || *result.TargetWeight != 100.0 {
			t.Errorf("target weight changed to %v", result.TargetWeight)
		}
		if result.ActualReps == nil || *result.ActualReps != 4 {
			t.Error("actual reps not updated")
		}
	})

	t.Run("update target fields when provided", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/1", map[string]any{
			"completed":    true,
			"actualReps":   4,
			"targetReps":   8,
			"targetWeight": 120.0,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.WorkoutLogExerciseSet
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.TargetReps == nil || *result.TargetReps != 8 {
			t.Errorf("target reps: expected 8, got %v", result.TargetReps)
		}
		if result.TargetWeight == nil || *result.TargetWeight != 120.0 {
			t.Errorf("target weight: expected 120, got %v", result.TargetWeight)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/999", map[string]any{
			"completed": true,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "PUT", "/api/user/workout-log-exercise-sets/1", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestProgressiveCompletion(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 2, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Test", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	// Complete first set — exercise should NOT be completed yet
	doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"completed": true, "actualReps": 5, "actualWeight": 100.0,
	})

	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Completed {
		t.Error("log should not be completed after 1 of 2 sets")
	}
	if log.Sections[0].Completed {
		t.Error("section should not be completed after 1 of 2 sets")
	}
	if log.Sections[0].Exercises[0].Completed {
		t.Error("exercise should not be completed after 1 of 2 sets")
	}

	// Complete second set — everything should propagate to completed
	doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"completed": true, "actualReps": 5, "actualWeight": 100.0,
	})

	w = doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	json.Unmarshal(w.Body.Bytes(), &log)
	if !log.Sections[0].Exercises[0].Completed {
		t.Error("exercise should be completed after all sets done")
	}
	if !log.Sections[0].Completed {
		t.Error("section should be completed after all exercises done")
	}
	if !log.Completed {
		t.Error("log should be completed after all sections done")
	}

	// Un-complete a set — should propagate false upward
	doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"completed": false, "actualReps": 5, "actualWeight": 100.0,
	})

	w = doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Sections[0].Exercises[0].Completed {
		t.Error("exercise should not be completed after un-completing a set")
	}
	if log.Sections[0].Completed {
		t.Error("section should not be completed after un-completing a set")
	}
	if log.Completed {
		t.Error("log should not be completed after un-completing a set")
	}
}

func TestDeleteSetPropagatesCompletion(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 2, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Test", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	// Complete only set 1 — exercise not yet complete
	doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"completed": true, "actualReps": 5, "actualWeight": 100.0,
	})

	// Delete the incomplete set 2 — now the only remaining set is completed,
	// so exercise/section/log should propagate to completed
	doJSON(r, "DELETE", "/api/user/workout-log-exercise-sets/2", nil)

	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)

	if !log.Sections[0].Exercises[0].Completed {
		t.Error("exercise should be completed after deleting the incomplete set")
	}
	if !log.Sections[0].Completed {
		t.Error("section should be completed")
	}
	if !log.Completed {
		t.Error("log should be completed")
	}
}

func TestDeleteWorkoutLogExerciseSet(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED", "sets": 3, "reps": 5,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Test", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/workout-log-exercise-sets/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
		// Verify only 2 sets remain
		w = doJSON(r, "GET", "/api/user/workout-log-exercise-sets?workoutLogExerciseId=1", nil)
		var result []models.WorkoutLogExerciseSet
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2 sets remaining, got %d", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/user/workout-log-exercise-sets/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
