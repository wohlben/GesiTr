package handlers

import (
	"net/http"

	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
)

func computeRecordValue(measurementType string, set *models.WorkoutLogExerciseSetEntity) (float64, bool) {
	switch measurementType {
	case "REP_BASED", "AMRAP":
		if set.ActualReps == nil || set.ActualWeight == nil || *set.ActualWeight <= 0 {
			return 0, false
		}
		return *set.ActualWeight * (1 + float64(*set.ActualReps)/30), true
	case "TIME_BASED", "TIME", "EMOM", "ROUNDS_FOR_TIME":
		if set.ActualDuration == nil {
			return 0, false
		}
		return float64(*set.ActualDuration), true
	case "DISTANCE_BASED", "DISTANCE":
		if set.ActualDistance == nil {
			return 0, false
		}
		return *set.ActualDistance, true
	default:
		return 0, false
	}
}

func ListUserRecords(c *gin.Context) {
	db := database.DB.Model(&models.UserRecordEntity{})

	if v := c.Query("userExerciseId"); v != "" {
		db = db.Where("user_exercise_id = ?", v)
	}
	if v := c.Query("owner"); v != "" {
		db = db.Joins("JOIN user_exercises ON user_exercises.id = user_records.user_exercise_id").
			Where("user_exercises.owner = ?", v)
	}

	var entities []models.UserRecordEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.UserRecord, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func GetUserRecord(c *gin.Context) {
	var entity models.UserRecordEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User record not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}
