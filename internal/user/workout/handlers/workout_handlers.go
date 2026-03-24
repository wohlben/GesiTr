package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/user/workout/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func preloadWorkout(db *gorm.DB) *gorm.DB {
	return db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Sections.Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	})
}

// ListWorkouts returns all workouts owned by the current user, each
// including its sections and section exercises. GET /api/user/workouts
func ListWorkouts(c *gin.Context) {
	userID := auth.GetUserID(c)
	db := database.DB.Model(&models.WorkoutEntity{}).Where("owner = ?", userID)

	var entities []models.WorkoutEntity
	if err := preloadWorkout(db).Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.Workout, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

// CreateWorkout creates an empty workout. Add sections via
// [CreateWorkoutSection] and exercises via [CreateWorkoutSectionExercise].
// POST /api/user/workouts
func CreateWorkout(c *gin.Context) {
	var dto models.Workout
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.WorkoutFromDTO(dto)
	entity.Owner = auth.GetUserID(c)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

// GetWorkout returns a workout with its full section and exercise tree.
// Returns 403 if the caller is not the owner. GET /api/user/workouts/:id
func GetWorkout(c *gin.Context) {
	var entity models.WorkoutEntity
	if err := preloadWorkout(database.DB).First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return
	}
	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

// UpdateWorkout updates workout metadata (name, notes). Sections and exercises
// are managed via their own endpoints. PUT /api/user/workouts/:id
func UpdateWorkout(c *gin.Context) {
	var existing models.WorkoutEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return
	}
	if existing.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var dto models.Workout
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.WorkoutFromDTO(dto)
	entity.ID = existing.ID
	entity.Owner = existing.Owner

	if err := database.DB.Save(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := preloadWorkout(database.DB).First(&entity, entity.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

// DeleteWorkout deletes a workout. DELETE /api/user/workouts/:id
func DeleteWorkout(c *gin.Context) {
	var entity models.WorkoutEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return
	}
	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
