package handlers

import (
	"net/http"

	"gesitr/internal/compendium/models"
	"gesitr/internal/database"

	"github.com/gin-gonic/gin"
)

func ListExerciseRelationships(c *gin.Context) {
	db := database.DB.Model(&models.ExerciseRelationshipEntity{})

	if v := c.Query("fromExerciseTemplateId"); v != "" {
		db = db.Where("from_exercise_template_id = ?", v)
	}
	if v := c.Query("toExerciseTemplateId"); v != "" {
		db = db.Where("to_exercise_template_id = ?", v)
	}
	if v := c.Query("relationshipType"); v != "" {
		db = db.Where("relationship_type = ?", v)
	}

	var entities []models.ExerciseRelationshipEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.ExerciseRelationship, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateExerciseRelationship(c *gin.Context) {
	var dto models.ExerciseRelationship
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.ExerciseRelationshipFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func DeleteExerciseRelationship(c *gin.Context) {
	if err := database.DB.Delete(&models.ExerciseRelationshipEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ExerciseRelationship not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
