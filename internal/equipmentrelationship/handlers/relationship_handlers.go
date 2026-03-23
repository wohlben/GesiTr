package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/equipmentrelationship/models"

	"github.com/gin-gonic/gin"
)

func ListEquipmentRelationships(c *gin.Context) {
	db := database.DB.Model(&models.EquipmentRelationshipEntity{})

	if v := c.Query("owner"); v != "" {
		db = db.Where("owner = ?", v)
	}
	if v := c.Query("fromEquipmentId"); v != "" {
		db = db.Where("from_equipment_id = ?", v)
	}
	if v := c.Query("toEquipmentId"); v != "" {
		db = db.Where("to_equipment_id = ?", v)
	}
	if v := c.Query("relationshipType"); v != "" {
		db = db.Where("relationship_type = ?", v)
	}

	var entities []models.EquipmentRelationshipEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.EquipmentRelationship, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateEquipmentRelationship(c *gin.Context) {
	var dto models.EquipmentRelationship
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.EquipmentRelationshipFromDTO(dto)
	entity.Owner = auth.GetUserID(c)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func DeleteEquipmentRelationship(c *gin.Context) {
	var entity models.EquipmentRelationshipEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "EquipmentRelationship not found"})
		return
	}

	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not the owner of this relationship"})
		return
	}

	if err := database.DB.Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
