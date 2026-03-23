package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	exercisemodels "gesitr/internal/exercise/models"
	workoutmodels "gesitr/internal/user/workout/models"
	"gesitr/internal/user/workoutlog/models"
)

func TestFullWorkoutToLogFlow(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// 1. Create an exercise
	w := doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Barbell Squat", "type": "STRENGTH", "technicalDifficulty": "intermediate",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create user exercise: status = %d, body = %s", w.Code, w.Body.String())
	}
	var userExercise exercisemodels.Exercise
	json.Unmarshal(w.Body.Bytes(), &userExercise)

	// 2. Create an exercise scheme for that user exercise
	w = doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId":      userExercise.ID,
		"measurementType": "REP_BASED",
		"sets":            5,
		"reps":            5,
		"weight":          140.0,
		"restBetweenSets": 180,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create scheme: status = %d, body = %s", w.Code, w.Body.String())
	}
	var scheme exercisemodels.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme)

	// 3. Create a second exercise + scheme for variety
	w = doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Bench Press", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create second user exercise: status = %d", w.Code)
	}
	var userExercise2 exercisemodels.Exercise
	json.Unmarshal(w.Body.Bytes(), &userExercise2)

	w = doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId":      userExercise2.ID,
		"measurementType": "REP_BASED",
		"sets":            4,
		"reps":            8,
		"weight":          80.0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create second scheme: status = %d", w.Code)
	}
	var scheme2 exercisemodels.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme2)

	// 4. Create a workout template
	w = doJSON(r, "POST", "/api/user/workouts", map[string]any{
		"owner": "alice", "name": "Strength Day", "date": "2026-03-07T10:00:00Z",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create workout: status = %d", w.Code)
	}
	var wkt workoutmodels.Workout
	json.Unmarshal(w.Body.Bytes(), &wkt)

	// 5. Add sections to the workout template (with restBetweenExercises)
	w = doJSON(r, "POST", "/api/user/workout-sections", map[string]any{
		"workoutId": wkt.ID, "type": "supplementary", "label": "Warmup", "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create warmup section: status = %d", w.Code)
	}

	w = doJSON(r, "POST", "/api/user/workout-sections", map[string]any{
		"workoutId": wkt.ID, "type": "main", "position": 1, "restBetweenExercises": 120,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create main section: status = %d", w.Code)
	}
	var mainSection workoutmodels.WorkoutSection
	json.Unmarshal(w.Body.Bytes(), &mainSection)

	// 6. Add exercises to the workout template sections
	w = doJSON(r, "POST", "/api/user/workout-section-exercises", map[string]any{
		"workoutSectionId": mainSection.ID, "exerciseSchemeId": scheme.ID, "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section exercise 1: status = %d", w.Code)
	}

	w = doJSON(r, "POST", "/api/user/workout-section-exercises", map[string]any{
		"workoutSectionId": mainSection.ID, "exerciseSchemeId": scheme2.ID, "position": 1,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section exercise 2: status = %d", w.Code)
	}

	// 7. Verify the workout template is fully loaded
	w = doJSON(r, "GET", "/api/user/workouts/"+itoa(wkt.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get workout: status = %d", w.Code)
	}
	json.Unmarshal(w.Body.Bytes(), &wkt)
	if len(wkt.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(wkt.Sections))
	}
	if len(wkt.Sections[1].Exercises) != 2 {
		t.Fatalf("expected 2 exercises in main section, got %d", len(wkt.Sections[1].Exercises))
	}
	if wkt.Sections[1].RestBetweenExercises == nil || *wkt.Sections[1].RestBetweenExercises != 120 {
		t.Error("workout section restBetweenExercises mismatch")
	}

	// 8. Create a workout log referencing the template
	w = doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Strength Day - March 7", "date": "2026-03-07T18:00:00Z",
		"workoutId": wkt.ID,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create workout log: status = %d", w.Code)
	}
	var workoutLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &workoutLog)
	if workoutLog.WorkoutID == nil || *workoutLog.WorkoutID != wkt.ID {
		t.Error("workout log workoutId mismatch")
	}
	if workoutLog.Status != models.WorkoutLogStatusPlanning {
		t.Errorf("expected planning status, got %s", workoutLog.Status)
	}

	// 9. Add sections to the log (mirroring the template but independent)
	w = doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": workoutLog.ID, "type": "main", "position": 0, "restBetweenExercises": 90,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create log section: status = %d", w.Code)
	}
	var logSection models.WorkoutLogSection
	json.Unmarshal(w.Body.Bytes(), &logSection)

	// 10. Add exercises to the log section — targets should be snapshotted from scheme
	w = doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": logSection.ID, "sourceExerciseSchemeId": scheme.ID, "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create log exercise 1: status = %d, body = %s", w.Code, w.Body.String())
	}
	var logExercise1 models.WorkoutLogExercise
	json.Unmarshal(w.Body.Bytes(), &logExercise1)

	// Verify target snapshot from scheme
	if logExercise1.TargetMeasurementType != "REP_BASED" {
		t.Errorf("expected REP_BASED, got %s", logExercise1.TargetMeasurementType)
	}
	// Exercise-level BreakAfterSeconds should come from section's RestBetweenExercises
	if logExercise1.BreakAfterSeconds == nil || *logExercise1.BreakAfterSeconds != 90 {
		t.Errorf("exercise breakAfterSeconds: expected 90, got %v", logExercise1.BreakAfterSeconds)
	}
	// Should have 5 auto-created sets with snapshotted targets
	if len(logExercise1.Sets) != 5 {
		t.Fatalf("expected 5 sets, got %d", len(logExercise1.Sets))
	}
	for i, s := range logExercise1.Sets {
		if s.TargetReps == nil || *s.TargetReps != 5 {
			t.Errorf("set %d: target reps not snapshotted correctly", i+1)
		}
		if s.TargetWeight == nil || *s.TargetWeight != 140.0 {
			t.Errorf("set %d: target weight not snapshotted correctly", i+1)
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

	w = doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": logSection.ID, "sourceExerciseSchemeId": scheme2.ID, "position": 1,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create log exercise 2: status = %d", w.Code)
	}
	var logExercise2 models.WorkoutLogExercise
	json.Unmarshal(w.Body.Bytes(), &logExercise2)

	// Second exercise should have 4 sets from scheme2
	if len(logExercise2.Sets) != 4 {
		t.Fatalf("expected 4 sets for exercise 2, got %d", len(logExercise2.Sets))
	}
	if logExercise2.Sets[0].TargetWeight == nil || *logExercise2.Sets[0].TargetWeight != 80.0 {
		t.Error("second exercise set target weight not snapshotted from scheme2")
	}

	// Start the workout log (planning -> in_progress)
	w = doJSON(r, "POST", "/api/user/workout-logs/"+itoa(workoutLog.ID)+"/start", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("start workout log: status = %d, body = %s", w.Code, w.Body.String())
	}

	// 11. Record actual performance — update individual sets
	for i, s := range logExercise1.Sets {
		w = doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/"+itoa(s.ID), map[string]any{
			"status": "finished", "actualReps": 5, "actualWeight": 140.0,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("update set %d: status = %d, body = %s", i+1, w.Code, w.Body.String())
		}
	}

	// Partial completion for exercise 2 — only complete 3 of 4 sets, with reduced reps
	for i, s := range logExercise2.Sets {
		if i >= 3 {
			break // skip last set (not attempted)
		}
		w = doJSON(r, "PATCH", "/api/user/workout-log-exercise-sets/"+itoa(s.ID), map[string]any{
			"status": "finished", "actualReps": 6, "actualWeight": 75.0,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("update exercise 2 set %d: status = %d", i+1, w.Code)
		}
	}

	// 12. Verify the full workout log with nested preloads (including sets)
	w = doJSON(r, "GET", "/api/user/workout-logs/"+itoa(workoutLog.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get workout log: status = %d", w.Code)
	}
	var fullLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &fullLog)

	if fullLog.WorkoutID == nil || *fullLog.WorkoutID != wkt.ID {
		t.Error("full log workoutId mismatch")
	}
	if len(fullLog.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(fullLog.Sections))
	}
	if fullLog.Sections[0].RestBetweenExercises == nil || *fullLog.Sections[0].RestBetweenExercises != 90 {
		t.Error("log section restBetweenExercises mismatch")
	}
	if len(fullLog.Sections[0].Exercises) != 2 {
		t.Fatalf("expected 2 exercises, got %d", len(fullLog.Sections[0].Exercises))
	}

	ex1 := fullLog.Sections[0].Exercises[0]
	ex2 := fullLog.Sections[0].Exercises[1]

	// Exercise 1: all 5 sets finished — exercise should be finished
	if len(ex1.Sets) != 5 {
		t.Fatalf("expected 5 sets for exercise 1, got %d", len(ex1.Sets))
	}
	if ex1.Status != models.WorkoutLogItemStatusFinished {
		t.Errorf("exercise 1 should be finished (all sets done), got %s", ex1.Status)
	}
	for i, s := range ex1.Sets {
		if s.Status != models.WorkoutLogItemStatusFinished {
			t.Errorf("exercise 1 set %d should be finished", i+1)
		}
		if s.ExerciseLog == nil || s.ExerciseLog.Weight == nil || *s.ExerciseLog.Weight != 140.0 {
			t.Errorf("exercise 1 set %d exerciseLog weight mismatch", i+1)
		}
		// Target fields must still be intact
		if s.TargetWeight == nil || *s.TargetWeight != 140.0 {
			t.Errorf("exercise 1 set %d target weight changed after update", i+1)
		}
	}

	// Exercise 2: 3 of 4 sets finished — exercise should still be in_progress
	if len(ex2.Sets) != 4 {
		t.Fatalf("expected 4 sets for exercise 2, got %d", len(ex2.Sets))
	}
	if ex2.Status == models.WorkoutLogItemStatusFinished {
		t.Error("exercise 2 should not be finished (1 set remaining)")
	}
	finishedCount := 0
	for _, s := range ex2.Sets {
		if s.Status == models.WorkoutLogItemStatusFinished {
			finishedCount++
			if s.ExerciseLog == nil || s.ExerciseLog.Reps == nil || *s.ExerciseLog.Reps != 6 {
				t.Error("exercise 2 finished set exerciseLog reps mismatch")
			}
		}
	}
	if finishedCount != 3 {
		t.Errorf("expected 3 finished sets for exercise 2, got %d", finishedCount)
	}
	// Last set should still be in_progress
	if ex2.Sets[3].Status != models.WorkoutLogItemStatusInProgress {
		t.Errorf("exercise 2 last set should be in_progress, got %s", ex2.Sets[3].Status)
	}

	// Section and log should not be finished (exercise 2 is incomplete)
	if fullLog.Sections[0].Status == models.WorkoutLogItemStatusFinished {
		t.Error("section should not be finished (exercise 2 incomplete)")
	}
	if fullLog.Status == models.WorkoutLogStatusFinished {
		t.Error("log should not be finished (section incomplete)")
	}

	// 13. Also create an ad-hoc log (no workout template)
	w = doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Quick Session", "date": "2026-03-08T10:00:00Z",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create ad-hoc log: status = %d", w.Code)
	}
	var adHocLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &adHocLog)
	if adHocLog.WorkoutID != nil {
		t.Error("ad-hoc log should have nil workoutId")
	}

	// 14. Verify filtering — only template-linked logs
	w = doJSON(r, "GET", "/api/user/workout-logs?workoutId="+itoa(wkt.ID), nil)
	var filteredLogs []models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &filteredLogs)
	if len(filteredLogs) != 1 || *filteredLogs[0].WorkoutID != wkt.ID {
		t.Errorf("workoutId filter: expected 1 log linked to workout %d, got %d", wkt.ID, len(filteredLogs))
	}

	// 15. Verify owner filtering returns all of alice's logs
	w = doJSON(r, "GET", "/api/user/workout-logs?owner=alice", nil)
	var aliceLogs []models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &aliceLogs)
	if len(aliceLogs) != 2 {
		t.Errorf("expected 2 logs for alice, got %d", len(aliceLogs))
	}
}

