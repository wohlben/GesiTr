package handlers

import (
	"context"

	workoutmodels "gesitr/internal/compendium/workout/models"
	"gesitr/internal/compendium/workoutgroup/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"

	"github.com/danielgtaylor/huma/v2"
)

// ListWorkoutGroups returns workout groups owned by the current user.
// GET /api/user/workout-groups
func ListWorkoutGroups(ctx context.Context, input *ListWorkoutGroupsInput) (*ListWorkoutGroupsOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	var entities []models.WorkoutGroupEntity
	if err := database.DB.Where("owner = ?", userID).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.WorkoutGroup, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutGroupsOutput{Body: dtos}, nil
}

// CreateWorkoutGroup creates a workout group. The caller must own the referenced workout.
// POST /api/user/workout-groups
func CreateWorkoutGroup(ctx context.Context, input *CreateWorkoutGroupInput) (*CreateWorkoutGroupOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	var workout workoutmodels.WorkoutEntity
	if err := database.DB.First(&workout, input.Body.WorkoutID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	if workout.Owner != userID {
		return nil, huma.Error403Forbidden("access denied")
	}

	entity := models.WorkoutGroupEntity{
		Name:      input.Body.Name,
		WorkoutID: input.Body.WorkoutID,
		Owner:     userID,
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error422UnprocessableEntity("a group already exists for this workout")
	}
	return &CreateWorkoutGroupOutput{Body: entity.ToDTO()}, nil
}

// GetWorkoutGroup returns a workout group. Visible if the caller is the owner or a member.
// GET /api/user/workout-groups/{id}
func GetWorkoutGroup(ctx context.Context, input *GetWorkoutGroupInput) (*GetWorkoutGroupOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	var entity models.WorkoutGroupEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout group not found")
	}

	if entity.Owner != userID {
		var membership models.WorkoutGroupMembershipEntity
		if err := database.DB.Where("group_id = ? AND user_id = ?", entity.ID, userID).First(&membership).Error; err != nil {
			return nil, huma.Error403Forbidden("access denied")
		}
	}

	return &GetWorkoutGroupOutput{Body: entity.ToDTO()}, nil
}

// UpdateWorkoutGroup updates a workout group name. Owner only.
// PUT /api/user/workout-groups/{id}
func UpdateWorkoutGroup(ctx context.Context, input *UpdateWorkoutGroupInput) (*UpdateWorkoutGroupOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	var existing models.WorkoutGroupEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout group not found")
	}
	if existing.Owner != userID {
		return nil, huma.Error403Forbidden("access denied")
	}

	existing.Name = input.Body.Name
	if err := database.DB.Save(&existing).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateWorkoutGroupOutput{Body: existing.ToDTO()}, nil
}

// DeleteWorkoutGroup deletes a workout group. Owner only. CASCADE deletes memberships.
// DELETE /api/user/workout-groups/{id}
func DeleteWorkoutGroup(ctx context.Context, input *DeleteWorkoutGroupInput) (*DeleteWorkoutGroupOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	var entity models.WorkoutGroupEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout group not found")
	}
	if entity.Owner != userID {
		return nil, huma.Error403Forbidden("access denied")
	}

	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
