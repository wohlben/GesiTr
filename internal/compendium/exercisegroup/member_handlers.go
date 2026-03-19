package exercisegroup

import (
	"net/http"

	"gesitr/internal/database"

	"github.com/gin-gonic/gin"
)

func ListExerciseGroupMembers(c *gin.Context) {
	db := database.DB.Model(&ExerciseGroupMemberEntity{})

	if v := c.Query("groupTemplateId"); v != "" {
		db = db.Where("group_template_id = ?", v)
	}
	if v := c.Query("exerciseTemplateId"); v != "" {
		db = db.Where("exercise_template_id = ?", v)
	}

	var entities []ExerciseGroupMemberEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]ExerciseGroupMember, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateExerciseGroupMember(c *gin.Context) {
	var dto ExerciseGroupMember
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := ExerciseGroupMemberFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func DeleteExerciseGroupMember(c *gin.Context) {
	if err := database.DB.Delete(&ExerciseGroupMemberEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ExerciseGroupMember not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
