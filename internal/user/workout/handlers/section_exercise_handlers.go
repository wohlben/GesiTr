package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	exercisemodels "gesitr/internal/exercise/models"
	"gesitr/internal/user/workout/models"

	"github.com/gin-gonic/gin"
)

// requireSectionOwner fetches the section, then checks workout ownership.
func requireSectionOwner(c *gin.Context, sectionID uint) bool {
	var section models.WorkoutSectionEntity
	if err := database.DB.First(&section, sectionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout section not found"})
		return false
	}
	return requireWorkoutOwner(c, section.WorkoutID)
}

// ListWorkoutSectionExercises returns section exercises owned by the current
// user. Filter by workoutSectionId query param.
// GET /api/user/workout-section-exercises
func ListWorkoutSectionExercises(c *gin.Context) {
	db := database.DB.Model(&models.WorkoutSectionExerciseEntity{})

	if v := c.Query("workoutSectionId"); v != "" {
		db = db.Where("workout_section_id = ?", v)
	}

	// Enforce ownership through workout section -> workout
	db = db.Where("workout_section_id IN (SELECT id FROM workout_sections WHERE workout_id IN (SELECT id FROM workouts WHERE owner = ? AND deleted_at IS NULL))", auth.GetUserID(c))

	var entities []models.WorkoutSectionExerciseEntity
	if err := db.Order("position").Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.WorkoutSectionExercise, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

// CreateWorkoutSectionExercise adds an exercise to a section. Requires a
// workoutSectionId (whose parent workout must be owned by the current user)
// and an exerciseSchemeId referencing an existing exercise scheme — see
// [gesitr/internal/exercise/handlers.CreateExerciseScheme].
// A workout and section must exist first — see [CreateWorkout] and
// [CreateWorkoutSection]. POST /api/user/workout-section-exercises
func CreateWorkoutSectionExercise(c *gin.Context) {
	var dto models.WorkoutSectionExercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !requireSectionOwner(c, dto.WorkoutSectionID) {
		return
	}

	var scheme exercisemodels.ExerciseSchemeEntity
	if err := database.DB.First(&scheme, dto.ExerciseSchemeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise scheme not found"})
		return
	}

	entity := models.WorkoutSectionExerciseFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

// DeleteWorkoutSectionExercise removes an exercise from a section.
// DELETE /api/user/workout-section-exercises/:id
func DeleteWorkoutSectionExercise(c *gin.Context) {
	var entity models.WorkoutSectionExerciseEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout section exercise not found"})
		return
	}
	if !requireSectionOwner(c, entity.WorkoutSectionID) {
		return
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
