package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	exercisemodels "gesitr/internal/compendium/exercise/models"
	workoutmodels "gesitr/internal/compendium/workout/models"
	"gesitr/internal/database"
	"gesitr/internal/user/workoutlog/models"
)

// TestCommitmentHappyPath exercises the full proposed → committed → in_progress → finished flow
// with a workout template, exercise schemes, and set tracking.
func TestCommitmentHappyPath(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// 1. Create exercise + scheme
	w := doJSONLog(t, r, "POST", "/api/exercises", map[string]any{
		"names": []string{"Deadlift"}, "type": "STRENGTH", "technicalDifficulty": "intermediate",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create exercise: %d", w.Code)
	}
	var exercise exercisemodels.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)

	w = doJSONLog(t, r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": exercise.ID, "measurementType": "REP_BASED",
		"sets": 3, "reps": 5, "weight": 180.0, "restBetweenSets": 180,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create scheme: %d", w.Code)
	}
	var scheme exercisemodels.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme)

	// 2. Create a workout template
	w = doJSONLog(t, r, "POST", "/api/workouts", map[string]any{"name": "Pull Day"})
	if w.Code != http.StatusCreated {
		t.Fatalf("create workout: %d", w.Code)
	}
	var workout workoutmodels.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)

	w = doJSONLog(t, r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": workout.ID, "type": "main", "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section: %d", w.Code)
	}
	var section workoutmodels.WorkoutSection
	json.Unmarshal(w.Body.Bytes(), &section)

	w = doJSONLog(t, r, "POST", "/api/workout-section-items", map[string]any{
		"workoutSectionId": section.ID, "type": "exercise", "exerciseSchemeId": scheme.ID, "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section item: %d", w.Code)
	}

	// 3. Create a proposed commitment for next week
	dueStart := time.Now().Add(7 * 24 * time.Hour).Truncate(time.Second)
	dueEnd := dueStart.Add(24 * time.Hour)

	w = doJSONLog(t, r, "POST", "/api/user/workout-logs", map[string]any{
		"name":      "Pull Day - Next Monday",
		"workoutId": workout.ID,
		"status":    "proposed",
		"dueStart":  dueStart.Format(time.RFC3339),
		"dueEnd":    dueEnd.Format(time.RFC3339),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create proposed log: %d", w.Code)
	}
	var proposedLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &proposedLog)

	if proposedLog.Status != models.WorkoutLogStatusProposed {
		t.Errorf("expected proposed, got %s", proposedLog.Status)
	}
	if proposedLog.DueStart == nil || proposedLog.DueEnd == nil {
		t.Fatal("due window should be set")
	}
	if proposedLog.WorkoutID == nil || *proposedLog.WorkoutID != workout.ID {
		t.Error("workoutId mismatch")
	}

	// 4. User configures exercise schemes (add section + exercise to the log)
	w = doJSONLog(t, r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": proposedLog.ID, "type": "main", "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create log section: %d", w.Code)
	}
	var logSection models.WorkoutLogSection
	json.Unmarshal(w.Body.Bytes(), &logSection)

	w = doJSONLog(t, r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": logSection.ID, "sourceExerciseSchemeId": scheme.ID, "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create log exercise: %d", w.Code)
	}
	var logExercise models.WorkoutLogExercise
	json.Unmarshal(w.Body.Bytes(), &logExercise)

	if len(logExercise.Sets) != 3 {
		t.Fatalf("expected 3 auto-created sets, got %d", len(logExercise.Sets))
	}

	// 5. Commit the proposal
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/"+itoa(proposedLog.ID)+"/commit", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("commit: %d", w.Code)
	}
	var committedLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &committedLog)

	if committedLog.Status != models.WorkoutLogStatusCommitted {
		t.Errorf("expected committed, got %s", committedLog.Status)
	}

	// 6. Start the workout (committed → in_progress)
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/"+itoa(proposedLog.ID)+"/start", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("start: %d", w.Code)
	}
	var startedLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &startedLog)

	if startedLog.Status != models.WorkoutLogStatusInProgress {
		t.Errorf("expected in_progress, got %s", startedLog.Status)
	}
	// Due window should be preserved
	if startedLog.DueStart == nil || startedLog.DueEnd == nil {
		t.Error("due window should persist after starting")
	}

	// 7. Complete all sets
	for i, s := range logExercise.Sets {
		w = doJSONLog(t, r, "PATCH", "/api/user/workout-log-exercise-sets/"+itoa(s.ID), map[string]any{
			"status": "finished", "actualReps": 5, "actualWeight": 180.0,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("finish set %d: %d", i+1, w.Code)
		}
	}

	// 8. Verify final state
	w = doJSONLog(t, r, "GET", "/api/user/workout-logs/"+itoa(proposedLog.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get final log: %d", w.Code)
	}
	var finalLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &finalLog)

	if finalLog.Sections[0].Exercises[0].Status != models.WorkoutLogItemStatusFinished {
		t.Errorf("exercise should be finished, got %s", finalLog.Sections[0].Exercises[0].Status)
	}
	if finalLog.DueStart == nil || finalLog.DueEnd == nil {
		t.Error("due window should still be on the finished log")
	}
}

// TestCommitmentSkipFlow tests the proposed → skipped path.
func TestCommitmentSkipFlow(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	// 1. Create a proposed commitment
	w := doJSONLog(t, r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "Skip Me",
		"status":   "proposed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create proposed: %d", w.Code)
	}

	// 2. Skip it
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/1/skip", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("skip: %d", w.Code)
	}
	var skippedLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &skippedLog)

	if skippedLog.Status != models.WorkoutLogStatusSkipped {
		t.Errorf("expected skipped, got %s", skippedLog.Status)
	}

	// 3. Skipped is terminal — cannot commit or start
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/1/commit", nil)
	if w.Code != http.StatusConflict {
		t.Errorf("commit after skip should be 409, got %d", w.Code)
	}

	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/1/start", nil)
	if w.Code != http.StatusConflict {
		t.Errorf("start after skip should be 409, got %d", w.Code)
	}
}

