package handlers

import (
	"net/http"
	"time"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/user/exerciselog/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListExerciseLogs(c *gin.Context) {
	db := database.DB.Model(&models.ExerciseLogEntity{}).
		Where("owner = ?", auth.GetUserID(c))

	if v := c.Query("exerciseId"); v != "" {
		db = db.Where("exercise_id = ?", v)
	}
	if v := c.Query("measurementType"); v != "" {
		db = db.Where("measurement_type = ?", v)
	}
	if v := c.Query("isRecord"); v == "true" {
		db = db.Where("is_record = ?", true)
	}
	if v := c.Query("from"); v != "" {
		db = db.Where("performed_at >= ?", v)
	}
	if v := c.Query("to"); v != "" {
		db = db.Where("performed_at <= ?", v)
	}

	var entities []models.ExerciseLogEntity
	if err := db.Order("performed_at DESC").Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.ExerciseLog, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateExerciseLog(c *gin.Context) {
	var dto models.ExerciseLog
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.ExerciseLogFromDTO(dto)
	entity.ID = 0
	entity.CreatedAt = time.Time{}
	entity.UpdatedAt = time.Time{}
	entity.Owner = auth.GetUserID(c)

	if entity.PerformedAt.IsZero() {
		entity.PerformedAt = time.Now()
	}

	value, ok := models.ComputeRecordValue(entity.MeasurementType, entity.Reps, entity.Weight, entity.Duration, entity.Distance)
	if ok {
		entity.RecordValue = value
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&entity).Error; err != nil {
			return err
		}
		return RecomputeRecord(tx, entity.ExerciseID, entity.MeasurementType)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload to get recomputed isRecord
	database.DB.First(&entity, entity.ID)
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetExerciseLog(c *gin.Context) {
	var entity models.ExerciseLogEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise log not found"})
		return
	}
	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateExerciseLog(c *gin.Context) {
	var existing models.ExerciseLogEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise log not found"})
		return
	}
	if existing.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var patch struct {
		Reps     *int     `json:"reps"`
		Weight   *float64 `json:"weight"`
		Duration *int     `json:"duration"`
		Distance *float64 `json:"distance"`
		Time     *int     `json:"time"`
	}
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if patch.Reps != nil {
		existing.Reps = patch.Reps
	}
	if patch.Weight != nil {
		existing.Weight = patch.Weight
	}
	if patch.Duration != nil {
		existing.Duration = patch.Duration
	}
	if patch.Distance != nil {
		existing.Distance = patch.Distance
	}
	if patch.Time != nil {
		existing.Time = patch.Time
	}

	value, ok := models.ComputeRecordValue(existing.MeasurementType, existing.Reps, existing.Weight, existing.Duration, existing.Distance)
	if ok {
		existing.RecordValue = value
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}
		return RecomputeRecord(tx, existing.ExerciseID, existing.MeasurementType)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload to get recomputed isRecord
	database.DB.First(&existing, existing.ID)
	c.JSON(http.StatusOK, existing.ToDTO())
}

func DeleteExerciseLog(c *gin.Context) {
	var existing models.ExerciseLogEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise log not found"})
		return
	}
	if existing.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	exerciseID := existing.ExerciseID
	measurementType := existing.MeasurementType

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&existing).Error; err != nil {
			return err
		}
		return RecomputeRecord(tx, exerciseID, measurementType)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// RecomputeRecord recalculates which ExerciseLog is the record for a given
// (exerciseId, measurementType) pair.
func RecomputeRecord(db *gorm.DB, exerciseID uint, measurementType string) error {
	// Clear all isRecord flags for this combo
	if err := db.Model(&models.ExerciseLogEntity{}).
		Where("exercise_id = ? AND measurement_type = ?", exerciseID, measurementType).
		Update("is_record", false).Error; err != nil {
		return err
	}

	// Find the entry with highest recordValue
	var best models.ExerciseLogEntity
	err := db.
		Where("exercise_id = ? AND measurement_type = ? AND record_value > 0", exerciseID, measurementType).
		Order("record_value DESC").
		First(&best).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	return db.Model(&best).Update("is_record", true).Error
}
