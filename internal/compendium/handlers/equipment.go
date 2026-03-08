package handlers

import (
	"net/http"

	"gesitr/internal/compendium/models"
	"gesitr/internal/database"

	"github.com/gin-gonic/gin"
)

func ListEquipment(c *gin.Context) {
	db := database.DB.Model(&models.EquipmentEntity{})

	if q := c.Query("q"); q != "" {
		pattern := "%" + q + "%"
		db = db.Where("name LIKE ? OR display_name LIKE ?", pattern, pattern)
	}
	if v := c.Query("category"); v != "" {
		db = db.Where("category = ?", v)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p := parsePagination(c)
	var entities []models.EquipmentEntity
	if err := applyPagination(db, p).Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.Equipment, len(entities))
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

func CreateEquipment(c *gin.Context) {
	var dto models.Equipment
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.EquipmentFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetEquipment(c *gin.Context) {
	var entity models.EquipmentEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateEquipment(c *gin.Context) {
	var existing models.EquipmentEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}

	var dto models.Equipment
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.EquipmentFromDTO(dto)
	entity.ID = existing.ID
	entity.Version = existing.Version

	if err := database.DB.Save(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func DeleteEquipment(c *gin.Context) {
	if err := database.DB.Delete(&models.EquipmentEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
