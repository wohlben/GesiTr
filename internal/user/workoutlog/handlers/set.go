package handlers

import (
	"context"
	"reflect"
	"time"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	exerciseloghandlers "gesitr/internal/user/exerciselog/handlers"
	exerciselogmodels "gesitr/internal/user/exerciselog/models"
	exerciseschememodels "gesitr/internal/user/exercisescheme/models"
	masteryHandlers "gesitr/internal/user/mastery/handlers"
	"gesitr/internal/user/workoutlog/models"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

// ListWorkoutLogExerciseSets returns workout log exercise sets owned by the current user.
// GET /api/user/workout-log-exercise-sets
//
// OpenAPI: /api/docs#/operations/ListWorkoutLogExerciseSets
func ListWorkoutLogExerciseSets(ctx context.Context, input *ListWorkoutLogExerciseSetsInput) (*ListWorkoutLogExerciseSetsOutput, error) {
	db := database.DB.Model(&models.WorkoutLogExerciseSetEntity{})

	if input.WorkoutLogExerciseID != "" {
		var exercise models.WorkoutLogExerciseEntity
		if err := database.DB.First(&exercise, input.WorkoutLogExerciseID).Error; err != nil {
			return nil, huma.Error404NotFound("Workout log exercise not found")
		}
		if err := requireLogOwner(ctx, exercise.WorkoutLogID); err != nil {
			return nil, err
		}
		db = db.Where("workout_log_exercise_id = ?", input.WorkoutLogExerciseID)
	} else {
		db = db.Joins("JOIN workout_logs ON workout_logs.id = workout_log_exercise_sets.workout_log_id").
			Where("workout_logs.owner = ?", humaconfig.GetUserID(ctx))
	}

	var entities []models.WorkoutLogExerciseSetEntity
	if err := db.Preload("ExerciseLog").Order("set_number").Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.WorkoutLogExerciseSet, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutLogExerciseSetsOutput{Body: dtos}, nil
}

// CreateWorkoutLogExerciseSet adds a set to a workout log exercise.
// POST /api/user/workout-log-exercise-sets
//
// OpenAPI: /api/docs#/operations/CreateWorkoutLogExerciseSet
func CreateWorkoutLogExerciseSet(ctx context.Context, input *CreateWorkoutLogExerciseSetInput) (*CreateWorkoutLogExerciseSetOutput, error) {
	var exercise models.WorkoutLogExerciseEntity
	if err := database.DB.First(&exercise, input.Body.WorkoutLogExerciseID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log exercise not found")
	}

	// Guard: parent log must be in planning or adhoc status (also checks ownership)
	log, err := requireLogStatus(ctx, exercise.WorkoutLogID, models.WorkoutLogStatusPlanning, models.WorkoutLogStatusAdhoc, models.WorkoutLogStatusProposed)
	if err != nil {
		return nil, err
	}

	// For adhoc logs, sets start in_progress immediately
	setStatus := models.WorkoutLogItemStatusPlanning
	if log.Status == models.WorkoutLogStatusAdhoc {
		setStatus = models.WorkoutLogItemStatusInProgress
	}

	entity := models.WorkoutLogExerciseSetEntity{
		WorkoutLogExerciseID: input.Body.WorkoutLogExerciseID,
		WorkoutLogID:         exercise.WorkoutLogID,
		SetNumber:            input.Body.SetNumber,
		Status:               setStatus,
		BreakAfterSeconds:    input.Body.BreakAfterSeconds,
		TargetReps:           input.Body.TargetReps,
		TargetWeight:         input.Body.TargetWeight,
		TargetDuration:       input.Body.TargetDuration,
		TargetDistance:       input.Body.TargetDistance,
		TargetTime:           input.Body.TargetTime,
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateWorkoutLogExerciseSetOutput{Body: entity.ToDTO()}, nil
}

// UpdateWorkoutLogExerciseSet partially updates a workout log exercise set.
// PATCH /api/user/workout-log-exercise-sets/{id}
//
// OpenAPI: /api/docs#/operations/UpdateWorkoutLogExerciseSet
func UpdateWorkoutLogExerciseSet(ctx context.Context, input *UpdateWorkoutLogExerciseSetInput) (*UpdateWorkoutLogExerciseSetOutput, error) {
	var existing models.WorkoutLogExerciseSetEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log exercise set not found")
	}

	// Owner check via parent log
	if err := requireLogOwner(ctx, existing.WorkoutLogID); err != nil {
		return nil, err
	}

	patch := input.Body

	if reflect.ValueOf(patch).IsZero() {
		return nil, huma.Error400BadRequest("patch body contains no updatable fields")
	}

	transitionToFinished := false

	// Status transition: validate against the state machine
	if patch.Status != "" && patch.Status != existing.Status {
		if err := existing.Status.TransitionTo(patch.Status); err != nil {
			return nil, huma.Error409Conflict(err.Error())
		}
		now := time.Now()
		existing.Status = patch.Status
		existing.StatusChangedAt = &now
		if patch.Status == models.WorkoutLogItemStatusFinished {
			transitionToFinished = true
		}
	}

	// Update target fields when provided
	if patch.TargetReps != nil {
		existing.TargetReps = patch.TargetReps
	}
	if patch.TargetWeight != nil {
		existing.TargetWeight = patch.TargetWeight
	}
	if patch.TargetDuration != nil {
		existing.TargetDuration = patch.TargetDuration
	}
	if patch.TargetDistance != nil {
		existing.TargetDistance = patch.TargetDistance
	}
	if patch.TargetTime != nil {
		existing.TargetTime = patch.TargetTime
	}

	if patch.BreakAfterSeconds != nil {
		existing.BreakAfterSeconds = patch.BreakAfterSeconds
	}

	hasActuals := patch.ActualReps != nil || patch.ActualWeight != nil ||
		patch.ActualDuration != nil || patch.ActualDistance != nil || patch.ActualTime != nil

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}

		if transitionToFinished || (existing.Status == models.WorkoutLogItemStatusFinished && hasActuals) {
			if err := createOrUpdateExerciseLog(tx, &existing, patch.ActualReps, patch.ActualWeight, patch.ActualDuration, patch.ActualDistance, patch.ActualTime); err != nil {
				return err
			}
		}

		return propagateStatus(tx, existing.WorkoutLogExerciseID)
	})
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Reload with ExerciseLog preloaded
	database.DB.Preload("ExerciseLog").First(&existing, existing.ID)
	return &UpdateWorkoutLogExerciseSetOutput{Body: existing.ToDTO()}, nil
}

