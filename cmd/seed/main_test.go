package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gesitr/internal/database"
	equipmentModels "gesitr/internal/equipment/models"
	fulfillmentModels "gesitr/internal/equipmentfulfillment/models"
	exerciseModels "gesitr/internal/exercise/models"
	relModels "gesitr/internal/exerciserelationship/models"
	profileModels "gesitr/internal/profile/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSeedTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&profileModels.UserProfileEntity{},
		&exerciseModels.ExerciseEntity{},
		&exerciseModels.ExerciseForce{},
		&exerciseModels.ExerciseMuscle{},
		&exerciseModels.ExerciseMeasurementParadigm{},
		&exerciseModels.ExerciseInstruction{},
		&exerciseModels.ExerciseImage{},
		&exerciseModels.ExerciseAlternativeName{},
		&equipmentModels.EquipmentEntity{},
		&exerciseModels.ExerciseEquipment{},
		&fulfillmentModels.FulfillmentEntity{},
		&relModels.ExerciseRelationshipEntity{},
		&exerciseModels.ExerciseHistoryEntity{},
		&equipmentModels.EquipmentHistoryEntity{},
	)
	database.DB = db

	// Create profiles for test owner values
	for _, id := range []string{"sinon", "s"} {
		db.Create(&profileModels.UserProfileEntity{ID: id, Name: id})
	}
}

func writeTempJSON(t *testing.T, base, dir, filename string, data any) {
	t.Helper()
	dirPath := filepath.Join(base, "data", dir)
	os.MkdirAll(dirPath, 0755)
	content, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dirPath, filename), content, 0644); err != nil {
		t.Fatal(err)
	}
}

func chdirTemp(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(origDir) })
	return tmpDir
}

// --- main ---

func TestMainFunction(t *testing.T) {
	tmpDir := chdirTemp(t)

	// Create minimal data for all 6 entity types
	ts := int64(1700000000)
	writeTempJSON(t, tmpDir, "compendium_equipments", "eq.json", map[string]any{
		"name": "bar", "displayName": "Bar", "description": "", "category": "free_weights",
		"imageUrl": nil, "templateId": "bar",
	})
	writeTempJSON(t, tmpDir, "compendium_equipment_fulfillment", "f.json", map[string]any{
		"createdBy": "sinon", "createdAt": ts, "equipmentTemplateId": "bar", "fulfillsEquipmentTemplateId": "bar",
	})
	writeTempJSON(t, tmpDir, "compendium_exercises", "ex.json", map[string]any{
		"name": "X", "type": "STRENGTH", "force": []string{}, "primaryMuscles": []string{},
		"secondaryMuscles": []string{}, "technicalDifficulty": "beginner", "bodyWeightScaling": 0.0,
		"suggestedMeasurementParadigms": []string{}, "description": "", "instructions": []string{},
		"images": []string{}, "alternativeNames": []string{}, "authorName": nil, "authorUrl": nil,
		"createdBy": "sinon", "createdAt": ts, "updatedAt": nil, "version": 0,
		"parentExerciseId": nil, "templateId": "x", "equipmentIds": []string{},
	})
	writeTempJSON(t, tmpDir, "compendium_relationships", "rel.json", map[string]any{
		"id": "x-y", "relationshipType": "similar", "strength": 0.5, "description": nil,
		"createdBy": "sinon", "createdAt": ts, "fromExerciseTemplateId": "x", "toExerciseTemplateId": "x",
	})
	// main() calls database.Init() which creates gesitr.db in cwd (tmpDir)
	main()

	// Verify all entities were seeded
	var eqCount, fCount, exCount, relCount int64
	database.DB.Model(&equipmentModels.EquipmentEntity{}).Count(&eqCount)
	database.DB.Model(&fulfillmentModels.FulfillmentEntity{}).Count(&fCount)
	database.DB.Model(&exerciseModels.ExerciseEntity{}).Count(&exCount)
	database.DB.Model(&relModels.ExerciseRelationshipEntity{}).Count(&relCount)

	if eqCount != 1 || fCount != 1 || exCount != 1 || relCount != 1 {
		t.Errorf("counts: eq=%d f=%d ex=%d rel=%d", eqCount, fCount, exCount, relCount)
	}

	// Verify history entries were created
	var eqHistCount, exHistCount int64
	database.DB.Model(&exerciseModels.ExerciseHistoryEntity{}).Count(&exHistCount)
	database.DB.Model(&equipmentModels.EquipmentHistoryEntity{}).Count(&eqHistCount)
	if exHistCount != 1 || eqHistCount != 1 {
		t.Errorf("history counts: exerciseHistory=%d equipmentHistory=%d", exHistCount, eqHistCount)
	}
}

// --- unixToTime ---

