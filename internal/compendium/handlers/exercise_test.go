package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"gesitr/internal/compendium/models"
	"gesitr/internal/database"

	"gorm.io/gorm"
)

func newExercisePayload(name, slug string) map[string]any {
	return map[string]any{
		"name": name, "slug": slug, "type": "STRENGTH",
		"technicalDifficulty": "beginner", "bodyWeightScaling": 0.5,
		"description": "test", "createdBy": "system", "version": 0,
		"force": []string{"PUSH"}, "primaryMuscles": []string{"CHEST"},
		"secondaryMuscles": []string{"TRICEPS"},
		"suggestedMeasurementParadigms": []string{"REP_BASED"},
		"instructions": []string{"Step 1", "Step 2"},
		"images": []string{"/img/a.jpg"},
		"alternativeNames": []string{"Alt Name"},
		"equipmentIds": []string{"barbell"},
	}
}

func TestListExercises(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("empty", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result []models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	// Seed exercises for filter tests
	doJSON(r, "POST", "/api/exercises", newExercisePayload("Bench Press", "bench-press"))

	cardio := map[string]any{
		"name": "Running", "slug": "running", "type": "CARDIO",
		"technicalDifficulty": "intermediate", "bodyWeightScaling": 1.0,
		"description": "run", "createdBy": "system", "version": 0,
		"force": []string{"DYNAMIC"}, "primaryMuscles": []string{"QUADS"},
		"secondaryMuscles": []string{"CALVES"},
		"suggestedMeasurementParadigms": []string{"DISTANCE"},
	}
	doJSON(r, "POST", "/api/exercises", cardio)

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises", nil)
		var result []models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by q (name)", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?q=Bench", nil)
		var result []models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("filter by q (alt name)", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?q=Alt+Name", nil)
		var result []models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("filter by type", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?type=CARDIO", nil)
		var result []models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Name != "Running" {
			t.Errorf("type filter: got %d results", len(result))
		}
	})

	t.Run("filter by difficulty", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?difficulty=beginner", nil)
		var result []models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Name != "Bench Press" {
			t.Errorf("difficulty filter: got %d results", len(result))
		}
	})

	t.Run("filter by force", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?force=PUSH", nil)
		var result []models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Name != "Bench Press" {
			t.Errorf("force filter: got %d results", len(result))
		}
	})

	t.Run("filter by muscle", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?muscle=CALVES", nil)
		var result []models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Name != "Running" {
			t.Errorf("muscle filter: got %d results", len(result))
		}
	})

	t.Run("filter by primaryMuscle", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?primaryMuscle=CHEST", nil)
		var result []models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if len(result) != 1 || result[0].Name != "Bench Press" {
			t.Errorf("primaryMuscle filter: got %d results", len(result))
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/exercises", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	t.Run("success with associations", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/exercises", newExercisePayload("Squat", "squat"))
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.ID == 0 || result.Name != "Squat" {
			t.Error("create response mismatch")
		}
		if len(result.Force) != 1 || result.Force[0] != "PUSH" {
			t.Errorf("Force = %v", result.Force)
		}
		if len(result.PrimaryMuscles) != 1 {
			t.Errorf("PrimaryMuscles = %v", result.PrimaryMuscles)
		}
		if len(result.SecondaryMuscles) != 1 {
			t.Errorf("SecondaryMuscles = %v", result.SecondaryMuscles)
		}
		if len(result.Instructions) != 2 {
			t.Errorf("Instructions len = %d", len(result.Instructions))
		}
		if len(result.Images) != 1 {
			t.Errorf("Images len = %d", len(result.Images))
		}
		if len(result.AlternativeNames) != 1 {
			t.Errorf("AlternativeNames len = %d", len(result.AlternativeNames))
		}
		if len(result.EquipmentIDs) != 1 {
			t.Errorf("EquipmentIDs len = %d", len(result.EquipmentIDs))
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "POST", "/api/exercises", "{invalid")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error on create", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "POST", "/api/exercises", newExercisePayload("X", "x"))
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})

	t.Run("reload error after create", func(t *testing.T) {
		setupTestDB(t)
		r = newRouter()
		callbackName := "test:fail_create_reload"
		database.DB.Callback().Query().Before("gorm:query").Register(callbackName, func(tx *gorm.DB) {
			// All queries fail - create doesn't use Query callbacks, but reload does
			_ = tx.AddError(fmt.Errorf("injected reload error"))
		})
		w := doJSON(r, "POST", "/api/exercises", newExercisePayload("ReloadFail", "reload-fail"))
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 for reload error, got %d", w.Code)
		}
		database.DB.Callback().Query().Remove(callbackName)
	})
}

func TestGetExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", newExercisePayload("Deadlift", "deadlift"))

	t.Run("found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var result models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Name != "Deadlift" {
			t.Error("get response mismatch")
		}
		// Verify preloads work
		if len(result.Force) == 0 {
			t.Error("expected preloaded forces")
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestUpdateExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", newExercisePayload("OHP", "ohp"))

	t.Run("success", func(t *testing.T) {
		updated := map[string]any{
			"name": "Overhead Press", "slug": "ohp", "type": "STRENGTH",
			"technicalDifficulty": "intermediate", "bodyWeightScaling": 0.0,
			"description": "updated", "createdBy": "system", "version": 0,
			"force": []string{"PUSH"}, "primaryMuscles": []string{"SHOULDERS"},
			"secondaryMuscles":              []string{"TRICEPS", "CHEST"},
			"suggestedMeasurementParadigms": []string{"REP_BASED", "AMRAP"},
			"instructions":                  []string{"New step 1"},
			"images":                        []string{"/img/new.jpg"},
			"alternativeNames":              []string{"Press", "Military Press"},
			"equipmentIds":                  []string{"barbell"},
		}
		w := doJSON(r, "PUT", "/api/exercises/1", updated)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Name != "Overhead Press" {
			t.Errorf("Name = %q", result.Name)
		}
		if result.Version != 1 {
			t.Errorf("Version = %d, want 1", result.Version)
		}
		if len(result.PrimaryMuscles) != 1 || result.PrimaryMuscles[0] != "SHOULDERS" {
			t.Errorf("PrimaryMuscles = %v", result.PrimaryMuscles)
		}
		if len(result.SecondaryMuscles) != 2 {
			t.Errorf("SecondaryMuscles = %v", result.SecondaryMuscles)
		}
		if len(result.Instructions) != 1 {
			t.Errorf("Instructions len = %d", len(result.Instructions))
		}
		if len(result.AlternativeNames) != 2 {
			t.Errorf("AlternativeNames len = %d", len(result.AlternativeNames))
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "PUT", "/api/exercises/999", newExercisePayload("X", "x"))
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		w := doRaw(r, "PUT", "/api/exercises/1", "{bad")
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("db error on update (slug conflict)", func(t *testing.T) {
		// Create a second exercise
		doJSON(r, "POST", "/api/exercises", newExercisePayload("Other", "other"))
		// Try to update it with the slug of the first exercise
		conflict := newExercisePayload("Conflict", "ohp")
		w := doJSON(r, "PUT", "/api/exercises/2", conflict)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 for slug conflict, got %d", w.Code)
		}
	})

	// Test transaction delete errors by dropping child tables one at a time
	deleteTargets := []struct {
		model any
		name  string
	}{
		{&models.ExerciseForce{}, "forces"},
		{&models.ExerciseMuscle{}, "muscles"},
		{&models.ExerciseMeasurementParadigm{}, "paradigms"},
		{&models.ExerciseInstruction{}, "instructions"},
		{&models.ExerciseImage{}, "images"},
		{&models.ExerciseAlternativeName{}, "alt_names"},
		{&models.ExerciseEquipment{}, "equipment"},
	}
	for _, target := range deleteTargets {
		t.Run("tx delete error "+target.name, func(t *testing.T) {
			database.DB.Migrator().DropTable(target.model)
			w := doJSON(r, "PUT", "/api/exercises/1", newExercisePayload("U", "ohp"))
			if w.Code != http.StatusInternalServerError {
				t.Errorf("expected 500, got %d", w.Code)
			}
			database.DB.AutoMigrate(target.model)
		})
	}

	// Test transaction insert errors by using GORM callbacks to inject errors
	insertTargets := []struct {
		table string
		name  string
	}{
		{"exercise_forces", "forces"},
		{"exercise_muscles", "muscles"},
		{"exercise_measurement_paradigms", "paradigms"},
		{"exercise_instructions", "instructions"},
		{"exercise_images", "images"},
		{"exercise_alternative_names", "alt_names"},
		{"exercise_equipments", "equipment"},
	}
	for _, target := range insertTargets {
		t.Run("tx insert error "+target.name, func(t *testing.T) {
			callbackName := "test:fail_" + target.name
			targetTable := target.table
			database.DB.Callback().Create().Before("gorm:create").Register(callbackName, func(tx *gorm.DB) {
				if tx.Statement.Table == targetTable {
					_ = tx.AddError(fmt.Errorf("injected error"))
				}
			})
			w := doJSON(r, "PUT", "/api/exercises/1", newExercisePayload("U2", "ohp"))
			if w.Code != http.StatusInternalServerError {
				t.Errorf("expected 500, got %d", w.Code)
			}
			database.DB.Callback().Create().Remove(callbackName)
		})
	}

	// Test reload error after successful transaction - drop exercises table
	// after the transaction succeeds by using a callback on query (reload phase)
	t.Run("reload error after update", func(t *testing.T) {
		callbackName := "test:fail_reload_update"
		queryCount := 0
		database.DB.Callback().Query().Before("gorm:query").Register(callbackName, func(tx *gorm.DB) {
			queryCount++
			// First query: lookup existing exercise. Second+ queries: reload after tx.
			// Fail on query 2+ to break the reload.
			if queryCount >= 2 {
				_ = tx.AddError(fmt.Errorf("injected reload error"))
			}
		})
		w := doJSON(r, "PUT", "/api/exercises/1", newExercisePayload("Reload", "ohp"))
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 for reload error, got %d", w.Code)
		}
		database.DB.Callback().Query().Remove(callbackName)
	})

	t.Run("db error first lookup", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "PUT", "/api/exercises/1", newExercisePayload("X", "x"))
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}

func TestDeleteExercise(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	doJSON(r, "POST", "/api/exercises", newExercisePayload("Curl", "curl"))

	t.Run("success", func(t *testing.T) {
		w := doJSON(r, "DELETE", "/api/exercises/1", nil)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "DELETE", "/api/exercises/1", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}
