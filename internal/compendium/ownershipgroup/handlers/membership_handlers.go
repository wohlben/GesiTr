package handlers

import (
	"context"

	"gesitr/internal/compendium/ownershipgroup/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"

	"github.com/danielgtaylor/huma/v2"
)

// requireGroupOwner fetches the group and checks that the caller is an owner-role member.
func requireGroupOwner(ctx context.Context, groupID uint) error {
	userID := humaconfig.GetUserID(ctx)
	var membership models.OwnershipGroupMembershipEntity
	err := database.DB.
		Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, models.RoleOwner).
		First(&membership).Error
	if err != nil {
		return huma.Error403Forbidden("access denied")
	}
	return nil
}

// ListOwnershipGroupMemberships returns memberships for an ownership group. Owner only.
// GET /api/ownership-groups/{id}/memberships
func ListOwnershipGroupMemberships(ctx context.Context, input *ListOwnershipGroupMembershipsInput) (*ListOwnershipGroupMembershipsOutput, error) {
	if err := requireGroupOwner(ctx, input.GroupID); err != nil {
		return nil, err
	}

	var entities []models.OwnershipGroupMembershipEntity
	if err := database.DB.Where("group_id = ?", input.GroupID).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.OwnershipGroupMembership, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListOwnershipGroupMembershipsOutput{Body: dtos}, nil
}

// CreateOwnershipGroupMembership adds a user to an ownership group. Owner only.
// POST /api/ownership-groups/{id}/memberships
func CreateOwnershipGroupMembership(ctx context.Context, input *CreateOwnershipGroupMembershipInput) (*CreateOwnershipGroupMembershipOutput, error) {
	if err := requireGroupOwner(ctx, input.GroupID); err != nil {
		return nil, err
	}

	// Prevent adding someone who is already the owner
	userID := humaconfig.GetUserID(ctx)
	if input.Body.UserID == userID {
		return nil, huma.Error422UnprocessableEntity("cannot add yourself — you are already the owner")
	}

	entity := models.OwnershipGroupMembershipEntity{
		GroupID: input.GroupID,
		UserID:  input.Body.UserID,
		Role:    models.RoleMember,
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error422UnprocessableEntity("user is already a member of this group")
	}
	return &CreateOwnershipGroupMembershipOutput{Body: entity.ToDTO()}, nil
}

// UpdateOwnershipGroupMembership changes a membership's role. Owner only.
// PUT /api/ownership-group-memberships/{id}
func UpdateOwnershipGroupMembership(ctx context.Context, input *UpdateOwnershipGroupMembershipInput) (*UpdateOwnershipGroupMembershipOutput, error) {
	var existing models.OwnershipGroupMembershipEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Membership not found")
	}

	if err := requireGroupOwner(ctx, existing.GroupID); err != nil {
		return nil, err
	}

	if input.Body.Role != models.RoleOwner && input.Body.Role != models.RoleAdmin && input.Body.Role != models.RoleMember {
		return nil, huma.Error422UnprocessableEntity("role must be 'owner', 'admin', or 'member'")
	}

	existing.Role = input.Body.Role
	if err := database.DB.Save(&existing).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateOwnershipGroupMembershipOutput{Body: existing.ToDTO()}, nil
}

// DeleteOwnershipGroupMembership removes a user from an ownership group. Owner only.
// DELETE /api/ownership-group-memberships/{id}
func DeleteOwnershipGroupMembership(ctx context.Context, input *DeleteOwnershipGroupMembershipInput) (*DeleteOwnershipGroupMembershipOutput, error) {
	var entity models.OwnershipGroupMembershipEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Membership not found")
	}

	if err := requireGroupOwner(ctx, entity.GroupID); err != nil {
		return nil, err
	}

	// Prevent removing the last owner
	if entity.Role == models.RoleOwner {
		var ownerCount int64
		database.DB.Model(&models.OwnershipGroupMembershipEntity{}).
			Where("group_id = ? AND role = ? AND id != ?", entity.GroupID, models.RoleOwner, entity.ID).
			Count(&ownerCount)
		if ownerCount == 0 {
			return nil, huma.Error422UnprocessableEntity("cannot remove the last owner")
		}
	}

	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
