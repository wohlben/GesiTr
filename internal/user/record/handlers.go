package record

import (
	"net/http"

	"gesitr/internal/database"

	"github.com/gin-gonic/gin"
)

func ListUserRecords(c *gin.Context) {
	db := database.DB.Model(&UserRecordEntity{})

	if v := c.Query("userExerciseId"); v != "" {
		db = db.Where("user_exercise_id = ?", v)
	}
	if v := c.Query("owner"); v != "" {
		db = db.Joins("JOIN user_exercises ON user_exercises.id = user_records.user_exercise_id").
			Where("user_exercises.owner = ?", v)
	}

	var entities []UserRecordEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]UserRecord, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func GetUserRecord(c *gin.Context) {
	var entity UserRecordEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User record not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}
