package handlers

import (
	"context"

	"gesitr/internal/compendium/exercise/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/compendium/ownershipgroup"
	masteryHandlers "gesitr/internal/user/mastery/handlers"

	"github.com/danielgtaylor/huma/v2"
)

// ListExerciseRelationships returns exercise relationships, optionally filtered
// by owner, fromExerciseId, toExerciseId, or relationshipType.
// GET /api/exercise-relationships
//
// OpenAPI: /api/docs#/operations/ListExerciseRelationships
func ListExerciseRelationships(ctx context.Context, input *ListExerciseRelationshipsInput) (*ListExerciseRelationshipsOutput, error) {
	db := database.DB.Model(&models.ExerciseRelationshipEntity{})

	if input.Owner != "" {
		userID := humaconfig.GetUserID(ctx)
		if input.Owner == "me" || input.Owner == userID {
			visibleGroups := ownershipgroup.VisibleGroupIDs(database.DB, userID)
			db = db.Where("ownership_group_id IN (?)", visibleGroups)
		} else {
			db = db.Where("ownership_group_id IN (SELECT group_id FROM ownership_group_memberships WHERE user_id = ? AND role = 'owner' AND deleted_at IS NULL)", input.Owner)
		}
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
// OpenAPI: /api/docs#/operations/CreateExerciseRelationship
func CreateExerciseRelationship(ctx context.Context, input *CreateExerciseRelationshipInput) (*CreateExerciseRelationshipOutput, error) {
	dto := models.ExerciseRelationship{
		RelationshipType: input.Body.RelationshipType,
		Strength:         input.Body.Strength,
		Description:      input.Body.Description,
		FromExerciseID:   input.Body.FromExerciseID,
		ToExerciseID:     input.Body.ToExerciseID,
	}

	entity := models.ExerciseRelationshipFromDTO(dto)
	userID := humaconfig.GetUserID(ctx)

	// Inherit ownership group from the "from" exercise.
	var fromExercise models.ExerciseEntity
	if err := database.DB.Select("ownership_group_id").First(&fromExercise, input.Body.FromExerciseID).Error; err != nil {
		return nil, huma.Error404NotFound("From exercise not found")
	}
	entity.OwnershipGroupID = fromExercise.OwnershipGroupID

	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	_ = masteryHandlers.RecalculateContributions(database.DB, userID, entity.FromExerciseID, entity.ToExerciseID)
	return &CreateExerciseRelationshipOutput{Body: entity.ToDTO()}, nil
}

// DeleteExerciseRelationship deletes an exercise relationship. Owner only.
// DELETE /api/exercise-relationships/:id
//
// OpenAPI: /api/docs#/operations/DeleteExerciseRelationship
func DeleteExerciseRelationship(ctx context.Context, input *DeleteExerciseRelationshipInput) (*DeleteExerciseRelationshipOutput, error) {
	var entity models.ExerciseRelationshipEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("ExerciseRelationship not found")
	}

	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, entity.OwnershipGroupID)
	if !access.CanModify() {
		return nil, huma.Error403Forbidden("not the owner of this relationship")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	_ = masteryHandlers.RecalculateContributions(database.DB, userID, entity.FromExerciseID, entity.ToExerciseID)
	return nil, nil
}
