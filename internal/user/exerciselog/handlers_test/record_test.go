package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	exerciselogmodels "gesitr/internal/user/exerciselog/models"
)

func TestRecordCreatedOnSetCompletion(t *testing.T) {
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

	// Start the workout
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish set 1 with 5 reps @ 100kg → e1RM = 100 * (1 + 5/30) = 116.667
	w := doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	// Verify exercise log was created with isRecord=true
	w = doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1&isRecord=true", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(logs))
	}
	expected := 100.0 * (1 + 5.0/30)
	if logs[0].RecordValue < expected-0.01 || logs[0].RecordValue > expected+0.01 {
		t.Errorf("expected e1RM ~%.3f, got %.3f", expected, logs[0].RecordValue)
	}
	if logs[0].MeasurementType != "REP_BASED" {
		t.Errorf("expected REP_BASED, got %s", logs[0].MeasurementType)
	}
	if !logs[0].IsRecord {
		t.Error("expected isRecord=true")
	}
}

func TestRecordUpdatedOnBetterPerformance(t *testing.T) {
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
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish set 1: 5 reps @ 100kg → e1RM = 116.667
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	// Finish set 2: 8 reps @ 100kg → e1RM = 126.667 (better)
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 8, "actualWeight": 100.0,
	})

	// Only one should be the record — the better one
	w := doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1&isRecord=true", nil)
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(logs))
	}
	expected := 100.0 * (1 + 8.0/30)
	if logs[0].RecordValue < expected-0.01 || logs[0].RecordValue > expected+0.01 {
		t.Errorf("expected e1RM ~%.3f, got %.3f", expected, logs[0].RecordValue)
	}

	// Both exercise logs should exist
	w = doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1", nil)
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 2 {
		t.Fatalf("expected 2 exercise logs total, got %d", len(logs))
	}
}

func TestRecordNotUpdatedOnWorsePerformance(t *testing.T) {
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
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish set 1: 8 reps @ 100kg → e1RM = 126.667
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 8, "actualWeight": 100.0,
	})

	// Finish set 2: 3 reps @ 100kg → e1RM = 110.0 (worse)
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 3, "actualWeight": 100.0,
	})

	w := doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1&isRecord=true", nil)
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(logs))
	}
	// The better one (set 1) should be the record
	expected := 100.0 * (1 + 8.0/30)
	if logs[0].RecordValue < expected-0.01 || logs[0].RecordValue > expected+0.01 {
		t.Errorf("expected e1RM ~%.3f (unchanged), got %.3f", expected, logs[0].RecordValue)
	}
}

func TestDifferentMeasurementTypes(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// TIME_BASED exercise
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "plank", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "TIME_BASED",
		"sets": 1, "duration": 60,
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
		"status": "finished", "actualDuration": 75,
	})

	w := doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1&isRecord=true", nil)
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(logs))
	}
	if logs[0].RecordValue != 75 {
		t.Errorf("expected duration 75, got %.1f", logs[0].RecordValue)
	}
	if logs[0].MeasurementType != "TIME_BASED" {
		t.Errorf("expected TIME_BASED, got %s", logs[0].MeasurementType)
	}
}

func TestDistanceBasedMeasurement(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "run", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "DISTANCE_BASED",
		"sets": 1, "distance": 5.0,
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
		"status": "finished", "actualDistance": 5.5,
	})

	w := doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1&isRecord=true", nil)
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(logs))
	}
	if logs[0].RecordValue != 5.5 {
		t.Errorf("expected distance 5.5, got %.1f", logs[0].RecordValue)
	}
	if logs[0].MeasurementType != "DISTANCE_BASED" {
		t.Errorf("expected DISTANCE_BASED, got %s", logs[0].MeasurementType)
	}
}

