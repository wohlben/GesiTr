package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/database"
	"gesitr/internal/user/models"
)

func TestStartCascade(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	// Start the log
	w := doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("start: status = %d, body = %s", w.Code, w.Body.String())
	}

	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)

	if log.Status != models.WorkoutLogStatusInProgress {
		t.Errorf("log status: expected in_progress, got %s", log.Status)
	}
	if log.StatusChangedAt == nil {
		t.Error("log statusChangedAt should be set")
	}
	if log.Date == nil {
		t.Error("log date should be set after start")
	}
	if log.Sections[0].Status != models.WorkoutLogStatusInProgress {
		t.Errorf("section status: expected in_progress, got %s", log.Sections[0].Status)
	}
	if log.Sections[0].Exercises[0].Status != models.WorkoutLogStatusInProgress {
		t.Errorf("exercise status: expected in_progress, got %s", log.Sections[0].Exercises[0].Status)
	}
	for i, s := range log.Sections[0].Exercises[0].Sets {
		if s.Status != models.WorkoutLogStatusInProgress {
			t.Errorf("set %d status: expected in_progress, got %s", i, s.Status)
		}
		if s.StatusChangedAt == nil {
			t.Errorf("set %d statusChangedAt should be set", i)
		}
	}
}

func TestStartWhenNotPlanning(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	// Start once
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Start again — should 409
	w := doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestAbandonCascade(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	// Start the log
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish the first set
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	// Abandon the log
	w := doJSON(r, "POST", "/api/user/workout-logs/1/abandon", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("abandon: status = %d, body = %s", w.Code, w.Body.String())
	}

	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)

	if log.Status != models.WorkoutLogStatusAborted {
		t.Errorf("log status: expected aborted, got %s", log.Status)
	}

	// The finished set should stay finished
	set1 := log.Sections[0].Exercises[0].Sets[0]
	if set1.Status != models.WorkoutLogStatusFinished {
		t.Errorf("set 1 should remain finished, got %s", set1.Status)
	}

	// The in_progress set should be aborted
	set2 := log.Sections[0].Exercises[0].Sets[1]
	if set2.Status != models.WorkoutLogStatusAborted {
		t.Errorf("set 2 should be aborted, got %s", set2.Status)
	}
}

func TestAbandonWhenTerminal(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish the only set — log should propagate to finished
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	// Try to abandon — should 409 because it's already finished
	w := doJSON(r, "POST", "/api/user/workout-logs/1/abandon", nil)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestUniquePlanningLog(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"owner": "alice", "name": "Template",
	})

	// First planning log — should succeed
	w := doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Log 1", "workoutId": 1,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("first log: status = %d", w.Code)
	}

	// Second planning log for same workout — should 409
	w = doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Log 2", "workoutId": 1,
	})
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate planning log, got %d", w.Code)
	}

	// Start the first log — it should no longer block new planning logs
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 1, "reps": 5, "weight": 100.0,
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Now a new planning log should be allowed
	w = doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Log 2", "workoutId": 1,
	})
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 after first log started, got %d", w.Code)
	}
}

func TestGuardCreateSectionWhenInProgress(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Try to create a section — should be blocked
	w := doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "supplementary", "position": 1,
	})
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestGuardDeleteExerciseWhenInProgress(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Try to delete exercise — should be blocked
	w := doJSON(r, "DELETE", "/api/user/workout-log-exercises/1", nil)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestPropagationAllFinished(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish both sets
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)

	if log.Status != models.WorkoutLogStatusFinished {
		t.Errorf("log should be finished, got %s", log.Status)
	}
	if log.Sections[0].Status != models.WorkoutLogStatusFinished {
		t.Errorf("section should be finished, got %s", log.Sections[0].Status)
	}
	if log.Sections[0].Exercises[0].Status != models.WorkoutLogStatusFinished {
		t.Errorf("exercise should be finished, got %s", log.Sections[0].Exercises[0].Status)
	}
}

func TestPropagationAnyAborted(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish set 1, abort set 2
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "aborted",
	})

	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)

	// When any child is aborted, parent should be aborted
	if log.Status != models.WorkoutLogStatusAborted {
		t.Errorf("log should be aborted (has aborted set), got %s", log.Status)
	}
	if log.Sections[0].Status != models.WorkoutLogStatusAborted {
		t.Errorf("section should be aborted, got %s", log.Sections[0].Status)
	}
	if log.Sections[0].Exercises[0].Status != models.WorkoutLogStatusAborted {
		t.Errorf("exercise should be aborted, got %s", log.Sections[0].Exercises[0].Status)
	}
}

