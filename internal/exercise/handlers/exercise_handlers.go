package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/exercise/models"
	"gesitr/internal/shared"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var exercisePreloads = []string{
	"Forces", "Muscles", "Paradigms", "Instructions", "Images", "AlternativeNames", "Equipment",
}

func preloadExercise(db *gorm.DB) *gorm.DB {
	for _, p := range exercisePreloads {
		db = db.Preload(p)
	}
	return db
}

// ListExercises returns exercises visible to the current user: their own
// exercises plus all public exercises. Filter by owner or public query params.
// GET /api/exercises
func ListExercises(c *gin.Context) {
	db := database.DB.Model(&models.ExerciseEntity{})

	userID := auth.GetUserID(c)
	if owner := c.Query("owner"); owner != "" {
		if owner == "me" {
			db = db.Where("owner = ?", userID)
		} else if owner == userID {
			db = db.Where("owner = ?", userID)
		} else {
			db = db.Where("owner = ? AND public = ?", owner, true)
		}
	} else {
		db = db.Where("owner = ? OR public = ?", userID, true)
	}
	if pub := c.Query("public"); pub == "true" {
		db = db.Where("public = ?", true)
	}

	if q := c.Query("q"); q != "" {
		pattern := "%" + q + "%"
		db = db.Where(
			"exercises.id IN (SELECT exercises.id FROM exercises LEFT JOIN exercise_alternative_names ON exercise_alternative_names.exercise_id = exercises.id WHERE exercises.name LIKE ? OR exercise_alternative_names.name LIKE ?)",
			pattern, pattern,
		)
	}
	if v := c.Query("type"); v != "" {
		db = db.Where("exercises.type = ?", v)
	}
	if v := c.Query("difficulty"); v != "" {
		db = db.Where("exercises.technical_difficulty = ?", v)
	}
	if v := c.Query("force"); v != "" {
		db = db.Where("exercises.id IN (SELECT exercise_id FROM exercise_forces WHERE force = ?)", v)
	}
	if v := c.Query("muscle"); v != "" {
		db = db.Where("exercises.id IN (SELECT exercise_id FROM exercise_muscles WHERE muscle = ?)", v)
	}
	if v := c.Query("primaryMuscle"); v != "" {
		db = db.Where("exercises.id IN (SELECT exercise_id FROM exercise_muscles WHERE muscle = ? AND is_primary = true)", v)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p := shared.ParsePagination(c)
	var entities []models.ExerciseEntity
	if err := preloadExercise(shared.ApplyPagination(db, p)).Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.Exercise, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, gin.H{
		"items":  dtos,
		"total":  total,
		"limit":  p.Limit,
		"offset": p.Offset,
	})
}

