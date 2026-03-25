package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/exerciserelationship/models"
	"gesitr/internal/humaconfig"

	"github.com/danielgtaylor/huma/v2"
)

// ListExerciseRelationships returns exercise relationships, optionally filtered
// by owner, fromExerciseId, toExerciseId, or relationshipType.
// GET /api/exercise-relationships
//
// OpenAPI: /api/docs#/operations/list-exercise-relationships
func ListExerciseRelationships(ctx context.Context, input *ListExerciseRelationshipsInput) (*ListExerciseRelationshipsOutput, error) {
	db := database.DB.Model(&models.ExerciseRelationshipEntity{})

	if input.Owner != "" {
		owner := input.Owner
		if owner == "me" {
			owner = humaconfig.GetUserID(ctx)
		}
		db = db.Where("owner = ?", owner)
	}
	if input.FromExerciseID != "" {
		db = db.Where("from_exercise_id = ?", input.FromExerciseID)
	}
	if input.ToExerciseID != "" {
		db = db.Where("to_exercise_id = ?", input.ToExerciseID)
	}
	if input.RelationshipType != "" {
		db = db.Where("relationship_type = ?", input.RelationshipType)
	}

	var entities []models.ExerciseRelationshipEntity
	if err := db.Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.ExerciseRelationship, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListExerciseRelationshipsOutput{Body: dtos}, nil
}

// CreateExerciseRelationship creates an exercise relationship owned by the current user.
// POST /api/exercise-relationships
//
// OpenAPI: /api/docs#/operations/create-exercise-relationship
func CreateExerciseRelationship(ctx context.Context, input *CreateExerciseRelationshipInput) (*CreateExerciseRelationshipOutput, error) {
	dto := models.ExerciseRelationship{
		RelationshipType: input.Body.RelationshipType,
		Strength:         input.Body.Strength,
		Description:      input.Body.Description,
		FromExerciseID:   input.Body.FromExerciseID,
		ToExerciseID:     input.Body.ToExerciseID,
	}

	entity := models.ExerciseRelationshipFromDTO(dto)
	entity.Owner = humaconfig.GetUserID(ctx)
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateExerciseRelationshipOutput{Body: entity.ToDTO()}, nil
}

// DeleteExerciseRelationship deletes an exercise relationship. Owner only.
// DELETE /api/exercise-relationships/:id
//
// OpenAPI: /api/docs#/operations/delete-exercise-relationship
func DeleteExerciseRelationship(ctx context.Context, input *DeleteExerciseRelationshipInput) (*DeleteExerciseRelationshipOutput, error) {
	var entity models.ExerciseRelationshipEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("ExerciseRelationship not found")
	}

	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("not the owner of this relationship")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