func createOrUpdateExerciseLog(db *gorm.DB, set *models.WorkoutLogExerciseSetEntity, reps *int, weight *float64, duration *int, distance *float64, tm *int) error {
	var logExercise models.WorkoutLogExerciseEntity
	if err := db.First(&logExercise, set.WorkoutLogExerciseID).Error; err != nil {
		return err
	}

	var scheme exerciseschememodels.ExerciseSchemeEntity
	if err := db.First(&scheme, logExercise.SourceExerciseSchemeID).Error; err != nil {
		return err
	}

	recordValue, _ := exerciselogmodels.ComputeRecordValue(logExercise.TargetMeasurementType, reps, weight, duration, distance)

	// Check if ExerciseLog already exists for this set
	var existing exerciselogmodels.ExerciseLogEntity
	err := db.Where("workout_log_exercise_set_id = ?", set.ID).First(&existing).Error

	var log models.WorkoutLogEntity
	db.Select("owner").First(&log, set.WorkoutLogID)

	if err == gorm.ErrRecordNotFound {
		schemeID := logExercise.SourceExerciseSchemeID
		exerciseLog := exerciselogmodels.ExerciseLogEntity{
			Owner:                   log.Owner,
			ExerciseID:              scheme.ExerciseID,
			MeasurementType:         logExercise.TargetMeasurementType,
			Reps:                    reps,
			Weight:                  weight,
			Duration:                duration,
			Distance:                distance,
			Time:                    tm,
			RecordValue:             recordValue,
			PerformedAt:             time.Now(),
			WorkoutLogExerciseSetID: &set.ID,
			SourceExerciseSchemeID:  &schemeID,
		}
		if err := db.Create(&exerciseLog).Error; err != nil {
			return err
		}
		_ = masteryHandlers.UpsertExperience(db, log.Owner, scheme.ExerciseID, reps)
		_ = masteryHandlers.UpsertEquipmentExperience(db, log.Owner, scheme.ExerciseID, reps)
		return exerciseloghandlers.RecomputeRecord(db, scheme.ExerciseID, logExercise.TargetMeasurementType)
	} else if err != nil {
		return err
	}

	// Update existing ExerciseLog
	if reps != nil {
		existing.Reps = reps
	}
	if weight != nil {
		existing.Weight = weight
	}
	if duration != nil {
		existing.Duration = duration
	}
	if distance != nil {
		existing.Distance = distance
	}
	if tm != nil {
		existing.Time = tm
	}
	existing.RecordValue = recordValue

	if err := db.Save(&existing).Error; err != nil {
		return err
	}
	return exerciseloghandlers.RecomputeRecord(db, existing.ExerciseID, existing.MeasurementType)
}

