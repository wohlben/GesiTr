package handlers

import (
	"context"

	"gesitr/internal/compendium/equipment/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/ownershipgroup"
	masteryHandlers "gesitr/internal/user/mastery/handlers"

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
		userID := humaconfig.GetUserID(ctx)
		if input.Owner == "me" || input.Owner == userID {
			visibleGroups := ownershipgroup.VisibleGroupIDs(database.DB, userID)
			db = db.Where("ownership_group_id IN (?)", visibleGroups)
		} else {
			db = db.Where("ownership_group_id IN (SELECT group_id FROM ownership_group_memberships WHERE user_id = ? AND role = 'owner' AND deleted_at IS NULL)", input.Owner)
		}
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
	userID := humaconfig.GetUserID(ctx)

	// Inherit ownership group from the "from" equipment.
	var fromEquipment models.EquipmentEntity
	if err := database.DB.Select("ownership_group_id").First(&fromEquipment, input.Body.FromEquipmentID).Error; err != nil {
		return nil, huma.Error404NotFound("From equipment not found")
	}
	entity.OwnershipGroupID = fromEquipment.OwnershipGroupID

	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	_ = masteryHandlers.RecalculateEquipmentContributions(database.DB, userID, entity.FromEquipmentID, entity.ToEquipmentID)
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

	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, entity.OwnershipGroupID)
	if !access.CanModify() {
		return nil, huma.Error403Forbidden("not the owner of this relationship")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	_ = masteryHandlers.RecalculateEquipmentContributions(database.DB, userID, entity.FromEquipmentID, entity.ToEquipmentID)
	return nil, nil
}
