package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	exercisemodels "gesitr/internal/exercise/models"
	workoutmodels "gesitr/internal/user/workout/models"
	workoutgroupmodels "gesitr/internal/user/workoutgroup/models"
	"gesitr/internal/user/workoutlog/models"
)

// TestGroupMemberStartFlow verifies that an invited group member can accept
// their invitation after creating exercise schemes, then start the workout
// and log sets using their own schemes.
func TestGroupMemberStartFlow(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// -- Setup: alice creates a public exercise --

	w := doJSONLog(t, r, "POST", "/api/exercises", map[string]any{
		"name": "Barbell Squat", "type": "STRENGTH",
		"technicalDifficulty": "intermediate", "public": true,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create exercise: status = %d", w.Code)
	}
	var exercise exercisemodels.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)

	// Alice creates her own scheme
	w = doJSONLog(t, r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": exercise.ID, "measurementType": "REP_BASED",
		"sets": 5, "reps": 5, "weight": 100.0, "restBetweenSets": 180,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create alice scheme: status = %d, body = %s", w.Code, w.Body.String())
	}
	var aliceScheme exercisemodels.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &aliceScheme)

	// -- Setup: alice creates workout with 1 exercise item --

	w = doJSONLog(t, r, "POST", "/api/user/workouts", map[string]any{"name": "Squat Day"})
	if w.Code != http.StatusCreated {
		t.Fatalf("create workout: status = %d", w.Code)
	}
	var wkt workoutmodels.Workout
	json.Unmarshal(w.Body.Bytes(), &wkt)

	w = doJSONLog(t, r, "POST", "/api/user/workout-sections", map[string]any{
		"workoutId": wkt.ID, "type": "main", "position": 0, "restBetweenExercises": 90,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section: status = %d", w.Code)
	}
	var section workoutmodels.WorkoutSection
	json.Unmarshal(w.Body.Bytes(), &section)

	w = doJSONLog(t, r, "POST", "/api/user/workout-section-items", map[string]any{
		"workoutSectionId": section.ID, "type": "exercise",
		"exerciseSchemeId": aliceScheme.ID, "position": 0,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create section item: status = %d", w.Code)
	}
	var sectionItem workoutmodels.WorkoutSectionItem
	json.Unmarshal(w.Body.Bytes(), &sectionItem)

	// -- Setup: alice creates workout group and invites bob --

	w = doJSONLog(t, r, "POST", "/api/user/workout-groups", map[string]any{
		"name": "Squat Crew", "workoutId": wkt.ID,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create group: status = %d", w.Code)
	}

	w = doJSONLog(t, r, "POST", "/api/user/workout-group-memberships", map[string]any{
		"groupId": 1, "userId": "bob", "role": "member",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("invite bob: status = %d, body = %s", w.Code, w.Body.String())
	}
	var membership workoutgroupmodels.WorkoutGroupMembership
	json.Unmarshal(w.Body.Bytes(), &membership)
	if membership.Role != workoutgroupmodels.WorkoutGroupRoleInvited {
		t.Errorf("expected role 'invited', got %q", membership.Role)
	}

	// -- Bob creates his own scheme for the workout section item --

	t.Log("=== Bob creates his own exercise scheme ===")
	w = doJSONLogAs(t, r, "POST", "/api/exercise-schemes", map[string]any{
		"exerciseId": exercise.ID, "measurementType": "REP_BASED",
		"sets": 3, "reps": 8, "weight": 70.0, "restBetweenSets": 120,
		"workoutSectionItemId": sectionItem.ID,
	}, "bob")
	if w.Code != http.StatusCreated {
		t.Fatalf("create bob scheme: status = %d, body = %s", w.Code, w.Body.String())
	}
	var bobScheme exercisemodels.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &bobScheme)

	// -- Bob accepts the invitation --

	t.Log("=== Bob accepts workout group invitation ===")
	w = doJSONLogAs(t, r, "POST", "/api/user/workouts/"+itoa(wkt.ID)+"/group/accept", nil, "bob")
	if w.Code != http.StatusOK {
		t.Fatalf("accept invitation: status = %d, body = %s", w.Code, w.Body.String())
	}

	// -- Bob starts the workout using his own scheme --

	t.Log("=== STEP 1: Bob creates workout log ===")
	w = doJSONLogAs(t, r, "POST", "/api/user/workout-logs", map[string]any{
		"name": "Squat Day - Bob", "workoutId": wkt.ID,
	}, "bob")
	if w.Code != http.StatusCreated {
		t.Fatalf("create log: status = %d, body = %s", w.Code, w.Body.String())
	}
	var workoutLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &workoutLog)

	t.Log("=== STEP 2: Bob creates workout log section ===")
	w = doJSONLogAs(t, r, "POST", "/api/user/workout-log-sections", map[string]any{
		"workoutLogId": workoutLog.ID, "type": "main", "position": 0, "restBetweenExercises": 90,
	}, "bob")
	if w.Code != http.StatusCreated {
		t.Fatalf("create log section: status = %d, body = %s", w.Code, w.Body.String())
	}
	var logSection models.WorkoutLogSection
	json.Unmarshal(w.Body.Bytes(), &logSection)

	t.Log("=== STEP 3: Bob creates workout log exercise using his own scheme ===")
	w = doJSONLogAs(t, r, "POST", "/api/user/workout-log-exercises", map[string]any{
		"workoutLogSectionId":    logSection.ID,
		"sourceExerciseSchemeId": bobScheme.ID,
		"position":               0,
	}, "bob")
	if w.Code != http.StatusCreated {
		t.Fatalf("create log exercise: status = %d, body = %s", w.Code, w.Body.String())
	}
	var logExercise models.WorkoutLogExercise
	json.Unmarshal(w.Body.Bytes(), &logExercise)

	// Bob's scheme is 3x8@70, so we expect 3 sets
	if len(logExercise.Sets) != 3 {
		t.Fatalf("expected 3 sets (from bob's scheme), got %d", len(logExercise.Sets))
	}
	if logExercise.Sets[0].TargetReps == nil || *logExercise.Sets[0].TargetReps != 8 {
		t.Errorf("expected targetReps=8, got %v", logExercise.Sets[0].TargetReps)
	}
	if logExercise.Sets[0].TargetWeight == nil || *logExercise.Sets[0].TargetWeight != 70.0 {
		t.Errorf("expected targetWeight=70, got %v", logExercise.Sets[0].TargetWeight)
	}

	t.Log("=== STEP 4: Bob starts the workout ===")
	w = doJSONLogAs(t, r, "POST", "/api/user/workout-logs/"+itoa(workoutLog.ID)+"/start", nil, "bob")
	if w.Code != http.StatusOK {
		t.Fatalf("start workout: status = %d, body = %s", w.Code, w.Body.String())
	}

	t.Log("=== STEP 5: Bob logs actual reps for each set ===")
	for _, set := range logExercise.Sets {
		w = doJSONLogAs(t, r, "PATCH", "/api/user/workout-log-exercise-sets/"+itoa(set.ID), map[string]any{
			"status":       "finished",
			"actualReps":   8,
			"actualWeight": 70.0,
		}, "bob")
		if w.Code != http.StatusOK {
			t.Fatalf("finish set %d: status = %d, body = %s", set.SetNumber, w.Code, w.Body.String())
		}
	}

	// -- Verify the final state --

	t.Log("=== VERIFY: Fetch full workout log ===")
	w = doJSONLogAs(t, r, "GET", "/api/user/workout-logs/"+itoa(workoutLog.ID), nil, "bob")
	if w.Code != http.StatusOK {
		t.Fatalf("get log: status = %d", w.Code)
	}
	var fullLog models.WorkoutLog
	json.Unmarshal(w.Body.Bytes(), &fullLog)

	// Log should be finished (all sets finished → exercise finished → section finished → log finished)
	if fullLog.Status != models.WorkoutLogStatusFinished {
		t.Errorf("expected log status 'finished', got %q", fullLog.Status)
	}
	if len(fullLog.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(fullLog.Sections))
	}
	if fullLog.Sections[0].Status != models.WorkoutLogItemStatusFinished {
		t.Errorf("expected section status 'finished', got %q", fullLog.Sections[0].Status)
	}
	if len(fullLog.Sections[0].Exercises) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(fullLog.Sections[0].Exercises))
	}

	ex := fullLog.Sections[0].Exercises[0]
	if ex.Status != models.WorkoutLogItemStatusFinished {
		t.Errorf("expected exercise status 'finished', got %q", ex.Status)
	}
	if len(ex.Sets) != 3 {
		t.Errorf("expected 3 sets, got %d", len(ex.Sets))
	}
	for i, s := range ex.Sets {
		if s.Status != models.WorkoutLogItemStatusFinished {
			t.Errorf("set %d: expected 'finished', got %q", i+1, s.Status)
		}
	}

	t.Log("=== Group member workout start flow completed successfully ===")
}
