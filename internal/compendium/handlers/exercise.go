package handlers

import (
	"net/http"

	"gesitr/internal/database"
	"gesitr/internal/compendium/models"

	"github.com/gin-gonic/gin"
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

	var entities []models.ExerciseEntity
	if err := preloadExercise(db).Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.Exercise, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateExercise(c *gin.Context) {
	var dto models.Exercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
	c.JSON(http.StatusCreated, entity.ToDTO())
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
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	var dto models.Exercise
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.ExerciseFromDTO(dto)
	entity.ID = existing.ID
	entity.Version = existing.Version

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

		// Save entity (triggers BeforeUpdate for version increment)
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

func DeleteExercise(c *gin.Context) {
	if err := database.DB.Delete(&models.ExerciseEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
