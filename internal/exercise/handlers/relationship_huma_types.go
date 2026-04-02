package handlers

import (
	"gesitr/internal/exercise/models"
)

// ExerciseRelationshipBody contains the client-provided fields for creating an exercise relationship.
type ExerciseRelationshipBody struct {
	RelationshipType models.ExerciseRelationshipType `json:"relationshipType" required:"true"`
	Strength         float64                         `json:"strength"`
	Description      *string                         `json:"description"`
	FromExerciseID   uint                            `json:"fromExerciseId" required:"true"`
	ToExerciseID     uint                            `json:"toExerciseId" required:"true"`
}

type ListExerciseRelationshipsInput struct {
	Owner            string `query:"owner" doc:"Filter by owner"`
	FromExerciseID   string `query:"fromExerciseId" doc:"Filter by source exercise ID"`
	ToExerciseID     string `query:"toExerciseId" doc:"Filter by target exercise ID"`
	RelationshipType string `query:"relationshipType" doc:"Filter by relationship type"`
}

type ListExerciseRelationshipsOutput struct {
	Body []models.ExerciseRelationship
}

type CreateExerciseRelationshipInput struct {
	Body ExerciseRelationshipBody
}

type CreateExerciseRelationshipOutput struct {
	Body models.ExerciseRelationship
}

type DeleteExerciseRelationshipInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseRelationshipOutput struct{}
