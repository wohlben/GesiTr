package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	equipmentModels "gesitr/internal/compendium/equipment/models"
	exerciseModels "gesitr/internal/compendium/exercise/models"
	workoutModels "gesitr/internal/compendium/workout/models"
	"gesitr/internal/database"
	exerciseSchemeModels "gesitr/internal/user/exercisescheme/models"

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
		&exerciseModels.ExerciseEntity{},
		&exerciseModels.ExerciseForce{},
		&exerciseModels.ExerciseMuscle{},
		&exerciseModels.ExerciseMeasurementParadigm{},
		&exerciseModels.ExerciseInstruction{},
		&exerciseModels.ExerciseImage{},
		&exerciseModels.ExerciseName{},
		&equipmentModels.EquipmentEntity{},
		&exerciseModels.ExerciseEquipment{},
		&equipmentModels.FulfillmentEntity{},
		&exerciseModels.ExerciseRelationshipEntity{},
		&exerciseModels.ExerciseHistoryEntity{},
		&equipmentModels.EquipmentHistoryEntity{},
		&workoutModels.WorkoutEntity{},
		&workoutModels.WorkoutSectionEntity{},
		&workoutModels.WorkoutSectionItemEntity{},
		&exerciseSchemeModels.ExerciseSchemeEntity{},
		&workoutModels.WorkoutHistoryEntity{},
	)
	database.DB = db
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
	writeTempJSON(t, tmpDir, "compendium_workouts", "w.json", map[string]any{
		"name": "Test Workout", "notes": nil, "createdBy": "sinon", "createdAt": ts, "version": 0,
		"sections": []map[string]any{
			{
				"type": "main", "label": nil, "position": 0, "restBetweenExercises": 60,
				"items": []map[string]any{
					{
						"type": "exercise", "position": 0, "exerciseTemplateId": "x",
						"scheme": map[string]any{"measurementType": "REP_BASED", "sets": 3, "restBetweenSets": 90},
					},
				},
			},
		},
	})
	// main() calls database.Init() which creates gesitr.db in cwd (tmpDir)
	main()

	// Verify all entities were seeded
	var eqCount, fCount, exCount, relCount, wCount int64
	database.DB.Model(&equipmentModels.EquipmentEntity{}).Count(&eqCount)
	database.DB.Model(&equipmentModels.FulfillmentEntity{}).Count(&fCount)
	database.DB.Model(&exerciseModels.ExerciseEntity{}).Count(&exCount)
	database.DB.Model(&exerciseModels.ExerciseRelationshipEntity{}).Count(&relCount)
	database.DB.Model(&workoutModels.WorkoutEntity{}).Count(&wCount)

	if eqCount != 1 || fCount != 1 || exCount != 1 || relCount != 1 || wCount != 1 {
		t.Errorf("counts: eq=%d f=%d ex=%d rel=%d w=%d", eqCount, fCount, exCount, relCount, wCount)
	}

	// Verify history entries were created
	var eqHistCount, exHistCount, wHistCount int64
	database.DB.Model(&exerciseModels.ExerciseHistoryEntity{}).Count(&exHistCount)
	database.DB.Model(&equipmentModels.EquipmentHistoryEntity{}).Count(&eqHistCount)
	database.DB.Model(&workoutModels.WorkoutHistoryEntity{}).Count(&wHistCount)
	if exHistCount != 1 || eqHistCount != 1 || wHistCount != 1 {
		t.Errorf("history counts: exercise=%d equipment=%d workout=%d", exHistCount, eqHistCount, wHistCount)
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
		database.DB.Model(&equipmentModels.FulfillmentEntity{}).Count(&count)
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
		database.DB.First(&ex) // only one exercise seeded
		var fc, mc, pc, ic, imgc, nc, eqc int64
		database.DB.Model(&exerciseModels.ExerciseForce{}).Where("exercise_id = ?", ex.ID).Count(&fc)
		database.DB.Model(&exerciseModels.ExerciseMuscle{}).Where("exercise_id = ?", ex.ID).Count(&mc)
		database.DB.Model(&exerciseModels.ExerciseMeasurementParadigm{}).Where("exercise_id = ?", ex.ID).Count(&pc)
		database.DB.Model(&exerciseModels.ExerciseInstruction{}).Where("exercise_id = ?", ex.ID).Count(&ic)
		database.DB.Model(&exerciseModels.ExerciseImage{}).Where("exercise_id = ?", ex.ID).Count(&imgc)
		database.DB.Model(&exerciseModels.ExerciseName{}).Where("exercise_id = ?", ex.ID).Count(&nc)
		database.DB.Model(&exerciseModels.ExerciseEquipment{}).Where("exercise_id = ?", ex.ID).Count(&eqc)
		// 2 names: "Squat" (from name) + "Back Squat" (from alternativeNames, deduplicated)
		if fc != 1 || mc != 2 || pc != 1 || ic != 2 || imgc != 1 || nc != 2 || eqc != 1 {
			t.Errorf("child counts: forces=%d muscles=%d paradigms=%d instr=%d img=%d names=%d eq=%d",
				fc, mc, pc, ic, imgc, nc, eqc)
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
		database.DB.Model(&exerciseModels.ExerciseRelationshipEntity{}).Count(&count)
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

// --- seedWorkouts ---

func TestSeedWorkouts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)

		// Seed exercises first so exerciseIDMap is populated
		equipmentIDMap = make(map[string]uint)
		ts := int64(1700000000)
		for _, name := range []string{"ex_a", "ex_b", "ex_c"} {
			writeTempJSON(t, tmpDir, "compendium_exercises", name+".json", map[string]any{
				"name": name, "type": "STRENGTH",
				"force": []string{}, "primaryMuscles": []string{}, "secondaryMuscles": []string{},
				"technicalDifficulty": "beginner", "bodyWeightScaling": 0.0,
				"suggestedMeasurementParadigms": []string{}, "description": "",
				"instructions": []string{}, "images": []string{},
				"alternativeNames": []string{}, "authorName": nil, "authorUrl": nil,
				"createdBy": "sinon", "createdAt": ts, "updatedAt": nil,
				"version": 0, "parentExerciseId": nil, "templateId": name,
				"equipmentIds": []string{},
			})
		}
		if err := seedExercises(); err != nil {
			t.Fatal(err)
		}

		sets := 3
		rest := 90
		writeTempJSON(t, tmpDir, "compendium_workouts", "w.json", map[string]any{
			"name": "Test Workout", "notes": "A test", "createdBy": "sinon", "createdAt": ts, "version": 0,
			"sections": []map[string]any{
				{
					"type": "main", "label": nil, "position": 0, "restBetweenExercises": 60,
					"items": []map[string]any{
						{"type": "exercise", "position": 0, "exerciseTemplateId": "ex_a",
							"scheme": map[string]any{"measurementType": "REP_BASED", "sets": sets, "restBetweenSets": rest}},
						{"type": "exercise", "position": 1, "exerciseTemplateId": "ex_b",
							"scheme": map[string]any{"measurementType": "REP_BASED", "sets": sets, "restBetweenSets": rest}},
						{"type": "exercise", "position": 2, "exerciseTemplateId": "ex_c",
							"scheme": map[string]any{"measurementType": "REP_BASED", "sets": sets, "restBetweenSets": rest}},
					},
				},
			},
		})

		if err := seedWorkouts(); err != nil {
			t.Fatal(err)
		}

		// Verify workout
		var workout workoutModels.WorkoutEntity
		if err := database.DB.Preload("Sections.Items").First(&workout).Error; err != nil {
			t.Fatal(err)
		}
		if workout.Name != "Test Workout" || workout.Owner != "sinon" || !workout.Public {
			t.Errorf("workout fields: name=%q owner=%q public=%v", workout.Name, workout.Owner, workout.Public)
		}

		// Verify section
		if len(workout.Sections) != 1 {
			t.Fatalf("expected 1 section, got %d", len(workout.Sections))
		}
		section := workout.Sections[0]
		if section.Type != "main" || *section.RestBetweenExercises != 60 {
			t.Errorf("section: type=%q rest=%v", section.Type, section.RestBetweenExercises)
		}

		// Verify items and schemes
		if len(section.Items) != 3 {
			t.Fatalf("expected 3 items, got %d", len(section.Items))
		}
		expectedExercises := []string{"ex_a", "ex_b", "ex_c"}
		for i, item := range section.Items {
			if item.Position != i {
				t.Errorf("item %d: position=%d", i, item.Position)
			}
			if item.ExerciseSchemeID == nil {
				t.Fatalf("item %d: nil ExerciseSchemeID", i)
			}

			var scheme exerciseSchemeModels.ExerciseSchemeEntity
			if err := database.DB.First(&scheme, *item.ExerciseSchemeID).Error; err != nil {
				t.Fatalf("item %d: scheme not found: %v", i, err)
			}
			if scheme.MeasurementType != "REP_BASED" {
				t.Errorf("item %d: measurement=%q", i, scheme.MeasurementType)
			}
			if scheme.Sets == nil || *scheme.Sets != 3 {
				t.Errorf("item %d: sets=%v", i, scheme.Sets)
			}
			if scheme.RestBetweenSets == nil || *scheme.RestBetweenSets != 90 {
				t.Errorf("item %d: restBetweenSets=%v", i, scheme.RestBetweenSets)
			}

			// Verify scheme points to correct exercise via exerciseIDMap
			expectedID := exerciseIDMap[expectedExercises[i]]
			if scheme.ExerciseID != expectedID {
				t.Errorf("item %d: exerciseID=%d, expected=%d", i, scheme.ExerciseID, expectedID)
			}

			// Verify back-link
			if scheme.WorkoutSectionItemID == nil || *scheme.WorkoutSectionItemID != item.ID {
				t.Errorf("item %d: scheme back-link=%v, expected=%d", i, scheme.WorkoutSectionItemID, item.ID)
			}
		}

		// Verify history
		var histCount int64
		database.DB.Model(&workoutModels.WorkoutHistoryEntity{}).Where("workout_id = ?", workout.ID).Count(&histCount)
		if histCount != 1 {
			t.Errorf("expected 1 history entry, got %d", histCount)
		}
	})

	t.Run("unknown exercise template", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		exerciseIDMap = make(map[string]uint)

		writeTempJSON(t, tmpDir, "compendium_workouts", "w.json", map[string]any{
			"name": "Bad", "notes": nil, "createdBy": "sinon", "createdAt": 0, "version": 0,
			"sections": []map[string]any{
				{
					"type": "main", "label": nil, "position": 0, "restBetweenExercises": nil,
					"items": []map[string]any{
						{"type": "exercise", "position": 0, "exerciseTemplateId": "nonexistent",
							"scheme": map[string]any{"measurementType": "REP_BASED", "sets": 1, "restBetweenSets": 60}},
					},
				},
			},
		})

		err := seedWorkouts()
		if err == nil {
			t.Error("expected error for unknown exercise template ID")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		setupSeedTestDB(t)
		tmpDir := chdirTemp(t)
		dirPath := filepath.Join(tmpDir, "data", "compendium_workouts")
		os.MkdirAll(dirPath, 0755)
		os.WriteFile(filepath.Join(dirPath, "bad.json"), []byte("{bad"), 0644)
		if err := seedWorkouts(); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("dir not found", func(t *testing.T) {
		setupSeedTestDB(t)
		chdirTemp(t)
		if err := seedWorkouts(); err == nil {
			t.Error("expected error")
		}
	})
}
