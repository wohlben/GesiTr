package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	equipmentModels "gesitr/internal/compendium/equipment/models"
	exerciseModels "gesitr/internal/compendium/exercise/models"
	workoutModels "gesitr/internal/compendium/workout/models"
	"gesitr/internal/database"
	"gesitr/internal/compendium/ownershipgroup"
	ownershipGroupModels "gesitr/internal/compendium/ownershipgroup/models"
	"gesitr/internal/shared"
	exerciseSchemeModels "gesitr/internal/user/exercisescheme/models"

	"gorm.io/gorm"
)

var equipmentIDMap map[string]uint
var exerciseIDMap map[string]uint

func main() {
	database.Init()
	database.DB.AutoMigrate(
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
		&exerciseSchemeModels.ExerciseSchemeSectionItemEntity{},
		&workoutModels.WorkoutHistoryEntity{},
		&ownershipGroupModels.OwnershipGroupEntity{},
		&ownershipGroupModels.OwnershipGroupMembershipEntity{},
	)

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Equipment", seedEquipment},
		{"Exercises", seedExercises},
		{"Fulfillments", seedFulfillments},
		{"ExerciseRelationships", seedExerciseRelationships},
		{"Workouts", seedWorkouts},
	}
	for _, s := range steps {
		if err := s.fn(); err != nil {
			log.Fatalf("Failed to seed %s: %v", s.name, err)
		}
	}

	// Create ownership groups for all seeded entities.
	if err := assignOwnershipGroups(); err != nil {
		log.Fatalf("Failed to assign ownership groups: %v", err)
	}
}

func readDir(dir string) ([][]byte, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read directory %s: %w", dir, err)
	}
	var results [][]byte
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read file %s: %w", entry.Name(), err)
		}
		results = append(results, data)
	}
	return results, nil
}

func unixToTime(ts *int64) time.Time {
	if ts == nil {
		return time.Now()
	}
	return time.Unix(*ts, 0)
}

// --- Equipment ---

type jsonEquipment struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"displayName"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	ImageUrl    *string `json:"imageUrl"`
	TemplateID  string  `json:"templateId"`
}

func seedEquipment() error {
	files, err := readDir("data/compendium_equipments")
	if err != nil {
		return err
	}
	var entities []equipmentModels.EquipmentEntity
	var equipTemplateIDs []string
	for _, data := range files {
		var j jsonEquipment
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse equipment JSON: %w", err)
		}
		equipTemplateIDs = append(equipTemplateIDs, j.TemplateID)
		entities = append(entities, equipmentModels.EquipmentEntity{
			Name:        j.Name,
			DisplayName: j.DisplayName,
			Description: j.Description,
			Category:    equipmentModels.EquipmentCategory(j.Category),
			ImageUrl:    j.ImageUrl,
			Public:      true,
		})
	}
	if err := database.DB.CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("insert equipment: %w", err)
	}

	var history []equipmentModels.EquipmentHistoryEntity
	for i := range entities {
		dto := entities[i].ToDTO()
		history = append(history, equipmentModels.EquipmentHistoryEntity{
			EquipmentID: entities[i].ID,
			Version:     0,
			Snapshot:    shared.SnapshotJSON(dto),
			ChangedAt:   entities[i].CreatedAt,
			ChangedBy:   "sinon",
		})
	}
	if err := database.DB.CreateInBatches(history, 100).Error; err != nil {
		return fmt.Errorf("insert equipment history: %w", err)
	}

	// Build equipmentIDMap for downstream seeders (templateID from seed JSON → DB ID)
	equipmentIDMap = make(map[string]uint)
	for i, eq := range entities {
		if equipTemplateIDs[i] != "" {
			equipmentIDMap[equipTemplateIDs[i]] = eq.ID
		}
	}

	log.Printf("Equipment: %d", len(entities))
	return nil
}

// --- Exercises ---

type jsonExercise struct {
	Name                          string   `json:"name"`
	Type                          string   `json:"type"`
	Force                         []string `json:"force"`
	PrimaryMuscles                []string `json:"primaryMuscles"`
	SecondaryMuscles              []string `json:"secondaryMuscles"`
	TechnicalDifficulty           string   `json:"technicalDifficulty"`
	BodyWeightScaling             float64  `json:"bodyWeightScaling"`
	SuggestedMeasurementParadigms []string `json:"suggestedMeasurementParadigms"`
	Description                   string   `json:"description"`
	Instructions                  []string `json:"instructions"`
	Images                        []string `json:"images"`
	AlternativeNames              []string `json:"alternativeNames"`
	AuthorName                    *string  `json:"authorName"`
	AuthorUrl                     *string  `json:"authorUrl"`
	CreatedBy                     string   `json:"createdBy"`
	CreatedAt                     *int64   `json:"createdAt"`
	UpdatedAt                     *int64   `json:"updatedAt"`
	Version                       int      `json:"version"`
	ParentExerciseID              *uint    `json:"parentExerciseId"`
	TemplateID                    *string  `json:"templateId"`
	EquipmentIDs                  []string `json:"equipmentIds"`
}

