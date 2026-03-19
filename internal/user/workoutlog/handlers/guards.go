package handlers

import (
	"net/http"
	"strconv"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/user/workoutlog/models"

	"github.com/gin-gonic/gin"
)

// parseUint converts a string to uint, returning 0 on failure.
func parseUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}

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
