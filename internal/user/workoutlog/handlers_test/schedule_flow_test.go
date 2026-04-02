package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	exercisemodels "gesitr/internal/compendium/exercise/models"
	workoutmodels "gesitr/internal/compendium/workout/models"
	"gesitr/internal/user/workoutlog/models"
	schedulemodels "gesitr/internal/user/workoutschedule/models"
)

// TestScheduleFixedDateFullFlow exercises the complete lifecycle:
// create exercises → create workout → create schedule → create period →
// create commitments (fixed_date) → list logs (triggers activation) →
// commit → start → complete sets → finished.
func TestScheduleFixedDateFullFlow(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	// 1. Create exercise + scheme
	w := doJSONLog(t, r, "POST", "/api/exercises", map[string]any{
		"name": "Barbell Row", "type": "STRENGTH", "technicalDifficulty": "intermediate",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create exercise: %d", w.Code)
	}
	var exercise exercisemodels.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)

	w = doJSONLog(t, r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": exercise.ID, "measurementType": "REP_BASED",
		"sets": 3, "reps": 8, "weight": 70.0, "restBetweenSets": 120,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create scheme: %d", w.Code)
	}
	var scheme exercisemodels.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme)

	// 2. Create a workout template
	w = doJSONLog(t, r, "POST", "/api/workouts", map[string]any{"name": "Back Day"})
	if w.Code != http.StatusCreated {
		t.Fatalf("create workout: %d", w.Code)
	}
	var workout workoutmodels.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)

	// 3. Create a schedule (proposed)
	startDate := time.Now().AddDate(0, 0, -3)
	w = doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId":     workout.ID,
		"startDate":     startDate.Format(time.RFC3339),
		"initialStatus": "proposed",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create schedule: %d, %s", w.Code, w.Body.String())
	}
	var schedule schedulemodels.WorkoutSchedule
	json.Unmarshal(w.Body.Bytes(), &schedule)

	// 4. Create the first period (started yesterday, ends in 6 days)
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	periodEnd := time.Now().AddDate(0, 0, 6).Truncate(24 * time.Hour)
	w = doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId":  schedule.ID,
		"periodStart": yesterday.Format(time.RFC3339),
		"periodEnd":   periodEnd.Format(time.RFC3339),
		"type":        "fixed_date",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create period: %d, %s", w.Code, w.Body.String())
	}
	var period schedulemodels.SchedulePeriod
	json.Unmarshal(w.Body.Bytes(), &period)

	// 5. Create commitments (toggle day 1 and day 3)
	day1 := yesterday
	day3 := yesterday.AddDate(0, 0, 2)
	for _, d := range []time.Time{day1, day3} {
		w = doJSONLog(t, r, "POST", "/api/user/schedule-commitments", map[string]any{
			"periodId": period.ID,
			"date":     d.Format(time.RFC3339),
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("create commitment: %d, %s", w.Code, w.Body.String())
		}
	}

	// 6. List workout logs — triggers activation
	w = doJSONLog(t, r, "GET", "/api/user/workout-logs", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("list logs: %d", w.Code)
	}
	var allLogs []models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &allLogs)

	// Should have 2 proposed logs from the schedule
	var proposedLogs []models.WorkoutLog
	for _, l := range allLogs {
		if l.ScheduleID != nil && *l.ScheduleID == schedule.ID {
			proposedLogs = append(proposedLogs, l)
		}
	}
	if len(proposedLogs) != 2 {
		t.Fatalf("expected 2 schedule-generated logs, got %d", len(proposedLogs))
	}

	firstLog := proposedLogs[0]
	if firstLog.Status != models.WorkoutLogStatusProposed {
		t.Errorf("expected proposed status, got %s", firstLog.Status)
	}
	if firstLog.Name != "Back Day" {
		t.Errorf("expected workout name, got %s", firstLog.Name)
	}

	// 7. Configure exercises on the proposed log
	logID := itoa(firstLog.ID)
	w = doJSONLog(t, r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": firstLog.ID, "type": "main", "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create log section: %d, %s", w.Code, w.Body.String())
	}
	var logSection models.WorkoutLogSection
	json.Unmarshal(w.Body.Bytes(), &logSection)

	w = doJSONLog(t, r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": logSection.ID, "sourceExerciseSchemeId": scheme.ID, "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create log exercise: %d, %s", w.Code, w.Body.String())
	}
	var logExercise models.WorkoutLogExercise
	json.Unmarshal(w.Body.Bytes(), &logExercise)

	// 8. Commit → Start → Complete
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/"+logID+"/commit", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("commit: %d", w.Code)
	}

	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/"+logID+"/start", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("start: %d", w.Code)
	}

	for i, s := range logExercise.Sets {
		w = doJSONLog(t, r, "PATCH", "/api/user/workout-log-exercise-sets/"+itoa(s.ID), map[string]any{
			"status": "finished", "actualReps": 8, "actualWeight": 70.0,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("finish set %d: %d", i+1, w.Code)
		}
	}

	// 9. Verify final state
	w = doJSONLog(t, r, "GET", "/api/user/workout-logs/"+logID, nil)
	var finalLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &finalLog)

	if finalLog.Sections[0].Exercises[0].Status != models.WorkoutLogItemStatusFinished {
		t.Errorf("exercise should be finished, got %s", finalLog.Sections[0].Exercises[0].Status)
	}
	if finalLog.ScheduleID == nil || *finalLog.ScheduleID != schedule.ID {
		t.Error("scheduleID should be preserved")
	}
}

