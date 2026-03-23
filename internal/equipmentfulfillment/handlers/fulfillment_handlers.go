package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/equipmentfulfillment/models"

	"github.com/gin-gonic/gin"
)

func ListFulfillments(c *gin.Context) {
	db := database.DB.Model(&models.FulfillmentEntity{})

	if v := c.Query("equipmentId"); v != "" {
		db = db.Where("equipment_id = ?", v)
	}
	if v := c.Query("fulfillsEquipmentId"); v != "" {
		db = db.Where("fulfills_equipment_id = ?", v)
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
	entity.Owner = auth.GetUserID(c)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func DeleteFulfillment(c *gin.Context) {
	var entity models.FulfillmentEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fulfillment not found"})
		return
	}

	if entity.Owner != auth.GetUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not the owner of this fulfillment"})
		return
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
