package handlers

import (
	"net/http"

	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
)

func ListUserExerciseSchemes(c *gin.Context) {
	db := database.DB.Model(&models.UserExerciseSchemeEntity{})

	if v := c.Query("userExerciseId"); v != "" {
		db = db.Where("user_exercise_id = ?", v)
	}
	if v := c.Query("measurementType"); v != "" {
		db = db.Where("measurement_type = ?", v)
	}

	var entities []models.UserExerciseSchemeEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.UserExerciseScheme, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateUserExerciseScheme(c *gin.Context) {
	var dto models.UserExerciseScheme
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the user exercise exists
	var userExercise models.UserExerciseEntity
	if err := database.DB.First(&userExercise, dto.UserExerciseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise not found"})
		return
	}

	entity := models.UserExerciseSchemeFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetUserExerciseScheme(c *gin.Context) {
	var entity models.UserExerciseSchemeEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise scheme not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateUserExerciseScheme(c *gin.Context) {
	var existing models.UserExerciseSchemeEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise scheme not found"})
		return
	}

	var dto models.UserExerciseScheme
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.UserExerciseSchemeFromDTO(dto)
	entity.ID = existing.ID
	entity.UserExerciseID = existing.UserExerciseID

	if err := database.DB.Save(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func DeleteUserExerciseScheme(c *gin.Context) {
	if err := database.DB.Delete(&models.UserExerciseSchemeEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise scheme not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
