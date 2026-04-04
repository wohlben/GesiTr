package handlers

import (
	"context"

	exercisemodels "gesitr/internal/compendium/exercise/models"
	"gesitr/internal/compendium/workout/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/ownershipgroup"

	"github.com/danielgtaylor/huma/v2"
)

// requireSectionRead fetches the section, then checks workout read access.
func requireSectionRead(ctx context.Context, sectionID uint) error {
	var section models.WorkoutSectionEntity
	if err := database.DB.First(&section, sectionID).Error; err != nil {
		return huma.Error404NotFound("Workout section not found")
	}
	return requireWorkoutRead(ctx, section.WorkoutID)
}

// requireSectionModify fetches the section, then checks workout modify access.
func requireSectionModify(ctx context.Context, sectionID uint) error {
	var section models.WorkoutSectionEntity
	if err := database.DB.First(&section, sectionID).Error; err != nil {
		return huma.Error404NotFound("Workout section not found")
	}
	return requireWorkoutModify(ctx, section.WorkoutID)
}

// ListWorkoutSectionItems returns section items owned by the current
// user. Filter by workoutSectionId query param.
// GET /api/workout-section-items
//
// OpenAPI: /api/docs#/operations/ListWorkoutSectionItems
func ListWorkoutSectionItems(ctx context.Context, input *ListWorkoutSectionItemsInput) (*ListWorkoutSectionItemsOutput, error) {
	db := database.DB.Model(&models.WorkoutSectionItemEntity{})

	if input.WorkoutSectionID != "" {
		db = db.Where("workout_section_id = ?", input.WorkoutSectionID)
	}

	// Include items for workouts the user can access via ownership group, public workouts, or workout group membership
	userID := humaconfig.GetUserID(ctx)
	visibleGroups := ownershipgroup.VisibleGroupIDs(database.DB, userID)
	db = db.Where(`workout_section_id IN (SELECT id FROM workout_sections WHERE
		workout_id IN (SELECT id FROM workouts WHERE (ownership_group_id IN (?) OR public = ?) AND deleted_at IS NULL)
		OR workout_id IN (SELECT wg.workout_id FROM workout_groups wg JOIN workout_group_memberships wgm ON wgm.group_id = wg.id WHERE wgm.user_id = ? AND wgm.deleted_at IS NULL AND wg.deleted_at IS NULL))`,
		visibleGroups, true, userID)

	var entities []models.WorkoutSectionItemEntity
	if err := db.Order("position").Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.WorkoutSectionItem, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutSectionItemsOutput{Body: dtos}, nil
}

// CreateWorkoutSectionItem adds an item to a section. The item type determines
// which reference field is required:
//   - "exercise": exerciseSchemeId is required
//   - "exercise_group": exerciseGroupId is required
//
// POST /api/workout-section-items
//
// OpenAPI: /api/docs#/operations/CreateWorkoutSectionItem
func CreateWorkoutSectionItem(ctx context.Context, input *CreateWorkoutSectionItemInput) (*CreateWorkoutSectionItemOutput, error) {
	if err := requireSectionModify(ctx, input.Body.WorkoutSectionID); err != nil {
		return nil, err
	}

	switch input.Body.Type {
	case models.WorkoutSectionItemTypeExercise:
		if input.Body.ExerciseID == nil {
			return nil, huma.Error422UnprocessableEntity("exerciseId is required for exercise type")
		}
		var exercise exercisemodels.ExerciseEntity
		if err := database.DB.First(&exercise, *input.Body.ExerciseID).Error; err != nil {
			return nil, huma.Error404NotFound("Exercise not found")
		}

	case models.WorkoutSectionItemTypeExerciseGroup:
		if input.Body.ExerciseGroupID == nil {
			return nil, huma.Error422UnprocessableEntity("exerciseGroupId is required for exercise_group type")
		}
		var group models.ExerciseGroupEntity
		if err := database.DB.First(&group, *input.Body.ExerciseGroupID).Error; err != nil {
			return nil, huma.Error404NotFound("Exercise group not found")
		}

	default:
		return nil, huma.Error422UnprocessableEntity("invalid item type: must be 'exercise' or 'exercise_group'")
	}

	entity := models.WorkoutSectionItemEntity{
		WorkoutSectionID: input.Body.WorkoutSectionID,
		Type:             input.Body.Type,
		ExerciseID:       input.Body.ExerciseID,
		ExerciseGroupID:  input.Body.ExerciseGroupID,
		Data:             input.Body.Data,
		Position:         input.Body.Position,
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateWorkoutSectionItemOutput{Body: entity.ToDTO()}, nil
}

// DeleteWorkoutSectionItem removes an item from a section.
// DELETE /api/workout-section-items/{id}
//
// OpenAPI: /api/docs#/operations/DeleteWorkoutSectionItem
func DeleteWorkoutSectionItem(ctx context.Context, input *DeleteWorkoutSectionItemInput) (*DeleteWorkoutSectionItemOutput, error) {
	var entity models.WorkoutSectionItemEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout section item not found")
	}
	if err := requireSectionModify(ctx, entity.WorkoutSectionID); err != nil {
		return nil, err
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
