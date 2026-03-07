package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/user/models"
)

func TestListUserExerciseSchemes(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create a user exercise to link schemes to
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "exerciseTemplateId": "bench-press", "compendiumVersion": 1,
	})

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.UserExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	sets3 := 3
	reps10 := 10
	weight60 := 60.0
	rest90 := 90
	duration30 := 30

	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": sets3, "reps": reps10, "weight": weight60, "restBetweenSets": rest90,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "TIME_BASED",
		"duration": duration30,
	})

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes", nil)
		var result []models.UserExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by userExerciseId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes?userExerciseId=1", nil)
		var result []models.UserExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by measurementType", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes?measurementType=REP_BASED", nil)
		var result []models.UserExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].MeasurementType != "REP_BASED" {
			t.Errorf("measurementType filter: got %d results", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/user/exercise-schemes", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateUserExerciseScheme(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create a user exercise first
	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "exerciseTemplateId": "squat", "compendiumVersion": 1,
	})

	t.Run("success with rep-based fields", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
			"userExerciseId": 1, "measurementType": "REP_BASED",
			"sets": 5, "reps": 5, "weight": 100.0, "restBetweenSets": 180,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.UserExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.MeasurementType != "REP_BASED" || *result.Sets != 5 || *result.Reps != 5 || *result.Weight != 100.0 {
			t.Error("create response mismatch")
		}
	})

	t.Run("success with time-based fields", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
			"userExerciseId": 1, "measurementType": "TIME_BASED",
			"duration": 60, "sets": 3, "restBetweenSets": 30,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.UserExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if *result.Duration != 60 || *result.Sets != 3 {
			t.Error("time-based create response mismatch")
		}
	})

	t.Run("success with distance-based fields", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
			"userExerciseId": 1, "measurementType": "DISTANCE_BASED",
			"distance": 5000.0, "targetTime": 1200,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.UserExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if *result.Distance != 5000.0 || *result.TargetTime != 1200 {
			t.Error("distance-based create response mismatch")
		}
	})

	t.Run("user exercise not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
			"userExerciseId": 999, "measurementType": "REP_BASED",
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/user/exercise-schemes", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
			"userExerciseId": 1, "measurementType": "REP_BASED",
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}

func TestGetUserExerciseScheme(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "exerciseTemplateId": "deadlift", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 3, "reps": 8,
	})

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.UserExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.MeasurementType != "REP_BASED" || *result.Sets != 3 {
			t.Error("get response mismatch")
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestUpdateUserExerciseScheme(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "exerciseTemplateId": "bench-press", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED",
		"sets": 3, "reps": 10, "weight": 60.0,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/exercise-schemes/1", map[string]any{
			"measurementType": "REP_BASED",
			"sets": 5, "reps": 5, "weight": 80.0, "restBetweenSets": 180,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.UserExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if *result.Sets != 5 || *result.Reps != 5 || *result.Weight != 80.0 || *result.RestBetweenSets != 180 {
			t.Errorf("update response mismatch: sets=%v reps=%v weight=%v rest=%v",
				result.Sets, result.Reps, result.Weight, result.RestBetweenSets)
		}
		if result.UserExerciseID != 1 {
			t.Errorf("userExerciseId should be preserved, got %d", result.UserExerciseID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/exercise-schemes/999", map[string]any{
			"measurementType": "REP_BASED",
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "PUT", "/api/user/exercise-schemes/1", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestDeleteUserExerciseScheme(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/user/exercises", map[string]any{
		"owner": "alice", "exerciseTemplateId": "row", "compendiumVersion": 1,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"userExerciseId": 1, "measurementType": "REP_BASED", "sets": 4, "reps": 12,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/exercise-schemes/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/user/exercise-schemes/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
