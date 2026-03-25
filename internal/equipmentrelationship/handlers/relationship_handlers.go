package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/equipmentrelationship/models"
	"gesitr/internal/humaconfig"

	"github.com/danielgtaylor/huma/v2"
)

// ListEquipmentRelationships returns equipment relationships, optionally filtered
// by owner, fromEquipmentId, toEquipmentId, or relationshipType.
// GET /api/equipment-relationships
//
// OpenAPI: /api/docs#/operations/ListEquipmentRelationships
func ListEquipmentRelationships(ctx context.Context, input *ListEquipmentRelationshipsInput) (*ListEquipmentRelationshipsOutput, error) {
	db := database.DB.Model(&models.EquipmentRelationshipEntity{})

	if input.Owner != "" {
		owner := input.Owner
		if owner == "me" {
			owner = humaconfig.GetUserID(ctx)
		}
		db = db.Where("owner = ?", owner)
	}
	if input.FromEquipmentID != "" {
		db = db.Where("from_equipment_id = ?", input.FromEquipmentID)
	}
	if input.ToEquipmentID != "" {
		db = db.Where("to_equipment_id = ?", input.ToEquipmentID)
	}
	if input.RelationshipType != "" {
		db = db.Where("relationship_type = ?", input.RelationshipType)
	}

	var entities []models.EquipmentRelationshipEntity
	if err := db.Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.EquipmentRelationship, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListEquipmentRelationshipsOutput{Body: dtos}, nil
}

// CreateEquipmentRelationship creates an equipment relationship owned by the current user.
// POST /api/equipment-relationships
//
// OpenAPI: /api/docs#/operations/CreateEquipmentRelationship
func CreateEquipmentRelationship(ctx context.Context, input *CreateEquipmentRelationshipInput) (*CreateEquipmentRelationshipOutput, error) {
	dto := models.EquipmentRelationship{
		RelationshipType: input.Body.RelationshipType,
		Strength:         input.Body.Strength,
		FromEquipmentID:  input.Body.FromEquipmentID,
		ToEquipmentID:    input.Body.ToEquipmentID,
	}

	entity := models.EquipmentRelationshipFromDTO(dto)
	entity.Owner = humaconfig.GetUserID(ctx)
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateEquipmentRelationshipOutput{Body: entity.ToDTO()}, nil
}

// DeleteEquipmentRelationship deletes an equipment relationship. Owner only.
// DELETE /api/equipment-relationships/:id
//
// OpenAPI: /api/docs#/operations/DeleteEquipmentRelationship
func DeleteEquipmentRelationship(ctx context.Context, input *DeleteEquipmentRelationshipInput) (*DeleteEquipmentRelationshipOutput, error) {
	var entity models.EquipmentRelationshipEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("EquipmentRelationship not found")
	}

	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("not the owner of this relationship")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
