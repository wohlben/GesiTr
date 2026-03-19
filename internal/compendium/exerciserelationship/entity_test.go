package exerciserelationship

import (
	"testing"

	"gesitr/internal/shared"
)

func TestExerciseRelationshipEntityTableName(t *testing.T) {
	if got := (ExerciseRelationshipEntity{}).TableName(); got != "exercise_relationships" {
		t.Errorf("TableName() = %q, want %q", got, "exercise_relationships")
	}
}

func TestExerciseRelationshipEntityToDTO(t *testing.T) {
	desc := "test desc"
	e := &ExerciseRelationshipEntity{
		BaseModel:              shared.BaseModel{ID: 5},
		RelationshipType:       ExerciseRelationshipTypeSimilar,
		Strength:               0.8,
		Description:            &desc,
		CreatedBy:              "system",
		FromExerciseTemplateID: "ex1",
		ToExerciseTemplateID:   "ex2",
	}
	dto := e.ToDTO()
	if dto.ID != 5 || dto.RelationshipType != ExerciseRelationshipTypeSimilar || dto.Strength != 0.8 {
		t.Error("ToDTO field mismatch")
	}
	if *dto.Description != "test desc" || dto.CreatedBy != "system" {
		t.Error("ToDTO field mismatch")
	}
	if dto.FromExerciseTemplateID != "ex1" || dto.ToExerciseTemplateID != "ex2" {
		t.Error("ToDTO field mismatch")
	}
}

func TestExerciseRelationshipFromDTO(t *testing.T) {
	desc := "desc"
	dto := ExerciseRelationship{
		BaseModel:              shared.BaseModel{ID: 3},
		RelationshipType:       ExerciseRelationshipTypeVariation,
		Strength:               0.5,
		Description:            &desc,
		CreatedBy:              "user",
		FromExerciseTemplateID: "a",
		ToExerciseTemplateID:   "b",
	}
	e := ExerciseRelationshipFromDTO(dto)
	if e.ID != 3 || e.RelationshipType != ExerciseRelationshipTypeVariation || e.Strength != 0.5 {
		t.Error("FromDTO field mismatch")
	}
	if *e.Description != "desc" || e.CreatedBy != "user" || e.FromExerciseTemplateID != "a" || e.ToExerciseTemplateID != "b" {
		t.Error("FromDTO field mismatch")
	}
}
