package handlers

import (
	"context"

	"gesitr/internal/compendium/ownershipgroup"
	"gesitr/internal/compendium/workout/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

// ListExerciseGroups returns exercise groups, optionally filtered by name search.
// GET /api/exercise-groups
//
// OpenAPI: /api/docs#/operations/ListExerciseGroups
func ListExerciseGroups(ctx context.Context, input *ListExerciseGroupsInput) (*ListExerciseGroupsOutput, error) {
	db := database.DB.Model(&models.ExerciseGroupEntity{})

	if input.Q != "" {
		pattern := "%" + input.Q + "%"
		db = db.Where("name IS NOT NULL AND name LIKE ?", pattern)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	p := input.ToPaginationParams()
	var entities []models.ExerciseGroupEntity
	if err := shared.ApplyPagination(db, p).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.ExerciseGroup, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListExerciseGroupsOutput{Body: humaconfig.PaginatedBody[models.ExerciseGroup]{
		Items: dtos, Total: total, Limit: p.Limit, Offset: p.Offset,
	}}, nil
}

// CreateExerciseGroup creates an exercise group owned by the current user.
// POST /api/exercise-groups
//
// OpenAPI: /api/docs#/operations/CreateExerciseGroup
func CreateExerciseGroup(ctx context.Context, input *CreateExerciseGroupInput) (*CreateExerciseGroupOutput, error) {
	dto := models.ExerciseGroup{
		Name: input.Body.Name,
	}

	entity := models.ExerciseGroupFromDTO(dto)
	userID := humaconfig.GetUserID(ctx)

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
	return &CreateExerciseGroupOutput{Body: entity.ToDTO()}, nil
}

// GetExerciseGroupPermissions returns the current user's permissions on an exercise group.
// GET /api/exercise-groups/:id/permissions
//
// OpenAPI: /api/docs#/operations/GetExerciseGroupPermissions
func GetExerciseGroupPermissions(ctx context.Context, input *GetExerciseGroupPermissionsInput) (*GetExerciseGroupPermissionsOutput, error) {
	var entity models.ExerciseGroupEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("ExerciseGroup not found")
	}
	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, entity.OwnershipGroupID)
	perms, _ := shared.ResolvePermissionsFromAccess(access, false)
	if perms == nil {
		perms = []shared.Permission{}
	}
	return &GetExerciseGroupPermissionsOutput{Body: shared.PermissionsResponse{Permissions: perms}}, nil
}

// GetExerciseGroup returns a single exercise group.
// GET /api/exercise-groups/:id
//
// OpenAPI: /api/docs#/operations/GetExerciseGroup
func GetExerciseGroup(ctx context.Context, input *GetExerciseGroupInput) (*GetExerciseGroupOutput, error) {
	var entity models.ExerciseGroupEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("ExerciseGroup not found")
	}
	return &GetExerciseGroupOutput{Body: entity.ToDTO()}, nil
}

// UpdateExerciseGroup updates an exercise group. Owner only.
// PUT /api/exercise-groups/:id
//
// OpenAPI: /api/docs#/operations/UpdateExerciseGroup
func UpdateExerciseGroup(ctx context.Context, input *UpdateExerciseGroupInput) (*UpdateExerciseGroupOutput, error) {
	var existing models.ExerciseGroupEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("ExerciseGroup not found")
	}

	access := ownershipgroup.CheckAccess(database.DB, humaconfig.GetUserID(ctx), existing.OwnershipGroupID)
	if !access.CanModify() {
		return nil, huma.Error403Forbidden("not the owner of this exercise group")
	}

	dto := models.ExerciseGroup{
		Name: input.Body.Name,
	}

	entity := models.ExerciseGroupFromDTO(dto)
	entity.ID = existing.ID
	entity.OwnershipGroupID = existing.OwnershipGroupID

	if err := database.DB.Save(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateExerciseGroupOutput{Body: entity.ToDTO()}, nil
}

// DeleteExerciseGroup deletes an exercise group. Owner only.
// DELETE /api/exercise-groups/:id
//
// OpenAPI: /api/docs#/operations/DeleteExerciseGroup
func DeleteExerciseGroup(ctx context.Context, input *DeleteExerciseGroupInput) (*DeleteExerciseGroupOutput, error) {
	var entity models.ExerciseGroupEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("ExerciseGroup not found")
	}

	access := ownershipgroup.CheckAccess(database.DB, humaconfig.GetUserID(ctx), entity.OwnershipGroupID)
	if !access.CanDelete() {
		return nil, huma.Error403Forbidden("not the owner of this exercise group")
	}

	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
