package exercise

import (
	"net/http"

	"gesitr/internal/database"

	"github.com/gin-gonic/gin"
)

func ListUserExercises(c *gin.Context) {
	db := database.DB.Model(&UserExerciseEntity{})

	if v := c.Query("owner"); v != "" {
		db = db.Where("owner = ?", v)
	}
	if v := c.Query("compendiumExerciseId"); v != "" {
		db = db.Where("compendium_exercise_id = ?", v)
	}

	var entities []UserExerciseEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]UserExercise, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateUserExercise(c *gin.Context) {
	var dto UserExercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := UserExerciseFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetUserExercise(c *gin.Context) {
	var entity UserExerciseEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func DeleteUserExercise(c *gin.Context) {
	if err := database.DB.Delete(&UserExerciseEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User exercise not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
