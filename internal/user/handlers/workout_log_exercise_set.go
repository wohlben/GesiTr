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

	// Update actual fields and completed
	existing.Completed = dto.Completed
	existing.ActualReps = dto.ActualReps
	existing.ActualWeight = dto.ActualWeight
	existing.ActualDuration = dto.ActualDuration
	existing.ActualDistance = dto.ActualDistance
	existing.ActualTime = dto.ActualTime

	// Update target fields and break when provided (non-nil in DTO)
	if dto.TargetReps != nil {
		existing.TargetReps = dto.TargetReps
	}
	if dto.TargetWeight != nil {
		existing.TargetWeight = dto.TargetWeight
	}
	if dto.TargetDuration != nil {
		existing.TargetDuration = dto.TargetDuration
	}
	if dto.TargetDistance != nil {
		existing.TargetDistance = dto.TargetDistance
	}
	if dto.TargetTime != nil {
		existing.TargetTime = dto.TargetTime
	}
	if dto.BreakAfterSeconds != nil {
		existing.BreakAfterSeconds = dto.BreakAfterSeconds
	}

	if err := database.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existing.Completed {
		maybeUpdateRecord(&existing)
	}

	propagateCompletion(existing.WorkoutLogExerciseID)

	c.JSON(http.StatusOK, existing.ToDTO())
}

func propagateCompletion(exerciseID uint) {
	// Check if all sets of this exercise are completed
	var exercise models.WorkoutLogExerciseEntity
	if err := database.DB.Preload("Sets").First(&exercise, exerciseID).Error; err != nil {
		return
	}

	allSetsCompleted := len(exercise.Sets) > 0
	for _, s := range exercise.Sets {
		if !s.Completed {
			allSetsCompleted = false
			break
		}
	}

	if exercise.Completed != allSetsCompleted {
		exercise.Completed = allSetsCompleted
		database.DB.Model(&exercise).Update("completed", allSetsCompleted)
	}

	propagateSectionCompletion(exercise.WorkoutLogSectionID)
}

func propagateSectionCompletion(sectionID uint) {
	// Check if all exercises of the section are completed
	var section models.WorkoutLogSectionEntity
	if err := database.DB.Preload("Exercises").First(&section, sectionID).Error; err != nil {
		return
	}

	allExercisesCompleted := len(section.Exercises) > 0
	for _, ex := range section.Exercises {
		if !ex.Completed {
			allExercisesCompleted = false
			break
		}
	}

	if section.Completed != allExercisesCompleted {
		section.Completed = allExercisesCompleted
		database.DB.Model(&section).Update("completed", allExercisesCompleted)
	}

	// Check if all sections of the log are completed
	var log models.WorkoutLogEntity
	if err := database.DB.Preload("Sections").First(&log, section.WorkoutLogID).Error; err != nil {
		return
	}

	allSectionsCompleted := len(log.Sections) > 0
	for _, s := range log.Sections {
		if !s.Completed {
			allSectionsCompleted = false
			break
		}
	}

	if log.Completed != allSectionsCompleted {
		log.Completed = allSectionsCompleted
		database.DB.Model(&log).Update("completed", allSectionsCompleted)
	}
}

func maybeUpdateRecord(set *models.WorkoutLogExerciseSetEntity) {
	var logExercise models.WorkoutLogExerciseEntity
	if err := database.DB.First(&logExercise, set.WorkoutLogExerciseID).Error; err != nil {
		return
	}

	var scheme models.UserExerciseSchemeEntity
	if err := database.DB.First(&scheme, logExercise.SourceExerciseSchemeID).Error; err != nil {
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
			ActualDistance:          set.ActualDistance,
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
	var existing models.WorkoutLogExerciseSetEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log exercise set not found"})
		return
	}

	exerciseID := existing.WorkoutLogExerciseID
	if err := database.DB.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	propagateCompletion(exerciseID)
	c.JSON(http.StatusNoContent, nil)
}
