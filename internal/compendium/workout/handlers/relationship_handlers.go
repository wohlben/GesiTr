package handlers

import (
	"context"

	"gesitr/internal/compendium/workout/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/ownershipgroup"

	"github.com/danielgtaylor/huma/v2"
)

// ListWorkoutRelationships returns workout relationships, optionally filtered.
func ListWorkoutRelationships(ctx context.Context, input *ListWorkoutRelationshipsInput) (*ListWorkoutRelationshipsOutput, error) {
	db := database.DB.Model(&models.WorkoutRelationshipEntity{})

	if input.Owner != "" {
		if input.Owner == "me" {
			userID := humaconfig.GetUserID(ctx)
			visibleGroups := ownershipgroup.VisibleGroupIDs(database.DB, userID)
			db = db.Where("ownership_group_id IN (?)", visibleGroups)
		} else {
			db = db.Where("ownership_group_id IN (SELECT group_id FROM ownership_group_memberships WHERE user_id = ? AND role = 'owner' AND deleted_at IS NULL)", input.Owner)
		}
	}
	if input.FromWorkoutID != "" {
		db = db.Where("from_workout_id = ?", input.FromWorkoutID)
	}
	if input.ToWorkoutID != "" {
		db = db.Where("to_workout_id = ?", input.ToWorkoutID)
	}
	if input.RelationshipType != "" {
		db = db.Where("relationship_type = ?", input.RelationshipType)
	}

	var entities []models.WorkoutRelationshipEntity
	if err := db.Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.WorkoutRelationship, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutRelationshipsOutput{Body: dtos}, nil
}

// CreateWorkoutRelationship creates a workout relationship owned by the current user.
func CreateWorkoutRelationship(ctx context.Context, input *CreateWorkoutRelationshipInput) (*CreateWorkoutRelationshipOutput, error) {
	// Look up the from-workout to inherit its ownership_group_id.
	var fromWorkout models.WorkoutEntity
	if err := database.DB.First(&fromWorkout, input.Body.FromWorkoutID).Error; err != nil {
		return nil, huma.Error404NotFound("From workout not found")
	}

	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, fromWorkout.OwnershipGroupID)
	if !access.CanModify() {
		return nil, huma.Error403Forbidden("access denied")
	}

	dto := models.WorkoutRelationship{
		RelationshipType: input.Body.RelationshipType,
		Strength:         input.Body.Strength,
		FromWorkoutID:    input.Body.FromWorkoutID,
		ToWorkoutID:      input.Body.ToWorkoutID,
	}

	entity := models.WorkoutRelationshipFromDTO(dto)
	entity.OwnershipGroupID = fromWorkout.OwnershipGroupID
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateWorkoutRelationshipOutput{Body: entity.ToDTO()}, nil
}

// DeleteWorkoutRelationship deletes a workout relationship. Owner only.
func DeleteWorkoutRelationship(ctx context.Context, input *DeleteWorkoutRelationshipInput) (*DeleteWorkoutRelationshipOutput, error) {
	var entity models.WorkoutRelationshipEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("WorkoutRelationship not found")
	}

	access := ownershipgroup.CheckAccess(database.DB, humaconfig.GetUserID(ctx), entity.OwnershipGroupID)
	if !access.CanDelete() {
		return nil, huma.Error403Forbidden("not the owner of this relationship")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