func propagateStatus(db *gorm.DB, exerciseID uint) error {
	var exercise models.WorkoutLogExerciseEntity
	if err := db.Preload("Sets").First(&exercise, exerciseID).Error; err != nil {
		return err
	}

	if len(exercise.Sets) == 0 {
		return nil
	}

	allTerminal := true
	anyAborted := false
	anySkipped := false
	anyPartiallyFinished := false
	allFinished := true
	allSkipped := true
	for _, s := range exercise.Sets {
		if !s.Status.IsTerminal() {
			allTerminal = false
			break
		}
		if s.Status == models.WorkoutLogItemStatusAborted {
			anyAborted = true
		}
		if s.Status == models.WorkoutLogItemStatusSkipped {
			anySkipped = true
		}
		if s.Status == models.WorkoutLogItemStatusPartiallyFinished {
			anyPartiallyFinished = true
		}
		if s.Status != models.WorkoutLogItemStatusFinished {
			allFinished = false
		}
		if s.Status != models.WorkoutLogItemStatusSkipped {
			allSkipped = false
		}
	}

	if !allTerminal {
		return nil
	}

	now := time.Now()
	var newStatus models.WorkoutLogItemStatus
	switch {
	case allFinished:
		newStatus = models.WorkoutLogItemStatusFinished
	case allSkipped:
		newStatus = models.WorkoutLogItemStatusSkipped
	case anyAborted:
		newStatus = models.WorkoutLogItemStatusAborted
	case anySkipped || anyPartiallyFinished:
		newStatus = models.WorkoutLogItemStatusPartiallyFinished
	default:
		newStatus = models.WorkoutLogItemStatusFinished
	}

	if exercise.Status != newStatus {
		if err := db.Model(&exercise).Updates(map[string]any{
			"status":            newStatus,
			"status_changed_at": now,
		}).Error; err != nil {
			return err
		}
	}

	// Stop propagation at exercise level for adhoc workouts
	var log models.WorkoutLogEntity
	if err := db.Select("status").First(&log, exercise.WorkoutLogID).Error; err != nil {
		return err
	}
	if log.Status == models.WorkoutLogStatusAdhoc {
		return nil
	}

	return propagateSectionStatus(db, exercise.WorkoutLogSectionID)
}

