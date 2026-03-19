package record_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/user/record"
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

	// Verify record was created
	w = doJSON(r, "GET", "/api/user/records?userExerciseId=1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var records []record.UserRecord
	json.Unmarshal(w.Body.Bytes(), &records)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	expected := 100.0 * (1 + 5.0/30)
	if records[0].RecordValue < expected-0.01 || records[0].RecordValue > expected+0.01 {
		t.Errorf("expected e1RM ~%.3f, got %.3f", expected, records[0].RecordValue)
	}
	if records[0].MeasurementType != "REP_BASED" {
		t.Errorf("expected REP_BASED, got %s", records[0].MeasurementType)
	}
	if records[0].WorkoutLogExerciseSetID != 1 {
		t.Errorf("expected set ID 1, got %d", records[0].WorkoutLogExerciseSetID)
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

	// Start the workout
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish set 1: 5 reps @ 100kg → e1RM = 116.667
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	// Finish set 2: 8 reps @ 100kg → e1RM = 100 * (1 + 8/30) = 126.667 (better)
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 8, "actualWeight": 100.0,
	})

	w := doJSON(r, "GET", "/api/user/records?userExerciseId=1", nil)
	var records []record.UserRecord
	json.Unmarshal(w.Body.Bytes(), &records)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	expected := 100.0 * (1 + 8.0/30)
	if records[0].RecordValue < expected-0.01 || records[0].RecordValue > expected+0.01 {
		t.Errorf("expected e1RM ~%.3f, got %.3f", expected, records[0].RecordValue)
	}
	if records[0].WorkoutLogExerciseSetID != 2 {
		t.Errorf("expected set ID 2, got %d", records[0].WorkoutLogExerciseSetID)
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

	// Start the workout
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish set 1: 8 reps @ 100kg → e1RM = 126.667
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 8, "actualWeight": 100.0,
	})

	// Finish set 2: 3 reps @ 100kg → e1RM = 100 * (1 + 3/30) = 110.0 (worse)
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 3, "actualWeight": 100.0,
	})

	w := doJSON(r, "GET", "/api/user/records?userExerciseId=1", nil)
	var records []record.UserRecord
	json.Unmarshal(w.Body.Bytes(), &records)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	expected := 100.0 * (1 + 8.0/30)
	if records[0].RecordValue < expected-0.01 || records[0].RecordValue > expected+0.01 {
		t.Errorf("expected e1RM ~%.3f (unchanged), got %.3f", expected, records[0].RecordValue)
	}
	if records[0].WorkoutLogExerciseSetID != 1 {
		t.Errorf("expected set ID 1 (unchanged), got %d", records[0].WorkoutLogExerciseSetID)
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

	// Update set without finishing — just update actuals without status change
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"actualReps": 5, "actualWeight": 100.0,
	})

	w := doJSON(r, "GET", "/api/user/records?userExerciseId=1", nil)
	var records []record.UserRecord
	json.Unmarshal(w.Body.Bytes(), &records)
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
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

	// Start the workout
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualDuration": 75,
	})

	w := doJSON(r, "GET", "/api/user/records?userExerciseId=1", nil)
	var records []record.UserRecord
	json.Unmarshal(w.Body.Bytes(), &records)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].RecordValue != 75 {
		t.Errorf("expected duration 75, got %.1f", records[0].RecordValue)
	}
	if records[0].MeasurementType != "TIME_BASED" {
		t.Errorf("expected TIME_BASED, got %s", records[0].MeasurementType)
	}

	// DISTANCE_BASED exercise
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "run", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 2, "measurementType": "DISTANCE_BASED",
		"sets": 1, "distance": 5.0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 2, "position": 1,
	})

	// Need to start a new log for the distance exercise since log 1 is already in_progress
	// and the new exercise was created while in_progress — this would be blocked.
	// Actually, the exercise was added while the log was already started.
	// The guard blocks creating exercises when log is NOT planning.
	// We need a separate log for this test.

	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Test2", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 2, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 2, "sourceExerciseSchemeId": 2, "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-logs/2/start", nil)

	// The set ID depends on how many sets were created.
	// Log 1: 1 set (TIME_BASED scheme with sets=1) = set ID 1
	// Log 2: 1 set (DISTANCE_BASED scheme with sets=1) = set ID 2
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualDistance": 5.5,
	})

	w = doJSON(r, "GET", "/api/user/records?userExerciseId=2", nil)
	json.Unmarshal(w.Body.Bytes(), &records)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].RecordValue != 5.5 {
		t.Errorf("expected distance 5.5, got %.1f", records[0].RecordValue)
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
	// Log exercise from scheme 1
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})
	// Log exercise from scheme 2
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 2, "position": 1,
	})

	// Start the workout
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish set from scheme 1: 5 reps @ 100kg → e1RM = 116.667
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	// Finish set from scheme 2: 10 reps @ 80kg → e1RM = 80 * (1 + 10/30) = 106.667 (worse)
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 10, "actualWeight": 80.0,
	})

	// Should have only 1 record for the exercise (not 2 for each scheme)
	w := doJSON(r, "GET", "/api/user/records?userExerciseId=1", nil)
	var records []record.UserRecord
	json.Unmarshal(w.Body.Bytes(), &records)
	if len(records) != 1 {
		t.Fatalf("expected 1 record (per exercise), got %d", len(records))
	}
	// The better one (scheme 1) should win
	expected := 100.0 * (1 + 5.0/30)
	if records[0].RecordValue < expected-0.01 || records[0].RecordValue > expected+0.01 {
		t.Errorf("expected e1RM ~%.3f, got %.3f", expected, records[0].RecordValue)
	}
}

func TestListAndGetRecordEndpoints(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "bob", "compendiumExerciseId": "bench", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 1, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 2, "measurementType": "REP_BASED",
		"sets": 1, "reps": 5, "weight": 60.0,
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

	// Start the workout
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish both
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 60.0,
	})

	t.Run("filter by userExerciseId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/records?userExerciseId=1", nil)
		var records []record.UserRecord
		json.Unmarshal(w.Body.Bytes(), &records)
		if len(records) != 1 {
			t.Fatalf("expected 1, got %d", len(records))
		}
		if records[0].UserExerciseID != 1 {
			t.Errorf("expected userExerciseId 1, got %d", records[0].UserExerciseID)
		}
	})

	t.Run("filter by owner", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/records?owner=alice", nil)
		var records []record.UserRecord
		json.Unmarshal(w.Body.Bytes(), &records)
		if len(records) != 1 {
			t.Fatalf("expected 1 record for alice, got %d", len(records))
		}
	})

	t.Run("get by ID", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/records/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var record record.UserRecord
		json.Unmarshal(w.Body.Bytes(), &record)
		if record.ID == 0 {
			t.Error("expected non-zero ID")
		}
	})

	t.Run("get not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/records/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}