func TestUnixToTime(t *testing.T) {
	t.Run("nil returns now-ish", func(t *testing.T) {
		before := time.Now().Add(-time.Second)
		result := unixToTime(nil)
		after := time.Now().Add(time.Second)
		if result.Before(before) || result.After(after) {
			t.Errorf("expected time near now, got %v", result)
		}
	})

	t.Run("non-nil converts", func(t *testing.T) {
		ts := int64(1700000000)
		result := unixToTime(&ts)
		if !result.Equal(time.Unix(1700000000, 0)) {
			t.Errorf("expected %v, got %v", time.Unix(1700000000, 0), result)
		}
	})
}

// --- readDir ---

func TestReadDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "a.json"), []byte(`{"a":1}`), 0644)
		os.WriteFile(filepath.Join(dir, "b.json"), []byte(`{"b":2}`), 0644)
		os.MkdirAll(filepath.Join(dir, "subdir"), 0755) // should be skipped

		results, err := readDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 2 {
			t.Errorf("expected 2 files, got %d", len(results))
		}
	})

	t.Run("dir not found", func(t *testing.T) {
		_, err := readDir("/nonexistent/path/xyz")
		if err == nil {
			t.Error("expected error for nonexistent directory")
		}
	})

	t.Run("file read error", func(t *testing.T) {
		dir := t.TempDir()
		// Create a file that can't be read
		path := filepath.Join(dir, "unreadable.json")
		os.WriteFile(path, []byte("data"), 0644)
		os.Chmod(path, 0000)
		t.Cleanup(func() { os.Chmod(path, 0644) })

		_, err := readDir(dir)
		if err == nil {
			t.Error("expected error for unreadable file")
		}
	})
}

// --- seedEquipment ---

func TestSeedEquipment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)

		writeTempJSON(t, tmpDir, "compendium_equipments", "barbell.json", map[string]any{
			"name": "barbell", "displayName": "Barbell", "description": "A bar",
			"category": "free_weights", "imageUrl": nil, "templateId": "barbell",
		})
		writeTempJSON(t, tmpDir, "compendium_equipments", "bench.json", map[string]any{
			"name": "bench", "displayName": "Bench", "description": "A bench",
			"category": "benches", "imageUrl": nil, "templateId": "bench",
		})

		if err := seedEquipment(); err != nil {
			t.Fatal(err)
		}

		var count int64
		database.DB.Model(&equipmentModels.EquipmentEntity{}).Count(&count)
		if count != 2 {
			t.Errorf("expected 2, got %d", count)
		}

		var eq equipmentModels.EquipmentEntity
		database.DB.Where("name = ?", "barbell").First(&eq)
		if eq.Name != "barbell" || eq.Category != "free_weights" || eq.Owner != "sinon" || !eq.Public {
			t.Errorf("field mismatch: %+v", eq)
		}

		var histCount int64
		database.DB.Model(&equipmentModels.EquipmentHistoryEntity{}).Count(&histCount)
		if histCount != 2 {
			t.Errorf("expected 2 equipment history entries, got %d", histCount)
		}

		// Verify equipmentIDMap was populated
		if len(equipmentIDMap) != 2 {
			t.Errorf("expected equipmentIDMap to have 2 entries, got %d", len(equipmentIDMap))
		}
	})

	t.Run("bad json", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		dirPath := filepath.Join(tmpDir, "data", "compendium_equipments")
		os.MkdirAll(dirPath, 0755)
		os.WriteFile(filepath.Join(dirPath, "bad.json"), []byte("{invalid"), 0644)

		err := seedEquipment()
		if err == nil {
			t.Error("expected error for bad JSON")
		}
	})

	t.Run("dir not found", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		_ = tmpDir // no data dir created

		err := seedEquipment()
		if err == nil {
			t.Error("expected error for missing directory")
		}
	})

	t.Run("db error", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		writeTempJSON(t, tmpDir, "compendium_equipments", "a.json", map[string]any{
			"name": "a", "displayName": "A", "description": "",
			"category": "other", "imageUrl": nil, "templateId": "a",
		})
		sqlDB, _ := database.DB.DB()
		sqlDB.Close()

		err := seedEquipment()
		if err == nil {
			t.Error("expected error for closed DB")
		}
	})
}

// --- seedFulfillments ---

