package models

import (
	"gesitr/internal/shared"
)

type ExerciseRelationshipEntity struct {
	shared.BaseModel
	RelationshipType ExerciseRelationshipType `gorm:"not null;uniqueIndex:idx_exercise_relationship"`
	Strength         float64                  `gorm:"not null"`
	Description      *string
	OwnershipGroupID uint `gorm:"uniqueIndex:idx_exercise_relationship"`
	FromExerciseID   uint `gorm:"not null;uniqueIndex:idx_exercise_relationship"`
	ToExerciseID     uint `gorm:"not null;uniqueIndex:idx_exercise_relationship"`
}

func (ExerciseRelationshipEntity) TableName() string { return "exercise_relationships" }

func (e *ExerciseRelationshipEntity) ToDTO() ExerciseRelationship {
	return ExerciseRelationship{
		BaseModel:        e.BaseModel,
		RelationshipType: e.RelationshipType,
		Strength:         e.Strength,
		Description:      e.Description,
		OwnershipGroupID: e.OwnershipGroupID,
		FromExerciseID:   e.FromExerciseID,
		ToExerciseID:     e.ToExerciseID,
	}
}

func ExerciseRelationshipFromDTO(dto ExerciseRelationship) ExerciseRelationshipEntity {
	return ExerciseRelationshipEntity{
		BaseModel:        dto.BaseModel,
		RelationshipType: dto.RelationshipType,
		Strength:         dto.Strength,
		Description:      dto.Description,
		OwnershipGroupID: dto.OwnershipGroupID,
		FromExerciseID:   dto.FromExerciseID,
		ToExerciseID:     dto.ToExerciseID,
	}
}
