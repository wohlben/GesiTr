package handlers

import (
	"context"
	"reflect"

	"gesitr/internal/database"
	exercisemodels "gesitr/internal/exercise/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/workoutlog/models"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

// ListWorkoutLogExercises returns workout log exercises owned by the current user.
// GET /api/user/workout-log-exercises
//
// OpenAPI: /api/docs#/operations/list-workout-log-exercises
func ListWorkoutLogExercises(ctx context.Context, input *ListWorkoutLogExercisesInput) (*ListWorkoutLogExercisesOutput, error) {
	db := database.DB.Model(&models.WorkoutLogExerciseEntity{})

	if input.WorkoutLogSectionID != "" {
		var section models.WorkoutLogSectionEntity
		if err := database.DB.First(&section, input.WorkoutLogSectionID).Error; err != nil {
			return nil, huma.Error404NotFound("Workout log section not found")
		}
		if err := requireLogOwner(ctx, section.WorkoutLogID); err != nil {
			return nil, err
		}
		db = db.Where("workout_log_section_id = ?", input.WorkoutLogSectionID)
	} else {
		db = db.Joins("JOIN workout_logs ON workout_logs.id = workout_log_exercises.workout_log_id").
			Where("workout_logs.owner = ?", humaconfig.GetUserID(ctx))
	}

	var entities []models.WorkoutLogExerciseEntity
	if err := db.Preload("Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).Order("position").Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.WorkoutLogExercise, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutLogExercisesOutput{Body: dtos}, nil
}

// CreateWorkoutLogExercise adds an exercise to a workout log section.
// POST /api/user/workout-log-exercises
//
// OpenAPI: /api/docs#/operations/create-workout-log-exercise
func CreateWorkoutLogExercise(ctx context.Context, input *CreateWorkoutLogExerciseInput) (*CreateWorkoutLogExerciseOutput, error) {
	var section models.WorkoutLogSectionEntity
	if err := database.DB.First(&section, input.Body.WorkoutLogSectionID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log section not found")
	}

	// Guard: parent log must be in planning or adhoc status
	log, err := requireLogStatus(ctx, section.WorkoutLogID, models.WorkoutLogStatusPlanning, models.WorkoutLogStatusAdhoc)
	if err != nil {
		return nil, err
	}

	var scheme exercisemodels.ExerciseSchemeEntity
	if err := database.DB.First(&scheme, input.Body.SourceExerciseSchemeID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise scheme not found")
	}

	breakAfter := section.RestBetweenExercises
	if input.Body.BreakAfterSeconds != nil {
		breakAfter = input.Body.BreakAfterSeconds
	}

	// For adhoc logs, exercises start in_progress immediately
	exerciseStatus := models.WorkoutLogItemStatusPlanning
	if log.Status == models.WorkoutLogStatusAdhoc {
		exerciseStatus = models.WorkoutLogItemStatusInProgress
	}

	entity := models.WorkoutLogExerciseEntity{
		WorkoutLogSectionID:    input.Body.WorkoutLogSectionID,
		WorkoutLogID:           section.WorkoutLogID,
		SourceExerciseSchemeID: input.Body.SourceExerciseSchemeID,
		Position:               input.Body.Position,
		Status:                 exerciseStatus,
		BreakAfterSeconds:      breakAfter,
		TargetMeasurementType:  scheme.MeasurementType,
		TargetTimePerRep:       scheme.TimePerRep,
	}

	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Auto-create set rows from the scheme
	// For adhoc logs, sets start in_progress immediately
	setStatus := models.WorkoutLogItemStatusPlanning
	if log.Status == models.WorkoutLogStatusAdhoc {
		setStatus = models.WorkoutLogItemStatusInProgress
	}

	numSets := 0
	if scheme.Sets != nil {
		numSets = *scheme.Sets
	}
	for i := 1; i <= numSets; i++ {
		set := models.WorkoutLogExerciseSetEntity{
			WorkoutLogExerciseID: entity.ID,
			WorkoutLogID:         section.WorkoutLogID,
			SetNumber:            i,
			Status:               setStatus,
			TargetReps:           scheme.Reps,
			TargetWeight:         scheme.Weight,
			TargetDuration:       scheme.Duration,
			TargetDistance:       scheme.Distance,
			TargetTime:           scheme.TargetTime,
		}
		// BreakAfterSeconds for sets 1..N-1, nil for the last set
		if i < numSets {
			set.BreakAfterSeconds = scheme.RestBetweenSets
		}
		if err := database.DB.Create(&set).Error; err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		entity.Sets = append(entity.Sets, set)
	}

	return &CreateWorkoutLogExerciseOutput{Body: entity.ToDTO()}, nil
}

// UpdateWorkoutLogExercise partially updates a workout log exercise.
// PATCH /api/user/workout-log-exercises/{id}
//
// OpenAPI: /api/docs#/operations/update-workout-log-exercise
func UpdateWorkoutLogExercise(ctx context.Context, input *UpdateWorkoutLogExerciseInput) (*UpdateWorkoutLogExerciseOutput, error) {
	var existing models.WorkoutLogExerciseEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log exercise not found")
	}

	if err := requireLogOwner(ctx, existing.WorkoutLogID); err != nil {
		return nil, err
	}

	patch := input.Body

	if reflect.ValueOf(patch).IsZero() {
		return nil, huma.Error400BadRequest("patch body contains no updatable fields")
	}

	if patch.Position != nil {
		existing.Position = *patch.Position
	}
	if patch.BreakAfterSeconds != nil {
		existing.BreakAfterSeconds = patch.BreakAfterSeconds
	}

	if err := database.DB.Save(&existing).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Reload with sets
	database.DB.Preload("Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).First(&existing, existing.ID)

	return &UpdateWorkoutLogExerciseOutput{Body: existing.ToDTO()}, nil
}

// DeleteWorkoutLogExercise deletes a workout log exercise.
// DELETE /api/user/workout-log-exercises/{id}
//
// OpenAPI: /api/docs#/operations/delete-workout-log-exercise
func DeleteWorkoutLogExercise(ctx context.Context, input *DeleteWorkoutLogExerciseInput) (*DeleteWorkoutLogExerciseOutput, error) {
	var existing models.WorkoutLogExerciseEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log exercise not found")
	}

	// Guard: parent log must be in planning or adhoc status
	if _, err := requireLogStatus(ctx, existing.WorkoutLogID, models.WorkoutLogStatusPlanning, models.WorkoutLogStatusAdhoc); err != nil {
		return nil, err
	}

	sectionID := existing.WorkoutLogSectionID
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&existing).Error; err != nil {
			return err
		}
		return propagateSectionStatus(tx, sectionID)
	})
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
