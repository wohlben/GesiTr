package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gesitr/internal/database"
	equipmentModels "gesitr/internal/equipment/models"
	fulfillmentModels "gesitr/internal/equipmentfulfillment/models"
	exerciseModels "gesitr/internal/exercise/models"
	groupModels "gesitr/internal/exercisegroup/models"
	relModels "gesitr/internal/exerciserelationship/models"
	profileModels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

var equipmentIDMap map[string]uint
var exerciseIDMap map[string]uint
var groupIDMap map[string]uint

func main() {
	database.Init()
	database.DB.AutoMigrate(
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
		&groupModels.ExerciseGroupEntity{},
		&groupModels.ExerciseGroupMemberEntity{},
		&exerciseModels.ExerciseHistoryEntity{},
		&equipmentModels.EquipmentHistoryEntity{},
	)

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Profile", seedProfile},
		{"Equipment", seedEquipment},
		{"Exercises", seedExercises},
		{"Fulfillments", seedFulfillments},
		{"ExerciseRelationships", seedExerciseRelationships},
		{"ExerciseGroups", seedExerciseGroups},
		{"ExerciseGroupMembers", seedExerciseGroupMembers},
	}
	for _, s := range steps {
		if err := s.fn(); err != nil {
			log.Fatalf("Failed to seed %s: %v", s.name, err)
		}
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

// --- Profile ---

func seedProfile() error {
	profile := profileModels.UserProfileEntity{
		ID:   "sinon",
		Name: "Sinon",
	}
	if err := database.DB.Create(&profile).Error; err != nil {
		return fmt.Errorf("insert profile: %w", err)
	}
	log.Printf("Profile: sinon")
	return nil
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
	for _, data := range files {
		var j jsonEquipment
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse equipment JSON: %w", err)
		}
		entities = append(entities, equipmentModels.EquipmentEntity{
			Name:        j.Name,
			DisplayName: j.DisplayName,
			Description: j.Description,
			Category:    equipmentModels.EquipmentCategory(j.Category),
			ImageUrl:    j.ImageUrl,
			TemplateID:  j.TemplateID,
			Owner:       "sinon",
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
			ChangedBy:   entities[i].Owner,
		})
	}
	if err := database.DB.CreateInBatches(history, 100).Error; err != nil {
		return fmt.Errorf("insert equipment history: %w", err)
	}

	// Build equipmentIDMap for downstream seeders
	var allEquipment []equipmentModels.EquipmentEntity
	database.DB.Find(&allEquipment)
	equipmentIDMap = make(map[string]uint)
	for _, eq := range allEquipment {
		equipmentIDMap[eq.TemplateID] = eq.ID
	}

	log.Printf("Equipment: %d", len(entities))
	return nil
}

// --- Exercises ---

