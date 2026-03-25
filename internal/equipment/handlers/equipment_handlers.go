package handlers

import (
	"context"
	"encoding/json"
	"time"

	"gesitr/internal/database"
	"gesitr/internal/equipment/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func equipmentDTOFromBody(body EquipmentBody) models.Equipment {
	return models.Equipment{
		Name:        body.Name,
		DisplayName: body.DisplayName,
		Description: body.Description,
		Category:    body.Category,
		ImageUrl:    body.ImageUrl,
		TemplateID:  body.TemplateID,
		Public:      body.Public,
	}
}

// ListEquipment returns equipment visible to the current user: their own
// equipment plus all public equipment. Filter by owner, public, or category
// query params. GET /api/equipment
//
// OpenAPI: /api/docs#/operations/list-equipment
func ListEquipment(ctx context.Context, input *ListEquipmentInput) (*ListEquipmentOutput, error) {
	db := database.DB.Model(&models.EquipmentEntity{})

	userID := humaconfig.GetUserID(ctx)
	if input.Owner != "" {
		if input.Owner == "me" || input.Owner == userID {
			db = db.Where("owner = ?", userID)
		} else {
			db = db.Where("owner = ? AND public = ?", input.Owner, true)
		}
	} else {
		db = db.Where("owner = ? OR public = ?", userID, true)
	}
	if input.Public == "true" {
		db = db.Where("public = ?", true)
	}

	if input.Q != "" {
		pattern := "%" + input.Q + "%"
		db = db.Where("name LIKE ? OR display_name LIKE ?", pattern, pattern)
	}
	if input.Category != "" {
		db = db.Where("category = ?", input.Category)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	p := input.ToPaginationParams()
	var entities []models.EquipmentEntity
	if err := shared.ApplyPagination(db, p).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.Equipment, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListEquipmentOutput{Body: humaconfig.PaginatedBody[models.Equipment]{
		Items: dtos, Total: total, Limit: p.Limit, Offset: p.Offset,
	}}, nil
}

// CreateEquipment creates equipment owned by the current user. Equipment can
// be referenced by exercises via their equipmentIds field — see
// [gesitr/internal/exercise/handlers.CreateExercise]. POST /api/equipment
//
// OpenAPI: /api/docs#/operations/create-equipment
func CreateEquipment(ctx context.Context, input *CreateEquipmentInput) (*CreateEquipmentOutput, error) {
	dto := equipmentDTOFromBody(input.Body)

	if dto.TemplateID == "" {
		dto.TemplateID = uuid.New().String()
	}

	entity := models.EquipmentFromDTO(dto)
	entity.Owner = humaconfig.GetUserID(ctx)
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resultDTO := entity.ToDTO()
	database.DB.Create(&models.EquipmentHistoryEntity{
		EquipmentID: entity.ID,
		Version:     resultDTO.Version,
		Snapshot:    shared.SnapshotJSON(resultDTO),
		ChangedAt:   time.Now(),
		ChangedBy:   resultDTO.Owner,
	})
	return &CreateEquipmentOutput{Body: resultDTO}, nil
}

// GetEquipmentPermissions returns the current user's permissions on equipment.
// See [gesitr/internal/shared.ResolvePermissions] for the permission model.
// GET /api/equipment/:id/permissions
//
// OpenAPI: /api/docs#/operations/get-equipment-permissions
func GetEquipmentPermissions(ctx context.Context, input *GetEquipmentPermissionsInput) (*GetEquipmentPermissionsOutput, error) {
	var entity models.EquipmentEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Equipment not found")
	}
	userID := humaconfig.GetUserID(ctx)
	perms, _ := shared.ResolvePermissions(userID, entity.Owner, entity.Public)
	if perms == nil {
		perms = []shared.Permission{}
	}
	return &GetEquipmentPermissionsOutput{Body: shared.PermissionsResponse{Permissions: perms}}, nil
}

// GetEquipment returns a single equipment item. Public equipment is visible
// to all users; private equipment is visible only to its owner.
// GET /api/equipment/:id
//
// OpenAPI: /api/docs#/operations/get-equipment
func GetEquipment(ctx context.Context, input *GetEquipmentInput) (*GetEquipmentOutput, error) {
	var entity models.EquipmentEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Equipment not found")
	}
	userID := humaconfig.GetUserID(ctx)
	perms, _ := shared.ResolvePermissions(userID, entity.Owner, entity.Public)
	if len(perms) == 0 {
		return nil, huma.Error403Forbidden("access denied")
	}
	return &GetEquipmentOutput{Body: entity.ToDTO()}, nil
}

// UpdateEquipment updates equipment. Creates a version history entry.
// Owner only — returns 403 for non-owners. PUT /api/equipment/:id
//
// OpenAPI: /api/docs#/operations/update-equipment
func UpdateEquipment(ctx context.Context, input *UpdateEquipmentInput) (*UpdateEquipmentOutput, error) {
	var existing models.EquipmentEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Equipment not found")
	}

	if existing.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}

	dto := equipmentDTOFromBody(input.Body)
	dto.Owner = existing.Owner

	oldDTO := existing.ToDTO()

	if !models.EquipmentChanged(oldDTO, dto) {
		return &UpdateEquipmentOutput{Body: oldDTO}, nil
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
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateEquipmentOutput{Body: entity.ToDTO()}, nil
}

// ListEquipmentVersions returns the version history for equipment. Each
// update via [UpdateEquipment] creates a new version entry.
// GET /api/equipment/:id/versions
//
// OpenAPI: /api/docs#/operations/list-equipment-versions
func ListEquipmentVersions(ctx context.Context, input *ListEquipmentVersionsInput) (*ListEquipmentVersionsOutput, error) {
	var entity models.EquipmentEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Equipment not found")
	}

	var history []models.EquipmentHistoryEntity
	if err := database.DB.Where("equipment_id = ?", entity.ID).Order("version ASC").Find(&history).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	entries := make([]shared.VersionEntry, len(history))
	for i := range history {
		entries[i] = history[i].ToVersionEntry()
	}
	return &ListEquipmentVersionsOutput{Body: entries}, nil
}

// GetEquipmentVersion returns a specific historical version of equipment
// by templateId and version number.
// GET /api/equipment/templates/:templateId/versions/:version
//
// OpenAPI: /api/docs#/operations/get-equipment-version
func GetEquipmentVersion(ctx context.Context, input *GetEquipmentVersionInput) (*GetEquipmentVersionOutput, error) {
	var history models.EquipmentHistoryEntity
	if err := database.DB.Where("json_extract(snapshot, '$.templateId') = ? AND version = ?", input.TemplateID, input.Version).First(&history).Error; err != nil {
		return nil, huma.Error404NotFound("Version not found")
	}

	var snap models.Equipment
	json.Unmarshal([]byte(history.Snapshot), &snap)
	userID := humaconfig.GetUserID(ctx)
	perms, _ := shared.ResolvePermissions(userID, snap.Owner, snap.Public)
	if len(perms) == 0 {
		return nil, huma.Error403Forbidden("access denied")
	}

	return &GetEquipmentVersionOutput{Body: history.ToVersionEntry()}, nil
}

// DeleteEquipment deletes equipment. Owner only.
// DELETE /api/equipment/:id
//
// OpenAPI: /api/docs#/operations/delete-equipment
func DeleteEquipment(ctx context.Context, input *DeleteEquipmentInput) (*DeleteEquipmentOutput, error) {
	var entity models.EquipmentEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Equipment not found")
	}
	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}
	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
