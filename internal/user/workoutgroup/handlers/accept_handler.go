package handlers

import (
	"context"
	"fmt"

	"gesitr/internal/database"
	exercisemodels "gesitr/internal/exercise/models"
	"gesitr/internal/humaconfig"
	workoutmodels "gesitr/internal/user/workout/models"
	"gesitr/internal/user/workoutgroup/models"

	"github.com/danielgtaylor/huma/v2"
)

type AcceptWorkoutGroupInvitationInput struct {
	ID uint `path:"id" doc:"Workout ID"`
}

type AcceptWorkoutGroupInvitationOutput struct {
	Body models.WorkoutGroupMembership
}

// AcceptWorkoutGroupInvitation promotes an "invited" membership to "member"
// after validating that the caller has created exercise schemes for all
// exercise-type items in the workout.
// POST /api/user/workouts/{id}/group/accept
func AcceptWorkoutGroupInvitation(ctx context.Context, input *AcceptWorkoutGroupInvitationInput) (*AcceptWorkoutGroupInvitationOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	// Find the workout group for this workout
	var group models.WorkoutGroupEntity
	if err := database.DB.Where("workout_id = ?", input.ID).First(&group).Error; err != nil {
		return nil, huma.Error404NotFound("No group exists for this workout")
	}

	// Find the caller's membership
	var membership models.WorkoutGroupMembershipEntity
	if err := database.DB.Where("group_id = ? AND user_id = ?", group.ID, userID).First(&membership).Error; err != nil {
		return nil, huma.Error404NotFound("You are not a member of this group")
	}

	if membership.Role != models.WorkoutGroupRoleInvited {
		return nil, huma.Error409Conflict("Invitation already accepted")
	}

	// Fetch all exercise-type items for this workout
	var items []workoutmodels.WorkoutSectionItemEntity
	if err := database.DB.
		Where("workout_section_id IN (SELECT id FROM workout_sections WHERE workout_id = ? AND deleted_at IS NULL)", input.ID).
		Where("type = ?", workoutmodels.WorkoutSectionItemTypeExercise).
		Find(&items).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// For each item, check that the user has a scheme
	var missingItemIDs []uint
	for _, item := range items {
		var count int64
		database.DB.Model(&exercisemodels.ExerciseSchemeEntity{}).
			Where("workout_section_item_id = ? AND owner = ?", item.ID, userID).
			Count(&count)
		if count == 0 {
			missingItemIDs = append(missingItemIDs, item.ID)
		}
	}

	if len(missingItemIDs) > 0 {
		return nil, huma.Error400BadRequest(
			fmt.Sprintf("missing exercise schemes for %d workout item(s)", len(missingItemIDs)),
		)
	}

	// Promote to member
	membership.Role = models.WorkoutGroupRoleMember
	if err := database.DB.Save(&membership).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &AcceptWorkoutGroupInvitationOutput{Body: membership.ToDTO()}, nil
}
