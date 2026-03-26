package workoutgroup

import (
	"gesitr/internal/database"
	"gesitr/internal/user/workoutgroup/models"
)

// WorkoutAccess describes a user's access level for a workout.
type WorkoutAccess struct {
	IsOwner  bool
	IsMember bool
	IsAdmin  bool
}

// CheckWorkoutAccess determines a user's access level for a given workout.
// It checks ownership first, then group membership.
func CheckWorkoutAccess(userID, workoutOwner string, workoutID uint) WorkoutAccess {
	if userID == workoutOwner {
		return WorkoutAccess{IsOwner: true}
	}

	var membership models.WorkoutGroupMembershipEntity
	err := database.DB.
		Joins("JOIN workout_groups ON workout_groups.id = workout_group_memberships.group_id AND workout_groups.deleted_at IS NULL").
		Where("workout_groups.workout_id = ? AND workout_group_memberships.user_id = ? AND workout_group_memberships.deleted_at IS NULL", workoutID, userID).
		First(&membership).Error

	if err != nil {
		return WorkoutAccess{}
	}

	return WorkoutAccess{
		IsMember: true,
		IsAdmin:  membership.Role == models.WorkoutGroupRoleAdmin,
	}
}

func (a WorkoutAccess) CanRead() bool {
	return a.IsOwner || a.IsMember || a.IsAdmin
}

func (a WorkoutAccess) CanModify() bool {
	return a.IsOwner || a.IsAdmin
}

func (a WorkoutAccess) CanDelete() bool {
	return a.IsOwner
}