type jsonExercise struct {
	Name                          string   `json:"name"`
	Slug                          string   `json:"slug"`
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
	for _, data := range files {
		var j jsonExercise
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse exercise JSON: %w", err)
		}
		templateID := ""
		if j.TemplateID != nil {
			templateID = *j.TemplateID
		}
		e := exerciseModels.ExerciseEntity{
			Name:                j.Name,
			Slug:                j.Slug,
			Type:                exerciseModels.ExerciseType(j.Type),
			TechnicalDifficulty: exerciseModels.TechnicalDifficulty(j.TechnicalDifficulty),
			BodyWeightScaling:   j.BodyWeightScaling,
			Description:         j.Description,
			AuthorName:          j.AuthorName,
			AuthorUrl:           j.AuthorUrl,
			Owner:               "sinon",
			Public:              true,
			Version:             j.Version,
			ParentExerciseID:    j.ParentExerciseID,
			TemplateID:          templateID,
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
		for _, name := range j.AlternativeNames {
			e.AlternativeNames = append(e.AlternativeNames, exerciseModels.ExerciseAlternativeName{Name: name})
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
			ChangedBy:  entities[i].Owner,
		})
	}
	if err := database.DB.CreateInBatches(history, 100).Error; err != nil {
		return fmt.Errorf("insert exercise history: %w", err)
	}

	// Build exerciseIDMap for downstream seeders
	var allExercises []exerciseModels.ExerciseEntity
	database.DB.Where("owner = ?", "sinon").Find(&allExercises)
	exerciseIDMap = make(map[string]uint)
	for _, ex := range allExercises {
		exerciseIDMap[ex.TemplateID] = ex.ID
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
	var entities []fulfillmentModels.FulfillmentEntity
	for _, data := range files {
		var j jsonFulfillment
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse fulfillment JSON: %w", err)
		}
		e := fulfillmentModels.FulfillmentEntity{
			EquipmentID:         equipmentIDMap[j.EquipmentTemplateID],
			FulfillsEquipmentID: equipmentIDMap[j.FulfillsEquipmentTemplateID],
			Owner:               "sinon",
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
	var entities []relModels.ExerciseRelationshipEntity
	for _, data := range files {
		var j jsonExerciseRelationship
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse relationship JSON: %w", err)
		}
		e := relModels.ExerciseRelationshipEntity{
			RelationshipType: relModels.ExerciseRelationshipType(j.RelationshipType),
			Strength:         j.Strength,
			Description:      j.Description,
			Owner:            "sinon",
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

// --- Exercise Groups ---

type jsonExerciseGroup struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedBy   string  `json:"createdBy"`
	CreatedAt   *int64  `json:"createdAt"`
	UpdatedAt   *int64  `json:"updatedAt"`
}

func seedExerciseGroups() error {
	files, err := readDir("data/compendium_exercise_groups")
	if err != nil {
		return err
	}
	var entities []groupModels.ExerciseGroupEntity
	for _, data := range files {
		var j jsonExerciseGroup
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse exercise group JSON: %w", err)
		}
		e := groupModels.ExerciseGroupEntity{
			TemplateID:  j.ID,
			Name:        j.Name,
			Description: j.Description,
			Owner:       "sinon",
		}
		e.CreatedAt = unixToTime(j.CreatedAt)
		if j.UpdatedAt != nil {
			e.UpdatedAt = time.Unix(*j.UpdatedAt, 0)
		} else {
			e.UpdatedAt = e.CreatedAt
		}
		entities = append(entities, e)
	}
	if err := database.DB.CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("insert exercise groups: %w", err)
	}

	// Build groupIDMap for downstream seeders
	var allGroups []groupModels.ExerciseGroupEntity
	database.DB.Where("owner = ?", "sinon").Find(&allGroups)
	groupIDMap = make(map[string]uint)
	for _, g := range allGroups {
		groupIDMap[g.TemplateID] = g.ID
	}

	log.Printf("ExerciseGroups: %d", len(entities))
	return nil
}

// --- Exercise Group Members ---

type jsonExerciseGroupMember struct {
	GroupID            string `json:"groupId"`
	ExerciseTemplateID string `json:"exerciseTemplateId"`
	AddedBy            string `json:"addedBy"`
	AddedAt            *int64 `json:"addedAt"`
}

func seedExerciseGroupMembers() error {
	files, err := readDir("data/compendium_exercise_group_members")
	if err != nil {
		return err
	}
	var entities []groupModels.ExerciseGroupMemberEntity
	for _, data := range files {
		var j jsonExerciseGroupMember
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse exercise group member JSON: %w", err)
		}
		e := groupModels.ExerciseGroupMemberEntity{
			GroupID:    groupIDMap[j.GroupID],
			ExerciseID: exerciseIDMap[j.ExerciseTemplateID],
			Owner:      "sinon",
		}
		e.CreatedAt = unixToTime(j.AddedAt)
		entities = append(entities, e)
	}
	if err := database.DB.CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("insert exercise group members: %w", err)
	}
	log.Printf("ExerciseGroupMembers: %d", len(entities))
	return nil
}
