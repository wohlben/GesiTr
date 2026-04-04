package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/user/exercisescheme/models"
)

func TestListExerciseSchemes(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create an exercise to link schemes to
	doJSON(r, "POST", "/api/exercises", newExercisePayload("Bench Press"))

	t.Run("empty list", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.ExerciseScheme
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
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": sets3, "reps": reps10, "weight": weight60, "restBetweenSets": rest90,
	})
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "TIME_BASED",
		"duration": duration30,
	})

	_ = sets3
	_ = reps10
	_ = weight60
	_ = rest90
	_ = duration30

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes", nil)
		var result []models.ExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by exerciseId", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes?exerciseId=1", nil)
		var result []models.ExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by measurementType", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes?measurementType=REP_BASED", nil)
		var result []models.ExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].MeasurementType != "REP_BASED" {
			t.Errorf("measurementType filter: got %d results", len(result))
		}
	})

	t.Run("scoped to owner", func(t *testing.T) {
		// Bob should see no schemes (they belong to testuser)
		w := doJSONAs(r, "GET", "/api/user/exercise-schemes", nil, "bob")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.ExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0 schemes for bob, got %d", len(result))
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

func TestCreateExerciseScheme(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create an exercise first
	doJSON(r, "POST", "/api/exercises", newExercisePayload("Squat"))

	t.Run("success with rep-based fields", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
			"exerciseId": 1, "measurementType": "REP_BASED",
			"sets": 5, "reps": 5, "weight": 100.0, "restBetweenSets": 180,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.ExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.MeasurementType != "REP_BASED" || *result.Sets != 5 || *result.Reps != 5 || *result.Weight != 100.0 {
			t.Error("create response mismatch")
		}
		if result.Owner != "testuser" {
			t.Errorf("Owner = %q, want testuser", result.Owner)
		}
	})

	t.Run("success with time-based fields", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
			"exerciseId": 1, "measurementType": "TIME_BASED",
			"duration": 60, "sets": 3, "restBetweenSets": 30,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.ExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if *result.Duration != 60 || *result.Sets != 3 {
			t.Error("time-based create response mismatch")
		}
	})

	t.Run("success with distance-based fields", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
			"exerciseId": 1, "measurementType": "DISTANCE_BASED",
			"distance": 5000.0, "targetTime": 1200,
		})
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.ExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if *result.Distance != 5000.0 || *result.TargetTime != 1200 {
			t.Error("distance-based create response mismatch")
		}
	})

	t.Run("exercise not found", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
			"exerciseId": 999, "measurementType": "REP_BASED",
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
			"exerciseId": 1, "measurementType": "REP_BASED",
		})
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}

func TestGetExerciseScheme(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", newExercisePayload("Deadlift"))
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": 3, "reps": 8,
	})

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.ExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.MeasurementType != "REP_BASED" || *result.Sets != 3 {
			t.Error("get response mismatch")
		}
	})

	t.Run("forbidden for non-owner", func(t *testing.T) {
		w := doJSONAs(r, "GET", "/api/user/exercise-schemes/1", nil, "bob")
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/user/exercise-schemes/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestUpdateExerciseScheme(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", newExercisePayload("Bench Press"))
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": 3, "reps": 10, "weight": 60.0,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/exercise-schemes/1", map[string]any{
			"exerciseId":      1,
			"measurementType": "REP_BASED",
			"sets":            5, "reps": 5, "weight": 80.0, "restBetweenSets": 180,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.ExerciseScheme
		json.Unmarshal(w.Body.Bytes(), &result)
		if *result.Sets != 5 || *result.Reps != 5 || *result.Weight != 80.0 || *result.RestBetweenSets != 180 {
			t.Errorf("update response mismatch: sets=%v reps=%v weight=%v rest=%v",
				result.Sets, result.Reps, result.Weight, result.RestBetweenSets)
		}
		if result.ExerciseID != 1 {
			t.Errorf("exerciseId should be preserved, got %d", result.ExerciseID)
		}
	})

	t.Run("forbidden for non-owner", func(t *testing.T) {
		w := doJSONAs(r, "PUT", "/api/user/exercise-schemes/1", map[string]any{
			"exerciseId":      1,
			"measurementType": "REP_BASED",
		}, "bob")
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/user/exercise-schemes/999", map[string]any{
			"exerciseId":      1,
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

func TestDeleteExerciseScheme(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", newExercisePayload("Row"))
	doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
		"exerciseId": 1, "measurementType": "REP_BASED", "sets": 4, "reps": 12,
	})

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/exercise-schemes/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("forbidden for non-owner", func(t *testing.T) {
		// Create a new scheme for this test
		doJSON(r, "POST", "/api/user/exercise-schemes", map[string]any{
			"exerciseId": 1, "measurementType": "REP_BASED", "sets": 4, "reps": 12,
		})
		w := doJSONAs(r, "DELETE", "/api/user/exercise-schemes/2", nil, "bob")
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/user/exercise-schemes/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
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