func TestGetWorkoutLogWithSectionsAndExercises(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Setup: exercise -> scheme -> workout log -> section -> exercise
	doJSON(r, "POST", "/api/exercises", map[string]any{
		"name": "Bench Press", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	doJSON(r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED", "sets": 3, "reps": 10,
	})
	doJSON(r, "POST", "/api/user/workout-logs", map[string]any{
		"owner": "alice", "name": "Full Session", "date": "2026-03-07T10:00:00Z",
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "supplementary", "label": "Warmup", "position": 0,
	})
	doJSON(r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": 1, "type": "main", "position": 1,
	})
	doJSON(r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId": 2, "sourceExerciseSchemeId": 1, "position": 0,
	})

	w := doJSON(r, "GET", "/api/user/workout-logs/1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}

	var result models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &result)

	if len(result.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(result.Sections))
	}
	if result.Sections[0].Position != 0 || result.Sections[1].Position != 1 {
		t.Error("sections not ordered by position")
	}
	if len(result.Sections[1].Exercises) != 1 {
		t.Fatalf("expected 1 exercise in main section, got %d", len(result.Sections[1].Exercises))
	}
	if result.Sections[1].Exercises[0].SourceExerciseSchemeID != 1 {
		t.Error("exercise scheme ID mismatch")
	}
}
