package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	exerciselogmodels "gesitr/internal/user/exerciselog/models"

	"github.com/gin-gonic/gin"
)

// createExerciseLogViaWorkflow finishes a set and returns the resulting exercise log ID.
func createExerciseLogViaWorkflow(t *testing.T, r *gin.Engine) uint {
	t.Helper()
	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Squat", "slug": "squat", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": 1, "reps": 5, "weight": 100.0,
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
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	w := doJSON(r, "GET", "/api/user/exercise-logs?exerciseId=1", nil)
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) == 0 {
		t.Fatal("expected at least 1 exercise log after workflow")
	}
	return logs[0].ID
}

func TestGetExerciseLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()
	logID := createExerciseLogViaWorkflow(t, r)

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "GET", fmt.Sprintf("/api/user/exercise-logs/%d", logID), nil)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var log exerciselogmodels.ExerciseLog
		json.Unmarshal(w.Body.Bytes(), &log)
		if log.ID != logID {
			t.Errorf("expected ID %d, got %d", logID, log.ID)
		}
		if log.ExerciseID != 1 {
			t.Errorf("expected exerciseId 1, got %d", log.ExerciseID)
		}
		if log.MeasurementType != "REP_BASED" {
			t.Errorf("expected REP_BASED, got %s", log.MeasurementType)
		}
		if log.Reps == nil || *log.Reps != 5 {
			t.Error("reps mismatch")
		}
		if log.Weight == nil || *log.Weight != 100.0 {
			t.Error("weight mismatch")
		}
		if log.Owner != "alice" {
			t.Errorf("expected owner alice, got %s", log.Owner)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-logs/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestCreateExerciseLogAdHoc(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Squat", "slug": "squat", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})

	t.Run("ad-hoc creation with record value", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercise-logs", map[string]any{
			"exerciseId":      1,
			"measurementType": "REP_BASED",
			"reps":            10,
			"weight":          80.0,
			"performedAt":     "2026-03-15T14:00:00Z",
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
		var log exerciselogmodels.ExerciseLog
		json.Unmarshal(w.Body.Bytes(), &log)
		if log.ExerciseID != 1 {
			t.Errorf("expected exerciseId 1, got %d", log.ExerciseID)
		}
		if log.WorkoutLogExerciseSetID != nil {
			t.Error("ad-hoc log should have nil workoutLogExerciseSetId")
		}
		if log.Owner != "alice" {
			t.Errorf("expected owner alice, got %s", log.Owner)
		}
		expected := 80.0 * (1 + 10.0/30)
		if log.RecordValue < expected-0.01 || log.RecordValue > expected+0.01 {
			t.Errorf("expected recordValue ~%.3f, got %.3f", expected, log.RecordValue)
		}
		if !log.IsRecord {
			t.Error("first entry should be isRecord=true")
		}
	})

	t.Run("default performedAt", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercise-logs", map[string]any{
			"exerciseId":      1,
			"measurementType": "REP_BASED",
			"reps":            5,
			"weight":          60.0,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
		var log exerciselogmodels.ExerciseLog
		json.Unmarshal(w.Body.Bytes(), &log)
		if log.PerformedAt.IsZero() {
			t.Error("performedAt should default to now, not be zero")
		}
	})
}

func TestUpdateExerciseLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()
	logID := createExerciseLogViaWorkflow(t, r)

	t.Run("update reps and weight", func(t *testing.T) {
		w := doJSON(r, "PATCH", fmt.Sprintf("/api/user/exercise-logs/%d", logID), map[string]any{
			"reps":   8,
			"weight": 105.0,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var log exerciselogmodels.ExerciseLog
		json.Unmarshal(w.Body.Bytes(), &log)
		if log.Reps == nil || *log.Reps != 8 {
			t.Error("reps not updated")
		}
		if log.Weight == nil || *log.Weight != 105.0 {
			t.Error("weight not updated")
		}
		expected := 105.0 * (1 + 8.0/30)
		if log.RecordValue < expected-0.01 || log.RecordValue > expected+0.01 {
			t.Errorf("expected recordValue ~%.3f, got %.3f", expected, log.RecordValue)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PATCH", "/api/user/exercise-logs/999", map[string]any{
			"reps": 10,
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestDeleteExerciseLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()
	logID := createExerciseLogViaWorkflow(t, r)

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", fmt.Sprintf("/api/user/exercise-logs/%d", logID), nil)
		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
		}

		// Verify it's gone
		w = doJSON(r, "GET", fmt.Sprintf("/api/user/exercise-logs/%d", logID), nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 after delete, got %d", w.Code)
		}

		// Verify list is empty
		w = doJSON(r, "GET", "/api/user/exercise-logs?exerciseId=1", nil)
		var logs []exerciselogmodels.ExerciseLog
		json.Unmarshal(w.Body.Bytes(), &logs)
		if len(logs) != 0 {
			t.Errorf("expected 0 logs after delete, got %d", len(logs))
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/exercise-logs/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestListExerciseLogsFilters(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Squat", "slug": "squat", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Plank", "slug": "plank", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})

	doJSON(r, "POST", "/api/user/exercise-logs", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"reps": 5, "weight": 100.0, "performedAt": "2026-03-10T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/exercise-logs", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"reps": 8, "weight": 100.0, "performedAt": "2026-03-15T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/exercise-logs", map[string]any{
		"exerciseId": 2, "measurementType": "TIME_BASED",
		"duration": 60, "performedAt": "2026-03-12T10:00:00Z",
	})

	t.Run("filter by exerciseId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-logs?exerciseId=1", nil)
		var logs []exerciselogmodels.ExerciseLog
		json.Unmarshal(w.Body.Bytes(), &logs)
		if len(logs) != 2 {
			t.Fatalf("expected 2, got %d", len(logs))
		}
	})

	t.Run("filter by measurementType", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-logs?measurementType=TIME_BASED", nil)
		var logs []exerciselogmodels.ExerciseLog
		json.Unmarshal(w.Body.Bytes(), &logs)
		if len(logs) != 1 {
			t.Fatalf("expected 1, got %d", len(logs))
		}
		if logs[0].MeasurementType != "TIME_BASED" {
			t.Errorf("expected TIME_BASED, got %s", logs[0].MeasurementType)
		}
	})

	t.Run("filter by isRecord", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-logs?isRecord=true", nil)
		var logs []exerciselogmodels.ExerciseLog
		json.Unmarshal(w.Body.Bytes(), &logs)
		if len(logs) != 2 {
			t.Fatalf("expected 2 records (one per exercise), got %d", len(logs))
		}
	})

	t.Run("filter by date range", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-logs?from=2026-03-11T00:00:00Z&to=2026-03-14T00:00:00Z", nil)
		var logs []exerciselogmodels.ExerciseLog
		json.Unmarshal(w.Body.Bytes(), &logs)
		if len(logs) != 1 {
			t.Fatalf("expected 1 log in date range, got %d", len(logs))
		}
		if logs[0].ExerciseID != 2 {
			t.Errorf("expected exercise 2 in range, got exercise %d", logs[0].ExerciseID)
		}
	})

	t.Run("ordered by performedAt DESC", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-logs?exerciseId=1", nil)
		var logs []exerciselogmodels.ExerciseLog
		json.Unmarshal(w.Body.Bytes(), &logs)
		if len(logs) != 2 {
			t.Fatalf("expected 2, got %d", len(logs))
		}
		if logs[0].PerformedAt.Before(logs[1].PerformedAt) {
			t.Error("expected descending order by performedAt")
		}
	})
}
