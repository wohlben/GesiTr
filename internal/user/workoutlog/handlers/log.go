package handlers

import (
	"context"
	"reflect"
	"time"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/workoutlog/models"
	workoutschedule "gesitr/internal/user/workoutschedule"
	workoutmodels "gesitr/internal/workout/models"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func preloadWorkoutLog(db *gorm.DB) *gorm.DB {
	return db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Sections.Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Sections.Exercises.Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).Preload("Sections.Exercises.Sets.ExerciseLog")
}

// ListWorkoutLogs returns workout logs owned by the current user.
// GET /api/user/workout-logs
//
// OpenAPI: /api/docs#/operations/ListWorkoutLogs
func ListWorkoutLogs(ctx context.Context, input *ListWorkoutLogsInput) (*ListWorkoutLogsOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	// Lazy generation: ensure schedule-derived logs exist for current/next window
	_ = workoutschedule.GenerateForUser(database.DB, userID, time.Now())

	db := database.DB.Model(&models.WorkoutLogEntity{}).
		Where("owner = ?", userID)

	if input.WorkoutID != "" {
		db = db.Where("workout_id = ?", input.WorkoutID)
	}
	if input.Status != "" {
		db = db.Where("status = ?", input.Status)
	}
	if input.PeriodID != "" {
		db = db.Where("period_id = ?", input.PeriodID)
	}

	var entities []models.WorkoutLogEntity
	if err := preloadWorkoutLog(db).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.WorkoutLog, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutLogsOutput{Body: dtos}, nil
}

// CreateWorkoutLog creates a workout log in planning status.
// POST /api/user/workout-logs
//
// OpenAPI: /api/docs#/operations/CreateWorkoutLog
func CreateWorkoutLog(ctx context.Context, input *CreateWorkoutLogInput) (*CreateWorkoutLogOutput, error) {
	status := models.WorkoutLogStatusPlanning
	if input.Body.Status != nil {
		s := models.WorkoutLogStatus(*input.Body.Status)
		switch s {
		case models.WorkoutLogStatusProposed, models.WorkoutLogStatusCommitted:
			status = s
		default:
			return nil, huma.Error400BadRequest("status must be 'proposed' or 'committed' (or omitted for 'planning')")
		}
	}

	isCommitment := status == models.WorkoutLogStatusProposed || status == models.WorkoutLogStatusCommitted
	if isCommitment {
		if input.Body.DueStart == nil || input.Body.DueEnd == nil {
			return nil, huma.Error400BadRequest("dueStart and dueEnd are required for proposed/committed logs")
		}
		if !input.Body.DueEnd.After(*input.Body.DueStart) {
			return nil, huma.Error400BadRequest("dueEnd must be after dueStart")
		}
	}

	if input.Body.WorkoutID != nil {
		var w workoutmodels.WorkoutEntity
		if err := database.DB.First(&w, *input.Body.WorkoutID).Error; err != nil {
			return nil, huma.Error404NotFound("Workout not found")
		}

		// Uniqueness: only one planning log per workout (commitments are exempt)
		if status == models.WorkoutLogStatusPlanning {
			var count int64
			database.DB.Model(&models.WorkoutLogEntity{}).
				Where("workout_id = ? AND status = ?", *input.Body.WorkoutID, models.WorkoutLogStatusPlanning).
				Count(&count)
			if count > 0 {
				return nil, huma.Error409Conflict("A planning log already exists for this workout")
			}
		}
	}

	entity := models.WorkoutLogEntity{
		WorkoutID: input.Body.WorkoutID,
		Name:      input.Body.Name,
		Notes:     input.Body.Notes,
		Date:      input.Body.Date,
		DueStart:  input.Body.DueStart,
		DueEnd:    input.Body.DueEnd,
	}
	entity.Owner = humaconfig.GetUserID(ctx)
	entity.Status = status
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateWorkoutLogOutput{Body: entity.ToDTO()}, nil
}

