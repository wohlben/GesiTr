package handlers

import (
	"net/http"

	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
)

func ListWorkoutSections(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutSectionEntity{})

	if v := c.Query("workoutId"); v != "" {
		db = db.Where("workout_id = ?", v)
	}

	var entities []models.WorkoutSectionEntity
	if err := db.Preload("Exercises").Order("position").Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.WorkoutSection, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateWorkoutSection(c *gin.Context) {
	var dto models.WorkoutSection
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var workout models.WorkoutEntity
	if err := database.DB.First(&workout, dto.WorkoutID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return
	}

	entity := models.WorkoutSectionFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetWorkoutSection(c *gin.Context) {
	var entity models.WorkoutSectionEntity
	if err := database.DB.Preload("Exercises").First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout section not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func DeleteWorkoutSection(c *gin.Context) {
	if err := database.DB.Delete(&models.WorkoutSectionEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout section not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
