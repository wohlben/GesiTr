package handlers

import (
	"context"

	"gesitr/internal/compendium/equipment/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/compendium/ownershipgroup"
	masteryHandlers "gesitr/internal/user/mastery/handlers"

	"github.com/danielgtaylor/huma/v2"
)

// ListFulfillments returns equipment fulfillments, optionally filtered by
// equipmentId or fulfillsEquipmentId.
// GET /api/fulfillments
//
// OpenAPI: /api/docs#/operations/ListFulfillments
func ListFulfillments(ctx context.Context, input *ListFulfillmentsInput) (*ListFulfillmentsOutput, error) {
	db := database.DB.Model(&models.FulfillmentEntity{})

	if input.EquipmentID != "" {
		db = db.Where("equipment_id = ?", input.EquipmentID)
	}
	if input.FulfillsEquipmentID != "" {
		db = db.Where("fulfills_equipment_id = ?", input.FulfillsEquipmentID)
	}

	var entities []models.FulfillmentEntity
	if err := db.Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.Fulfillment, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListFulfillmentsOutput{Body: dtos}, nil
}

// CreateFulfillment creates a fulfillment owned by the current user.
// POST /api/fulfillments
//
// OpenAPI: /api/docs#/operations/CreateFulfillment
func CreateFulfillment(ctx context.Context, input *CreateFulfillmentInput) (*CreateFulfillmentOutput, error) {
	dto := models.Fulfillment{
		EquipmentID:         input.Body.EquipmentID,
		FulfillsEquipmentID: input.Body.FulfillsEquipmentID,
	}

	entity := models.FulfillmentFromDTO(dto)
	userID := humaconfig.GetUserID(ctx)

	// Inherit ownership group from the equipment.
	var parent models.EquipmentEntity
	if err := database.DB.Select("ownership_group_id").First(&parent, input.Body.EquipmentID).Error; err != nil {
		return nil, huma.Error404NotFound("Equipment not found")
	}
	entity.OwnershipGroupID = parent.OwnershipGroupID

	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	_ = masteryHandlers.RecalculateEquipmentContributions(database.DB, userID, entity.EquipmentID, entity.FulfillsEquipmentID)
	return &CreateFulfillmentOutput{Body: entity.ToDTO()}, nil
}

// DeleteFulfillment deletes a fulfillment. Owner only.
// DELETE /api/fulfillments/:id
//
// OpenAPI: /api/docs#/operations/DeleteFulfillment
func DeleteFulfillment(ctx context.Context, input *DeleteFulfillmentInput) (*DeleteFulfillmentOutput, error) {
	var entity models.FulfillmentEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Fulfillment not found")
	}

	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, entity.OwnershipGroupID)
	if !access.CanModify() {
		return nil, huma.Error403Forbidden("not the owner of this fulfillment")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	_ = masteryHandlers.RecalculateEquipmentContributions(database.DB, userID, entity.EquipmentID, entity.FulfillsEquipmentID)
	return nil, nil
}
