package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	exercisemodels "gesitr/internal/exercise/models"
	workoutmodels "gesitr/internal/user/workout/models"
	"gesitr/internal/user/workoutlog/models"
)

func TestWorkoutStartFlow(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// -- Setup: create user exercises and schemes --

	w := doJSONLog(t, r, "POST", "/api/exercises", map[string]any{
		"name": "Barbell Squat", "type": "STRENGTH", "technicalDifficulty": "intermediate",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create user exercise 1: status = %d", w.Code)
	}
	var userExercise1 exercisemodels.Exercise
	json.Unmarshal(w.Body.Bytes(), &userExercise1)

	w = doJSONLog(t, r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId":      userExercise1.ID,
		"measurementType": "REP_BASED",
		"sets":            5,
		"reps":            5,
		"weight":          140.0,
		"restBetweenSets": 180,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create scheme 1: status = %d, body = %s", w.Code, w.Body.String())
	}
	var scheme1 exercisemodels.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme1)

	w = doJSONLog(t, r, "POST", "/api/exercises", map[string]any{
		"name": "Bench Press", "type": "STRENGTH", "technicalDifficulty": "beginner",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create user exercise 2: status = %d", w.Code)
	}
	var userExercise2 exercisemodels.Exercise
	json.Unmarshal(w.Body.Bytes(), &userExercise2)

	w = doJSONLog(t, r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId":      userExercise2.ID,
		"measurementType": "REP_BASED",
		"sets":            4,
		"reps":            8,
		"weight":          80.0,
		"restBetweenSets": 120,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create scheme 2: status = %d, body = %s", w.Code, w.Body.String())
	}
	var scheme2 exercisemodels.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme2)

	// -- Setup: create workout template with sections and exercises --

	w = doJSONLog(t, r, "POST", "/api/workouts", map[string]any{
		"name": "Strength Day",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create workout: status = %d", w.Code)
	}
	var wkt workoutmodels.Workout
	json.Unmarshal(w.Body.Bytes(), &wkt)

	w = doJSONLog(t, r, "POST", "/api/workout-sections", map[string]any{
		"workoutId": wkt.ID, "type": "main", "position": 0, "restBetweenExercises": 90,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section: status = %d", w.Code)
	}
	var section workoutmodels.WorkoutSection
	json.Unmarshal(w.Body.Bytes(), &section)

	w = doJSONLog(t, r, "POST", "/api/workout-section-items", map[string]any{
		"workoutSectionId": section.ID, "type": "exercise", "exerciseSchemeId": scheme1.ID, "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section exercise 1: status = %d", w.Code)
	}

	w = doJSONLog(t, r, "POST", "/api/workout-section-items", map[string]any{
		"workoutSectionId": section.ID, "type": "exercise", "exerciseSchemeId": scheme2.ID, "position": 1,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section exercise 2: status = %d", w.Code)
	}

	// Fetch the full workout template to see what the frontend has to work with
	w = doJSONLog(t, r, "GET", "/api/workouts/"+itoa(wkt.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get workout template: status = %d", w.Code)
	}
	json.Unmarshal(w.Body.Bytes(), &wkt)
	t.Logf("=== Workout template loaded: %d section(s), %d item(s) in section 0 ===",
		len(wkt.Sections), len(wkt.Sections[0].Items))

	// -- Step 1: POST /api/user/workout-logs -- create the log --

	t.Log("=== STEP 1: Create workout log (no date during planning) ===")
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs", map[string]any{
		"name":      "Strength Day - March 14",
		"workoutId": wkt.ID,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("STEP 1 FAILED: create workout log: status = %d", w.Code)
	}
	var workoutLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &workoutLog)
	if workoutLog.Status != models.WorkoutLogStatusPlanning {
		t.Errorf("expected planning status, got %s", workoutLog.Status)
	}
	if workoutLog.Date != nil {
		t.Errorf("expected nil date during planning, got %v", workoutLog.Date)
	}

	// -- Step 2: POST /api/user/workout-log-sections -- create each section --

	t.Log("=== STEP 2: Create workout log sections ===")
	templateSection := wkt.Sections[0]
	w = doJSONLog(t, r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId":         workoutLog.ID,
		"type":                 templateSection.Type,
		"position":             templateSection.Position,
		"restBetweenExercises": templateSection.RestBetweenExercises,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("STEP 2 FAILED: create log section: status = %d, body = %s", w.Code, w.Body.String())
	}
	var logSection models.WorkoutLogSection
	json.Unmarshal(w.Body.Bytes(), &logSection)

	// -- Step 3: POST /api/user/workout-log-exercises -- create each exercise --

	t.Log("=== STEP 3: Create workout log exercises ===")
	var logExercises []models.WorkoutLogExercise
	for i, templateItem := range templateSection.Items {
		t.Logf("--- Creating log exercise %d (scheme %d) ---", i, *templateItem.ExerciseSchemeID)
		w = doJSONLog(t, r, "POST", "/api/user/workout-log-exercises", map[string]any{
			"workoutLogSectionId":    logSection.ID,
			"sourceExerciseSchemeId": *templateItem.ExerciseSchemeID,
			"position":               templateItem.Position,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("STEP 3 FAILED: create log exercise %d: status = %d, body = %s", i, w.Code, w.Body.String())
		}
		var logExercise models.WorkoutLogExercise
		json.Unmarshal(w.Body.Bytes(), &logExercise)
		logExercises = append(logExercises, logExercise)
		t.Logf("  -> Created exercise with %d auto-generated sets", len(logExercise.Sets))
	}

	// -- Step 4: PATCH /api/user/workout-log-exercise-sets/:id -- optional target overrides --

	t.Log("=== STEP 4: Override targets on first exercise's first set ===")
	if len(logExercises) > 0 && len(logExercises[0].Sets) > 0 {
		firstSet := logExercises[0].Sets[0]
		w = doJSONLog(t, r, "PATCH", "/api/user/workout-log-exercise-sets/"+itoa(firstSet.ID), map[string]any{
			"targetReps":   3,
			"targetWeight": 160.0,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("STEP 4 FAILED: update set targets: status = %d, body = %s", w.Code, w.Body.String())
		}
	}

	// -- Step 5: Start the workout log --

	t.Log("=== STEP 5: Start workout log ===")
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs/"+itoa(workoutLog.ID)+"/start", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("STEP 5 FAILED: start workout log: status = %d, body = %s", w.Code, w.Body.String())
	}
	var startedLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &startedLog)
	if startedLog.Status != models.WorkoutLogStatusInProgress {
		t.Errorf("expected in_progress status, got %s", startedLog.Status)
	}
	if startedLog.Date == nil {
		t.Error("expected date to be set after start")
	}

	// -- Verify: GET the full nested workout log --

	t.Log("=== VERIFY: Fetch full workout log ===")
	w = doJSONLog(t, r, "GET", "/api/user/workout-logs/"+itoa(workoutLog.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("VERIFY FAILED: get workout log: status = %d", w.Code)
	}
	var fullLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &fullLog)

	// Assertions on the final structure
	if len(fullLog.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(fullLog.Sections))
	}
	if len(fullLog.Sections[0].Exercises) != 2 {
		t.Fatalf("expected 2 exercises, got %d", len(fullLog.Sections[0].Exercises))
	}

	ex1 := fullLog.Sections[0].Exercises[0]
	ex2 := fullLog.Sections[0].Exercises[1]

	if len(ex1.Sets) != 5 {
		t.Errorf("exercise 1: expected 5 sets, got %d", len(ex1.Sets))
	}
	if len(ex2.Sets) != 4 {
		t.Errorf("exercise 2: expected 4 sets, got %d", len(ex2.Sets))
	}

	// All should be in_progress after start
	if fullLog.Status != models.WorkoutLogStatusInProgress {
		t.Errorf("log status: expected in_progress, got %s", fullLog.Status)
	}
	if fullLog.Sections[0].Status != models.WorkoutLogItemStatusInProgress {
		t.Errorf("section status: expected in_progress, got %s", fullLog.Sections[0].Status)
	}

	// Verify the target override from step 4 was applied
	if ex1.Sets[0].TargetReps == nil || *ex1.Sets[0].TargetReps != 3 {
		t.Errorf("exercise 1 set 1: expected overridden targetReps=3, got %v", ex1.Sets[0].TargetReps)
	}
	if ex1.Sets[0].TargetWeight == nil || *ex1.Sets[0].TargetWeight != 160.0 {
		t.Errorf("exercise 1 set 1: expected overridden targetWeight=160, got %v", ex1.Sets[0].TargetWeight)
	}

	// Non-overridden sets should retain original targets
	if ex1.Sets[1].TargetReps == nil || *ex1.Sets[1].TargetReps != 5 {
		t.Errorf("exercise 1 set 2: expected original targetReps=5, got %v", ex1.Sets[1].TargetReps)
	}

	t.Log("=== Workout start flow completed successfully ===")
}
