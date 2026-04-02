package handlers

import (
	"context"

	"gesitr/internal/compendium/workout/models"
	"gesitr/internal/compendium/workoutgroup"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"

	"github.com/danielgtaylor/huma/v2"
)

// requireWorkoutRead fetches the workout and checks that the caller has read access
// (owner or group member).
func requireWorkoutRead(ctx context.Context, workoutID uint) error {
	var workout models.WorkoutEntity
	if err := database.DB.First(&workout, workoutID).Error; err != nil {
		return huma.Error404NotFound("Workout not found")
	}
	if !canReadWorkout(humaconfig.GetUserID(ctx), &workout) {
		return huma.Error403Forbidden("access denied")
	}
	return nil
}

// requireWorkoutModify fetches the workout and checks that the caller has modify access
// (owner or group admin).
func requireWorkoutModify(ctx context.Context, workoutID uint) error {
	var workout models.WorkoutEntity
	if err := database.DB.First(&workout, workoutID).Error; err != nil {
		return huma.Error404NotFound("Workout not found")
	}
	access := workoutgroup.CheckWorkoutAccess(humaconfig.GetUserID(ctx), workout.Owner, workoutID)
	if !access.CanModify() {
		return huma.Error403Forbidden("access denied")
	}
	return nil
}

// ListWorkoutSections returns sections owned by the current user. Filter by
// workoutId query param to get sections for a specific workout.
// GET /api/workout-sections
//
// OpenAPI: /api/docs#/operations/ListWorkoutSections
func ListWorkoutSections(ctx context.Context, input *ListWorkoutSectionsInput) (*ListWorkoutSectionsOutput, error) {
	db := database.DB.Model(&models.WorkoutSectionEntity{})

	if input.WorkoutID != "" {
		db = db.Where("workout_id = ?", input.WorkoutID)
	}

	// Include workouts the user owns, public workouts, or group membership workouts
	userID := humaconfig.GetUserID(ctx)
	db = db.Where(`workout_id IN (SELECT id FROM workouts WHERE (owner = ? OR public = ?) AND deleted_at IS NULL)
		OR workout_id IN (SELECT wg.workout_id FROM workout_groups wg JOIN workout_group_memberships wgm ON wgm.group_id = wg.id WHERE wgm.user_id = ? AND wgm.deleted_at IS NULL AND wg.deleted_at IS NULL)`,
		userID, true, userID)

	var entities []models.WorkoutSectionEntity
	if err := db.Preload("Items").Order("position").Find(&entities).Error; err != nil {
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
// POST /api/workout-sections
//
// OpenAPI: /api/docs#/operations/CreateWorkoutSection
func CreateWorkoutSection(ctx context.Context, input *CreateWorkoutSectionInput) (*CreateWorkoutSectionOutput, error) {
	if err := requireWorkoutModify(ctx, input.Body.WorkoutID); err != nil {
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
// GET /api/workout-sections/{id}
//
// OpenAPI: /api/docs#/operations/GetWorkoutSection
func GetWorkoutSection(ctx context.Context, input *GetWorkoutSectionInput) (*GetWorkoutSectionOutput, error) {
	var entity models.WorkoutSectionEntity
	if err := database.DB.Preload("Items").First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout section not found")
	}
	if err := requireWorkoutRead(ctx, entity.WorkoutID); err != nil {
		return nil, err
	}
	return &GetWorkoutSectionOutput{Body: entity.ToDTO()}, nil
}

// DeleteWorkoutSection removes a section from its workout.
// DELETE /api/workout-sections/{id}
//
// OpenAPI: /api/docs#/operations/DeleteWorkoutSection
func DeleteWorkoutSection(ctx context.Context, input *DeleteWorkoutSectionInput) (*DeleteWorkoutSectionOutput, error) {
	var entity models.WorkoutSectionEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout section not found")
	}
	if err := requireWorkoutModify(ctx, entity.WorkoutID); err != nil {
		return nil, err
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