// CreateExercise creates an exercise owned by the current user. The exercise
// can reference equipment via equipmentIds — equipment must already exist
// (see [gesitr/internal/equipment/handlers.CreateEquipment]). To use this
// exercise in a workout, create an exercise scheme via [CreateExerciseScheme].
// POST /api/exercises
func CreateExercise(c *gin.Context) {
	var dto models.Exercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default templateId to a UUID if not provided
	if dto.TemplateID == "" {
		dto.TemplateID = uuid.New().String()
	}

	entity := models.ExerciseFromDTO(dto)
	entity.Owner = auth.GetUserID(c)
	entity.Public = dto.Public
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with associations
	if err := preloadExercise(database.DB).First(&entity, entity.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resultDTO := entity.ToDTO()
	database.DB.Create(&models.ExerciseHistoryEntity{
		ExerciseID: entity.ID,
		Version:    resultDTO.Version,
		Snapshot:   shared.SnapshotJSON(resultDTO),
		ChangedAt:  time.Now(),
		ChangedBy:  resultDTO.Owner,
	})
	c.JSON(http.StatusCreated, resultDTO)
}

// GetExercisePermissions returns the current user's permissions on an exercise.
// See [gesitr/internal/shared.ResolvePermissions] for the permission model.
// GET /api/exercises/:id/permissions
func GetExercisePermissions(c *gin.Context) {
	var entity models.ExerciseEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	userID := auth.GetUserID(c)
	perms, _ := shared.ResolvePermissions(userID, entity.Owner, entity.Public)
	if perms == nil {
		perms = []shared.Permission{}
	}
	c.JSON(http.StatusOK, shared.PermissionsResponse{Permissions: perms})
}

// GetExercise returns a single exercise. Public exercises are visible to all
// users; private exercises are visible only to their owner.
// GET /api/exercises/:id
func GetExercise(c *gin.Context) {
	var entity models.ExerciseEntity
	if err := preloadExercise(database.DB).First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	userID := auth.GetUserID(c)
	perms, _ := shared.ResolvePermissions(userID, entity.Owner, entity.Public)
	if len(perms) == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

// UpdateExercise updates an exercise. Creates a version history entry.
// Owner only — returns 403 for non-owners. PUT /api/exercises/:id
func UpdateExercise(c *gin.Context) {
	var existing models.ExerciseEntity
	if err := preloadExercise(database.DB).First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	if existing.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var dto models.Exercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldDTO := existing.ToDTO()

	// If nothing changed, return the existing entity without bumping version
	if !models.ExerciseChanged(oldDTO, dto) {
		c.JSON(http.StatusOK, oldDTO)
		return
	}

	entity := models.ExerciseFromDTO(dto)
	entity.ID = existing.ID
	entity.Owner = existing.Owner
	entity.Version = existing.Version + 1

	// Stash child records and clear from entity so Save doesn't try to upsert them
	forces := entity.Forces
	muscles := entity.Muscles
	paradigms := entity.Paradigms
	instructions := entity.Instructions
	images := entity.Images
	altNames := entity.AlternativeNames
	equipment := entity.Equipment
	entity.Forces = nil
	entity.Muscles = nil
	entity.Paradigms = nil
	entity.Instructions = nil
	entity.Images = nil
	entity.AlternativeNames = nil
	entity.Equipment = nil

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Delete old child records
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseForce{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseMuscle{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseMeasurementParadigm{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseInstruction{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseImage{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseAlternativeName{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exercise_id = ?", entity.ID).Delete(&models.ExerciseEquipment{}).Error; err != nil {
			return err
		}

		// Save entity
		if err := tx.Save(&entity).Error; err != nil {
			return err
		}

		// Insert new child records
		for i := range forces {
			forces[i].ExerciseID = entity.ID
			if err := tx.Create(&forces[i]).Error; err != nil {
				return err
			}
		}
		for i := range muscles {
			muscles[i].ExerciseID = entity.ID
			if err := tx.Create(&muscles[i]).Error; err != nil {
				return err
			}
		}
		for i := range paradigms {
			paradigms[i].ExerciseID = entity.ID
			if err := tx.Create(&paradigms[i]).Error; err != nil {
				return err
			}
		}
		for i := range instructions {
			instructions[i].ExerciseID = entity.ID
			if err := tx.Create(&instructions[i]).Error; err != nil {
				return err
			}
		}
		for i := range images {
			images[i].ExerciseID = entity.ID
			if err := tx.Create(&images[i]).Error; err != nil {
				return err
			}
		}
		for i := range altNames {
			altNames[i].ExerciseID = entity.ID
			if err := tx.Create(&altNames[i]).Error; err != nil {
				return err
			}
		}
		for i := range equipment {
			equipment[i].ExerciseID = entity.ID
			if err := tx.Create(&equipment[i]).Error; err != nil {
				return err
			}
		}

		// Reconstruct entity for snapshot
		entity.Forces = forces
		entity.Muscles = muscles
		entity.Paradigms = paradigms
		entity.Instructions = instructions
		entity.Images = images
		entity.AlternativeNames = altNames
		entity.Equipment = equipment

		resultDTO := entity.ToDTO()
		if err := tx.Create(&models.ExerciseHistoryEntity{
			ExerciseID: entity.ID,
			Version:    resultDTO.Version,
			Snapshot:   shared.SnapshotJSON(resultDTO),
			ChangedAt:  time.Now(),
			ChangedBy:  resultDTO.Owner,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with associations
	if err := preloadExercise(database.DB).First(&entity, entity.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

// ListExerciseVersions returns the version history for an exercise. Each
// update via [UpdateExercise] creates a new version entry.
// GET /api/exercises/:id/versions
func ListExerciseVersions(c *gin.Context) {
	var entity models.ExerciseEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	var history []models.ExerciseHistoryEntity
	if err := database.DB.Where("exercise_id = ?", entity.ID).Order("version ASC").Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entries := make([]shared.VersionEntry, len(history))
	for i := range history {
		entries[i] = history[i].ToVersionEntry()
	}
	c.JSON(http.StatusOK, entries)
}

// GetExerciseVersion returns a specific historical version of an exercise
// by templateId and version number.
// GET /api/exercises/templates/:templateId/versions/:version
func GetExerciseVersion(c *gin.Context) {
	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
		return
	}

	// Query the history table directly — the snapshot contains the full
	// exercise, so this works even after the exercise has been deleted.
	var history models.ExerciseHistoryEntity
	if err := database.DB.Where("json_extract(snapshot, '$.templateId') = ? AND version = ?", c.Param("templateId"), version).First(&history).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	// Parse snapshot for permission check.
	var snap models.Exercise
	json.Unmarshal([]byte(history.Snapshot), &snap)
	userID := auth.GetUserID(c)
	perms, _ := shared.ResolvePermissions(userID, snap.Owner, snap.Public)
	if len(perms) == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, history.ToVersionEntry())
}

// DeleteExercise deletes an exercise. Owner only.
// DELETE /api/exercises/:id
func DeleteExercise(c *gin.Context) {
	var entity models.ExerciseEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	if err := database.DB.Unscoped().Select(clause.Associations).Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