func propagateSectionStatus(db *gorm.DB, sectionID uint) error {
	var section models.WorkoutLogSectionEntity
	if err := db.Preload("Exercises").First(&section, sectionID).Error; err != nil {
		return err
	}

	if len(section.Exercises) == 0 {
		return nil
	}

	allTerminal := true
	anyAborted := false
	anySkipped := false
	anyPartiallyFinished := false
	allFinished := true
	allSkipped := true
	for _, ex := range section.Exercises {
		if !ex.Status.IsTerminal() {
			allTerminal = false
			break
		}
		if ex.Status == models.WorkoutLogItemStatusAborted {
			anyAborted = true
		}
		if ex.Status == models.WorkoutLogItemStatusSkipped {
			anySkipped = true
		}
		if ex.Status == models.WorkoutLogItemStatusPartiallyFinished {
			anyPartiallyFinished = true
		}
		if ex.Status != models.WorkoutLogItemStatusFinished {
			allFinished = false
		}
		if ex.Status != models.WorkoutLogItemStatusSkipped {
			allSkipped = false
		}
	}

	if !allTerminal {
		return nil
	}

	now := time.Now()
	var newSectionStatus models.WorkoutLogItemStatus
	switch {
	case allFinished:
		newSectionStatus = models.WorkoutLogItemStatusFinished
	case allSkipped:
		newSectionStatus = models.WorkoutLogItemStatusSkipped
	case anyAborted:
		newSectionStatus = models.WorkoutLogItemStatusAborted
	case anySkipped || anyPartiallyFinished:
		newSectionStatus = models.WorkoutLogItemStatusPartiallyFinished
	default:
		newSectionStatus = models.WorkoutLogItemStatusFinished
	}

	if section.Status != newSectionStatus {
		if err := db.Model(&section).Updates(map[string]any{
			"status":            newSectionStatus,
			"status_changed_at": now,
		}).Error; err != nil {
			return err
		}
	}

	// Propagate to log
	var log models.WorkoutLogEntity
	if err := db.Preload("Sections").First(&log, section.WorkoutLogID).Error; err != nil {
		return err
	}

	if len(log.Sections) == 0 {
		return nil
	}

	allTerminal = true
	anyAborted = false
	allFinished = true
	for _, s := range log.Sections {
		if !s.Status.IsTerminal() {
			allTerminal = false
			break
		}
		if s.Status == models.WorkoutLogItemStatusAborted {
			anyAborted = true
		}
		if s.Status != models.WorkoutLogItemStatusFinished {
			allFinished = false
		}
	}

	if !allTerminal {
		return nil
	}

	var newLogStatus models.WorkoutLogStatus
	switch {
	case allFinished:
		newLogStatus = models.WorkoutLogStatusFinished
	case anyAborted:
		newLogStatus = models.WorkoutLogStatusAborted
	default:
		newLogStatus = models.WorkoutLogStatusPartiallyFinished
	}

	if log.Status != newLogStatus {
		if err := db.Model(&log).Updates(map[string]any{
			"status":            newLogStatus,
			"status_changed_at": now,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

// DeleteWorkoutLogExerciseSet deletes a workout log exercise set.
// DELETE /api/user/workout-log-exercise-sets/{id}
//
// OpenAPI: /api/docs#/operations/DeleteWorkoutLogExerciseSet
func DeleteWorkoutLogExerciseSet(ctx context.Context, input *DeleteWorkoutLogExerciseSetInput) (*DeleteWorkoutLogExerciseSetOutput, error) {
	var existing models.WorkoutLogExerciseSetEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log exercise set not found")
	}

	// Guard: parent log must be in planning or adhoc status (also checks ownership)
	if _, err := requireLogStatus(ctx, existing.WorkoutLogID, models.WorkoutLogStatusPlanning, models.WorkoutLogStatusAdhoc, models.WorkoutLogStatusProposed); err != nil {
		return nil, err
	}

	exerciseID := existing.WorkoutLogExerciseID
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&existing).Error; err != nil {
			return err
		}
		return propagateStatus(tx, exerciseID)
	})
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
