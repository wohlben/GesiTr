package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/exercise/models"

	"github.com/gin-gonic/gin"
)

func ListExerciseSchemes(c *gin.Context) {
	db := database.DB.Model(&models.ExerciseSchemeEntity{}).
		Where("owner = ?", auth.GetUserID(c))

	if v := c.Query("exerciseId"); v != "" {
		db = db.Where("exercise_id = ?", v)
	}
	if v := c.Query("measurementType"); v != "" {
		db = db.Where("measurement_type = ?", v)
	}

	var entities []models.ExerciseSchemeEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.ExerciseScheme, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateExerciseScheme(c *gin.Context) {
	var dto models.ExerciseScheme
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the exercise exists
	var exercise models.ExerciseEntity
	if err := database.DB.First(&exercise, dto.ExerciseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	entity := models.ExerciseSchemeFromDTO(dto)
	entity.Owner = auth.GetUserID(c)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetExerciseScheme(c *gin.Context) {
	var entity models.ExerciseSchemeEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise scheme not found"})
		return
	}
	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateExerciseScheme(c *gin.Context) {
	var existing models.ExerciseSchemeEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise scheme not found"})
		return
	}
	if existing.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var dto models.ExerciseScheme
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.ExerciseSchemeFromDTO(dto)
	entity.ID = existing.ID
	entity.Owner = existing.Owner
	entity.ExerciseID = existing.ExerciseID

	if err := database.DB.Save(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func DeleteExerciseScheme(c *gin.Context) {
	var entity models.ExerciseSchemeEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise scheme not found"})
		return
	}
	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
