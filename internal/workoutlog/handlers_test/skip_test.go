package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/workoutlog/models"
)

func TestSkipSet(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Squat", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": 2, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Skip the first set
	w := doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "skipped",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("skip set: status = %d, body = %s", w.Code, w.Body.String())
	}
	var result models.WorkoutLogExerciseSet
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.Status != models.WorkoutLogItemStatusSkipped {
		t.Errorf("expected skipped status, got %s", result.Status)
	}
}

func TestPropagationAllSkipped(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Squat", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": 2, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Skip both sets
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{"status": "skipped"})
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{"status": "skipped"})

	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)

	// Exercise should be skipped
	if log.Sections[0].Exercises[0].Status != models.WorkoutLogItemStatusSkipped {
		t.Errorf("exercise should be skipped, got %s", log.Sections[0].Exercises[0].Status)
	}
	// Section should be skipped
	if log.Sections[0].Status != models.WorkoutLogItemStatusSkipped {
		t.Errorf("section should be skipped, got %s", log.Sections[0].Status)
	}
	// Log should be partially_finished (no skipped status for logs)
	if log.Status != models.WorkoutLogStatusPartiallyFinished {
		t.Errorf("log should be partially_finished, got %s", log.Status)
	}
}

func TestPropagationMixedSkipFinish(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Squat", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": 2, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish set 1, skip set 2
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "skipped",
	})

	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)

	// Exercise should be partially_finished (mix of finished and skipped)
	if log.Sections[0].Exercises[0].Status != models.WorkoutLogItemStatusPartiallyFinished {
		t.Errorf("exercise should be partially_finished, got %s", log.Sections[0].Exercises[0].Status)
	}
	// Section should be partially_finished
	if log.Sections[0].Status != models.WorkoutLogItemStatusPartiallyFinished {
		t.Errorf("section should be partially_finished, got %s", log.Sections[0].Status)
	}
	// Log should be partially_finished
	if log.Status != models.WorkoutLogStatusPartiallyFinished {
		t.Errorf("log should be partially_finished, got %s", log.Status)
	}
}

func TestAbandonPreservesSkipped(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Squat", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": 2, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Skip set 1
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "skipped",
	})

	// Abandon the log
	w := doJSON(r, "POST", "/api/user/workout-logs/1/abandon", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("abandon: status = %d, body = %s", w.Code, w.Body.String())
	}

	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)

	// Skipped set should remain skipped
	if log.Sections[0].Exercises[0].Sets[0].Status != models.WorkoutLogItemStatusSkipped {
		t.Errorf("set 1 should remain skipped, got %s", log.Sections[0].Exercises[0].Sets[0].Status)
	}
	// In-progress set should be aborted
	if log.Sections[0].Exercises[0].Sets[1].Status != models.WorkoutLogItemStatusAborted {
		t.Errorf("set 2 should be aborted, got %s", log.Sections[0].Exercises[0].Sets[1].Status)
	}
}