func TestDeleteInProgressLog(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Try to delete an in-progress log — should 409
	w := doJSON(r, "DELETE", "/api/user/workout-logs/1", nil)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestMultiSectionPropagation(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	// Section 1
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})
	// Section 2
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "supplementary", "position": 1,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 2, "sourceExerciseSchemeId": 1, "position": 0,
	})

	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Finish the set in section 1 — section 1 should be finished, section 2 and log should not
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/1", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)

	if log.Sections[0].Status != models.WorkoutLogStatusFinished {
		t.Errorf("section 1 should be finished, got %s", log.Sections[0].Status)
	}
	if log.Sections[1].Status != models.WorkoutLogStatusInProgress {
		t.Errorf("section 2 should still be in_progress, got %s", log.Sections[1].Status)
	}
	if log.Status != models.WorkoutLogStatusInProgress {
		t.Errorf("log should still be in_progress, got %s", log.Status)
	}

	// Finish the set in section 2 — everything should propagate to finished
	doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/2", map[string]any{
		"status": "finished", "actualReps": 5, "actualWeight": 100.0,
	})

	w = doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	json.Unmarshal(w.Body.Bytes(), &log)

	if log.Sections[1].Status != models.WorkoutLogStatusFinished {
		t.Errorf("section 2 should be finished, got %s", log.Sections[1].Status)
	}
	if log.Status != models.WorkoutLogStatusFinished {
		t.Errorf("log should be finished, got %s", log.Status)
	}
}

func TestUpdateLogPreservesStatus(t *testing.T) {
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
		"owner": "alice", "name": "Test",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})

	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	// Try to change status via generic PATCH — should be ignored
	w := doJSON(r, "PATCH", "/api/user/workout-logs/1", map[string]any{
		"name": "Updated", "status": "planning",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	var log models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &log)
	if log.Status != models.WorkoutLogStatusInProgress {
		t.Errorf("status should still be in_progress, got %s", log.Status)
	}
	if log.Name != "Updated" {
		t.Errorf("name should be updated, got %s", log.Name)
	}
}

func TestListWorkoutLogsStatusFilter(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "squat", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 1, "reps": 5, "weight": 100.0,
	})

	// Create two logs
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Log 1",
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Log 2",
	})

	// Start log 1
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 1, "sourceExerciseSchemeId": 1, "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-logs/1/start", nil)

	t.Run("filter planning", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-logs?status=planning", nil)
		var result []models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Fatalf("expected 1 planning log, got %d", len(result))
		}
		if result[0].Name != "Log 2" {
			t.Errorf("expected Log 2, got %s", result[0].Name)
		}
	})

	t.Run("filter in_progress", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-logs?status=in_progress", nil)
		var result []models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Fatalf("expected 1 in_progress log, got %d", len(result))
		}
		if result[0].Name != "Log 1" {
			t.Errorf("expected Log 1, got %s", result[0].Name)
		}
	})
}

func TestOwnerAuthorization(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create a log as alice (via fallback user)
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"name": "Alice's Log",
	})

	// Manually insert a log owned by bob (bypass handler to simulate another user)
	database.DB.Create(&models.WorkoutLogEntity{
		Owner:  "bob",
		Name:   "Bob's Log",
		Status: models.WorkoutLogStatusPlanning,
	})

	t.Run("cannot get another user's log", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-logs/2", nil)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("cannot update another user's log", func(t *testing.T) {
		w := doJSON(r, "PATCH", "/api/user/workout-logs/2", map[string]any{
			"name": "Hacked",
		})
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("cannot delete another user's log", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/workout-logs/2", nil)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("cannot start another user's log", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/workout-logs/2/start", nil)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("list only shows own logs", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/workout-logs", nil)
		var result []models.WorkoutLog
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Fatalf("expected 1 log (only alice's), got %d", len(result))
		}
		if result[0].Owner != "alice" {
			t.Errorf("expected alice, got %s", result[0].Owner)
		}
	})
}
