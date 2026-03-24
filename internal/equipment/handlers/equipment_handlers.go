package handlers

import (
	"net/http"
	"strconv"
	"time"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/equipment/models"
	"gesitr/internal/shared"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListEquipment returns equipment visible to the current user: their own
// equipment plus all public equipment. Filter by owner, public, or category
// query params. GET /api/equipment
func ListEquipment(c *gin.Context) {
	db := database.DB.Model(&models.EquipmentEntity{})

	userID := auth.GetUserID(c)
	if owner := c.Query("owner"); owner != "" {
		if owner == "me" || owner == userID {
			db = db.Where("owner = ?", userID)
		} else {
			db = db.Where("owner = ? AND public = ?", owner, true)
		}
	} else {
		db = db.Where("owner = ? OR public = ?", userID, true)
	}
	if pub := c.Query("public"); pub == "true" {
		db = db.Where("public = ?", true)
	}

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

	p := shared.ParsePagination(c)
	var entities []models.EquipmentEntity
	if err := shared.ApplyPagination(db, p).Find(&entities).Error; err != nil {
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

// CreateEquipment creates equipment owned by the current user. Equipment can
// be referenced by exercises via their equipmentIds field — see
// [gesitr/internal/exercise/handlers.CreateExercise]. POST /api/equipment
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
	entity.Owner = auth.GetUserID(c)
	if err := database.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resultDTO := entity.ToDTO()
	database.DB.Create(&models.EquipmentHistoryEntity{
		EquipmentID: entity.ID,
		Version:     resultDTO.Version,
		Snapshot:    shared.SnapshotJSON(resultDTO),
		ChangedAt:   time.Now(),
		ChangedBy:   resultDTO.Owner,
	})
	c.JSON(http.StatusCreated, resultDTO)
}

// GetEquipmentPermissions returns the current user's permissions on equipment.
// See [gesitr/internal/shared.ResolvePermissions] for the permission model.
// GET /api/equipment/:id/permissions
func GetEquipmentPermissions(c *gin.Context) {
	var entity models.EquipmentEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}
	userID := auth.GetUserID(c)
	perms, _ := shared.ResolvePermissions(userID, entity.Owner, entity.Public)
	if perms == nil {
		perms = []shared.Permission{}
	}
	c.JSON(http.StatusOK, shared.PermissionsResponse{Permissions: perms})
}

// GetEquipment returns a single equipment item. Public equipment is visible
// to all users; private equipment is visible only to its owner.
// GET /api/equipment/:id
func GetEquipment(c *gin.Context) {
	var entity models.EquipmentEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}
	userID := auth.GetUserID(c)
	perms, _ := shared.ResolvePermissions(userID, entity.Owner, entity.Public)
	if len(perms) == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

// UpdateEquipment updates equipment. Creates a version history entry.
// Owner only. PUT /api/equipment/:id
func UpdateEquipment(c *gin.Context) {
	var existing models.EquipmentEntity
	if err := database.DB.First(&existing, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}

	userID := auth.GetUserID(c)
	if existing.Owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not the owner"})
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
	entity.Owner = existing.Owner
	entity.Version = existing.Version + 1

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&entity).Error; err != nil {
			return err
		}
		resultDTO := entity.ToDTO()
		return tx.Create(&models.EquipmentHistoryEntity{
			EquipmentID: entity.ID,
			Version:     resultDTO.Version,
			Snapshot:    shared.SnapshotJSON(resultDTO),
			ChangedAt:   time.Now(),
			ChangedBy:   resultDTO.Owner,
		}).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.ToDTO())
}

// ListEquipmentVersions returns the version history for equipment. Each
// update via [UpdateEquipment] creates a new version entry.
// GET /api/equipment/:id/versions
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

	entries := make([]shared.VersionEntry, len(history))
	for i := range history {
		entries[i] = history[i].ToVersionEntry()
	}
	c.JSON(http.StatusOK, entries)
}

// GetEquipmentVersion returns a specific historical version of equipment
// by templateId and version number.
// GET /api/equipment/templates/:templateId/versions/:version
func GetEquipmentVersion(c *gin.Context) {
	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
		return
	}

	var entity models.EquipmentEntity
	if err := database.DB.Where("template_id = ?", c.Param("templateId")).First(&entity).Error; err != nil {
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

// DeleteEquipment deletes equipment. Owner only.
// DELETE /api/equipment/:id
func DeleteEquipment(c *gin.Context) {
	var entity models.EquipmentEntity
	if err := database.DB.First(&entity, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}

	userID := auth.GetUserID(c)
	if entity.Owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not the owner"})
		return
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
