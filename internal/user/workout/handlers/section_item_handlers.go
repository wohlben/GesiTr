package handlers

import (
	"context"

	"gesitr/internal/database"
	exercisemodels "gesitr/internal/exercise/models"
	exercisegroupmodels "gesitr/internal/exercisegroup/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/workout/models"

	"github.com/danielgtaylor/huma/v2"
)

// requireSectionOwner fetches the section, then checks workout ownership.
func requireSectionOwner(ctx context.Context, sectionID uint) error {
	var section models.WorkoutSectionEntity
	if err := database.DB.First(&section, sectionID).Error; err != nil {
		return huma.Error404NotFound("Workout section not found")
	}
	return requireWorkoutOwner(ctx, section.WorkoutID)
}

// ListWorkoutSectionItems returns section items owned by the current
// user. Filter by workoutSectionId query param.
// GET /api/user/workout-section-items
//
// OpenAPI: /api/docs#/operations/ListWorkoutSectionItems
func ListWorkoutSectionItems(ctx context.Context, input *ListWorkoutSectionItemsInput) (*ListWorkoutSectionItemsOutput, error) {
	db := database.DB.Model(&models.WorkoutSectionItemEntity{})

	if input.WorkoutSectionID != "" {
		db = db.Where("workout_section_id = ?", input.WorkoutSectionID)
	}

	// Enforce ownership through workout section -> workout
	db = db.Where("workout_section_id IN (SELECT id FROM workout_sections WHERE workout_id IN (SELECT id FROM workouts WHERE owner = ? AND deleted_at IS NULL))", humaconfig.GetUserID(ctx))

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
// POST /api/user/workout-section-items
//
// OpenAPI: /api/docs#/operations/CreateWorkoutSectionItem
func CreateWorkoutSectionItem(ctx context.Context, input *CreateWorkoutSectionItemInput) (*CreateWorkoutSectionItemOutput, error) {
	if err := requireSectionOwner(ctx, input.Body.WorkoutSectionID); err != nil {
		return nil, err
	}

	switch input.Body.Type {
	case models.WorkoutSectionItemTypeExercise:
		if input.Body.ExerciseSchemeID == nil {
			return nil, huma.Error422UnprocessableEntity("exerciseSchemeId is required for exercise type")
		}
		var scheme exercisemodels.ExerciseSchemeEntity
		if err := database.DB.First(&scheme, *input.Body.ExerciseSchemeID).Error; err != nil {
			return nil, huma.Error404NotFound("Exercise scheme not found")
		}

	case models.WorkoutSectionItemTypeExerciseGroup:
		if input.Body.ExerciseGroupID == nil {
			return nil, huma.Error422UnprocessableEntity("exerciseGroupId is required for exercise_group type")
		}
		var group exercisegroupmodels.ExerciseGroupEntity
		if err := database.DB.First(&group, *input.Body.ExerciseGroupID).Error; err != nil {
			return nil, huma.Error404NotFound("Exercise group not found")
		}

	default:
		return nil, huma.Error422UnprocessableEntity("invalid item type: must be 'exercise' or 'exercise_group'")
	}

	entity := models.WorkoutSectionItemEntity{
		WorkoutSectionID: input.Body.WorkoutSectionID,
		Type:             input.Body.Type,
		ExerciseSchemeID: input.Body.ExerciseSchemeID,
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
// DELETE /api/user/workout-section-items/{id}
//
// OpenAPI: /api/docs#/operations/DeleteWorkoutSectionItem
func DeleteWorkoutSectionItem(ctx context.Context, input *DeleteWorkoutSectionItemInput) (*DeleteWorkoutSectionItemOutput, error) {
	var entity models.WorkoutSectionItemEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout section item not found")
	}
	if err := requireSectionOwner(ctx, entity.WorkoutSectionID); err != nil {
		return nil, err
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
