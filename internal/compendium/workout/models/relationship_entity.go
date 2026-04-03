package models

import (
	"gesitr/internal/shared"
)

type WorkoutRelationshipEntity struct {
	shared.BaseModel
	RelationshipType WorkoutRelationshipType `gorm:"not null;uniqueIndex:idx_workout_relationship"`
	Strength         float64                 `gorm:"not null"`
	Owner            string                  `gorm:"not null;uniqueIndex:idx_workout_relationship"`
	FromWorkoutID    uint                    `gorm:"not null;uniqueIndex:idx_workout_relationship"`
	ToWorkoutID      uint                    `gorm:"not null;uniqueIndex:idx_workout_relationship"`
}

func (WorkoutRelationshipEntity) TableName() string { return "workout_relationships" }

func (e *WorkoutRelationshipEntity) ToDTO() WorkoutRelationship {
	return WorkoutRelationship{
		BaseModel:        e.BaseModel,
		RelationshipType: e.RelationshipType,
		Strength:         e.Strength,
		Owner:            e.Owner,
		FromWorkoutID:    e.FromWorkoutID,
		ToWorkoutID:      e.ToWorkoutID,
	}
}

func WorkoutRelationshipFromDTO(dto WorkoutRelationship) WorkoutRelationshipEntity {
	return WorkoutRelationshipEntity{
		BaseModel:        dto.BaseModel,
		RelationshipType: dto.RelationshipType,
		Strength:         dto.Strength,
		Owner:            dto.Owner,
		FromWorkoutID:    dto.FromWorkoutID,
		ToWorkoutID:      dto.ToWorkoutID,
	}
}