// TestScheduleDeleteOrphansLogs verifies that deleting a schedule orphans
// its workout logs rather than deleting them.
func TestScheduleDeleteOrphansLogs(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	doJSONLog(t, r, "POST", "/api/workouts", map[string]any{"name": "Orphan Test"})

	// Create schedule + period + commitment
	startDate := time.Now().AddDate(0, 0, -3)
	w := doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 1, "startDate": startDate.Format(time.RFC3339),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create schedule: %d, %s", w.Code, w.Body.String())
	}
	var schedule schedulemodels.WorkoutSchedule
	json.Unmarshal(w.Body.Bytes(), &schedule)

	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	periodEnd := time.Now().AddDate(0, 0, 6).Truncate(24 * time.Hour)
	doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId": schedule.ID, "periodStart": yesterday.Format(time.RFC3339), "periodEnd": periodEnd.Format(time.RFC3339), "type": "fixed_date",
	})
	doJSONLog(t, r, "POST", "/api/user/schedule-commitments", map[string]any{
		"periodId": 1, "date": yesterday.Format(time.RFC3339),
	})

	// Trigger activation
	doJSONLog(t, r, "GET", "/api/user/workout-logs", nil)

	// Verify logs exist
	w = doJSONLog(t, r, "GET", "/api/user/workout-logs", nil)
	var logsBefore []models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &logsBefore)
	if len(logsBefore) == 0 {
		t.Fatal("expected logs before deletion")
	}

	// Delete schedule
	w = doJSONLog(t, r, "DELETE", "/api/user/workout-schedules/"+itoa(schedule.ID), nil)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete: expected 204, got %d", w.Code)
	}

	// Logs should survive with schedule_id = null
	w = doJSONLog(t, r, "GET", "/api/user/workout-logs", nil)
	var logsAfter []models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &logsAfter)

	if len(logsAfter) != len(logsBefore) {
		t.Errorf("logs should survive: had %d, now %d", len(logsBefore), len(logsAfter))
	}
	for _, l := range logsAfter {
		if l.ScheduleID != nil {
			t.Error("schedule_id should be null after deletion")
		}
	}
}

// TestScheduleCRUD tests basic CRUD operations.
func TestScheduleCRUD(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	doJSONLog(t, r, "POST", "/api/workouts", map[string]any{"name": "CRUD Test"})

	// Create
	w := doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 1,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create: %d, %s", w.Code, w.Body.String())
	}
	var created schedulemodels.WorkoutSchedule
	json.Unmarshal(w.Body.Bytes(), &created)

	// Get
	w = doJSONLog(t, r, "GET", "/api/user/workout-schedules/"+itoa(created.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get: %d", w.Code)
	}

	// List
	w = doJSONLog(t, r, "GET", "/api/user/workout-schedules?workoutId=1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("list: %d", w.Code)
	}
	var listed []schedulemodels.WorkoutSchedule
	json.Unmarshal(w.Body.Bytes(), &listed)
	if len(listed) != 1 {
		t.Errorf("expected 1, got %d", len(listed))
	}

	// Delete
	w = doJSONLog(t, r, "DELETE", "/api/user/workout-schedules/"+itoa(created.ID), nil)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete: expected 204, got %d", w.Code)
	}
}

// TestScheduleValidation tests that invalid creation is rejected.
func TestScheduleValidation(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	doJSONLog(t, r, "POST", "/api/workouts", map[string]any{"name": "Validation"})

	// Invalid initialStatus
	w := doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 1, "initialStatus": "finished",
	})
	if w.Code != http.StatusBadRequest {
		t.Errorf("invalid initialStatus should be 400, got %d", w.Code)
	}
}

// TestScheduleIdempotentGeneration verifies listing logs multiple times
// doesn't create duplicate schedule-generated logs.
func TestScheduleIdempotentGeneration(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)
	r := newRouter()

	doJSONLog(t, r, "POST", "/api/workouts", map[string]any{"name": "Idempotent"})

	startDate := time.Now().AddDate(0, 0, -3)
	doJSONLog(t, r, "POST", "/api/user/workout-schedules", map[string]any{
		"workoutId": 1, "startDate": startDate.Format(time.RFC3339),
	})

	// Create period + commitment
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	periodEnd := time.Now().AddDate(0, 0, 6).Truncate(24 * time.Hour)
	doJSONLog(t, r, "POST", "/api/user/schedule-periods", map[string]any{
		"scheduleId": 1, "periodStart": yesterday.Format(time.RFC3339), "periodEnd": periodEnd.Format(time.RFC3339), "type": "fixed_date",
	})
	doJSONLog(t, r, "POST", "/api/user/schedule-commitments", map[string]any{
		"periodId": 1, "date": yesterday.Format(time.RFC3339),
	})

	// List logs 3 times
	for i := 0; i < 3; i++ {
		doJSONLog(t, r, "GET", "/api/user/workout-logs", nil)
	}

	w := doJSONLog(t, r, "GET", "/api/user/workout-logs", nil)
	var logs []models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &logs)

	if len(logs) != 1 {
		t.Errorf("expected exactly 1 log (idempotent), got %d", len(logs))
	}
}
