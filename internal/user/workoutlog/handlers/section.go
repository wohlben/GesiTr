package handlers

import (
	"context"
	"encoding/json"
	"reflect"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	workoutmodels "gesitr/internal/user/workout/models"
	"gesitr/internal/user/workoutlog/models"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

// ListWorkoutLogSections returns workout log sections owned by the current user.
// GET /api/user/workout-log-sections
//
// OpenAPI: /api/docs#/operations/list-workout-log-sections
func ListWorkoutLogSections(ctx context.Context, input *ListWorkoutLogSectionsInput) (*ListWorkoutLogSectionsOutput, error) {
	db := database.DB.Model(&models.WorkoutLogSectionEntity{})

	if input.WorkoutLogID != "" {
		if err := requireLogOwner(ctx, parseUint(input.WorkoutLogID)); err != nil {
			return nil, err
		}
		db = db.Where("workout_log_id = ?", input.WorkoutLogID)
	} else {
		db = db.Joins("JOIN workout_logs ON workout_logs.id = workout_log_sections.workout_log_id").
			Where("workout_logs.owner = ?", humaconfig.GetUserID(ctx))
	}

	var entities []models.WorkoutLogSectionEntity
	if err := db.Preload("Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).Order("position").Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.WorkoutLogSection, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutLogSectionsOutput{Body: dtos}, nil
}

// CreateWorkoutLogSection adds a section to a workout log.
// POST /api/user/workout-log-sections
//
// OpenAPI: /api/docs#/operations/create-workout-log-section
func CreateWorkoutLogSection(ctx context.Context, input *CreateWorkoutLogSectionInput) (*CreateWorkoutLogSectionOutput, error) {
	var dto models.WorkoutLogSection
	if err := json.Unmarshal(input.RawBody, &dto); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	log, err := requireLogStatus(ctx, dto.WorkoutLogID, models.WorkoutLogStatusPlanning, models.WorkoutLogStatusAdhoc)
	if err != nil {
		return nil, err
	}

	entity := models.WorkoutLogSectionFromDTO(dto)
	// For adhoc logs, sections start in_progress immediately
	if log.Status == models.WorkoutLogStatusAdhoc {
		entity.Status = models.WorkoutLogItemStatusInProgress
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateWorkoutLogSectionOutput{Body: entity.ToDTO()}, nil
}

// GetWorkoutLogSection returns a single workout log section.
// GET /api/user/workout-log-sections/{id}
//
// OpenAPI: /api/docs#/operations/get-workout-log-section
func GetWorkoutLogSection(ctx context.Context, input *GetWorkoutLogSectionInput) (*GetWorkoutLogSectionOutput, error) {
	var entity models.WorkoutLogSectionEntity
	if err := database.DB.Preload("Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log section not found")
	}
	if err := requireLogOwner(ctx, entity.WorkoutLogID); err != nil {
		return nil, err
	}
	return &GetWorkoutLogSectionOutput{Body: entity.ToDTO()}, nil
}

// UpdateWorkoutLogSection partially updates a workout log section.
// PATCH /api/user/workout-log-sections/{id}
//
// OpenAPI: /api/docs#/operations/update-workout-log-section
func UpdateWorkoutLogSection(ctx context.Context, input *UpdateWorkoutLogSectionInput) (*UpdateWorkoutLogSectionOutput, error) {
	var existing models.WorkoutLogSectionEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log section not found")
	}

	if err := requireLogOwner(ctx, existing.WorkoutLogID); err != nil {
		return nil, err
	}

	var patch struct {
		Type                 *workoutmodels.WorkoutSectionType `json:"type"`
		Label                *string                           `json:"label"`
		RestBetweenExercises *int                              `json:"restBetweenExercises"`
		Position             *int                              `json:"position"`
	}
	if err := json.Unmarshal(input.RawBody, &patch); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	if reflect.ValueOf(patch).IsZero() {
		return nil, huma.Error400BadRequest("patch body contains no updatable fields")
	}

	if patch.Type != nil {
		existing.Type = *patch.Type
	}
	if patch.Label != nil {
		existing.Label = patch.Label
	}
	if patch.RestBetweenExercises != nil {
		existing.RestBetweenExercises = patch.RestBetweenExercises
	}
	if patch.Position != nil {
		existing.Position = *patch.Position
	}

	if err := database.DB.Save(&existing).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Reload with exercises and sets
	database.DB.Preload("Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).First(&existing, existing.ID)

	return &UpdateWorkoutLogSectionOutput{Body: existing.ToDTO()}, nil
}

// DeleteWorkoutLogSection deletes a workout log section.
// DELETE /api/user/workout-log-sections/{id}
//
// OpenAPI: /api/docs#/operations/delete-workout-log-section
func DeleteWorkoutLogSection(ctx context.Context, input *DeleteWorkoutLogSectionInput) (*DeleteWorkoutLogSectionOutput, error) {
	var existing models.WorkoutLogSectionEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log section not found")
	}

	if _, err := requireLogStatus(ctx, existing.WorkoutLogID, models.WorkoutLogStatusPlanning, models.WorkoutLogStatusAdhoc); err != nil {
		return nil, err
	}

	if err := database.DB.Delete(&existing).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
