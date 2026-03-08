package handlers

import (
	"net/http"
	"strconv"
	"time"

	"gesitr/internal/compendium/models"
	"gesitr/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

	if dto.TemplateID == "" {
		dto.TemplateID = uuid.New().String()
	}

	entity := models.EquipmentFromDTO(dto)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resultDTO := entity.ToDTO()
	database.DB.Create(&models.EquipmentHistoryEntity{
		EquipmentID: entity.ID,
		Version:     resultDTO.Version,
		Snapshot:    models.SnapshotJSON(resultDTO),
		ChangedAt:   time.Now(),
		ChangedBy:   resultDTO.CreatedBy,
	})
	c.JSON(http.StatusCreated, resultDTO)
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

	oldDTO := existing.ToDTO()

	if !models.EquipmentChanged(oldDTO, dto) {
		c.JSON(http.StatusOK, oldDTO)
		return
	}

	entity := models.EquipmentFromDTO(dto)
	entity.ID = existing.ID
	entity.Version = existing.Version + 1

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&entity).Error; err != nil {
			return err
		}
		resultDTO := entity.ToDTO()
		return tx.Create(&models.EquipmentHistoryEntity{
			EquipmentID: entity.ID,
			Version:     resultDTO.Version,
			Snapshot:    models.SnapshotJSON(resultDTO),
			ChangedAt:   time.Now(),
			ChangedBy:   resultDTO.CreatedBy,
		}).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

func ListEquipmentVersions(c *gin.Context) {
	var entity models.EquipmentEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}

	var history []models.EquipmentHistoryEntity
	if err := database.DB.Where("equipment_id = ?", entity.ID).Order("version ASC").Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entries := make([]models.VersionEntry, len(history))
	for i := range history {
		entries[i] = history[i].ToVersionEntry()
	}
	c.JSON(http.StatusOK, entries)
}

func GetEquipmentVersion(c *gin.Context) {
	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
		return
	}

	var entity models.EquipmentEntity
	if err := database.DB.Unscoped().Where("template_id = ?", c.Param("templateId")).First(&entity).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}

	var history models.EquipmentHistoryEntity
	if err := database.DB.Where("equipment_id = ? AND version = ?", entity.ID, version).First(&history).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	c.JSON(http.StatusOK, history.ToVersionEntry())
}

func DeleteEquipment(c *gin.Context) {
	if err := database.DB.Delete(&models.EquipmentEntity{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