func seedExercises() error {
	files, err := readDir("data/compendium_exercises")
	if err != nil {
		return err
	}
	var entities []exerciseModels.ExerciseEntity
	var templateIDs []string
	for _, data := range files {
		var j jsonExercise
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse exercise JSON: %w", err)
		}
		templateID := ""
		if j.TemplateID != nil {
			templateID = *j.TemplateID
		}
		templateIDs = append(templateIDs, templateID)
		e := exerciseModels.ExerciseEntity{
			Type:                exerciseModels.ExerciseType(j.Type),
			TechnicalDifficulty: exerciseModels.TechnicalDifficulty(j.TechnicalDifficulty),
			BodyWeightScaling:   j.BodyWeightScaling,
			Description:         j.Description,
			AuthorName:          j.AuthorName,
			AuthorUrl:           j.AuthorUrl,
			Public:              true,
			Version:             j.Version,
			ParentExerciseID:    j.ParentExerciseID,
		}
		e.CreatedAt = unixToTime(j.CreatedAt)
		if j.UpdatedAt != nil {
			e.UpdatedAt = time.Unix(*j.UpdatedAt, 0)
		} else {
			e.UpdatedAt = e.CreatedAt
		}

		for _, f := range j.Force {
			e.Forces = append(e.Forces, exerciseModels.ExerciseForce{Force: exerciseModels.Force(f)})
		}
		for _, m := range j.PrimaryMuscles {
			e.Muscles = append(e.Muscles, exerciseModels.ExerciseMuscle{Muscle: exerciseModels.Muscle(m), IsPrimary: true})
		}
		for _, m := range j.SecondaryMuscles {
			e.Muscles = append(e.Muscles, exerciseModels.ExerciseMuscle{Muscle: exerciseModels.Muscle(m), IsPrimary: false})
		}
		for _, p := range j.SuggestedMeasurementParadigms {
			e.Paradigms = append(e.Paradigms, exerciseModels.ExerciseMeasurementParadigm{Paradigm: exerciseModels.MeasurementParadigm(p)})
		}
		for i, text := range j.Instructions {
			e.Instructions = append(e.Instructions, exerciseModels.ExerciseInstruction{Position: i, Text: text})
		}
		for i, path := range j.Images {
			e.Images = append(e.Images, exerciseModels.ExerciseImage{Position: i, Path: path})
		}
		// Merge name + alternativeNames into a single deduplicated names list.
		seen := map[string]bool{}
		pos := 0
		addName := func(n string) {
			if n != "" && !seen[n] {
				seen[n] = true
				e.Names = append(e.Names, exerciseModels.ExerciseName{Position: pos, Name: n})
				pos++
			}
		}
		addName(j.Name)
		for _, n := range j.AlternativeNames {
			addName(n)
		}
		for _, tid := range j.EquipmentIDs {
			if id, ok := equipmentIDMap[tid]; ok {
				e.Equipment = append(e.Equipment, exerciseModels.ExerciseEquipment{EquipmentID: id})
			}
		}

		entities = append(entities, e)
	}
	if err := database.DB.CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("insert exercises: %w", err)
	}

	var history []exerciseModels.ExerciseHistoryEntity
	for i := range entities {
		dto := entities[i].ToDTO()
		history = append(history, exerciseModels.ExerciseHistoryEntity{
			ExerciseID: entities[i].ID,
			Version:    0,
			Snapshot:   shared.SnapshotJSON(dto),
			ChangedAt:  entities[i].CreatedAt,
			ChangedBy:  "sinon",
		})
	}
	if err := database.DB.CreateInBatches(history, 100).Error; err != nil {
		return fmt.Errorf("insert exercise history: %w", err)
	}

	// Build exerciseIDMap for downstream seeders (templateID from seed JSON → DB ID)
	exerciseIDMap = make(map[string]uint)
	for i, ex := range entities {
		if templateIDs[i] != "" {
			exerciseIDMap[templateIDs[i]] = ex.ID
		}
	}

	log.Printf("Exercises: %d", len(entities))
	return nil
}

// --- Fulfillments ---

