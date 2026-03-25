package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/workout/models"

	"github.com/danielgtaylor/huma/v2"
)

// requireWorkoutOwner fetches the workout by ID and checks ownership.
// Returns an error if the workout is not found or the caller is not the owner.
func requireWorkoutOwner(ctx context.Context, workoutID uint) error {
	var workout models.WorkoutEntity
	if err := database.DB.First(&workout, workoutID).Error; err != nil {
		return huma.Error404NotFound("Workout not found")
	}
	if workout.Owner != humaconfig.GetUserID(ctx) {
		return huma.Error403Forbidden("access denied")
	}
	return nil
}

// ListWorkoutSections returns sections owned by the current user. Filter by
// workoutId query param to get sections for a specific workout.
// GET /api/user/workout-sections
//
// OpenAPI: /api/docs#/operations/list-workout-sections
func ListWorkoutSections(ctx context.Context, input *ListWorkoutSectionsInput) (*ListWorkoutSectionsOutput, error) {
	db := database.DB.Model(&models.WorkoutSectionEntity{})

	if input.WorkoutID != "" {
		db = db.Where("workout_id = ?", input.WorkoutID)
	}

	// Join through workout to enforce ownership
	db = db.Where("workout_id IN (SELECT id FROM workouts WHERE owner = ? AND deleted_at IS NULL)", humaconfig.GetUserID(ctx))

	var entities []models.WorkoutSectionEntity
	if err := db.Preload("Exercises").Order("position").Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.WorkoutSection, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutSectionsOutput{Body: dtos}, nil
}

// CreateWorkoutSection adds a section to a workout. Requires a workoutId
// referencing a workout owned by the current user. A workout must exist
// before sections can be added — see [CreateWorkout].
// POST /api/user/workout-sections
//
// OpenAPI: /api/docs#/operations/create-workout-section
func CreateWorkoutSection(ctx context.Context, input *CreateWorkoutSectionInput) (*CreateWorkoutSectionOutput, error) {
	if err := requireWorkoutOwner(ctx, input.Body.WorkoutID); err != nil {
		return nil, err
	}

	entity := models.WorkoutSectionEntity{
		WorkoutID:            input.Body.WorkoutID,
		Type:                 input.Body.Type,
		Label:                input.Body.Label,
		Position:             input.Body.Position,
		RestBetweenExercises: input.Body.RestBetweenExercises,
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateWorkoutSectionOutput{Body: entity.ToDTO()}, nil
}

// GetWorkoutSection returns a single section with its exercises.
// GET /api/user/workout-sections/{id}
//
// OpenAPI: /api/docs#/operations/get-workout-section
func GetWorkoutSection(ctx context.Context, input *GetWorkoutSectionInput) (*GetWorkoutSectionOutput, error) {
	var entity models.WorkoutSectionEntity
	if err := database.DB.Preload("Exercises").First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout section not found")
	}
	if err := requireWorkoutOwner(ctx, entity.WorkoutID); err != nil {
		return nil, err
	}
	return &GetWorkoutSectionOutput{Body: entity.ToDTO()}, nil
}

// DeleteWorkoutSection removes a section from its workout.
// DELETE /api/user/workout-sections/{id}
//
// OpenAPI: /api/docs#/operations/delete-workout-section
func DeleteWorkoutSection(ctx context.Context, input *DeleteWorkoutSectionInput) (*DeleteWorkoutSectionOutput, error) {
	var entity models.WorkoutSectionEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout section not found")
	}
	if err := requireWorkoutOwner(ctx, entity.WorkoutID); err != nil {
		return nil, err
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
