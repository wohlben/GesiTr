package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
)

// doJSONLog wraps doJSON and logs the request body and response with pretty-printed JSON.
func doJSONLog(t *testing.T, r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	if body != nil {
		reqJSON, _ := json.MarshalIndent(body, "  ", "  ")
		t.Logf(">>> %s %s\n  Request body:\n  %s", method, path, reqJSON)
	} else {
		t.Logf(">>> %s %s (no body)", method, path)
	}

	w := doJSON(r, method, path, body)

	var pretty json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &pretty); err == nil {
		respJSON, _ := json.MarshalIndent(pretty, "  ", "  ")
		t.Logf("<<< %d\n  Response body:\n  %s", w.Code, respJSON)
	} else {
		t.Logf("<<< %d\n  Response body (raw): %s", w.Code, w.Body.String())
	}

	return w
}

func TestWorkoutStartFlow(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// ── Setup: create user exercises and schemes ──

	w := doJSONLog(t, r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "barbell-squat", "compendiumVersion": 1,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create user exercise 1: status = %d", w.Code)
	}
	var userExercise1 models.UserExercise
	json.Unmarshal(w.Body.Bytes(), &userExercise1)

	w = doJSONLog(t, r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId":  userExercise1.ID,
		"measurementType": "REP_BASED",
		"sets":            5,
		"reps":            5,
		"weight":          140.0,
		"restBetweenSets": 180,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create scheme 1: status = %d, body = %s", w.Code, w.Body.String())
	}
	var scheme1 models.UserExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme1)

	w = doJSONLog(t, r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "compendiumExerciseId": "bench-press", "compendiumVersion": 1,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create user exercise 2: status = %d", w.Code)
	}
	var userExercise2 models.UserExercise
	json.Unmarshal(w.Body.Bytes(), &userExercise2)

	w = doJSONLog(t, r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId":  userExercise2.ID,
		"measurementType": "REP_BASED",
		"sets":            4,
		"reps":            8,
		"weight":          80.0,
		"restBetweenSets": 120,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create scheme 2: status = %d, body = %s", w.Code, w.Body.String())
	}
	var scheme2 models.UserExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme2)

	// ── Setup: create workout template with sections and exercises ──

	w = doJSONLog(t, r, "POST", "/api/user/workouts", map[string]any{
		"owner": "alice", "name": "Strength Day",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create workout: status = %d", w.Code)
	}
	var workout models.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)

	w = doJSONLog(t, r, "POST", "/api/user/workout-sections", map[string]any{
		"workoutId": workout.ID, "type": "main", "position": 0, "restBetweenExercises": 90,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section: status = %d", w.Code)
	}
	var section models.WorkoutSection
	json.Unmarshal(w.Body.Bytes(), &section)

	w = doJSONLog(t, r, "POST", "/api/user/workout-section-exercises", map[string]any{
		"workoutSectionId": section.ID, "userExerciseSchemeId": scheme1.ID, "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section exercise 1: status = %d", w.Code)
	}

	w = doJSONLog(t, r, "POST", "/api/user/workout-section-exercises", map[string]any{
		"workoutSectionId": section.ID, "userExerciseSchemeId": scheme2.ID, "position": 1,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section exercise 2: status = %d", w.Code)
	}

	// Fetch the full workout template to see what the frontend has to work with
	w = doJSONLog(t, r, "GET", "/api/user/workouts/"+itoa(workout.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get workout template: status = %d", w.Code)
	}
	json.Unmarshal(w.Body.Bytes(), &workout)
	t.Logf("=== Workout template loaded: %d section(s), %d exercise(s) in section 0 ===",
		len(workout.Sections), len(workout.Sections[0].Exercises))

	// ── Step 1: POST /api/user/workout-logs — create the log ──

	t.Log("=== STEP 1: Create workout log (no date during planning) ===")
	w = doJSONLog(t, r, "POST", "/api/user/workout-logs", map[string]any{
		"owner":     "alice",
		"name":      "Strength Day - March 14",
		"workoutId": workout.ID,
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

	// ── Step 2: POST /api/user/workout-log-sections — create each section ──

	t.Log("=== STEP 2: Create workout log sections ===")
	templateSection := workout.Sections[0]
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

	// ── Step 3: POST /api/user/workout-log-exercises — create each exercise ──

	t.Log("=== STEP 3: Create workout log exercises ===")
	var logExercises []models.WorkoutLogExercise
	for i, templateExercise := range templateSection.Exercises {
		t.Logf("--- Creating log exercise %d (scheme %d) ---", i, templateExercise.UserExerciseSchemeID)
		w = doJSONLog(t, r, "POST", "/api/user/workout-log-exercises", map[string]any{
			"workoutLogSectionId":    logSection.ID,
			"sourceExerciseSchemeId": templateExercise.UserExerciseSchemeID,
			"position":               templateExercise.Position,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("STEP 3 FAILED: create log exercise %d: status = %d, body = %s", i, w.Code, w.Body.String())
		}
		var logExercise models.WorkoutLogExercise
		json.Unmarshal(w.Body.Bytes(), &logExercise)
		logExercises = append(logExercises, logExercise)
		t.Logf("  → Created exercise with %d auto-generated sets", len(logExercise.Sets))
	}

	// ── Step 4: PUT /api/user/workout-log-exercise-sets/:id — optional target overrides ──

	t.Log("=== STEP 4: Override targets on first exercise's first set ===")
	if len(logExercises) > 0 && len(logExercises[0].Sets) > 0 {
		firstSet := logExercises[0].Sets[0]
		w = doJSONLog(t, r, "PUT", "/api/user/workout-log-exercise-sets/"+itoa(firstSet.ID), map[string]any{
			"targetReps":   3,
			"targetWeight": 160.0,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("STEP 4 FAILED: update set targets: status = %d, body = %s", w.Code, w.Body.String())
		}
	}

	// ── Step 5: Start the workout log ──

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

	// ── Verify: GET the full nested workout log ──

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
	if fullLog.Sections[0].Status != models.WorkoutLogStatusInProgress {
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
