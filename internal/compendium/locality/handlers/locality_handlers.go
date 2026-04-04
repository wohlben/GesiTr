package handlers

import (
	"context"

	"gesitr/internal/compendium/locality/models"
	"gesitr/internal/compendium/ownershipgroup"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

// ListLocalities returns localities visible to the current user: their own
// localities plus all public localities.
// GET /api/localities
//
// OpenAPI: /api/docs#/operations/ListLocalities
func ListLocalities(ctx context.Context, input *ListLocalitiesInput) (*ListLocalitiesOutput, error) {
	db := database.DB.Model(&models.LocalityEntity{})

	userID := humaconfig.GetUserID(ctx)
	visibleGroups := ownershipgroup.VisibleGroupIDs(database.DB, userID)
	if input.Owner != "" {
		if input.Owner == "me" || input.Owner == userID {
			db = db.Where("ownership_group_id IN (?)", visibleGroups)
		} else {
			db = db.Where("ownership_group_id IN (SELECT group_id FROM ownership_group_memberships WHERE user_id = ? AND role = 'owner' AND deleted_at IS NULL) AND public = ?", input.Owner, true)
		}
	} else {
		db = db.Where("ownership_group_id IN (?) OR public = ?", visibleGroups, true)
	}

	if input.Public == "true" {
		db = db.Where("public = ?", true)
	} else if input.Public == "false" {
		db = db.Where("public = ?", false)
	}

	if input.Q != "" {
		pattern := "%" + input.Q + "%"
		db = db.Where("name LIKE ?", pattern)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	p := input.ToPaginationParams()
	var entities []models.LocalityEntity
	if err := shared.ApplyPagination(db, p).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.Locality, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListLocalitiesOutput{Body: humaconfig.PaginatedBody[models.Locality]{
		Items: dtos, Total: total, Limit: p.Limit, Offset: p.Offset,
	}}, nil
}

// CreateLocality creates a locality owned by the current user.
// POST /api/localities
//
// OpenAPI: /api/docs#/operations/CreateLocality
func CreateLocality(ctx context.Context, input *CreateLocalityInput) (*CreateLocalityOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	entity := models.LocalityEntity{
		Name:   input.Body.Name,
		Public: input.Body.Public,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&entity).Error; err != nil {
			return err
		}

		groupID, err := ownershipgroup.CreateGroupForEntity(tx, userID)
		if err != nil {
			return err
		}
		entity.OwnershipGroupID = groupID
		return tx.Model(&entity).Update("ownership_group_id", groupID).Error
	})

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &CreateLocalityOutput{Body: entity.ToDTO()}, nil
}

// GetLocalityPermissions returns the current user's permissions on a locality.
// GET /api/localities/:id/permissions
//
// OpenAPI: /api/docs#/operations/GetLocalityPermissions
func GetLocalityPermissions(ctx context.Context, input *GetLocalityPermissionsInput) (*GetLocalityPermissionsOutput, error) {
	var entity models.LocalityEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Locality not found")
	}
	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, entity.OwnershipGroupID)
	perms, _ := shared.ResolvePermissionsFromAccess(access, entity.Public)
	if perms == nil {
		perms = []shared.Permission{}
	}
	return &GetLocalityPermissionsOutput{Body: shared.PermissionsResponse{Permissions: perms}}, nil
}

// GetLocality returns a single locality. Public localities are visible to all
// users; private localities are visible only to their owner.
// GET /api/localities/:id
//
// OpenAPI: /api/docs#/operations/GetLocality
func GetLocality(ctx context.Context, input *GetLocalityInput) (*GetLocalityOutput, error) {
	var entity models.LocalityEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Locality not found")
	}
	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, entity.OwnershipGroupID)
	perms, _ := shared.ResolvePermissionsFromAccess(access, entity.Public)
	if len(perms) == 0 {
		return nil, huma.Error403Forbidden("access denied")
	}
	return &GetLocalityOutput{Body: entity.ToDTO()}, nil
}

// UpdateLocality updates a locality. Owner only.
// PUT /api/localities/:id
//
// OpenAPI: /api/docs#/operations/UpdateLocality
func UpdateLocality(ctx context.Context, input *UpdateLocalityInput) (*UpdateLocalityOutput, error) {
	var existing models.LocalityEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Locality not found")
	}

	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, existing.OwnershipGroupID)
	if !access.CanModify() {
		return nil, huma.Error403Forbidden("access denied")
	}

	existing.Name = input.Body.Name
	existing.Public = input.Body.Public

	if err := database.DB.Save(&existing).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateLocalityOutput{Body: existing.ToDTO()}, nil
}

// DeleteLocality deletes a locality and its availability entries. Owner only.
// DELETE /api/localities/:id
//
// OpenAPI: /api/docs#/operations/DeleteLocality
func DeleteLocality(ctx context.Context, input *DeleteLocalityInput) (*DeleteLocalityOutput, error) {
	var entity models.LocalityEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Locality not found")
	}
	access := ownershipgroup.CheckAccess(database.DB, humaconfig.GetUserID(ctx), entity.OwnershipGroupID)
	if !access.CanDelete() {
		return nil, huma.Error403Forbidden("access denied")
	}

	// Delete availability entries first, then the locality.
	database.DB.Where("locality_id = ?", entity.ID).Delete(&models.LocalityAvailabilityEntity{})
	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
