package workoutgroup

import (
	"gesitr/internal/database"
	"gesitr/internal/workoutgroup/models"
)

// WorkoutAccess describes a user's access level for a workout.
type WorkoutAccess struct {
	IsOwner        bool
	IsMember       bool
	IsAdmin        bool
	GroupName      string
	MembershipRole string
}

// CheckWorkoutAccess determines a user's access level for a given workout.
// It checks ownership first, then group membership.
func CheckWorkoutAccess(userID, workoutOwner string, workoutID uint) WorkoutAccess {
	if userID == workoutOwner {
		return WorkoutAccess{IsOwner: true}
	}

	var result struct {
		Role      models.WorkoutGroupRole
		GroupName string
	}
	err := database.DB.
		Table("workout_group_memberships").
		Select("workout_group_memberships.role, workout_groups.name AS group_name").
		Joins("JOIN workout_groups ON workout_groups.id = workout_group_memberships.group_id AND workout_groups.deleted_at IS NULL").
		Where("workout_groups.workout_id = ? AND workout_group_memberships.user_id = ? AND workout_group_memberships.deleted_at IS NULL", workoutID, userID).
		Scan(&result).Error

	if err != nil || result.Role == "" {
		return WorkoutAccess{}
	}

	return WorkoutAccess{
		IsMember:       true,
		IsAdmin:        result.Role == models.WorkoutGroupRoleAdmin,
		GroupName:      result.GroupName,
		MembershipRole: string(result.Role),
	}
}

// GroupInfoForWorkouts returns group name and membership role keyed by workout ID.
func GroupInfoForWorkouts(userID string, workoutIDs []uint) map[uint]WorkoutAccess {
	if len(workoutIDs) == 0 {
		return nil
	}

	var results []struct {
		WorkoutID uint
		GroupName string
		Role      models.WorkoutGroupRole
	}
	database.DB.
		Table("workout_group_memberships").
		Select("workout_groups.workout_id, workout_groups.name AS group_name, workout_group_memberships.role").
		Joins("JOIN workout_groups ON workout_groups.id = workout_group_memberships.group_id AND workout_groups.deleted_at IS NULL").
		Where("workout_groups.workout_id IN ? AND workout_group_memberships.user_id = ? AND workout_group_memberships.deleted_at IS NULL", workoutIDs, userID).
		Scan(&results)

	m := make(map[uint]WorkoutAccess, len(results))
	for _, r := range results {
		m[r.WorkoutID] = WorkoutAccess{
			IsMember:       true,
			IsAdmin:        r.Role == models.WorkoutGroupRoleAdmin,
			GroupName:      r.GroupName,
			MembershipRole: string(r.Role),
		}
	}
	return m
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
