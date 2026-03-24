package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/exercisegroup/models"
	"gesitr/internal/shared"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListExerciseGroups(c *gin.Context) {
	db := database.DB.Model(&models.ExerciseGroupEntity{})

	if q := c.Query("q"); q != "" {
		pattern := "%" + q + "%"
		db = db.Where("name LIKE ?", pattern)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p := shared.ParsePagination(c)
	var entities []models.ExerciseGroupEntity
	if err := shared.ApplyPagination(db, p).Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.ExerciseGroup, len(entities))
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

func CreateExerciseGroup(c *gin.Context) {
	var dto models.ExerciseGroup
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dto.TemplateID == "" {
		dto.TemplateID = uuid.New().String()
	}

	entity := models.ExerciseGroupFromDTO(dto)
	entity.Owner = auth.GetUserID(c)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetExerciseGroupPermissions(c *gin.Context) {
	var entity models.ExerciseGroupEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ExerciseGroup not found"})
		return
	}
	userID := auth.GetUserID(c)
	perms, _ := shared.ResolvePermissions(userID, entity.Owner, false)
	if perms == nil {
		perms = []shared.Permission{}
	}
	c.JSON(http.StatusOK, shared.PermissionsResponse{Permissions: perms})
}

func GetExerciseGroup(c *gin.Context) {
	var entity models.ExerciseGroupEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ExerciseGroup not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateExerciseGroup(c *gin.Context) {
	var existing models.ExerciseGroupEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ExerciseGroup not found"})
		return
	}

	if existing.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not the owner of this exercise group"})
		return
	}

	var dto models.ExerciseGroup
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.ExerciseGroupFromDTO(dto)
	entity.ID = existing.ID
	entity.Owner = existing.Owner

	if err := database.DB.Save(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func DeleteExerciseGroup(c *gin.Context) {
	var entity models.ExerciseGroupEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ExerciseGroup not found"})
		return
	}

	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not the owner of this exercise group"})
		return
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
