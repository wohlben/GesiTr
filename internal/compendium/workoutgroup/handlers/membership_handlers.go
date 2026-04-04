package handlers

import (
	"context"

	"gesitr/internal/compendium/workoutgroup/models"
	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/compendium/ownershipgroup"

	"github.com/danielgtaylor/huma/v2"
)

// requireGroupOwner fetches the group and checks that the caller is the owner
// via ownership group access.
func requireGroupOwner(ctx context.Context, groupID uint) (*models.WorkoutGroupEntity, error) {
	var group models.WorkoutGroupEntity
	if err := database.DB.First(&group, groupID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout group not found")
	}
	userID := humaconfig.GetUserID(ctx)
	access := ownershipgroup.CheckAccess(database.DB, userID, group.OwnershipGroupID)
	if !access.CanModify() {
		return nil, huma.Error403Forbidden("access denied")
	}
	return &group, nil
}

// ListWorkoutGroupMemberships returns memberships for a group. Owner only.
// GET /api/user/workout-group-memberships
func ListWorkoutGroupMemberships(ctx context.Context, input *ListWorkoutGroupMembershipsInput) (*ListWorkoutGroupMembershipsOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	visibleGroups := ownershipgroup.VisibleGroupIDs(database.DB, userID)
	db := database.DB.Model(&models.WorkoutGroupMembershipEntity{})

	if input.GroupID != "" {
		db = db.Where("group_id = ?", input.GroupID)
	}

	// Only show memberships for groups the caller can access via ownership groups
	db = db.Where("group_id IN (SELECT id FROM workout_groups WHERE ownership_group_id IN (?) AND deleted_at IS NULL)", visibleGroups)

	var entities []models.WorkoutGroupMembershipEntity
	if err := db.Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.WorkoutGroupMembership, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutGroupMembershipsOutput{Body: dtos}, nil
}

// CreateWorkoutGroupMembership adds a user to a workout group. Owner only.
// POST /api/user/workout-group-memberships
func CreateWorkoutGroupMembership(ctx context.Context, input *CreateWorkoutGroupMembershipInput) (*CreateWorkoutGroupMembershipOutput, error) {
	_, err := requireGroupOwner(ctx, input.Body.GroupID)
	if err != nil {
		return nil, err
	}

	// Prevent owner from adding themselves
	callerID := humaconfig.GetUserID(ctx)
	if input.Body.UserID == callerID {
		return nil, huma.Error422UnprocessableEntity("the workout owner does not need a membership")
	}

	// New memberships always start as "invited"
	entity := models.WorkoutGroupMembershipEntity{
		GroupID: input.Body.GroupID,
		UserID:  input.Body.UserID,
		Role:    models.WorkoutGroupRoleInvited,
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error422UnprocessableEntity("user is already a member of this group")
	}
	return &CreateWorkoutGroupMembershipOutput{Body: entity.ToDTO()}, nil
}

// UpdateWorkoutGroupMembership changes a membership's role. Owner only.
// PUT /api/user/workout-group-memberships/{id}
func UpdateWorkoutGroupMembership(ctx context.Context, input *UpdateWorkoutGroupMembershipInput) (*UpdateWorkoutGroupMembershipOutput, error) {
	var existing models.WorkoutGroupMembershipEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Membership not found")
	}

	if _, err := requireGroupOwner(ctx, existing.GroupID); err != nil {
		return nil, err
	}

	if input.Body.Role != models.WorkoutGroupRoleInvited && input.Body.Role != models.WorkoutGroupRoleMember && input.Body.Role != models.WorkoutGroupRoleAdmin {
		return nil, huma.Error422UnprocessableEntity("role must be 'invited', 'member', or 'admin'")
	}

	existing.Role = input.Body.Role
	if err := database.DB.Save(&existing).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateWorkoutGroupMembershipOutput{Body: existing.ToDTO()}, nil
}

// DeleteWorkoutGroupMembership removes a user from a workout group. Owner only.
// DELETE /api/user/workout-group-memberships/{id}
func DeleteWorkoutGroupMembership(ctx context.Context, input *DeleteWorkoutGroupMembershipInput) (*DeleteWorkoutGroupMembershipOutput, error) {
	var entity models.WorkoutGroupMembershipEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Membership not found")
	}

	if _, err := requireGroupOwner(ctx, entity.GroupID); err != nil {
		return nil, err
	}

	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
