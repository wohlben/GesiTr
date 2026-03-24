package handlers

import (
	"gesitr/internal/exerciserelationship/models"
)

type ListExerciseRelationshipsInput struct {
	Owner            string `query:"owner" doc:"Filter by owner"`
	FromExerciseID   string `query:"fromExerciseId" doc:"Filter by source exercise ID"`
	ToExerciseID     string `query:"toExerciseId" doc:"Filter by target exercise ID"`
	RelationshipType string `query:"relationshipType" doc:"Filter by relationship type"`
}

type ListExerciseRelationshipsOutput struct {
	Body []models.ExerciseRelationship
}

// RawBody skips huma's automatic validation — the ExerciseRelationship DTO is shared
// between request and response and has server-set fields (id, createdAt, etc.)
// that aren't present in create requests.
type CreateExerciseRelationshipInput struct {
	RawBody []byte
}

type CreateExerciseRelationshipOutput struct {
	Body models.ExerciseRelationship
}

type DeleteExerciseRelationshipInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseRelationshipOutput struct{}
