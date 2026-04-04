package models

import "gesitr/internal/shared"

type WorkoutRelationship struct {
	shared.BaseModel `tstype:",extends"`
	RelationshipType WorkoutRelationshipType `json:"relationshipType"`
	Strength         float64                 `json:"strength"`
	OwnershipGroupID uint                    `json:"ownershipGroupId"`
	FromWorkoutID    uint                    `json:"fromWorkoutId"`
	ToWorkoutID      uint                    `json:"toWorkoutId"`
}

type WorkoutRelationshipType string

const (
	WorkoutRelationshipTypeEquivalent WorkoutRelationshipType = "equivalent"
	WorkoutRelationshipTypeForked     WorkoutRelationshipType = "forked"
)
