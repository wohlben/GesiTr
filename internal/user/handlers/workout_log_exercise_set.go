package handlers

import (
	"net/http"
	"time"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/user/models"

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
	entity.Status = models.WorkoutLogStatusPlanning
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

	var dto models.WorkoutLogExerciseSet
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Status transition: validate against the state machine
	if dto.Status != "" && dto.Status != existing.Status {
		if err := existing.Status.TransitionTo(dto.Status); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		existing.Status = dto.Status
		existing.StatusChangedAt = &now
	}

	// Update actual fields
	existing.ActualReps = dto.ActualReps
	existing.ActualWeight = dto.ActualWeight
	existing.ActualDuration = dto.ActualDuration
	existing.ActualDistance = dto.ActualDistance
	existing.ActualTime = dto.ActualTime

	// Update target fields and break when provided (non-nil in DTO)
	if dto.TargetReps != nil {
		existing.TargetReps = dto.TargetReps
	}
	if dto.TargetWeight != nil {
		existing.TargetWeight = dto.TargetWeight
	}
	if dto.TargetDuration != nil {
		existing.TargetDuration = dto.TargetDuration
	}
	if dto.TargetDistance != nil {
		existing.TargetDistance = dto.TargetDistance
	}
	if dto.TargetTime != nil {
		existing.TargetTime = dto.TargetTime
	}
	if dto.BreakAfterSeconds != nil {
		existing.BreakAfterSeconds = dto.BreakAfterSeconds
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}
		if existing.Status == models.WorkoutLogStatusFinished {
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

	// Check if all sets are terminal
	allTerminal := true
	anyAborted := false
	for _, s := range exercise.Sets {
		if !s.Status.IsTerminal() {
			allTerminal = false
			break
		}
		if s.Status == models.WorkoutLogStatusAborted {
			anyAborted = true
		}
	}

	if !allTerminal {
		return nil
	}

	now := time.Now()
	newStatus := models.WorkoutLogStatusFinished
	if anyAborted {
		newStatus = models.WorkoutLogStatusAborted
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
	for _, ex := range section.Exercises {
		if !ex.Status.IsTerminal() {
			allTerminal = false
			break
		}
		if ex.Status == models.WorkoutLogStatusAborted {
			anyAborted = true
		}
	}

	if !allTerminal {
		return nil
	}

	now := time.Now()
	newStatus := models.WorkoutLogStatusFinished
	if anyAborted {
		newStatus = models.WorkoutLogStatusAborted
	}

	if section.Status != newStatus {
		if err := db.Model(&section).Updates(map[string]any{
			"status":            newStatus,
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
	for _, s := range log.Sections {
		if !s.Status.IsTerminal() {
			allTerminal = false
			break
		}
		if s.Status == models.WorkoutLogStatusAborted {
			anyAborted = true
		}
	}

	if !allTerminal {
		return nil
	}

	newStatus = models.WorkoutLogStatusFinished
	if anyAborted {
		newStatus = models.WorkoutLogStatusAborted
	}

	if log.Status != newStatus {
		if err := db.Model(&log).Updates(map[string]any{
			"status":            newStatus,
			"status_changed_at": now,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func maybeUpdateRecord(db *gorm.DB, set *models.WorkoutLogExerciseSetEntity) {
	var logExercise models.WorkoutLogExerciseEntity
	if err := db.First(&logExercise, set.WorkoutLogExerciseID).Error; err != nil {
		return
	}

	var scheme models.UserExerciseSchemeEntity
	if err := db.First(&scheme, logExercise.SourceExerciseSchemeID).Error; err != nil {
		return
	}

	value, ok := computeRecordValue(logExercise.TargetMeasurementType, set)
	if !ok {
		return
	}

	var existing models.UserRecordEntity
	err := db.
		Where("user_exercise_id = ? AND measurement_type = ?", scheme.UserExerciseID, logExercise.TargetMeasurementType).
		First(&existing).Error

	if err != nil {
		// No existing record — create
		db.Create(&models.UserRecordEntity{
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
