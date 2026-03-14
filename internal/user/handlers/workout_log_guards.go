package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
)

// requireOwner checks that the authenticated user matches the given owner.
func requireOwner(c *gin.Context, owner string) bool {
	if owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return false
	}
	return true
}

// requireLogOwner fetches the log by ID and checks that the authenticated user owns it.
func requireLogOwner(c *gin.Context, logID uint) bool {
	var log models.WorkoutLogEntity
	if err := database.DB.Select("owner").First(&log, logID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log not found"})
		return false
	}
	return requireOwner(c, log.Owner)
}

// requireLogStatus fetches the log by ID, checks ownership, and checks its status is in the allowed set.
// Returns the log entity and true if OK, otherwise writes an error response and returns false.
func requireLogStatus(c *gin.Context, logID uint, allowed ...models.WorkoutLogStatus) (models.WorkoutLogEntity, bool) {
	var log models.WorkoutLogEntity
	if err := database.DB.First(&log, logID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout log not found"})
		return log, false
	}
	if !requireOwner(c, log.Owner) {
		return log, false
	}
	for _, s := range allowed {
		if log.Status == s {
			return log, true
		}
	}
	c.JSON(http.StatusConflict, gin.H{"error": "Workout log status is " + string(log.Status) + ", operation not allowed"})
	return log, false
}

// getLogIDFromSection returns the parent log ID for a given section ID.
func getLogIDFromSection(sectionID uint) (uint, error) {
	var section models.WorkoutLogSectionEntity
	if err := database.DB.Select("workout_log_id").First(&section, sectionID).Error; err != nil {
		return 0, err
	}
	return section.WorkoutLogID, nil
}

// getLogIDFromExercise returns the parent log ID for a given exercise ID
// by traversing exercise → section → log.
func getLogIDFromExercise(exerciseID uint) (uint, error) {
	var exercise models.WorkoutLogExerciseEntity
	if err := database.DB.Select("workout_log_section_id").First(&exercise, exerciseID).Error; err != nil {
		return 0, err
	}
	return getLogIDFromSection(exercise.WorkoutLogSectionID)
}
