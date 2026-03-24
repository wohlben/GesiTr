package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/user/workout/models"

	"github.com/gin-gonic/gin"
)

// requireWorkoutOwner fetches the workout by ID and checks ownership.
func requireWorkoutOwner(c *gin.Context, workoutID uint) bool {
	var workout models.WorkoutEntity
	if err := database.DB.First(&workout, workoutID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return false
	}
	if workout.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return false
	}
	return true
}

// ListWorkoutSections returns sections owned by the current user. Filter by
// workoutId query param to get sections for a specific workout.
// GET /api/user/workout-sections
func ListWorkoutSections(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutSectionEntity{})

	if v := c.Query("workoutId"); v != "" {
		db = db.Where("workout_id = ?", v)
	}

	// Join through workout to enforce ownership
	db = db.Where("workout_id IN (SELECT id FROM workouts WHERE owner = ? AND deleted_at IS NULL)", auth.GetUserID(c))

	var entities []models.WorkoutSectionEntity
	if err := db.Preload("Exercises").Order("position").Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.WorkoutSection, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

// CreateWorkoutSection adds a section to a workout. Requires a workoutId
// referencing a workout owned by the current user. A workout must exist
// before sections can be added — see [CreateWorkout].
// POST /api/user/workout-sections
func CreateWorkoutSection(c *gin.Context) {
	var dto models.WorkoutSection
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !requireWorkoutOwner(c, dto.WorkoutID) {
		return
	}

	entity := models.WorkoutSectionFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

// GetWorkoutSection returns a single section with its exercises.
// GET /api/user/workout-sections/:id
func GetWorkoutSection(c *gin.Context) {
	var entity models.WorkoutSectionEntity
	if err := database.DB.Preload("Exercises").First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout section not found"})
		return
	}
	if !requireWorkoutOwner(c, entity.WorkoutID) {
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

// DeleteWorkoutSection removes a section from its workout.
// DELETE /api/user/workout-sections/:id
func DeleteWorkoutSection(c *gin.Context) {
	var entity models.WorkoutSectionEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout section not found"})
		return
	}
	if !requireWorkoutOwner(c, entity.WorkoutID) {
		return
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
