package handlers

import (
	"context"

	"gesitr/internal/compendium/locality/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/compendium/ownershipgroup"

	"github.com/danielgtaylor/huma/v2"
)

// ListLocalityAvailabilities returns availability entries for a locality.
// GET /api/locality-availabilities
//
// OpenAPI: /api/docs#/operations/ListLocalityAvailabilities
func ListLocalityAvailabilities(ctx context.Context, input *ListLocalityAvailabilitiesInput) (*ListLocalityAvailabilitiesOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	visibleGroups := ownershipgroup.VisibleGroupIDs(database.DB, userID)

	db := database.DB.Model(&models.LocalityAvailabilityEntity{}).
		Select("locality_availabilities.*, equipment.display_name as equipment_name").
		Joins("LEFT JOIN equipment ON equipment.id = locality_availabilities.equipment_id").
		Where("locality_availabilities.ownership_group_id IN (?)", visibleGroups)

	if input.LocalityID != "" {
		db = db.Where("locality_id = ?", input.LocalityID)
	}
	if input.EquipmentID != "" {
		db = db.Where("equipment_id = ?", input.EquipmentID)
	}
	if input.Available == "true" {
		db = db.Where("available = ?", true)
	} else if input.Available == "false" {
		db = db.Where("available = ?", false)
	}

	var entities []models.LocalityAvailabilityEntity
	if err := db.Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.LocalityAvailability, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListLocalityAvailabilitiesOutput{Body: dtos}, nil
}

// CreateLocalityAvailability adds equipment to a locality.
// POST /api/locality-availabilities
//
// OpenAPI: /api/docs#/operations/CreateLocalityAvailability
func CreateLocalityAvailability(ctx context.Context, input *CreateLocalityAvailabilityInput) (*CreateLocalityAvailabilityOutput, error) {
	var locality models.LocalityEntity
	if err := database.DB.First(&locality, input.Body.LocalityID).Error; err != nil {
		return nil, huma.Error404NotFound("Locality not found")
	}
	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, locality.OwnershipGroupID)
	if !access.CanModify() {
		return nil, huma.Error403Forbidden("access denied")
	}

	available := true
	if input.Body.Available != nil {
		available = *input.Body.Available
	}

	entity := models.LocalityAvailabilityEntity{
		LocalityID:       input.Body.LocalityID,
		EquipmentID:      input.Body.EquipmentID,
		Available:        available,
		OwnershipGroupID: locality.OwnershipGroupID,
	}

	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateLocalityAvailabilityOutput{Body: entity.ToDTO()}, nil
}

// UpdateLocalityAvailability toggles the availability of equipment at a locality.
// PUT /api/locality-availabilities/:id
//
// OpenAPI: /api/docs#/operations/UpdateLocalityAvailability
func UpdateLocalityAvailability(ctx context.Context, input *UpdateLocalityAvailabilityInput) (*UpdateLocalityAvailabilityOutput, error) {
	var entity models.LocalityAvailabilityEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Locality availability not found")
	}
	access := ownershipgroup.CheckAccess(database.DB, humaconfig.GetUserID(ctx), entity.OwnershipGroupID)
	if !access.CanModify() {
		return nil, huma.Error403Forbidden("access denied")
	}

	entity.Available = input.Body.Available
	if err := database.DB.Save(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateLocalityAvailabilityOutput{Body: entity.ToDTO()}, nil
}

// DeleteLocalityAvailability removes equipment from a locality.
// DELETE /api/locality-availabilities/:id
//
// OpenAPI: /api/docs#/operations/DeleteLocalityAvailability
func DeleteLocalityAvailability(ctx context.Context, input *DeleteLocalityAvailabilityInput) (*DeleteLocalityAvailabilityOutput, error) {
	var entity models.LocalityAvailabilityEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Locality availability not found")
	}
	access := ownershipgroup.CheckAccess(database.DB, humaconfig.GetUserID(ctx), entity.OwnershipGroupID)
	if !access.CanDelete() {
		return nil, huma.Error403Forbidden("access denied")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