// GetWorkoutLog returns a single workout log with its full tree.
// GET /api/user/workout-logs/{id}
//
// OpenAPI: /api/docs#/operations/GetWorkoutLog
func GetWorkoutLog(ctx context.Context, input *GetWorkoutLogInput) (*GetWorkoutLogOutput, error) {
	var entity models.WorkoutLogEntity
	if err := preloadWorkoutLog(database.DB).First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log not found")
	}
	if err := requireOwner(ctx, entity.Owner); err != nil {
		return nil, err
	}
	return &GetWorkoutLogOutput{Body: entity.ToDTO()}, nil
}

// UpdateWorkoutLog partially updates a workout log (name, notes).
// PATCH /api/user/workout-logs/{id}
//
// OpenAPI: /api/docs#/operations/UpdateWorkoutLog
func UpdateWorkoutLog(ctx context.Context, input *UpdateWorkoutLogInput) (*UpdateWorkoutLogOutput, error) {
	var existing models.WorkoutLogEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log not found")
	}
	if err := requireOwner(ctx, existing.Owner); err != nil {
		return nil, err
	}

	patch := input.Body

	if reflect.ValueOf(patch).IsZero() {
		return nil, huma.Error400BadRequest("patch body contains no updatable fields")
	}

	if patch.Name != nil {
		existing.Name = *patch.Name
	}
	if patch.Notes != nil {
		existing.Notes = patch.Notes
	}

	if err := database.DB.Save(&existing).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	if err := preloadWorkoutLog(database.DB).First(&existing, existing.ID).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateWorkoutLogOutput{Body: existing.ToDTO()}, nil
}

// DeleteWorkoutLog deletes a workout log. Only planning logs can be deleted.
// DELETE /api/user/workout-logs/{id}
//
// OpenAPI: /api/docs#/operations/DeleteWorkoutLog
func DeleteWorkoutLog(ctx context.Context, input *DeleteWorkoutLogInput) (*DeleteWorkoutLogOutput, error) {
	var existing models.WorkoutLogEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log not found")
	}
	if err := requireOwner(ctx, existing.Owner); err != nil {
		return nil, err
	}
	if existing.Status != models.WorkoutLogStatusPlanning && existing.Status != models.WorkoutLogStatusProposed && existing.Status != models.WorkoutLogStatusCommitted {
		return nil, huma.Error409Conflict("can only delete logs in planning, proposed, or committed status")
	}
	if err := database.DB.Delete(&existing).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}

