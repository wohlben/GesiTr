package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"gesitr/internal/compendium/exercise/models"
	"gesitr/internal/database"
	"gesitr/internal/shared"

	"gorm.io/gorm"
)

func newExercisePayload(name, templateID string) map[string]any {
	return map[string]any{
		"name": name, "templateId": templateID, "type": "STRENGTH",
		"technicalDifficulty": "beginner", "bodyWeightScaling": 0.5,
		"description": "test", "createdBy": "system", "version": 0,
		"force": []string{"PUSH"}, "primaryMuscles": []string{"CHEST"},
		"secondaryMuscles":              []string{"TRICEPS"},
		"suggestedMeasurementParadigms": []string{"REP_BASED"},
		"instructions":                  []string{"Step 1", "Step 2"},
		"images":                        []string{"/img/a.jpg"},
		"alternativeNames":              []string{"Alt Name"},
		"equipmentIds":                  []string{"barbell"},
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
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
		if page.Total != 0 {
			t.Errorf("expected total 0, got %d", page.Total)
		}
	})

	// Seed exercises for filter tests
	doJSON(r, "POST", "/api/exercises", newExercisePayload("Bench Press", "bench-press"))

	cardio := map[string]any{
		"name": "Running", "templateId": "running", "type": "CARDIO",
		"technicalDifficulty": "intermediate", "bodyWeightScaling": 1.0,
		"description": "run", "createdBy": "system", "version": 0,
		"force": []string{"DYNAMIC"}, "primaryMuscles": []string{"QUADS"},
		"secondaryMuscles":              []string{"CALVES"},
		"suggestedMeasurementParadigms": []string{"DISTANCE"},
	}
	doJSON(r, "POST", "/api/exercises", cardio)

	t.Run("list all", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
		if page.Total != 2 {
			t.Errorf("expected total 2, got %d", page.Total)
		}
	})

	t.Run("filter by q (name)", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?q=Bench", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("filter by q (alt name)", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?q=Alt+Name", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("filter by type", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?type=CARDIO", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 || result[0].Name != "Running" {
			t.Errorf("type filter: got %d results", len(result))
		}
	})

	t.Run("filter by difficulty", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?difficulty=beginner", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 || result[0].Name != "Bench Press" {
			t.Errorf("difficulty filter: got %d results", len(result))
		}
	})

	t.Run("filter by force", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?force=PUSH", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 || result[0].Name != "Bench Press" {
			t.Errorf("force filter: got %d results", len(result))
		}
	})

	t.Run("filter by muscle", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?muscle=CALVES", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 || result[0].Name != "Running" {
			t.Errorf("muscle filter: got %d results", len(result))
		}
	})

	t.Run("filter by primaryMuscle", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?primaryMuscle=CHEST", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 || result[0].Name != "Bench Press" {
			t.Errorf("primaryMuscle filter: got %d results", len(result))
		}
	})

	t.Run("pagination limit", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?limit=1", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 {
			t.Errorf("expected 1 item, got %d", len(result))
		}
		if page.Total != 2 {
			t.Errorf("expected total 2, got %d", page.Total)
		}
		if page.Limit != 1 {
			t.Errorf("expected limit 1, got %d", page.Limit)
		}
	})

	t.Run("pagination offset", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises?limit=1&offset=1", nil)
		var page paginatedJSON
		json.Unmarshal(w.Body.Bytes(), &page)
		var result []models.Exercise
		json.Unmarshal(page.Items, &result)
		if len(result) != 1 {
			t.Errorf("expected 1 item, got %d", len(result))
		}
		if page.Offset != 1 {
			t.Errorf("expected offset 1, got %d", page.Offset)
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

	t.Run("auto-generates slug from name", func(t *testing.T) {
		w := doJSON(r, "POST", "/api/exercises", newExercisePayload("Bench Press", "bench-press"))
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Slug != "bench-press" {
			t.Errorf("Slug = %q, want %q", result.Slug, "bench-press")
		}
	})

	t.Run("auto-generates slug when not provided", func(t *testing.T) {
		payload := map[string]any{
			"name": "Overhead Press", "templateId": "overhead-press", "type": "STRENGTH",
			"technicalDifficulty": "beginner", "bodyWeightScaling": 0.0,
			"description": "test", "createdBy": "system",
		}
		w := doJSON(r, "POST", "/api/exercises", payload)
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Slug != "overhead-press" {
			t.Errorf("Slug = %q, want %q", result.Slug, "overhead-press")
		}
	})

	t.Run("defaults templateId to UUID when not provided", func(t *testing.T) {
		payload := map[string]any{
			"name": "No Template", "type": "STRENGTH",
			"technicalDifficulty": "beginner", "bodyWeightScaling": 0.0,
			"description": "test", "createdBy": "system",
		}
		w := doJSON(r, "POST", "/api/exercises", payload)
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.TemplateID == "" {
			t.Error("expected auto-generated templateId, got empty string")
		}
		if len(result.TemplateID) != 36 {
			t.Errorf("expected UUID-length templateId, got %q (len %d)", result.TemplateID, len(result.TemplateID))
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
			"name": "Overhead Press", "templateId": "ohp", "type": "STRENGTH",
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

	t.Run("no version bump when unchanged", func(t *testing.T) {
		// Re-PUT the same data that's currently in the DB (from the "success" test above)
		same := map[string]any{
			"name": "Overhead Press", "templateId": "ohp", "type": "STRENGTH",
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
		w := doJSON(r, "PUT", "/api/exercises/1", same)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Version != 1 {
			t.Errorf("Version = %d, want 1 (should not have bumped)", result.Version)
		}
	})

	t.Run("no extra history on unchanged update", func(t *testing.T) {
		var count int64
		database.DB.Model(&models.ExerciseHistoryEntity{}).Where("exercise_id = ?", 1).Count(&count)
		// version 0 (create) + version 1 (first update) = 2 records; no-op update above should not add more
		if count != 2 {
			t.Errorf("expected 2 history records, got %d", count)
		}
	})

	t.Run("successive updates accumulate history", func(t *testing.T) {
		// Second real update -> version 2
		w := doJSON(r, "PUT", "/api/exercises/1", map[string]any{
			"name": "OHP v2", "templateId": "ohp", "type": "STRENGTH",
			"technicalDifficulty": "intermediate", "bodyWeightScaling": 0.0,
			"description": "v2", "createdBy": "system", "version": 0,
			"force": []string{"PUSH"}, "primaryMuscles": []string{"SHOULDERS"},
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var result models.Exercise
		json.Unmarshal(w.Body.Bytes(), &result)
		if result.Version != 2 {
			t.Errorf("Version = %d, want 2", result.Version)
		}

		var count int64
		database.DB.Model(&models.ExerciseHistoryEntity{}).Where("exercise_id = ?", 1).Count(&count)
		if count != 3 {
			t.Errorf("expected 3 history records (v0, v1, v2), got %d", count)
		}
	})

	t.Run("history insert error rolls back", func(t *testing.T) {
		callbackName := "test:fail_history_insert"
		database.DB.Callback().Create().Before("gorm:create").Register(callbackName, func(tx *gorm.DB) {
			if tx.Statement.Table == "exercise_history" {
				_ = tx.AddError(fmt.Errorf("injected history error"))
			}
		})
		w := doJSON(r, "PUT", "/api/exercises/1", map[string]any{
			"name": "OHP v3 fail", "templateId": "ohp", "type": "STRENGTH",
			"technicalDifficulty": "advanced", "bodyWeightScaling": 0.0,
			"description": "v3 fail", "createdBy": "system",
			"force": []string{"PUSH"}, "primaryMuscles": []string{"SHOULDERS"},
		})
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 for history insert error, got %d", w.Code)
		}
		// Version should NOT have bumped since the whole transaction rolled back
		var entity models.ExerciseEntity
		database.DB.First(&entity, 1)
		if entity.Version != 2 {
			t.Errorf("Version = %d after rolled-back tx, want 2", entity.Version)
		}
		database.DB.Callback().Create().Remove(callbackName)
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

	t.Run("db error on update (templateId conflict)", func(t *testing.T) {
		// Create a second exercise
		doJSON(r, "POST", "/api/exercises", newExercisePayload("Other", "other"))
		// Try to update it with the templateId of the first exercise
		conflict := newExercisePayload("Conflict", "ohp")
		w := doJSON(r, "PUT", "/api/exercises/2", conflict)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 for templateId conflict, got %d", w.Code)
		}
	})

	// Test transaction delete errors by injecting errors via callbacks
	deleteTargets := []struct {
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
	for _, target := range deleteTargets {
		t.Run("tx delete error "+target.name, func(t *testing.T) {
			callbackName := "test:fail_del_" + target.name
			targetTable := target.table
			database.DB.Callback().Delete().Before("gorm:delete").Register(callbackName, func(tx *gorm.DB) {
				if tx.Statement.Table == targetTable {
					_ = tx.AddError(fmt.Errorf("injected delete error"))
				}
			})
			w := doJSON(r, "PUT", "/api/exercises/1", newExercisePayload("U", "ohp"))
			if w.Code != http.StatusInternalServerError {
				t.Errorf("expected 500, got %d", w.Code)
			}
			database.DB.Callback().Delete().Remove(callbackName)
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

	// Test reload error after successful transaction by injecting query errors
	// after the initial preloaded load (1 First + 7 preloads = 8 queries).
	// Reload queries start at query 9+.
	t.Run("reload error after update", func(t *testing.T) {
		callbackName := "test:fail_reload_update"
		queryCount := 0
		database.DB.Callback().Query().Before("gorm:query").Register(callbackName, func(tx *gorm.DB) {
			queryCount++
			if queryCount >= 9 {
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

func TestListExerciseVersions(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create exercise (v0) and update it twice (v1, v2)
	doJSON(r, "POST", "/api/exercises", newExercisePayload("Press", "press"))
	doJSON(r, "PUT", "/api/exercises/1", map[string]any{
		"name": "Press v1", "templateId": "press", "type": "STRENGTH",
		"technicalDifficulty": "intermediate", "bodyWeightScaling": 0.0,
		"description": "v1", "createdBy": "system",
		"force": []string{"PUSH"}, "primaryMuscles": []string{"CHEST"},
	})
	doJSON(r, "PUT", "/api/exercises/1", map[string]any{
		"name": "Press v2", "templateId": "press", "type": "STRENGTH",
		"technicalDifficulty": "advanced", "bodyWeightScaling": 0.0,
		"description": "v2", "createdBy": "system",
		"force": []string{"PUSH"}, "primaryMuscles": []string{"CHEST"},
	})

	t.Run("returns all versions ordered", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises/1/versions", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var entries []shared.VersionEntry
		json.Unmarshal(w.Body.Bytes(), &entries)
		if len(entries) != 3 {
			t.Fatalf("expected 3 versions, got %d", len(entries))
		}
		if entries[0].Version != 0 || entries[1].Version != 1 || entries[2].Version != 2 {
			t.Errorf("versions = %d, %d, %d", entries[0].Version, entries[1].Version, entries[2].Version)
		}
		for i, e := range entries {
			if e.ChangedBy != "system" {
				t.Errorf("entries[%d].ChangedBy = %q", i, e.ChangedBy)
			}
			if e.ChangedAt.IsZero() {
				t.Errorf("entries[%d].ChangedAt is zero", i)
			}
		}
	})

	t.Run("snapshot contains correct data", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises/1/versions", nil)
		var entries []shared.VersionEntry
		json.Unmarshal(w.Body.Bytes(), &entries)

		// Check v0 snapshot
		var v0 models.Exercise
		json.Unmarshal(entries[0].Snapshot, &v0)
		if v0.Name != "Press" {
			t.Errorf("v0 name = %q, want Press", v0.Name)
		}

		// Check v2 snapshot
		var v2 models.Exercise
		json.Unmarshal(entries[2].Snapshot, &v2)
		if v2.Name != "Press v2" {
			t.Errorf("v2 name = %q, want Press v2", v2.Name)
		}
		if v2.TechnicalDifficulty != "advanced" {
			t.Errorf("v2 difficulty = %q, want advanced", v2.TechnicalDifficulty)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises/999/versions", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		closeDB(t)
		w := doJSON(r, "GET", "/api/exercises/1/versions", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404 (db closed), got %d", w.Code)
		}
	})
}

func TestGetExerciseVersion(t *testing.T) {
	setupTestDB(t)
	r := newRouter()

	// Create exercise (v0) and update it (v1)
	doJSON(r, "POST", "/api/exercises", newExercisePayload("Press", "press"))
	doJSON(r, "PUT", "/api/exercises/1", map[string]any{
		"name": "Press v1", "templateId": "press", "type": "STRENGTH",
		"technicalDifficulty": "intermediate", "bodyWeightScaling": 0.0,
		"description": "v1", "createdBy": "system",
		"force": []string{"PUSH"}, "primaryMuscles": []string{"CHEST"},
	})

	t.Run("returns specific version", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises/templates/press/versions/0", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var entry shared.VersionEntry
		json.Unmarshal(w.Body.Bytes(), &entry)
		if entry.Version != 0 {
			t.Errorf("version = %d, want 0", entry.Version)
		}
		var snapshot models.Exercise
		json.Unmarshal(entry.Snapshot, &snapshot)
		if snapshot.Name != "Press" {
			t.Errorf("snapshot name = %q, want Press", snapshot.Name)
		}
	})

	t.Run("returns v1", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises/templates/press/versions/1", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var entry shared.VersionEntry
		json.Unmarshal(w.Body.Bytes(), &entry)
		if entry.Version != 1 {
			t.Errorf("version = %d, want 1", entry.Version)
		}
		var snapshot models.Exercise
		json.Unmarshal(entry.Snapshot, &snapshot)
		if snapshot.Name != "Press v1" {
			t.Errorf("snapshot name = %q, want Press v1", snapshot.Name)
		}
	})

	t.Run("works for soft-deleted exercises", func(t *testing.T) {
		doJSON(r, "DELETE", "/api/exercises/1", nil)
		w := doJSON(r, "GET", "/api/exercises/templates/press/versions/0", nil)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200 for soft-deleted exercise version, got %d", w.Code)
		}
	})

	t.Run("template not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises/templates/nonexistent/versions/0", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("version not found", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises/templates/press/versions/999", nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		w := doJSON(r, "GET", "/api/exercises/templates/press/versions/abc", nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
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
