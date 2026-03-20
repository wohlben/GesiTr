package handlers

import (
	"net/http"
	"reflect"
	"time"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	workoutmodels "gesitr/internal/user/workout/models"
	"gesitr/internal/user/workoutlog/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func preloadWorkoutLog(db *gorm.DB) *gorm.DB {
	return db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Sections.Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Sections.Exercises.Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("set_number")
	})
}

func ListWorkoutLogs(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutLogEntity{}).
		Where("owner = ?", auth.GetUserID(c))

	if v := c.Query("workoutId"); v != "" {
		db = db.Where("workout_id = ?", v)
	}
	if v := c.Query("status"); v != "" {
		db = db.Where("status = ?", v)
	}

	var entities []models.WorkoutLogEntity
	if err := preloadWorkoutLog(db).Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.WorkoutLog, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateWorkoutLog(c *gin.Context) {
	var dto models.WorkoutLog
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dto.WorkoutID != nil {
		var w workoutmodels.WorkoutEntity
		if err := database.DB.First(&w, *dto.WorkoutID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
			return
		}

		// Uniqueness: only one planning log per workout
		var count int64
		database.DB.Model(&models.WorkoutLogEntity{}).
			Where("workout_id = ? AND status = ?", *dto.WorkoutID, models.WorkoutLogStatusPlanning).
			Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "A planning log already exists for this workout"})
			return
		}
	}

	entity := models.WorkoutLogFromDTO(dto)
	entity.Owner = auth.GetUserID(c)
	entity.Status = models.WorkoutLogStatusPlanning
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetWorkoutLog(c *gin.Context) {
	var entity models.WorkoutLogEntity
	if err := preloadWorkoutLog(database.DB).First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log not found"})
		return
	}
	if !requireOwner(c, entity.Owner) {
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateWorkoutLog(c *gin.Context) {
	var existing models.WorkoutLogEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log not found"})
		return
	}
	if !requireOwner(c, existing.Owner) {
		return
	}

	var patch struct {
		Name  *string `json:"name"`
		Notes *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reflect.ValueOf(patch).IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "patch body contains no updatable fields"})
		return
	}

	if patch.Name != nil {
		existing.Name = *patch.Name
	}
	if patch.Notes != nil {
		existing.Notes = patch.Notes
	}

	if err := database.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := preloadWorkoutLog(database.DB).First(&existing, existing.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing.ToDTO())
}

func DeleteWorkoutLog(c *gin.Context) {
	var existing models.WorkoutLogEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log not found"})
		return
	}
	if !requireOwner(c, existing.Owner) {
		return
	}
	if existing.Status != models.WorkoutLogStatusPlanning {
		c.JSON(http.StatusConflict, gin.H{"error": "can only delete logs in planning status"})
		return
	}
	if err := database.DB.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func StartWorkoutLog(c *gin.Context) {
	var entity models.WorkoutLogEntity
	if err := preloadWorkoutLog(database.DB).First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log not found"})
		return
	}
	if !requireOwner(c, entity.Owner) {
		return
	}

	if err := entity.Status.TransitionTo(models.WorkoutLogStatusInProgress); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Update log
		if err := tx.Model(&entity).Updates(map[string]any{
			"status":            models.WorkoutLogStatusInProgress,
			"status_changed_at": now,
			"date":              now,
		}).Error; err != nil {
			return err
		}

		// Update all sections
		for i := range entity.Sections {
			if err := tx.Model(&entity.Sections[i]).Updates(map[string]any{
				"status":            models.WorkoutLogItemStatusInProgress,
				"status_changed_at": now,
			}).Error; err != nil {
				return err
			}

			// Update all exercises
			for j := range entity.Sections[i].Exercises {
				if err := tx.Model(&entity.Sections[i].Exercises[j]).Updates(map[string]any{
					"status":            models.WorkoutLogItemStatusInProgress,
					"status_changed_at": now,
				}).Error; err != nil {
					return err
				}

				// Update all sets
				for k := range entity.Sections[i].Exercises[j].Sets {
					if err := tx.Model(&entity.Sections[i].Exercises[j].Sets[k]).Updates(map[string]any{
						"status":            models.WorkoutLogItemStatusInProgress,
						"status_changed_at": now,
					}).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload
	if err := preloadWorkoutLog(database.DB).First(&entity, entity.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func AbandonWorkoutLog(c *gin.Context) {
	var entity models.WorkoutLogEntity
	if err := preloadWorkoutLog(database.DB).First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log not found"})
		return
	}
	if !requireOwner(c, entity.Owner) {
		return
	}

	if err := entity.Status.TransitionTo(models.WorkoutLogStatusAborted); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Update log to aborted
		if err := tx.Model(&entity).Updates(map[string]any{
			"status":            models.WorkoutLogStatusAborted,
			"status_changed_at": now,
		}).Error; err != nil {
			return err
		}

		// Cascade to children — only change in_progress ones (preserve finished/skipped)
		for i := range entity.Sections {
			sec := &entity.Sections[i]
			if sec.Status == models.WorkoutLogItemStatusInProgress {
				if err := tx.Model(sec).Updates(map[string]any{
					"status":            models.WorkoutLogItemStatusAborted,
					"status_changed_at": now,
				}).Error; err != nil {
					return err
				}
			}

			for j := range sec.Exercises {
				ex := &sec.Exercises[j]
				if ex.Status == models.WorkoutLogItemStatusInProgress {
					if err := tx.Model(ex).Updates(map[string]any{
						"status":            models.WorkoutLogItemStatusAborted,
						"status_changed_at": now,
					}).Error; err != nil {
						return err
					}
				}

				for k := range ex.Sets {
					set := &ex.Sets[k]
					if set.Status == models.WorkoutLogItemStatusInProgress {
						if err := tx.Model(set).Updates(map[string]any{
							"status":            models.WorkoutLogItemStatusAborted,
							"status_changed_at": now,
						}).Error; err != nil {
							return err
						}
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload
	if err := preloadWorkoutLog(database.DB).First(&entity, entity.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}