func TestUpdateExerciseLogShiftsRecord(t *testing.T) {
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
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Set 1: 5 reps @ 100kg → e1RM = 116.667
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})
	// Set 2: 8 reps @ 100kg → e1RM = 126.667 (record)
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 8, "actualWeight": 100.0,
	})

	// Get the exercise log IDs
	w := doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1", nil)
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}

	// Find the record (set 2, e1RM=126.667)
	var recordID uint
	for _, l := range logs {
		if l.IsRecord {
			recordID = l.ID
		}
	}
	if recordID == 0 {
		t.Fatal("no record found")
	}

	// Correct the record downward: change to 2 reps @ 100kg → e1RM = 106.667
	// Now set 1 (116.667) should become the new record
	w = doJSON(r, "PATCH", fmt.Sprintf("/api/user/exercise-logs/%d", recordID), map[string]any{
		"reps": 2,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("update failed: %d %s", w.Code, w.Body.String())
	}

	// Verify record shifted to set 1
	w = doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1&isRecord=true", nil)
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 record after update, got %d", len(logs))
	}
	expected := 100.0 * (1 + 5.0/30) // set 1's value
	if logs[0].RecordValue < expected-0.01 || logs[0].RecordValue > expected+0.01 {
		t.Errorf("expected record to shift to e1RM ~%.3f, got %.3f", expected, logs[0].RecordValue)
	}
}

func TestDeleteExerciseLogShiftsRecord(t *testing.T) {
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
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Set 1: 5 reps @ 100kg → e1RM = 116.667
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})
	// Set 2: 8 reps @ 100kg → e1RM = 126.667 (record)
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 8, "actualWeight": 100.0,
	})

	// Find the record
	w := doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1&isRecord=true", nil)
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(logs))
	}
	recordID := logs[0].ID

	// Delete the record entry
	w = doJSON(r, "DELETE", fmt.Sprintf("/api/user/exercise-logs/%d", recordID), nil)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete failed: %d %s", w.Code, w.Body.String())
	}

	// The remaining entry should now be the record
	w = doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1&isRecord=true", nil)
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 record after delete, got %d", len(logs))
	}
	expected := 100.0 * (1 + 5.0/30)
	if logs[0].RecordValue < expected-0.01 || logs[0].RecordValue > expected+0.01 {
		t.Errorf("expected record to shift to e1RM ~%.3f, got %.3f", expected, logs[0].RecordValue)
	}
}

func TestDeleteLastExerciseLogClearsRecord(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})

	// Create a single ad-hoc log
	w := doJSON(r, "POST", "/api/user/exercise-logs", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"reps": 5, "weight": 100.0, "performedAt": "2026-03-10T10:00:00Z",
	})
	var log exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &log)

	// Delete it
	doJSON(r, "DELETE", fmt.Sprintf("/api/user/exercise-logs/%d", log.ID), nil)

	// No records should remain
	w = doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1&isRecord=true", nil)
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 0 {
		t.Errorf("expected 0 records after deleting last entry, got %d", len(logs))
	}
}

func TestNoRecordWhenNotFinished(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
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

	// Skip the set instead of finishing — no ExerciseLog should be created
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "skipped",
	})

	w := doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1", nil)
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 0 {
		t.Errorf("expected 0 exercise logs for skipped set, got %d", len(logs))
	}
}

func TestPerExerciseNotPerScheme(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	// Two schemes for the same exercise
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 1, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 1, "reps": 10, "weight": 80.0,
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
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 2, "position": 1,
	})
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish set from scheme 1: 5 reps @ 100kg → e1RM = 116.667
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	// Finish set from scheme 2: 10 reps @ 80kg → e1RM = 106.667 (worse)
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 10, "actualWeight": 80.0,
	})

	// Should have only 1 record for the exercise
	w := doJSON(r, "GET", "/api/user/exercise-logs?userExerciseId=1&isRecord=true", nil)
	var logs []exerciselogmodels.ExerciseLog
	json.Unmarshal(w.Body.Bytes(), &logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 record (per exercise), got %d", len(logs))
	}
	expected := 100.0 * (1 + 5.0/30)
	if logs[0].RecordValue < expected-0.01 || logs[0].RecordValue > expected+0.01 {
		t.Errorf("expected e1RM ~%.3f, got %.3f", expected, logs[0].RecordValue)
	}
}
