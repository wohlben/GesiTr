package handlers

import (
	"context"

	"gesitr/internal/compendium/workout/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/compendium/ownershipgroup"

	"github.com/danielgtaylor/huma/v2"
)

// ListExerciseGroupMembers returns exercise group members, optionally filtered
// by groupId or exerciseId.
// GET /api/exercise-group-members
//
// OpenAPI: /api/docs#/operations/ListExerciseGroupMembers
func ListExerciseGroupMembers(ctx context.Context, input *ListExerciseGroupMembersInput) (*ListExerciseGroupMembersOutput, error) {
	db := database.DB.Model(&models.ExerciseGroupMemberEntity{})

	if input.GroupID != "" {
		db = db.Where("group_id = ?", input.GroupID)
	}
	if input.ExerciseID != "" {
		db = db.Where("exercise_id = ?", input.ExerciseID)
	}

	var entities []models.ExerciseGroupMemberEntity
	if err := db.Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.ExerciseGroupMember, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListExerciseGroupMembersOutput{Body: dtos}, nil
}

// CreateExerciseGroupMember creates an exercise group member owned by the current user.
// POST /api/exercise-group-members
//
// OpenAPI: /api/docs#/operations/CreateExerciseGroupMember
func CreateExerciseGroupMember(ctx context.Context, input *CreateExerciseGroupMemberInput) (*CreateExerciseGroupMemberOutput, error) {
	// Look up the parent exercise group to inherit its ownership_group_id.
	var group models.ExerciseGroupEntity
	if err := database.DB.First(&group, input.Body.GroupID).Error; err != nil {
		return nil, huma.Error404NotFound("ExerciseGroup not found")
	}

	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, group.OwnershipGroupID)
	if !access.CanModify() {
		return nil, huma.Error403Forbidden("not the owner of this exercise group")
	}

	dto := models.ExerciseGroupMember{
		GroupID:    input.Body.GroupID,
		ExerciseID: input.Body.ExerciseID,
	}

	entity := models.ExerciseGroupMemberFromDTO(dto)
	entity.OwnershipGroupID = group.OwnershipGroupID
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateExerciseGroupMemberOutput{Body: entity.ToDTO()}, nil
}

// DeleteExerciseGroupMember deletes an exercise group member. Owner only.
// DELETE /api/exercise-group-members/:id
//
// OpenAPI: /api/docs#/operations/DeleteExerciseGroupMember
func DeleteExerciseGroupMember(ctx context.Context, input *DeleteExerciseGroupMemberInput) (*DeleteExerciseGroupMemberOutput, error) {
	var entity models.ExerciseGroupMemberEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("ExerciseGroupMember not found")
	}

	access := ownershipgroup.CheckAccess(database.DB, humaconfig.GetUserID(ctx), entity.OwnershipGroupID)
	if !access.CanDelete() {
		return nil, huma.Error403Forbidden("not the owner of this group member")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