// TestCommitmentBrokenFlow tests the committed → broken path via the ticker logic.
func TestCommitmentBrokenFlow(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// 1. Create a committed log with due_end in the past (simulating time passing)
	pastStart := time.Now().Add(-48 * time.Hour)
	pastEnd := time.Now().Add(-1 * time.Hour)

	w := doJSONLog(t, r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "Overdue Workout",
		"status":   "committed",
		"dueStart": pastStart.Format(time.RFC3339),
		"dueEnd":   pastEnd.Format(time.RFC3339),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create committed: %d", w.Code)
	}

	// 2. Run the ticker logic (same UPDATE query as the goroutine)
	now := time.Now()
	result := database.DB.Model(&models.WorkoutLogEntity{}).
		Where("status = ? AND due_end < ?", models.WorkoutLogStatusCommitted, now).
		Updates(map[string]any{
			"status":            models.WorkoutLogStatusBroken,
			"status_changed_at": now,
		})
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.RowsAffected != 1 {
		t.Errorf("expected 1 row affected, got %d", result.RowsAffected)
	}

	// 3. Verify via API
	w = doJSONLog(t, r, "GET", "/api/user/workout-logs/1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get: %d", w.Code)
	}
	var brokenLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &brokenLog)

	if brokenLog.Status != models.WorkoutLogStatusBroken {
		t.Errorf("expected broken, got %s", brokenLog.Status)
	}

	// 4. Broken is terminal — cannot start
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/1/start", nil)
	if w.Code != http.StatusConflict {
		t.Errorf("start after broken should be 409, got %d", w.Code)
	}
}

// TestMultipleCommitmentsCoexist verifies that multiple proposed/committed logs
// can exist for the same workout simultaneously.
func TestMultipleCommitmentsCoexist(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// Create a workout
	w := doJSONLog(t, r, "POST", "/api/workouts", map[string]any{"name": "Leg Day"})
	if w.Code != http.StatusCreated {
		t.Fatalf("create workout: %d", w.Code)
	}
	var workout workoutmodels.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)

	// Create 3 commitments for different weeks
	for i := 0; i < 3; i++ {
		dueStart := time.Now().Add(time.Duration(i*7*24) * time.Hour)
		dueEnd := dueStart.Add(24 * time.Hour)

		w = doJSONLog(t, r, "POST", "/api/user/workout-logs", map[string]any{
			"name":      "Leg Day commitment",
			"workoutId": workout.ID,
			"status":    "proposed",
			"dueStart":  dueStart.Format(time.RFC3339),
			"dueEnd":    dueEnd.Format(time.RFC3339),
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("create commitment %d: %d, %s", i+1, w.Code, w.Body.String())
		}
	}

	// Verify all 3 exist
	w = doJSONLog(t, r, "GET", "/api/user/workout-logs?workoutId="+itoa(workout.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("list: %d", w.Code)
	}
	var logs []models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &logs)

	if len(logs) != 3 {
		t.Errorf("expected 3 commitment logs, got %d", len(logs))
	}

	// Skip one, commit one, leave one proposed
	doJSONLog(t, r, "POST", "/api/user/workout-logs/1/skip", nil)
	doJSONLog(t, r, "POST", "/api/user/workout-logs/2/commit", nil)

	// Filter by status
	w = doJSONLog(t, r, "GET", "/api/user/workout-logs?status=proposed", nil)
	var proposed []models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &proposed)
	if len(proposed) != 1 {
		t.Errorf("expected 1 proposed log, got %d", len(proposed))
	}

	w = doJSONLog(t, r, "GET", "/api/user/workout-logs?status=committed", nil)
	var committed []models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &committed)
	if len(committed) != 1 {
		t.Errorf("expected 1 committed log, got %d", len(committed))
	}

	w = doJSONLog(t, r, "GET", "/api/user/workout-logs?status=skipped", nil)
	var skipped []models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &skipped)
	if len(skipped) != 1 {
		t.Errorf("expected 1 skipped log, got %d", len(skipped))
	}
}

// TestCommitmentInvalidTransitions verifies that forbidden state transitions are rejected.
func TestCommitmentInvalidTransitions(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	dueStart := time.Now().Add(24 * time.Hour)
	dueEnd := time.Now().Add(48 * time.Hour)

	// Create a committed log
	doJSONLog(t, r, "POST", "/api/user/workout-logs", map[string]any{
		"name":     "Committed",
		"status":   "committed",
		"dueStart": dueStart.Format(time.RFC3339),
		"dueEnd":   dueEnd.Format(time.RFC3339),
	})

	// Cannot skip a committed log (only proposed → skipped is valid)
	w := doJSONLog(t, r, "POST", "/api/user/workout-logs/1/skip", nil)
	if w.Code != http.StatusConflict {
		t.Errorf("skip committed should be 409, got %d", w.Code)
	}

	// Cannot commit an already-committed log
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/1/commit", nil)
	if w.Code != http.StatusConflict {
		t.Errorf("commit committed should be 409, got %d", w.Code)
	}

	// Cannot abandon a committed log (abandon is in_progress → aborted)
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/1/abandon", nil)
	if w.Code != http.StatusConflict {
		t.Errorf("abandon committed should be 409, got %d", w.Code)
	}
}