type jsonFulfillment struct {
	CreatedBy                   string `json:"createdBy"`
	CreatedAt                   *int64 `json:"createdAt"`
	EquipmentTemplateID         string `json:"equipmentTemplateId"`
	FulfillsEquipmentTemplateID string `json:"fulfillsEquipmentTemplateId"`
}

func seedFulfillments() error {
	files, err := readDir("data/compendium_equipment_fulfillment")
	if err != nil {
		return err
	}
	var entities []equipmentModels.FulfillmentEntity
	for _, data := range files {
		var j jsonFulfillment
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse fulfillment JSON: %w", err)
		}
		e := equipmentModels.FulfillmentEntity{
			EquipmentID:         equipmentIDMap[j.EquipmentTemplateID],
			FulfillsEquipmentID: equipmentIDMap[j.FulfillsEquipmentTemplateID],
		}
		e.CreatedAt = unixToTime(j.CreatedAt)
		entities = append(entities, e)
	}
	if err := database.DB.CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("insert fulfillments: %w", err)
	}
	log.Printf("Fulfillments: %d", len(entities))
	return nil
}

// --- Exercise Relationships ---

type jsonExerciseRelationship struct {
	RelationshipType       string  `json:"relationshipType"`
	Strength               float64 `json:"strength"`
	Description            *string `json:"description"`
	CreatedBy              string  `json:"createdBy"`
	CreatedAt              *int64  `json:"createdAt"`
	FromExerciseTemplateID string  `json:"fromExerciseTemplateId"`
	ToExerciseTemplateID   string  `json:"toExerciseTemplateId"`
}

func seedExerciseRelationships() error {
	files, err := readDir("data/compendium_relationships")
	if err != nil {
		return err
	}
	var entities []exerciseModels.ExerciseRelationshipEntity
	for _, data := range files {
		var j jsonExerciseRelationship
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse relationship JSON: %w", err)
		}
		e := exerciseModels.ExerciseRelationshipEntity{
			RelationshipType: exerciseModels.ExerciseRelationshipType(j.RelationshipType),
			Strength:         j.Strength,
			Description:      j.Description,
			FromExerciseID:   exerciseIDMap[j.FromExerciseTemplateID],
			ToExerciseID:     exerciseIDMap[j.ToExerciseTemplateID],
		}
		e.CreatedAt = unixToTime(j.CreatedAt)
		entities = append(entities, e)
	}
	if err := database.DB.CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("insert exercise relationships: %w", err)
	}
	log.Printf("ExerciseRelationships: %d", len(entities))
	return nil
}

// --- Workouts ---

type jsonWorkoutScheme struct {
	MeasurementType string `json:"measurementType"`
	Sets            *int   `json:"sets"`
	Reps            *int   `json:"reps"`
	RestBetweenSets *int   `json:"restBetweenSets"`
}

type jsonWorkoutItem struct {
	Type               string            `json:"type"`
	Position           int               `json:"position"`
	ExerciseTemplateID string            `json:"exerciseTemplateId"`
	Scheme             jsonWorkoutScheme `json:"scheme"`
}

type jsonWorkoutSection struct {
	Type                 string            `json:"type"`
	Label                *string           `json:"label"`
	Position             int               `json:"position"`
	RestBetweenExercises *int              `json:"restBetweenExercises"`
	Items                []jsonWorkoutItem `json:"items"`
}

type jsonWorkout struct {
	Name      string               `json:"name"`
	Notes     *string              `json:"notes"`
	CreatedBy string               `json:"createdBy"`
	CreatedAt *int64               `json:"createdAt"`
	Version   int                  `json:"version"`
	Sections  []jsonWorkoutSection `json:"sections"`
}

