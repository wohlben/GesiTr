package handlers

import (
	"context"
	"encoding/json"

	"gesitr/internal/database"
	exercisemodels "gesitr/internal/exercise/models"
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

// ListWorkoutSectionExercises returns section exercises owned by the current
// user. Filter by workoutSectionId query param.
// GET /api/user/workout-section-exercises
//
// OpenAPI: /api/docs#/operations/list-workout-section-exercises
func ListWorkoutSectionExercises(ctx context.Context, input *ListWorkoutSectionExercisesInput) (*ListWorkoutSectionExercisesOutput, error) {
	db := database.DB.Model(&models.WorkoutSectionExerciseEntity{})

	if input.WorkoutSectionID != "" {
		db = db.Where("workout_section_id = ?", input.WorkoutSectionID)
	}

	// Enforce ownership through workout section -> workout
	db = db.Where("workout_section_id IN (SELECT id FROM workout_sections WHERE workout_id IN (SELECT id FROM workouts WHERE owner = ? AND deleted_at IS NULL))", humaconfig.GetUserID(ctx))

	var entities []models.WorkoutSectionExerciseEntity
	if err := db.Order("position").Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.WorkoutSectionExercise, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutSectionExercisesOutput{Body: dtos}, nil
}

// CreateWorkoutSectionExercise adds an exercise to a section. Requires a
// workoutSectionId (whose parent workout must be owned by the current user)
// and an exerciseSchemeId referencing an existing exercise scheme — see
// [gesitr/internal/exercise/handlers.CreateExerciseScheme].
// A workout and section must exist first — see [CreateWorkout] and
// [CreateWorkoutSection]. POST /api/user/workout-section-exercises
//
// OpenAPI: /api/docs#/operations/create-workout-section-exercise
func CreateWorkoutSectionExercise(ctx context.Context, input *CreateWorkoutSectionExerciseInput) (*CreateWorkoutSectionExerciseOutput, error) {
	var dto models.WorkoutSectionExercise
	if err := json.Unmarshal(input.RawBody, &dto); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	if err := requireSectionOwner(ctx, dto.WorkoutSectionID); err != nil {
		return nil, err
	}

	var scheme exercisemodels.ExerciseSchemeEntity
	if err := database.DB.First(&scheme, dto.ExerciseSchemeID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise scheme not found")
	}

	entity := models.WorkoutSectionExerciseFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateWorkoutSectionExerciseOutput{Body: entity.ToDTO()}, nil
}

// DeleteWorkoutSectionExercise removes an exercise from a section.
// DELETE /api/user/workout-section-exercises/{id}
//
// OpenAPI: /api/docs#/operations/delete-workout-section-exercise
func DeleteWorkoutSectionExercise(ctx context.Context, input *DeleteWorkoutSectionExerciseInput) (*DeleteWorkoutSectionExerciseOutput, error) {
	var entity models.WorkoutSectionExerciseEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout section exercise not found")
	}
	if err := requireSectionOwner(ctx, entity.WorkoutSectionID); err != nil {
		return nil, err
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