// StartWorkoutLog transitions a workout log from planning to in-progress.
// POST /api/user/workout-logs/{id}/start
//
// OpenAPI: /api/docs#/operations/StartWorkoutLog
func StartWorkoutLog(ctx context.Context, input *StartWorkoutLogInput) (*StartWorkoutLogOutput, error) {
	var entity models.WorkoutLogEntity
	if err := preloadWorkoutLog(database.DB).First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log not found")
	}
	if err := requireOwner(ctx, entity.Owner); err != nil {
		return nil, err
	}

	if err := entity.Status.TransitionTo(models.WorkoutLogStatusInProgress); err != nil {
		return nil, huma.Error409Conflict(err.Error())
	}

	now := time.Now()

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Update log
		if err := tx.Model(&entity).Updates(map[string]any{
			"status":            models.WorkoutLogStatusInProgress,
			"status_changed_at": now,
			"date":              now,
		}).Error; err != nil {
			return err
		}

		// Update all sections
		for i := range entity.Sections {
			if err := tx.Model(&entity.Sections[i]).Updates(map[string]any{
				"status":            models.WorkoutLogItemStatusInProgress,
				"status_changed_at": now,
			}).Error; err != nil {
				return err
			}

			// Update all exercises
			for j := range entity.Sections[i].Exercises {
				if err := tx.Model(&entity.Sections[i].Exercises[j]).Updates(map[string]any{
					"status":            models.WorkoutLogItemStatusInProgress,
					"status_changed_at": now,
				}).Error; err != nil {
					return err
				}

				// Update all sets
				for k := range entity.Sections[i].Exercises[j].Sets {
					if err := tx.Model(&entity.Sections[i].Exercises[j].Sets[k]).Updates(map[string]any{
						"status":            models.WorkoutLogItemStatusInProgress,
						"status_changed_at": now,
					}).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Reload
	if err := preloadWorkoutLog(database.DB).First(&entity, entity.ID).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &StartWorkoutLogOutput{Body: entity.ToDTO()}, nil
}

// StartAdhocWorkoutLog creates and immediately starts an ad-hoc workout log.
// POST /api/user/workout-logs/adhoc
//
// OpenAPI: /api/docs#/operations/StartAdhocWorkoutLog
func StartAdhocWorkoutLog(ctx context.Context, input *StartAdhocWorkoutLogInput) (*StartAdhocWorkoutLogOutput, error) {
	now := time.Now()
	adhocLabel := "Adhoc"

	entity := models.WorkoutLogEntity{
		Owner:  humaconfig.GetUserID(ctx),
		Status: models.WorkoutLogStatusAdhoc,
		Name:   "Ad-hoc Workout",
		Date:   &now,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&entity).Error; err != nil {
			return err
		}

		section := models.WorkoutLogSectionEntity{
			WorkoutLogID: entity.ID,
			Type:         workoutmodels.WorkoutSectionTypeMain,
			Label:        &adhocLabel,
			Position:     0,
			Status:       models.WorkoutLogItemStatusInProgress,
		}
		if err := tx.Create(&section).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Reload with preloading
	if err := preloadWorkoutLog(database.DB).First(&entity, entity.ID).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &StartAdhocWorkoutLogOutput{Body: entity.ToDTO()}, nil
}

// FinishWorkoutLog finishes an adhoc workout log, deriving statuses.
// POST /api/user/workout-logs/{id}/finish
//
// OpenAPI: /api/docs#/operations/FinishWorkoutLog
func FinishWorkoutLog(ctx context.Context, input *FinishWorkoutLogInput) (*FinishWorkoutLogOutput, error) {
	var entity models.WorkoutLogEntity
	if err := preloadWorkoutLog(database.DB).First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log not found")
	}
	if err := requireOwner(ctx, entity.Owner); err != nil {
		return nil, err
	}

	if entity.Status != models.WorkoutLogStatusAdhoc {
		return nil, huma.Error409Conflict("can only finish adhoc workout logs")
	}

	now := time.Now()

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// First: skip any remaining in_progress sets
		for i := range entity.Sections {
			for j := range entity.Sections[i].Exercises {
				for k := range entity.Sections[i].Exercises[j].Sets {
					set := &entity.Sections[i].Exercises[j].Sets[k]
					if set.Status == models.WorkoutLogItemStatusInProgress {
						if err := tx.Model(set).Updates(map[string]any{
							"status":            models.WorkoutLogItemStatusSkipped,
							"status_changed_at": now,
						}).Error; err != nil {
							return err
						}
						set.Status = models.WorkoutLogItemStatusSkipped
					}
				}
			}
		}

		// Derive exercise statuses from sets
		for i := range entity.Sections {
			for j := range entity.Sections[i].Exercises {
				ex := &entity.Sections[i].Exercises[j]
				if ex.Status.IsTerminal() {
					continue
				}
				newStatus := deriveStatusFromChildren(extractItemStatuses(ex.Sets))
				if err := tx.Model(ex).Updates(map[string]any{
					"status":            newStatus,
					"status_changed_at": now,
				}).Error; err != nil {
					return err
				}
				ex.Status = newStatus
			}
		}

		// Derive section statuses from exercises
		for i := range entity.Sections {
			sec := &entity.Sections[i]
			if sec.Status.IsTerminal() {
				continue
			}
			newStatus := deriveStatusFromExercises(sec.Exercises)
			if err := tx.Model(sec).Updates(map[string]any{
				"status":            newStatus,
				"status_changed_at": now,
			}).Error; err != nil {
				return err
			}
			sec.Status = newStatus
		}

		// Derive log status from sections
		logStatus := deriveLogStatusFromSections(entity.Sections)
		if err := tx.Model(&entity).Updates(map[string]any{
			"status":            logStatus,
			"status_changed_at": now,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Reload
	if err := preloadWorkoutLog(database.DB).First(&entity, entity.ID).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &FinishWorkoutLogOutput{Body: entity.ToDTO()}, nil
}

// deriveStatusFromChildren derives the aggregate status from a list of item statuses.
func deriveStatusFromChildren(statuses []models.WorkoutLogItemStatus) models.WorkoutLogItemStatus {
	if len(statuses) == 0 {
		return models.WorkoutLogItemStatusSkipped
	}

	allFinished := true
	allSkipped := true
	anyAborted := false
	anySkipped := false
	anyPartiallyFinished := false

	for _, s := range statuses {
		if s != models.WorkoutLogItemStatusFinished {
			allFinished = false
		}
		if s != models.WorkoutLogItemStatusSkipped {
			allSkipped = false
		}
		if s == models.WorkoutLogItemStatusAborted {
			anyAborted = true
		}
		if s == models.WorkoutLogItemStatusSkipped {
			anySkipped = true
		}
		if s == models.WorkoutLogItemStatusPartiallyFinished {
			anyPartiallyFinished = true
		}
	}

	switch {
	case allFinished:
		return models.WorkoutLogItemStatusFinished
	case allSkipped:
		return models.WorkoutLogItemStatusSkipped
	case anyAborted:
		return models.WorkoutLogItemStatusAborted
	case anySkipped || anyPartiallyFinished:
		return models.WorkoutLogItemStatusPartiallyFinished
	default:
		return models.WorkoutLogItemStatusFinished
	}
}

func extractItemStatuses(sets []models.WorkoutLogExerciseSetEntity) []models.WorkoutLogItemStatus {
	statuses := make([]models.WorkoutLogItemStatus, len(sets))
	for i, s := range sets {
		statuses[i] = s.Status
	}
	return statuses
}

func deriveStatusFromExercises(exercises []models.WorkoutLogExerciseEntity) models.WorkoutLogItemStatus {
	if len(exercises) == 0 {
		return models.WorkoutLogItemStatusSkipped
	}
	statuses := make([]models.WorkoutLogItemStatus, len(exercises))
	for i, ex := range exercises {
		statuses[i] = ex.Status
	}
	return deriveStatusFromChildren(statuses)
}

func deriveLogStatusFromSections(sections []models.WorkoutLogSectionEntity) models.WorkoutLogStatus {
	if len(sections) == 0 {
		return models.WorkoutLogStatusFinished
	}

	allFinished := true
	anyAborted := false
	for _, s := range sections {
		if s.Status != models.WorkoutLogItemStatusFinished {
			allFinished = false
		}
		if s.Status == models.WorkoutLogItemStatusAborted {
			anyAborted = true
		}
	}

	switch {
	case allFinished:
		return models.WorkoutLogStatusFinished
	case anyAborted:
		return models.WorkoutLogStatusAborted
	default:
		return models.WorkoutLogStatusPartiallyFinished
	}
}

// AbandonWorkoutLog aborts a workout log and cascades to children.
// POST /api/user/workout-logs/{id}/abandon
//
// OpenAPI: /api/docs#/operations/AbandonWorkoutLog
func AbandonWorkoutLog(ctx context.Context, input *AbandonWorkoutLogInput) (*AbandonWorkoutLogOutput, error) {
	var entity models.WorkoutLogEntity
	if err := preloadWorkoutLog(database.DB).First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log not found")
	}
	if err := requireOwner(ctx, entity.Owner); err != nil {
		return nil, err
	}

	if err := entity.Status.TransitionTo(models.WorkoutLogStatusAborted); err != nil {
		return nil, huma.Error409Conflict(err.Error())
	}

	now := time.Now()

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Update log to aborted
		if err := tx.Model(&entity).Updates(map[string]any{
			"status":            models.WorkoutLogStatusAborted,
			"status_changed_at": now,
		}).Error; err != nil {
			return err
		}

		// Cascade to children — only change in_progress ones (preserve finished/skipped)
		for i := range entity.Sections {
			sec := &entity.Sections[i]
			if sec.Status == models.WorkoutLogItemStatusInProgress {
				if err := tx.Model(sec).Updates(map[string]any{
					"status":            models.WorkoutLogItemStatusAborted,
					"status_changed_at": now,
				}).Error; err != nil {
					return err
				}
			}

			for j := range sec.Exercises {
				ex := &sec.Exercises[j]
				if ex.Status == models.WorkoutLogItemStatusInProgress {
					if err := tx.Model(ex).Updates(map[string]any{
						"status":            models.WorkoutLogItemStatusAborted,
						"status_changed_at": now,
					}).Error; err != nil {
						return err
					}
				}

				for k := range ex.Sets {
					set := &ex.Sets[k]
					if set.Status == models.WorkoutLogItemStatusInProgress {
						if err := tx.Model(set).Updates(map[string]any{
							"status":            models.WorkoutLogItemStatusAborted,
							"status_changed_at": now,
						}).Error; err != nil {
							return err
						}
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Reload
	if err := preloadWorkoutLog(database.DB).First(&entity, entity.ID).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &AbandonWorkoutLogOutput{Body: entity.ToDTO()}, nil
}

// SkipWorkoutLog transitions a proposed workout log to skipped.
// POST /api/user/workout-logs/{id}/skip
//
// OpenAPI: /api/docs#/operations/SkipWorkoutLog
func SkipWorkoutLog(ctx context.Context, input *SkipWorkoutLogInput) (*SkipWorkoutLogOutput, error) {
	var entity models.WorkoutLogEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log not found")
	}
	if err := requireOwner(ctx, entity.Owner); err != nil {
		return nil, err
	}

	if err := entity.Status.TransitionTo(models.WorkoutLogStatusSkipped); err != nil {
		return nil, huma.Error409Conflict(err.Error())
	}

	now := time.Now()
	if err := database.DB.Model(&entity).Updates(map[string]any{
		"status":            models.WorkoutLogStatusSkipped,
		"status_changed_at": now,
	}).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	if err := preloadWorkoutLog(database.DB).First(&entity, entity.ID).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &SkipWorkoutLogOutput{Body: entity.ToDTO()}, nil
}

// CommitWorkoutLog transitions a proposed workout log to committed.
// POST /api/user/workout-logs/{id}/commit
//
// OpenAPI: /api/docs#/operations/CommitWorkoutLog
func CommitWorkoutLog(ctx context.Context, input *CommitWorkoutLogInput) (*CommitWorkoutLogOutput, error) {
	var entity models.WorkoutLogEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout log not found")
	}
	if err := requireOwner(ctx, entity.Owner); err != nil {
		return nil, err
	}

	if err := entity.Status.TransitionTo(models.WorkoutLogStatusCommitted); err != nil {
		return nil, huma.Error409Conflict(err.Error())
	}

	now := time.Now()
	if err := database.DB.Model(&entity).Updates(map[string]any{
		"status":            models.WorkoutLogStatusCommitted,
		"status_changed_at": now,
	}).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	if err := preloadWorkoutLog(database.DB).First(&entity, entity.ID).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CommitWorkoutLogOutput{Body: entity.ToDTO()}, nil
}
