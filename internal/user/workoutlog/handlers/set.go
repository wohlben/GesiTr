package handlers

import (
	"net/http"
	"reflect"
	"time"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	userexercise "gesitr/internal/user/exercise"
	"gesitr/internal/user/record"
	"gesitr/internal/user/workoutlog/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListWorkoutLogExerciseSets(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutLogExerciseSetEntity{})

	if v := c.Query("workoutLogExerciseId"); v != "" {
		var exercise models.WorkoutLogExerciseEntity
		if err := database.DB.First(&exercise, v).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise not found"})
			return
		}
		if !requireLogOwner(c, exercise.WorkoutLogID) {
			return
		}
		db = db.Where("workout_log_exercise_id = ?", v)
	} else {
		db = db.Joins("JOIN workout_logs ON workout_logs.id = workout_log_exercise_sets.workout_log_id").
			Where("workout_logs.owner = ?", auth.GetUserID(c))
	}

	var entities []models.WorkoutLogExerciseSetEntity
	if err := db.Order("set_number").Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.WorkoutLogExerciseSet, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateWorkoutLogExerciseSet(c *gin.Context) {
	var dto models.WorkoutLogExerciseSet
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exercise models.WorkoutLogExerciseEntity
	if err := database.DB.First(&exercise, dto.WorkoutLogExerciseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise not found"})
		return
	}

	// Guard: parent log must be in planning status (also checks ownership)
	if _, ok := requireLogStatus(c, exercise.WorkoutLogID, models.WorkoutLogStatusPlanning); !ok {
		return
	}

	entity := models.WorkoutLogExerciseSetFromDTO(dto)
	entity.ID = 0
	entity.CreatedAt = time.Time{}
	entity.UpdatedAt = time.Time{}
	entity.WorkoutLogID = exercise.WorkoutLogID
	entity.Status = models.WorkoutLogItemStatusPlanning
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func UpdateWorkoutLogExerciseSet(c *gin.Context) {
	var existing models.WorkoutLogExerciseSetEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise set not found"})
		return
	}

	// Owner check via parent log
	if !requireLogOwner(c, existing.WorkoutLogID) {
		return
	}

	var patch struct {
		Status            models.WorkoutLogItemStatus `json:"status"`
		BreakAfterSeconds *int                        `json:"breakAfterSeconds"`
		TargetReps        *int                        `json:"targetReps"`
		TargetWeight      *float64                    `json:"targetWeight"`
		TargetDuration    *int                        `json:"targetDuration"`
		TargetDistance    *float64                    `json:"targetDistance"`
		TargetTime        *int                        `json:"targetTime"`
		ActualReps        *int                        `json:"actualReps"`
		ActualWeight      *float64                    `json:"actualWeight"`
		ActualDuration    *int                        `json:"actualDuration"`
		ActualDistance    *float64                    `json:"actualDistance"`
		ActualTime        *int                        `json:"actualTime"`
	}
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reflect.ValueOf(patch).IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "patch body contains no updatable fields"})
		return
	}

	// Status transition: validate against the state machine
	if patch.Status != "" && patch.Status != existing.Status {
		if err := existing.Status.TransitionTo(patch.Status); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		existing.Status = patch.Status
		existing.StatusChangedAt = &now
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

	// Update actual fields when provided
	if patch.ActualReps != nil {
		existing.ActualReps = patch.ActualReps
	}
	if patch.ActualWeight != nil {
		existing.ActualWeight = patch.ActualWeight
	}
	if patch.ActualDuration != nil {
		existing.ActualDuration = patch.ActualDuration
	}
	if patch.ActualDistance != nil {
		existing.ActualDistance = patch.ActualDistance
	}
	if patch.ActualTime != nil {
		existing.ActualTime = patch.ActualTime
	}

	if patch.BreakAfterSeconds != nil {
		existing.BreakAfterSeconds = patch.BreakAfterSeconds
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}
		if existing.Status == models.WorkoutLogItemStatusFinished {
			maybeUpdateRecord(tx, &existing)
		}
		return propagateStatus(tx, existing.WorkoutLogExerciseID)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing.ToDTO())
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

func computeRecordValue(measurementType string, set *models.WorkoutLogExerciseSetEntity) (float64, bool) {
	switch measurementType {
	case "REP_BASED", "AMRAP":
		if set.ActualReps == nil || set.ActualWeight == nil || *set.ActualWeight <= 0 {
			return 0, false
		}
		return *set.ActualWeight * (1 + float64(*set.ActualReps)/30), true
	case "TIME_BASED", "TIME", "EMOM", "ROUNDS_FOR_TIME":
		if set.ActualDuration == nil {
			return 0, false
		}
		return float64(*set.ActualDuration), true
	case "DISTANCE_BASED", "DISTANCE":
		if set.ActualDistance == nil {
			return 0, false
		}
		return *set.ActualDistance, true
	default:
		return 0, false
	}
}

func maybeUpdateRecord(db *gorm.DB, set *models.WorkoutLogExerciseSetEntity) {
	var logExercise models.WorkoutLogExerciseEntity
	if err := db.First(&logExercise, set.WorkoutLogExerciseID).Error; err != nil {
		return
	}

	var scheme userexercise.UserExerciseSchemeEntity
	if err := db.First(&scheme, logExercise.SourceExerciseSchemeID).Error; err != nil {
		return
	}

	value, ok := computeRecordValue(logExercise.TargetMeasurementType, set)
	if !ok {
		return
	}

	var existing record.UserRecordEntity
	err := db.
		Where("user_exercise_id = ? AND measurement_type = ?", scheme.UserExerciseID, logExercise.TargetMeasurementType).
		First(&existing).Error

	if err != nil {
		// No existing record — create
		db.Create(&record.UserRecordEntity{
			UserExerciseID:          scheme.UserExerciseID,
			MeasurementType:         logExercise.TargetMeasurementType,
			RecordValue:             value,
			ActualReps:              set.ActualReps,
			ActualWeight:            set.ActualWeight,
			ActualDuration:          set.ActualDuration,
			ActualDistance:          set.ActualDistance,
			ActualTime:              set.ActualTime,
			WorkoutLogExerciseSetID: set.ID,
		})
		return
	}

	if value > existing.RecordValue {
		existing.RecordValue = value
		existing.ActualReps = set.ActualReps
		existing.ActualWeight = set.ActualWeight
		existing.ActualDuration = set.ActualDuration
		existing.ActualDistance = set.ActualDistance
		existing.ActualTime = set.ActualTime
		existing.WorkoutLogExerciseSetID = set.ID
		db.Save(&existing)
	}
}

func DeleteWorkoutLogExerciseSet(c *gin.Context) {
	var existing models.WorkoutLogExerciseSetEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise set not found"})
		return
	}

	// Guard: parent log must be in planning status (also checks ownership)
	if _, ok := requireLogStatus(c, existing.WorkoutLogID, models.WorkoutLogStatusPlanning); !ok {
		return
	}

	exerciseID := existing.WorkoutLogExerciseID
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&existing).Error; err != nil {
			return err
		}
		return propagateStatus(tx, exerciseID)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
