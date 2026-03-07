package handlers

import (
	"net/http"

	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
)

func ListUserExercises(c *gin.Context) {
	db := database.DB.Model(&models.UserExerciseEntity{})

	if v := c.Query("owner"); v != "" {
		db = db.Where("owner = ?", v)
	}
	if v := c.Query("exerciseTemplateId"); v != "" {
		db = db.Where("exercise_template_id = ?", v)
	}

	var entities []models.UserExerciseEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.UserExercise, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateUserExercise(c *gin.Context) {
	var dto models.UserExercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.UserExerciseFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetUserExercise(c *gin.Context) {
	var entity models.UserExerciseEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func DeleteUserExercise(c *gin.Context) {
	if err := database.DB.Delete(&models.UserExerciseEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
