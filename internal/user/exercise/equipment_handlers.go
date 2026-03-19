package exercise

import (
	"net/http"

	"gesitr/internal/database"

	"github.com/gin-gonic/gin"
)

func ListUserEquipment(c *gin.Context) {
	db := database.DB.Model(&UserEquipmentEntity{})

	if v := c.Query("owner"); v != "" {
		db = db.Where("owner = ?", v)
	}
	if v := c.Query("compendiumEquipmentId"); v != "" {
		db = db.Where("compendium_equipment_id = ?", v)
	}

	var entities []UserEquipmentEntity
	if err := db.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]UserEquipment, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	c.JSON(http.StatusOK, dtos)
}

func CreateUserEquipment(c *gin.Context) {
	var dto UserEquipment
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := UserEquipmentFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entity.ToDTO())
}

func GetUserEquipment(c *gin.Context) {
	var entity UserEquipmentEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User equipment not found"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func DeleteUserEquipment(c *gin.Context) {
	if err := database.DB.Delete(&UserEquipmentEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User equipment not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
