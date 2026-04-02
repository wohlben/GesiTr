package handlers

import (
	"gesitr/internal/workoutrelationship/models"
)

type WorkoutRelationshipBody struct {
	RelationshipType models.WorkoutRelationshipType `json:"relationshipType" required:"true"`
	Strength         float64                        `json:"strength"`
	FromWorkoutID    uint                           `json:"fromWorkoutId" required:"true"`
	ToWorkoutID      uint                           `json:"toWorkoutId" required:"true"`
}

type ListWorkoutRelationshipsInput struct {
	Owner            string `query:"owner" doc:"Filter by owner"`
	FromWorkoutID    string `query:"fromWorkoutId" doc:"Filter by source workout ID"`
	ToWorkoutID      string `query:"toWorkoutId" doc:"Filter by target workout ID"`
	RelationshipType string `query:"relationshipType" doc:"Filter by relationship type"`
}

type ListWorkoutRelationshipsOutput struct {
	Body []models.WorkoutRelationship
}

type CreateWorkoutRelationshipInput struct {
	Body WorkoutRelationshipBody
}

type CreateWorkoutRelationshipOutput struct {
	Body models.WorkoutRelationship
}

type DeleteWorkoutRelationshipInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutRelationshipOutput struct{}
