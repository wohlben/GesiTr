package workout

import (
	"net/http"

	"gesitr/internal/database"
	userexercise "gesitr/internal/user/exercise"

	"github.com/gin-gonic/gin"
)

func ListWorkoutSectionExercises(c *gin.Context) {
	db := database.DB.Model(&WorkoutSectionExerciseEntity{})

	if v := c.Query("workoutSectionId"); v != "" {
		db = db.Where("workout_section_id = ?", v)
	}

	var entities []WorkoutSectionExerciseEntity
	if err := db.Order("position").Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]WorkoutSectionExercise, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateWorkoutSectionExercise(c *gin.Context) {
	var dto WorkoutSectionExercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var section WorkoutSectionEntity
	if err := database.DB.First(&section, dto.WorkoutSectionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout section not found"})
		return
	}

	var scheme userexercise.UserExerciseSchemeEntity
	if err := database.DB.First(&scheme, dto.UserExerciseSchemeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise scheme not found"})
		return
	}

	entity := WorkoutSectionExerciseFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func DeleteWorkoutSectionExercise(c *gin.Context) {
	if err := database.DB.Delete(&WorkoutSectionExerciseEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout section exercise not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