func TestSeedFulfillments(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)

		// Seed equipment first so equipmentIDMap is populated
		writeTempJSON(t, tmpDir, "compendium_equipments", "a.json", map[string]any{
			"name": "a", "displayName": "A", "description": "",
			"category": "other", "imageUrl": nil, "templateId": "a",
		})
		writeTempJSON(t, tmpDir, "compendium_equipments", "b.json", map[string]any{
			"name": "b", "displayName": "B", "description": "",
			"category": "other", "imageUrl": nil, "templateId": "b",
		})
		if err := seedEquipment(); err != nil {
			t.Fatal(err)
		}

		ts := int64(1700000000)
		writeTempJSON(t, tmpDir, "compendium_equipment_fulfillment", "a~b.json", map[string]any{
			"createdBy": "sinon", "createdAt": ts,
			"equipmentTemplateId": "a", "fulfillsEquipmentTemplateId": "b",
		})

		if err := seedFulfillments(); err != nil {
			t.Fatal(err)
		}

		var count int64
		database.DB.Model(&fulfillmentModels.FulfillmentEntity{}).Count(&count)
		if count != 1 {
			t.Errorf("expected 1, got %d", count)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		dirPath := filepath.Join(tmpDir, "data", "compendium_equipment_fulfillment")
		os.MkdirAll(dirPath, 0755)
		os.WriteFile(filepath.Join(dirPath, "bad.json"), []byte("{bad"), 0644)

		if err := seedFulfillments(); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("dir not found", func(t *testing.T) {
		setupSeedTestDB(t)
		chdirTemp(t)
		if err := seedFulfillments(); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("db error", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		// Initialize equipmentIDMap so JSON parsing works
		equipmentIDMap = map[string]uint{"a": 1, "b": 2}
		writeTempJSON(t, tmpDir, "compendium_equipment_fulfillment", "a.json", map[string]any{
			"createdBy": "sinon", "createdAt": 0, "equipmentTemplateId": "a", "fulfillsEquipmentTemplateId": "b",
		})
		sqlDB, _ := database.DB.DB()
		sqlDB.Close()
		if err := seedFulfillments(); err == nil {
			t.Error("expected error")
		}
	})
}

// --- seedExercises ---

func TestSeedExercises(t *testing.T) {
	t.Run("success with updatedAt", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)

		// Seed equipment first so equipmentIDMap is populated
		writeTempJSON(t, tmpDir, "compendium_equipments", "barbell.json", map[string]any{
			"name": "barbell", "displayName": "Barbell", "description": "",
			"category": "free_weights", "imageUrl": nil, "templateId": "barbell",
		})
		if err := seedEquipment(); err != nil {
			t.Fatal(err)
		}

		ts := int64(1700000000)
		updTs := int64(1700001000)
		writeTempJSON(t, tmpDir, "compendium_exercises", "squat.json", map[string]any{
			"name": "Squat", "type": "STRENGTH",
			"force": []string{"PUSH"}, "primaryMuscles": []string{"QUADS"},
			"secondaryMuscles":    []string{"GLUTES"},
			"technicalDifficulty": "intermediate", "bodyWeightScaling": 1.0,
			"suggestedMeasurementParadigms": []string{"REP_BASED"},
			"description":                   "A squat", "instructions": []string{"Go down", "Go up"},
			"images": []string{"/img/squat.jpg"}, "alternativeNames": []string{"Back Squat"},
			"authorName": nil, "authorUrl": nil,
			"createdBy": "sinon", "createdAt": ts, "updatedAt": updTs,
			"version": 0, "parentExerciseId": nil, "templateId": "squat",
			"equipmentIds": []string{"barbell"},
		})

		if err := seedExercises(); err != nil {
			t.Fatal(err)
		}

		var count int64
		database.DB.Model(&exerciseModels.ExerciseEntity{}).Count(&count)
		if count != 1 {
			t.Errorf("expected 1, got %d", count)
		}

		// Verify child records
		var ex exerciseModels.ExerciseEntity
		database.DB.Where("name = ?", "Squat").First(&ex)
		var fc, mc, pc, ic, imgc, alc, eqc int64
		database.DB.Model(&exerciseModels.ExerciseForce{}).Where("exercise_id = ?", ex.ID).Count(&fc)
		database.DB.Model(&exerciseModels.ExerciseMuscle{}).Where("exercise_id = ?", ex.ID).Count(&mc)
		database.DB.Model(&exerciseModels.ExerciseMeasurementParadigm{}).Where("exercise_id = ?", ex.ID).Count(&pc)
		database.DB.Model(&exerciseModels.ExerciseInstruction{}).Where("exercise_id = ?", ex.ID).Count(&ic)
		database.DB.Model(&exerciseModels.ExerciseImage{}).Where("exercise_id = ?", ex.ID).Count(&imgc)
		database.DB.Model(&exerciseModels.ExerciseAlternativeName{}).Where("exercise_id = ?", ex.ID).Count(&alc)
		database.DB.Model(&exerciseModels.ExerciseEquipment{}).Where("exercise_id = ?", ex.ID).Count(&eqc)
		if fc != 1 || mc != 2 || pc != 1 || ic != 2 || imgc != 1 || alc != 1 || eqc != 1 {
			t.Errorf("child counts: forces=%d muscles=%d paradigms=%d instr=%d img=%d alt=%d eq=%d",
				fc, mc, pc, ic, imgc, alc, eqc)
		}

		var histCount int64
		database.DB.Model(&exerciseModels.ExerciseHistoryEntity{}).Where("exercise_id = ?", ex.ID).Count(&histCount)
		if histCount != 1 {
			t.Errorf("expected 1 exercise history entry, got %d", histCount)
		}
	})

	t.Run("success with nil updatedAt", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)

		// Initialize equipmentIDMap (empty, no equipment needed)
		equipmentIDMap = make(map[string]uint)

		ts := int64(1700000000)
		writeTempJSON(t, tmpDir, "compendium_exercises", "curl.json", map[string]any{
			"name": "Curl", "type": "STRENGTH",
			"force": []string{}, "primaryMuscles": []string{},
			"secondaryMuscles":    []string{},
			"technicalDifficulty": "beginner", "bodyWeightScaling": 0.0,
			"suggestedMeasurementParadigms": []string{}, "description": "",
			"instructions": []string{}, "images": []string{},
			"alternativeNames": []string{}, "authorName": nil, "authorUrl": nil,
			"createdBy": "sinon", "createdAt": ts, "updatedAt": nil,
			"version": 0, "parentExerciseId": nil, "templateId": "curl",
			"equipmentIds": []string{},
		})

		if err := seedExercises(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		dirPath := filepath.Join(tmpDir, "data", "compendium_exercises")
		os.MkdirAll(dirPath, 0755)
		os.WriteFile(filepath.Join(dirPath, "bad.json"), []byte("{bad"), 0644)

		if err := seedExercises(); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("dir not found", func(t *testing.T) {
		setupSeedTestDB(t)
		chdirTemp(t)
		if err := seedExercises(); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("db error", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		// Initialize equipmentIDMap
		equipmentIDMap = make(map[string]uint)
		ts := int64(1700000000)
		writeTempJSON(t, tmpDir, "compendium_exercises", "x.json", map[string]any{
			"name": "X", "type": "STRENGTH",
			"force": []string{}, "primaryMuscles": []string{}, "secondaryMuscles": []string{},
			"technicalDifficulty": "beginner", "bodyWeightScaling": 0.0,
			"suggestedMeasurementParadigms": []string{}, "description": "",
			"instructions": []string{}, "images": []string{},
			"alternativeNames": []string{}, "authorName": nil, "authorUrl": nil,
			"createdBy": "sinon", "createdAt": ts, "updatedAt": nil,
			"version": 0, "parentExerciseId": nil, "templateId": "x",
			"equipmentIds": []string{},
		})
		sqlDB, _ := database.DB.DB()
		sqlDB.Close()
		if err := seedExercises(); err == nil {
			t.Error("expected error")
		}
	})
}

// --- seedExerciseRelationships ---

func TestSeedExerciseRelationships(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)

		// Initialize exerciseIDMap with test data
		exerciseIDMap = map[string]uint{"a": 1, "b": 2}

		ts := int64(1700000000)
		writeTempJSON(t, tmpDir, "compendium_relationships", "rel1.json", map[string]any{
			"id": "a-b-similar", "relationshipType": "similar", "strength": 0.8,
			"description": nil, "createdBy": "sinon", "createdAt": ts,
			"fromExerciseTemplateId": "a", "toExerciseTemplateId": "b",
		})

		if err := seedExerciseRelationships(); err != nil {
			t.Fatal(err)
		}

		var count int64
		database.DB.Model(&relModels.ExerciseRelationshipEntity{}).Count(&count)
		if count != 1 {
			t.Errorf("expected 1, got %d", count)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		dirPath := filepath.Join(tmpDir, "data", "compendium_relationships")
		os.MkdirAll(dirPath, 0755)
		os.WriteFile(filepath.Join(dirPath, "bad.json"), []byte("{bad"), 0644)
		if err := seedExerciseRelationships(); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("dir not found", func(t *testing.T) {
		setupSeedTestDB(t)
		chdirTemp(t)
		if err := seedExerciseRelationships(); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("db error", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		// Initialize exerciseIDMap
		exerciseIDMap = map[string]uint{"a": 1, "b": 2}
		writeTempJSON(t, tmpDir, "compendium_relationships", "r.json", map[string]any{
			"id": "x", "relationshipType": "similar", "strength": 0.5,
			"description": nil, "createdBy": "sinon", "createdAt": 0,
			"fromExerciseTemplateId": "a", "toExerciseTemplateId": "b",
		})
		sqlDB, _ := database.DB.DB()
		sqlDB.Close()
		if err := seedExerciseRelationships(); err == nil {
			t.Error("expected error")
		}
	})
}