func seedWorkouts() error {
	files, err := readDir("data/compendium_workouts")
	if err != nil {
		return err
	}
	count := 0
	for _, data := range files {
		var j jsonWorkout
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse workout JSON: %w", err)
		}

		workout := workoutModels.WorkoutEntity{
			Name:    j.Name,
			Notes:   j.Notes,
			Public:  true,
			Version: j.Version,
		}
		workout.CreatedAt = unixToTime(j.CreatedAt)
		if err := database.DB.Create(&workout).Error; err != nil {
			return fmt.Errorf("insert workout %q: %w", j.Name, err)
		}

		for _, js := range j.Sections {
			section := workoutModels.WorkoutSectionEntity{
				WorkoutID:            workout.ID,
				Type:                 workoutModels.WorkoutSectionType(js.Type),
				Label:                js.Label,
				Position:             js.Position,
				RestBetweenExercises: js.RestBetweenExercises,
			}
			if err := database.DB.Create(&section).Error; err != nil {
				return fmt.Errorf("insert section for workout %q: %w", j.Name, err)
			}

			for _, ji := range js.Items {
				exerciseID, ok := exerciseIDMap[ji.ExerciseTemplateID]
				if !ok {
					return fmt.Errorf("unknown exercise template ID %q in workout %q", ji.ExerciseTemplateID, j.Name)
				}

				item := workoutModels.WorkoutSectionItemEntity{
					WorkoutSectionID: section.ID,
					Type:             workoutModels.WorkoutSectionItemType(ji.Type),
					ExerciseID:       &exerciseID,
					Position:         ji.Position,
				}
				if err := database.DB.Create(&item).Error; err != nil {
					return fmt.Errorf("insert section item for exercise %q: %w", ji.ExerciseTemplateID, err)
				}

				scheme := exerciseSchemeModels.ExerciseSchemeEntity{
					Owner:           "sinon",
					ExerciseID:      exerciseID,
					MeasurementType: ji.Scheme.MeasurementType,
					Sets:            ji.Scheme.Sets,
					Reps:            ji.Scheme.Reps,
					RestBetweenSets: ji.Scheme.RestBetweenSets,
				}
				if err := database.DB.Create(&scheme).Error; err != nil {
					return fmt.Errorf("insert scheme for exercise %q: %w", ji.ExerciseTemplateID, err)
				}

				// Link the scheme to the section item via join table
				link := exerciseSchemeModels.ExerciseSchemeSectionItemEntity{
					ExerciseSchemeID:     scheme.ID,
					WorkoutSectionItemID: item.ID,
					Owner:                "sinon",
				}
				if err := database.DB.Create(&link).Error; err != nil {
					return fmt.Errorf("insert scheme-section-item link for exercise %q: %w", ji.ExerciseTemplateID, err)
				}
			}
		}

		// Create version 0 history snapshot.
		if err := database.DB.Preload("Sections.Items").First(&workout, workout.ID).Error; err != nil {
			return fmt.Errorf("reload workout %q: %w", j.Name, err)
		}
		dto := workout.ToDTO()
		history := workoutModels.WorkoutHistoryEntity{
			WorkoutID: workout.ID,
			Version:   0,
			Snapshot:  shared.SnapshotJSON(dto),
			ChangedAt: workout.CreatedAt,
			ChangedBy: "sinon",
		}
		if err := database.DB.Create(&history).Error; err != nil {
			return fmt.Errorf("insert workout history for %q: %w", j.Name, err)
		}

		count++
	}
	log.Printf("Workouts: %d", count)
	return nil
}

// assignOwnershipGroups creates an ownership group (owned by "sinon") for each
// top-level seeded entity and propagates the group to sub-entities.
func assignOwnershipGroups() error {
	topLevelTables := []string{"equipment", "exercises", "workouts"}
	subEntityTables := []struct {
		table    string
		parentFK string
		parent   string
	}{
		{"fulfillments", "equipment_id", "equipment"},
		{"exercise_relationships", "from_exercise_id", "exercises"},
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		for _, table := range topLevelTables {
			var ids []uint
			if err := tx.Table(table).
				Where("ownership_group_id = 0 OR ownership_group_id IS NULL").
				Where("deleted_at IS NULL").
				Pluck("id", &ids).Error; err != nil {
				return fmt.Errorf("query %s: %w", table, err)
			}
			for _, id := range ids {
				groupID, err := ownershipgroup.CreateGroupForEntity(tx, "sinon")
				if err != nil {
					return fmt.Errorf("create group for %s id=%d: %w", table, id, err)
				}
				if err := tx.Table(table).Where("id = ?", id).Update("ownership_group_id", groupID).Error; err != nil {
					return fmt.Errorf("update %s id=%d: %w", table, id, err)
				}
			}
			log.Printf("OwnershipGroups: assigned %d groups in %s", len(ids), table)
		}

		for _, se := range subEntityTables {
			result := tx.Exec(fmt.Sprintf(`
				UPDATE %s SET ownership_group_id = (
					SELECT %s.ownership_group_id FROM %s
					WHERE %s.id = %s.%s
				)
				WHERE (ownership_group_id IS NULL OR ownership_group_id = 0)
				AND deleted_at IS NULL
			`, se.table, se.parent, se.parent, se.parent, se.table, se.parentFK))
			if result.Error != nil {
				return fmt.Errorf("propagate to %s: %w", se.table, result.Error)
			}
			log.Printf("OwnershipGroups: propagated %d rows in %s", result.RowsAffected, se.table)
		}

		return nil
	})
}
