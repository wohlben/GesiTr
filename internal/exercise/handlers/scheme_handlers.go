package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/exercise/models"
	"gesitr/internal/shared"

	"github.com/gin-gonic/gin"
)

// ListExerciseSchemes returns schemes the current user has access to: their
// own schemes plus schemes linked to public exercises. Filter by exerciseId
// or measurementType query params. GET /api/exercise-schemes
func ListExerciseSchemes(c *gin.Context) {
	userID := auth.GetUserID(c)
	// FIXME: subquery doesn't scale — replace with a join or denormalize visibility
	db := database.DB.Model(&models.ExerciseSchemeEntity{}).
		Where("owner = ? OR exercise_id IN (SELECT id FROM exercises WHERE public = ? AND deleted_at IS NULL)", userID, true)

	if v := c.Query("exerciseId"); v != "" {
		db = db.Where("exercise_id = ?", v)
	}
	if v := c.Query("measurementType"); v != "" {
		db = db.Where("measurement_type = ?", v)
	}

	var entities []models.ExerciseSchemeEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.ExerciseScheme, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

// CreateExerciseScheme creates an exercise scheme — a user-specific
// configuration of an exercise (sets, reps, measurement type). Requires
// an exerciseId referencing an existing exercise (see [CreateExercise]).
// Schemes are referenced when adding exercises to workouts via
// [gesitr/internal/user/workout/handlers.CreateWorkoutSectionExercise].
// POST /api/exercise-schemes
func CreateExerciseScheme(c *gin.Context) {
	var dto models.ExerciseScheme
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the exercise exists
	var exercise models.ExerciseEntity
	if err := database.DB.First(&exercise, dto.ExerciseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	entity := models.ExerciseSchemeFromDTO(dto)
	entity.Owner = auth.GetUserID(c)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

// GetExerciseScheme returns a single exercise scheme. Access is determined by
// the linked exercise's visibility — if the user can see the exercise, they
// can see its schemes. GET /api/exercise-schemes/:id
func GetExerciseScheme(c *gin.Context) {
	var entity models.ExerciseSchemeEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise scheme not found"})
		return
	}
	var exercise models.ExerciseEntity
	if err := database.DB.First(&exercise, entity.ExerciseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	userID := auth.GetUserID(c)
	perms, _ := shared.ResolvePermissions(userID, exercise.Owner, exercise.Public)
	if len(perms) == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

// UpdateExerciseScheme updates a scheme's configuration. The exerciseId
// cannot be changed. PUT /api/exercise-schemes/:id
func UpdateExerciseScheme(c *gin.Context) {
	var existing models.ExerciseSchemeEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise scheme not found"})
		return
	}
	if existing.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var dto models.ExerciseScheme
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.ExerciseSchemeFromDTO(dto)
	entity.ID = existing.ID
	entity.Owner = existing.Owner
	entity.ExerciseID = existing.ExerciseID

	if err := database.DB.Save(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

// DeleteExerciseScheme deletes an exercise scheme. Owner only.
// DELETE /api/exercise-schemes/:id
func DeleteExerciseScheme(c *gin.Context) {
	var entity models.ExerciseSchemeEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise scheme not found"})
		return
	}
	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
