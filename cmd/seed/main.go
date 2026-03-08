package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gesitr/internal/compendium/models"
	"gesitr/internal/database"
)

func main() {
	database.Init()
	database.DB.AutoMigrate(
		&models.ExerciseEntity{},
		&models.ExerciseForce{},
		&models.ExerciseMuscle{},
		&models.ExerciseMeasurementParadigm{},
		&models.ExerciseInstruction{},
		&models.ExerciseImage{},
		&models.ExerciseAlternativeName{},
		&models.EquipmentEntity{},
		&models.ExerciseEquipment{},
		&models.FulfillmentEntity{},
		&models.ExerciseRelationshipEntity{},
		&models.ExerciseGroupEntity{},
		&models.ExerciseGroupMemberEntity{},
		&models.ExerciseHistoryEntity{},
		&models.EquipmentHistoryEntity{},
	)

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Equipment", seedEquipment},
		{"Fulfillments", seedFulfillments},
		{"Exercises", seedExercises},
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
	var entities []models.EquipmentEntity
	for _, data := range files {
		var j jsonEquipment
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse equipment JSON: %w", err)
		}
		entities = append(entities, models.EquipmentEntity{
			Name:        j.Name,
			DisplayName: j.DisplayName,
			Description: j.Description,
			Category:    models.EquipmentCategory(j.Category),
			ImageUrl:    j.ImageUrl,
			TemplateID:  j.TemplateID,
			CreatedBy:   "system",
		})
	}
	if err := database.DB.CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("insert equipment: %w", err)
	}

	var history []models.EquipmentHistoryEntity
	for i := range entities {
		dto := entities[i].ToDTO()
		history = append(history, models.EquipmentHistoryEntity{
			EquipmentID: entities[i].ID,
			Version:     0,
			Snapshot:    models.SnapshotJSON(dto),
			ChangedAt:   entities[i].CreatedAt,
			ChangedBy:   entities[i].CreatedBy,
		})
	}
	if err := database.DB.CreateInBatches(history, 100).Error; err != nil {
		return fmt.Errorf("insert equipment history: %w", err)
	}

	log.Printf("Equipment: %d", len(entities))
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
	var entities []models.FulfillmentEntity
	for _, data := range files {
		var j jsonFulfillment
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse fulfillment JSON: %w", err)
		}
		e := models.FulfillmentEntity{
			EquipmentTemplateID:         j.EquipmentTemplateID,
			FulfillsEquipmentTemplateID: j.FulfillsEquipmentTemplateID,
			CreatedBy:                   j.CreatedBy,
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
	var entities []models.ExerciseEntity
	for _, data := range files {
		var j jsonExercise
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse exercise JSON: %w", err)
		}
		e := models.ExerciseEntity{
			Name:                j.Name,
			Slug:                j.Slug,
			Type:                models.ExerciseType(j.Type),
			TechnicalDifficulty: models.TechnicalDifficulty(j.TechnicalDifficulty),
			BodyWeightScaling:   j.BodyWeightScaling,
			Description:         j.Description,
			AuthorName:          j.AuthorName,
			AuthorUrl:           j.AuthorUrl,
			CreatedBy:           j.CreatedBy,
			Version:             j.Version,
			ParentExerciseID:    j.ParentExerciseID,
			TemplateID:          j.TemplateID,
		}
		e.CreatedAt = unixToTime(j.CreatedAt)
		if j.UpdatedAt != nil {
			e.UpdatedAt = time.Unix(*j.UpdatedAt, 0)
		} else {
			e.UpdatedAt = e.CreatedAt
		}

		for _, f := range j.Force {
			e.Forces = append(e.Forces, models.ExerciseForce{Force: models.Force(f)})
		}
		for _, m := range j.PrimaryMuscles {
			e.Muscles = append(e.Muscles, models.ExerciseMuscle{Muscle: models.Muscle(m), IsPrimary: true})
		}
		for _, m := range j.SecondaryMuscles {
			e.Muscles = append(e.Muscles, models.ExerciseMuscle{Muscle: models.Muscle(m), IsPrimary: false})
		}
		for _, p := range j.SuggestedMeasurementParadigms {
			e.Paradigms = append(e.Paradigms, models.ExerciseMeasurementParadigm{Paradigm: models.MeasurementParadigm(p)})
		}
		for i, text := range j.Instructions {
			e.Instructions = append(e.Instructions, models.ExerciseInstruction{Position: i, Text: text})
		}
		for i, path := range j.Images {
			e.Images = append(e.Images, models.ExerciseImage{Position: i, Path: path})
		}
		for _, name := range j.AlternativeNames {
			e.AlternativeNames = append(e.AlternativeNames, models.ExerciseAlternativeName{Name: name})
		}
		for _, tid := range j.EquipmentIDs {
			e.Equipment = append(e.Equipment, models.ExerciseEquipment{EquipmentTemplateID: tid})
		}

		entities = append(entities, e)
	}
	if err := database.DB.CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("insert exercises: %w", err)
	}

	var history []models.ExerciseHistoryEntity
	for i := range entities {
		dto := entities[i].ToDTO()
		history = append(history, models.ExerciseHistoryEntity{
			ExerciseID: entities[i].ID,
			Version:    0,
			Snapshot:   models.SnapshotJSON(dto),
			ChangedAt:  entities[i].CreatedAt,
			ChangedBy:  entities[i].CreatedBy,
		})
	}
	if err := database.DB.CreateInBatches(history, 100).Error; err != nil {
		return fmt.Errorf("insert exercise history: %w", err)
	}

	log.Printf("Exercises: %d", len(entities))
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
	var entities []models.ExerciseRelationshipEntity
	for _, data := range files {
		var j jsonExerciseRelationship
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse relationship JSON: %w", err)
		}
		e := models.ExerciseRelationshipEntity{
			RelationshipType:       models.ExerciseRelationshipType(j.RelationshipType),
			Strength:               j.Strength,
			Description:            j.Description,
			CreatedBy:              j.CreatedBy,
			FromExerciseTemplateID: j.FromExerciseTemplateID,
			ToExerciseTemplateID:   j.ToExerciseTemplateID,
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
	var entities []models.ExerciseGroupEntity
	for _, data := range files {
		var j jsonExerciseGroup
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse exercise group JSON: %w", err)
		}
		e := models.ExerciseGroupEntity{
			TemplateID:  j.ID,
			Name:        j.Name,
			Description: j.Description,
			CreatedBy:   j.CreatedBy,
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
	var entities []models.ExerciseGroupMemberEntity
	for _, data := range files {
		var j jsonExerciseGroupMember
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("parse exercise group member JSON: %w", err)
		}
		e := models.ExerciseGroupMemberEntity{
			GroupTemplateID:    j.GroupID,
			ExerciseTemplateID: j.ExerciseTemplateID,
			AddedBy:            j.AddedBy,
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
