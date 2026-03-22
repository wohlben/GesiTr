package handlers

import (
	"net/http"
	"reflect"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	userexercisemodels "gesitr/internal/user/exercise/models"
	"gesitr/internal/user/workoutlog/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListWorkoutLogExercises(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutLogExerciseEntity{})

	if v := c.Query("workoutLogSectionId"); v != "" {
		var section models.WorkoutLogSectionEntity
		if err := database.DB.First(&section, v).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout log section not found"})
			return
		}
		if !requireLogOwner(c, section.WorkoutLogID) {
			return
		}
		db = db.Where("workout_log_section_id = ?", v)
	} else {
		db = db.Joins("JOIN workout_logs ON workout_logs.id = workout_log_exercises.workout_log_id").
			Where("workout_logs.owner = ?", auth.GetUserID(c))
	}

	var entities []models.WorkoutLogExerciseEntity
	if err := db.Preload("Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).Order("position").Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.WorkoutLogExercise, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateWorkoutLogExercise(c *gin.Context) {
	var dto models.WorkoutLogExercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var section models.WorkoutLogSectionEntity
	if err := database.DB.First(&section, dto.WorkoutLogSectionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log section not found"})
		return
	}

	// Guard: parent log must be in planning or adhoc status
	log, ok := requireLogStatus(c, section.WorkoutLogID, models.WorkoutLogStatusPlanning, models.WorkoutLogStatusAdhoc)
	if !ok {
		return
	}

	var scheme userexercisemodels.UserExerciseSchemeEntity
	if err := database.DB.First(&scheme, dto.SourceExerciseSchemeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise scheme not found"})
		return
	}

	breakAfter := section.RestBetweenExercises
	if dto.BreakAfterSeconds != nil {
		breakAfter = dto.BreakAfterSeconds
	}

	// For adhoc logs, exercises start in_progress immediately
	exerciseStatus := models.WorkoutLogItemStatusPlanning
	if log.Status == models.WorkoutLogStatusAdhoc {
		exerciseStatus = models.WorkoutLogItemStatusInProgress
	}

	entity := models.WorkoutLogExerciseEntity{
		WorkoutLogSectionID:    dto.WorkoutLogSectionID,
		WorkoutLogID:           section.WorkoutLogID,
		SourceExerciseSchemeID: dto.SourceExerciseSchemeID,
		Position:               dto.Position,
		Status:                 exerciseStatus,
		BreakAfterSeconds:      breakAfter,
		TargetMeasurementType:  scheme.MeasurementType,
		TargetTimePerRep:       scheme.TimePerRep,
	}

	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		entity.Sets = append(entity.Sets, set)
	}

	c.JSON(http.StatusCreated, entity.ToDTO())
}

func UpdateWorkoutLogExercise(c *gin.Context) {
	var existing models.WorkoutLogExerciseEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise not found"})
		return
	}

	if !requireLogOwner(c, existing.WorkoutLogID) {
		return
	}

	var patch struct {
		Position          *int `json:"position"`
		BreakAfterSeconds *int `json:"breakAfterSeconds"`
	}
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reflect.ValueOf(patch).IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "patch body contains no updatable fields"})
		return
	}

	if patch.Position != nil {
		existing.Position = *patch.Position
	}
	if patch.BreakAfterSeconds != nil {
		existing.BreakAfterSeconds = patch.BreakAfterSeconds
	}

	if err := database.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with sets
	database.DB.Preload("Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).First(&existing, existing.ID)

	c.JSON(http.StatusOK, existing.ToDTO())
}

func DeleteWorkoutLogExercise(c *gin.Context) {
	var existing models.WorkoutLogExerciseEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise not found"})
		return
	}

	// Guard: parent log must be in planning or adhoc status
	if _, ok := requireLogStatus(c, existing.WorkoutLogID, models.WorkoutLogStatusPlanning, models.WorkoutLogStatusAdhoc); !ok {
		return
	}

	sectionID := existing.WorkoutLogSectionID
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&existing).Error; err != nil {
			return err
		}
		return propagateSectionStatus(tx, sectionID)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
