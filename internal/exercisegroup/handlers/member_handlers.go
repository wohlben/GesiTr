package handlers

import (
	"context"
	"encoding/json"

	"gesitr/internal/database"
	"gesitr/internal/exercisegroup/models"
	"gesitr/internal/humaconfig"

	"github.com/danielgtaylor/huma/v2"
)

// ListExerciseGroupMembers returns exercise group members, optionally filtered
// by groupId or exerciseId.
// GET /api/exercise-group-members
//
// OpenAPI: /api/docs#/operations/list-exercise-group-members
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
// OpenAPI: /api/docs#/operations/create-exercise-group-member
func CreateExerciseGroupMember(ctx context.Context, input *CreateExerciseGroupMemberInput) (*CreateExerciseGroupMemberOutput, error) {
	var dto models.ExerciseGroupMember
	if err := json.Unmarshal(input.RawBody, &dto); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	entity := models.ExerciseGroupMemberFromDTO(dto)
	entity.Owner = humaconfig.GetUserID(ctx)
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateExerciseGroupMemberOutput{Body: entity.ToDTO()}, nil
}

// DeleteExerciseGroupMember deletes an exercise group member. Owner only.
// DELETE /api/exercise-group-members/:id
//
// OpenAPI: /api/docs#/operations/delete-exercise-group-member
func DeleteExerciseGroupMember(ctx context.Context, input *DeleteExerciseGroupMemberInput) (*DeleteExerciseGroupMemberOutput, error) {
	var entity models.ExerciseGroupMemberEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("ExerciseGroupMember not found")
	}

	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("not the owner of this group member")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
