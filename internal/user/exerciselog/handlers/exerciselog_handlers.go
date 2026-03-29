package handlers

import (
	"context"
	"time"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/exerciselog/models"
	masteryHandlers "gesitr/internal/user/mastery/handlers"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

// ListExerciseLogs returns exercise logs owned by the current user, ordered
// by performedAt DESC. Supports filtering by exerciseId, measurementType,
// isRecord, from, and to. GET /api/user/exercise-logs
//
// OpenAPI: /api/docs#/operations/ListExerciseLogs
func ListExerciseLogs(ctx context.Context, input *ListExerciseLogsInput) (*ListExerciseLogsOutput, error) {
	db := database.DB.Model(&models.ExerciseLogEntity{}).
		Where("owner = ?", humaconfig.GetUserID(ctx))

	if input.ExerciseID != "" {
		db = db.Where("exercise_id = ?", input.ExerciseID)
	}
	if input.MeasurementType != "" {
		db = db.Where("measurement_type = ?", input.MeasurementType)
	}
	if input.IsRecord == "true" {
		db = db.Where("is_record = ?", true)
	}
	if input.From != "" {
		db = db.Where("performed_at >= ?", input.From)
	}
	if input.To != "" {
		db = db.Where("performed_at <= ?", input.To)
	}

	var entities []models.ExerciseLogEntity
	if err := db.Order("performed_at DESC").Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.ExerciseLog, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListExerciseLogsOutput{Body: dtos}, nil
}

// CreateExerciseLog creates an exercise log entry. Sets owner, computes
// RecordValue, and recomputes which entry is the record for the
// (exerciseId, measurementType) pair. POST /api/user/exercise-logs
//
// OpenAPI: /api/docs#/operations/CreateExerciseLog
func CreateExerciseLog(ctx context.Context, input *CreateExerciseLogInput) (*CreateExerciseLogOutput, error) {
	entity := models.ExerciseLogEntity{
		ExerciseID:              input.Body.ExerciseID,
		MeasurementType:         input.Body.MeasurementType,
		Reps:                    input.Body.Reps,
		Weight:                  input.Body.Weight,
		Duration:                input.Body.Duration,
		Distance:                input.Body.Distance,
		Time:                    input.Body.Time,
		PerformedAt:             input.Body.PerformedAt,
		WorkoutLogExerciseSetID: input.Body.WorkoutLogExerciseSetID,
		SourceExerciseSchemeID:  input.Body.SourceExerciseSchemeID,
	}
	entity.Owner = humaconfig.GetUserID(ctx)

	if entity.PerformedAt.IsZero() {
		entity.PerformedAt = time.Now()
	}

	value, ok := models.ComputeRecordValue(entity.MeasurementType, entity.Reps, entity.Weight, entity.Duration, entity.Distance)
	if ok {
		entity.RecordValue = value
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&entity).Error; err != nil {
			return err
		}
		_ = masteryHandlers.UpsertExperience(tx, entity.Owner, entity.ExerciseID, entity.Reps)
		return RecomputeRecord(tx, entity.ExerciseID, entity.MeasurementType)
	})
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Reload to get recomputed isRecord
	database.DB.First(&entity, entity.ID)
	return &CreateExerciseLogOutput{Body: entity.ToDTO()}, nil
}

// GetExerciseLog returns a single exercise log entry.
// GET /api/user/exercise-logs/{id}
//
// OpenAPI: /api/docs#/operations/GetExerciseLog
func GetExerciseLog(ctx context.Context, input *GetExerciseLogInput) (*GetExerciseLogOutput, error) {
	var entity models.ExerciseLogEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise log not found")
	}
	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}
	return &GetExerciseLogOutput{Body: entity.ToDTO()}, nil
}

// UpdateExerciseLog partially updates an exercise log entry (PATCH).
// Recomputes RecordValue and recomputes which entry is the record.
// PATCH /api/user/exercise-logs/{id}
//
// OpenAPI: /api/docs#/operations/UpdateExerciseLog
func UpdateExerciseLog(ctx context.Context, input *UpdateExerciseLogInput) (*UpdateExerciseLogOutput, error) {
	var existing models.ExerciseLogEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise log not found")
	}
	if existing.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}

	patch := input.Body

	if patch.Reps != nil {
		existing.Reps = patch.Reps
	}
	if patch.Weight != nil {
		existing.Weight = patch.Weight
	}
	if patch.Duration != nil {
		existing.Duration = patch.Duration
	}
	if patch.Distance != nil {
		existing.Distance = patch.Distance
	}
	if patch.Time != nil {
		existing.Time = patch.Time
	}

	value, ok := models.ComputeRecordValue(existing.MeasurementType, existing.Reps, existing.Weight, existing.Duration, existing.Distance)
	if ok {
		existing.RecordValue = value
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}
		return RecomputeRecord(tx, existing.ExerciseID, existing.MeasurementType)
	})
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Reload to get recomputed isRecord
	database.DB.First(&existing, existing.ID)
	return &UpdateExerciseLogOutput{Body: existing.ToDTO()}, nil
}

// DeleteExerciseLog deletes an exercise log entry. Recomputes record
// afterwards. DELETE /api/user/exercise-logs/{id}
//
// OpenAPI: /api/docs#/operations/DeleteExerciseLog
func DeleteExerciseLog(ctx context.Context, input *DeleteExerciseLogInput) (*DeleteExerciseLogOutput, error) {
	var existing models.ExerciseLogEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise log not found")
	}
	if existing.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}

	exerciseID := existing.ExerciseID
	measurementType := existing.MeasurementType

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&existing).Error; err != nil {
			return err
		}
		return RecomputeRecord(tx, exerciseID, measurementType)
	})
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return nil, nil
}

// RecomputeRecord recalculates which ExerciseLog is the record for a given
// (exerciseId, measurementType) pair.
func RecomputeRecord(db *gorm.DB, exerciseID uint, measurementType string) error {
	// Clear all isRecord flags for this combo
	if err := db.Model(&models.ExerciseLogEntity{}).
		Where("exercise_id = ? AND measurement_type = ?", exerciseID, measurementType).
		Update("is_record", false).Error; err != nil {
		return err
	}

	// Find the entry with highest recordValue
	var best models.ExerciseLogEntity
	err := db.
		Where("exercise_id = ? AND measurement_type = ? AND record_value > 0", exerciseID, measurementType).
		Order("record_value DESC").
		First(&best).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	return db.Model(&best).Update("is_record", true).Error
}
