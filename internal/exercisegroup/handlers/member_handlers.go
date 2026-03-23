package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/exercisegroup/models"

	"github.com/gin-gonic/gin"
)

func ListExerciseGroupMembers(c *gin.Context) {
	db := database.DB.Model(&models.ExerciseGroupMemberEntity{})

	if v := c.Query("groupId"); v != "" {
		db = db.Where("group_id = ?", v)
	}
	if v := c.Query("exerciseId"); v != "" {
		db = db.Where("exercise_id = ?", v)
	}

	var entities []models.ExerciseGroupMemberEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.ExerciseGroupMember, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateExerciseGroupMember(c *gin.Context) {
	var dto models.ExerciseGroupMember
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.ExerciseGroupMemberFromDTO(dto)
	entity.Owner = auth.GetUserID(c)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func DeleteExerciseGroupMember(c *gin.Context) {
	var entity models.ExerciseGroupMemberEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ExerciseGroupMember not found"})
		return
	}

	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not the owner of this group member"})
		return
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
