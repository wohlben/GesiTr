package handlers

import (
	"context"
	"encoding/json"

	"gesitr/internal/database"
	"gesitr/internal/equipmentfulfillment/models"
	"gesitr/internal/humaconfig"

	"github.com/danielgtaylor/huma/v2"
)

// ListFulfillments returns equipment fulfillments, optionally filtered by
// equipmentId or fulfillsEquipmentId.
// GET /api/fulfillments
//
// OpenAPI: /api/docs#/operations/list-fulfillments
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
// OpenAPI: /api/docs#/operations/create-fulfillment
func CreateFulfillment(ctx context.Context, input *CreateFulfillmentInput) (*CreateFulfillmentOutput, error) {
	var dto models.Fulfillment
	if err := json.Unmarshal(input.RawBody, &dto); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	entity := models.FulfillmentFromDTO(dto)
	entity.Owner = humaconfig.GetUserID(ctx)
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateFulfillmentOutput{Body: entity.ToDTO()}, nil
}

// DeleteFulfillment deletes a fulfillment. Owner only.
// DELETE /api/fulfillments/:id
//
// OpenAPI: /api/docs#/operations/delete-fulfillment
func DeleteFulfillment(ctx context.Context, input *DeleteFulfillmentInput) (*DeleteFulfillmentOutput, error) {
	var entity models.FulfillmentEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Fulfillment not found")
	}

	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("not the owner of this fulfillment")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
