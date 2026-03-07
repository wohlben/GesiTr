package handlers

import (
	"net/http"

	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListWorkoutLogExercises(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutLogExerciseEntity{})

	if v := c.Query("workoutLogSectionId"); v != "" {
		db = db.Where("workout_log_section_id = ?", v)
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

	var scheme models.UserExerciseSchemeEntity
	if err := database.DB.First(&scheme, dto.UserExerciseSchemeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise scheme not found"})
		return
	}

	entity := models.WorkoutLogExerciseEntity{
		WorkoutLogSectionID:   dto.WorkoutLogSectionID,
		UserExerciseSchemeID:  dto.UserExerciseSchemeID,
		Position:              dto.Position,
		TargetMeasurementType: scheme.MeasurementType,
		TargetRestBetweenSets: scheme.RestBetweenSets,
		TargetTimePerRep:      scheme.TimePerRep,
	}

	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-create set rows from the scheme
	numSets := 0
	if scheme.Sets != nil {
		numSets = *scheme.Sets
	}
	for i := 1; i <= numSets; i++ {
		set := models.WorkoutLogExerciseSetEntity{
			WorkoutLogExerciseID: entity.ID,
			SetNumber:            i,
			Completed:            false,
			TargetReps:           scheme.Reps,
			TargetWeight:         scheme.Weight,
			TargetDuration:       scheme.Duration,
			TargetDistance:        scheme.Distance,
			TargetTime:           scheme.TargetTime,
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

	var dto models.WorkoutLogExercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only allow updating position
	existing.Position = dto.Position

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
	if err := database.DB.Delete(&models.WorkoutLogExerciseEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
