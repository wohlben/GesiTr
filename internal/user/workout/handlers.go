package workout

import (
	"net/http"

	"gesitr/internal/database"

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

func ListWorkouts(c *gin.Context) {
	db := database.DB.Model(&WorkoutEntity{})

	if v := c.Query("owner"); v != "" {
		db = db.Where("owner = ?", v)
	}

	var entities []WorkoutEntity
	if err := preloadWorkout(db).Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]Workout, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateWorkout(c *gin.Context) {
	var dto Workout
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := WorkoutFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetWorkout(c *gin.Context) {
	var entity WorkoutEntity
	if err := preloadWorkout(database.DB).First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateWorkout(c *gin.Context) {
	var existing WorkoutEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return
	}

	var dto Workout
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := WorkoutFromDTO(dto)
	entity.ID = existing.ID

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

func DeleteWorkout(c *gin.Context) {
	if err := database.DB.Delete(&WorkoutEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
