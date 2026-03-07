package handlers

import (
	"net/http"

	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func preloadWorkoutLog(db *gorm.DB) *gorm.DB {
	return db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Sections.Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Sections.Exercises.Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	})
}

func ListWorkoutLogs(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutLogEntity{})

	if v := c.Query("owner"); v != "" {
		db = db.Where("owner = ?", v)
	}
	if v := c.Query("workoutId"); v != "" {
		db = db.Where("workout_id = ?", v)
	}

	var entities []models.WorkoutLogEntity
	if err := preloadWorkoutLog(db).Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.WorkoutLog, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateWorkoutLog(c *gin.Context) {
	var dto models.WorkoutLog
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dto.WorkoutID != nil {
		var workout models.WorkoutEntity
		if err := database.DB.First(&workout, *dto.WorkoutID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
			return
		}
	}

	entity := models.WorkoutLogFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetWorkoutLog(c *gin.Context) {
	var entity models.WorkoutLogEntity
	if err := preloadWorkoutLog(database.DB).First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateWorkoutLog(c *gin.Context) {
	var existing models.WorkoutLogEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log not found"})
		return
	}

	var dto models.WorkoutLog
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.WorkoutLogFromDTO(dto)
	entity.ID = existing.ID

	if err := database.DB.Save(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := preloadWorkoutLog(database.DB).First(&entity, entity.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func DeleteWorkoutLog(c *gin.Context) {
	if err := database.DB.Delete(&models.WorkoutLogEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
