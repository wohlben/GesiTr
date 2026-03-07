package handlers

import (
	"net/http"

	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
)

func ListWorkoutLogExerciseSets(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutLogExerciseSetEntity{})

	if v := c.Query("workoutLogExerciseId"); v != "" {
		db = db.Where("workout_log_exercise_id = ?", v)
	}

	var entities []models.WorkoutLogExerciseSetEntity
	if err := db.Order("set_number").Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.WorkoutLogExerciseSet, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateWorkoutLogExerciseSet(c *gin.Context) {
	var dto models.WorkoutLogExerciseSet
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exercise models.WorkoutLogExerciseEntity
	if err := database.DB.First(&exercise, dto.WorkoutLogExerciseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise not found"})
		return
	}

	entity := models.WorkoutLogExerciseSetFromDTO(dto)
	entity.ID = 0
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func UpdateWorkoutLogExerciseSet(c *gin.Context) {
	var existing models.WorkoutLogExerciseSetEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise set not found"})
		return
	}

	var dto models.WorkoutLogExerciseSet
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only allow updating actual fields and completed — preserve target fields
	existing.Completed = dto.Completed
	existing.ActualReps = dto.ActualReps
	existing.ActualWeight = dto.ActualWeight
	existing.ActualDuration = dto.ActualDuration
	existing.ActualDistance = dto.ActualDistance
	existing.ActualTime = dto.ActualTime

	if err := database.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existing.Completed {
		maybeUpdateRecord(&existing)
	}

	c.JSON(http.StatusOK, existing.ToDTO())
}

func maybeUpdateRecord(set *models.WorkoutLogExerciseSetEntity) {
	var logExercise models.WorkoutLogExerciseEntity
	if err := database.DB.First(&logExercise, set.WorkoutLogExerciseID).Error; err != nil {
		return
	}

	var scheme models.UserExerciseSchemeEntity
	if err := database.DB.First(&scheme, logExercise.UserExerciseSchemeID).Error; err != nil {
		return
	}

	value, ok := computeRecordValue(logExercise.TargetMeasurementType, set)
	if !ok {
		return
	}

	var existing models.UserRecordEntity
	err := database.DB.
		Where("user_exercise_id = ? AND measurement_type = ?", scheme.UserExerciseID, logExercise.TargetMeasurementType).
		First(&existing).Error

	if err != nil {
		// No existing record — create
		database.DB.Create(&models.UserRecordEntity{
			UserExerciseID:          scheme.UserExerciseID,
			MeasurementType:         logExercise.TargetMeasurementType,
			RecordValue:             value,
			ActualReps:              set.ActualReps,
			ActualWeight:            set.ActualWeight,
			ActualDuration:          set.ActualDuration,
			ActualDistance:           set.ActualDistance,
			ActualTime:              set.ActualTime,
			WorkoutLogExerciseSetID: set.ID,
		})
		return
	}

	if value > existing.RecordValue {
		existing.RecordValue = value
		existing.ActualReps = set.ActualReps
		existing.ActualWeight = set.ActualWeight
		existing.ActualDuration = set.ActualDuration
		existing.ActualDistance = set.ActualDistance
		existing.ActualTime = set.ActualTime
		existing.WorkoutLogExerciseSetID = set.ID
		database.DB.Save(&existing)
	}
}

func DeleteWorkoutLogExerciseSet(c *gin.Context) {
	if err := database.DB.Delete(&models.WorkoutLogExerciseSetEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise set not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
