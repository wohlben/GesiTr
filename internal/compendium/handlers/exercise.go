package handlers

import (
	"net/http"
	"strconv"
	"time"

	"gesitr/internal/compendium/models"
	"gesitr/internal/database"
	"gesitr/internal/slug"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func ListExercises(c *gin.Context) {
	db := database.DB.Model(&models.ExerciseEntity{})

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

	p := parsePagination(c)
	var entities []models.ExerciseEntity
	if err := preloadExercise(applyPagination(db, p)).Find(&entities).Error; err != nil {
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

func CreateExercise(c *gin.Context) {
	var dto models.Exercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auto-generate slug from name if not provided
	if dto.Slug == "" {
		dto.Slug = slug.Generate(dto.Name)
	}

	// Default templateId to a UUID if not provided
	if dto.TemplateID == "" {
		dto.TemplateID = uuid.New().String()
	}

	entity := models.ExerciseFromDTO(dto)
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
		Snapshot:   models.SnapshotJSON(resultDTO),
		ChangedAt:  time.Now(),
		ChangedBy:  resultDTO.CreatedBy,
	})
	c.JSON(http.StatusCreated, resultDTO)
}

func GetExercise(c *gin.Context) {
	var entity models.ExerciseEntity
	if err := preloadExercise(database.DB).First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateExercise(c *gin.Context) {
	var existing models.ExerciseEntity
	if err := preloadExercise(database.DB).First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
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
			Snapshot:   models.SnapshotJSON(resultDTO),
			ChangedAt:  time.Now(),
			ChangedBy:  resultDTO.CreatedBy,
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

	entries := make([]models.VersionEntry, len(history))
	for i := range history {
		entries[i] = history[i].ToVersionEntry()
	}
	c.JSON(http.StatusOK, entries)
}

func GetExerciseVersion(c *gin.Context) {
	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
		return
	}

	var entity models.ExerciseEntity
	if err := database.DB.Unscoped().Where("template_id = ?", c.Param("templateId")).First(&entity).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	var history models.ExerciseHistoryEntity
	if err := database.DB.Where("exercise_id = ? AND version = ?", entity.ID, version).First(&history).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	c.JSON(http.StatusOK, history.ToVersionEntry())
}

func DeleteExercise(c *gin.Context) {
	if err := database.DB.Delete(&models.ExerciseEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
