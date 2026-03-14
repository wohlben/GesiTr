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
		if result.Status != models.WorkoutLogStatusPlanning {
			t.Errorf("expected planning status, got %s", result.Status)
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

	// Start the workout log
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	t.Run("finish a set", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/1", map[string]any{
			"status": "finished", "actualReps": 5, "actualWeight": 100.0,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.WorkoutLogExerciseSet
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Status != models.WorkoutLogStatusFinished {
			t.Errorf("expected finished status, got %s", result.Status)
		}
		if result.ActualReps == nil || *result.ActualReps != 5 {
			t.Error("actual reps mismatch")
		}
		if result.ActualWeight == nil || *result.ActualWeight != 100.0 {
			t.Error("actual weight mismatch")
		}
	})

	t.Run("update preserves target fields when not provided", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/2", map[string]any{
			"status":     "finished",
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

	t.Run("cannot transition from finished", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/1", map[string]any{
			"status": "in_progress",
		})
		if w.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/999", map[string]any{
			"status": "finished",
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

	// Start the workout
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish first set — exercise should NOT be finished yet
	doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Status == models.WorkoutLogStatusFinished {
		t.Error("log should not be finished after 1 of 2 sets")
	}
	if log.Sections[0].Status == models.WorkoutLogStatusFinished {
		t.Error("section should not be finished after 1 of 2 sets")
	}
	if log.Sections[0].Exercises[0].Status == models.WorkoutLogStatusFinished {
		t.Error("exercise should not be finished after 1 of 2 sets")
	}

	// Finish second set — everything should propagate to finished
	doJSON(r, "PUT", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	w = doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Sections[0].Exercises[0].Status != models.WorkoutLogStatusFinished {
		t.Error("exercise should be finished after all sets done")
	}
	if log.Sections[0].Status != models.WorkoutLogStatusFinished {
		t.Error("section should be finished after all exercises done")
	}
	if log.Status != models.WorkoutLogStatusFinished {
		t.Error("log should be finished after all sections done")
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
