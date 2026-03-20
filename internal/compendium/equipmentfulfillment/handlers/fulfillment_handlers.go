package handlers

import (
	"net/http"

	"gesitr/internal/compendium/equipmentfulfillment/models"
	"gesitr/internal/database"

	"github.com/gin-gonic/gin"
)

func ListFulfillments(c *gin.Context) {
	db := database.DB.Model(&models.FulfillmentEntity{})

	if v := c.Query("equipmentTemplateId"); v != "" {
		db = db.Where("equipment_template_id = ?", v)
	}
	if v := c.Query("fulfillsEquipmentTemplateId"); v != "" {
		db = db.Where("fulfills_equipment_template_id = ?", v)
	}

	var entities []models.FulfillmentEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.Fulfillment, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateFulfillment(c *gin.Context) {
	var dto models.Fulfillment
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.FulfillmentFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func DeleteFulfillment(c *gin.Context) {
	if err := database.DB.Delete(&models.FulfillmentEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fulfillment not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
