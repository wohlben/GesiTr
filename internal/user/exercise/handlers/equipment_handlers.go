package handlers

import (
	"net/http"

	"gesitr/internal/database"
	"gesitr/internal/user/exercise/models"

	"github.com/gin-gonic/gin"
)

func ListUserEquipment(c *gin.Context) {
	db := database.DB.Model(&models.UserEquipmentEntity{})

	if v := c.Query("owner"); v != "" {
		db = db.Where("owner = ?", v)
	}
	if v := c.Query("compendiumEquipmentId"); v != "" {
		db = db.Where("compendium_equipment_id = ?", v)
	}

	var entities []models.UserEquipmentEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]models.UserEquipment, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateUserEquipment(c *gin.Context) {
	var dto models.UserEquipment
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := models.UserEquipmentFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetUserEquipment(c *gin.Context) {
	var entity models.UserEquipmentEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User equipment not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func DeleteUserEquipment(c *gin.Context) {
	if err := database.DB.Delete(&models.UserEquipmentEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User equipment not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
