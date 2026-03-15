package handlers

import (
	"net/http"
	"reflect"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListWorkoutLogSections(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutLogSectionEntity{})

	if v := c.Query("workoutLogId"); v != "" {
		if !requireLogOwner(c, parseUint(v)) {
			return
		}
		db = db.Where("workout_log_id = ?", v)
	} else {
		db = db.Joins("JOIN workout_logs ON workout_logs.id = workout_log_sections.workout_log_id").
			Where("workout_logs.owner = ?", auth.GetUserID(c))
	}

	var entities []models.WorkoutLogSectionEntity
	if err := db.Preload("Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).Order("position").Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.WorkoutLogSection, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateWorkoutLogSection(c *gin.Context) {
	var dto models.WorkoutLogSection
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, ok := requireLogStatus(c, dto.WorkoutLogID, models.WorkoutLogStatusPlanning); !ok {
		return
	}

	entity := models.WorkoutLogSectionFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetWorkoutLogSection(c *gin.Context) {
	var entity models.WorkoutLogSectionEntity
	if err := database.DB.Preload("Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log section not found"})
		return
	}
	if !requireLogOwner(c, entity.WorkoutLogID) {
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateWorkoutLogSection(c *gin.Context) {
	var existing models.WorkoutLogSectionEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log section not found"})
		return
	}

	if !requireLogOwner(c, existing.WorkoutLogID) {
		return
	}

	var patch struct {
		Type                 *models.WorkoutSectionType `json:"type"`
		Label                *string                    `json:"label"`
		RestBetweenExercises *int                       `json:"restBetweenExercises"`
	}
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reflect.ValueOf(patch).IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "patch body contains no updatable fields"})
		return
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

	if err := database.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with exercises and sets
	database.DB.Preload("Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	}).First(&existing, existing.ID)

	c.JSON(http.StatusOK, existing.ToDTO())
}

func DeleteWorkoutLogSection(c *gin.Context) {
	var existing models.WorkoutLogSectionEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log section not found"})
		return
	}

	if _, ok := requireLogStatus(c, existing.WorkoutLogID, models.WorkoutLogStatusPlanning); !ok {
		return
	}

	if err := database.DB.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
