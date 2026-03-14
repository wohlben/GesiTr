package handlers

import (
	"net/http"

	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListWorkoutLogSections(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutLogSectionEntity{})

	if v := c.Query("workoutLogId"); v != "" {
		db = db.Where("workout_log_id = ?", v)
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
	c.JSON(http.StatusOK, entity.ToDTO())
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
