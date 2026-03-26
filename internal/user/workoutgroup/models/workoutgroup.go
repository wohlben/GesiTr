package models

import "gesitr/internal/shared"

type WorkoutGroupRole string

const (
	WorkoutGroupRoleInvited WorkoutGroupRole = "invited"
	WorkoutGroupRoleMember  WorkoutGroupRole = "member"
	WorkoutGroupRoleAdmin   WorkoutGroupRole = "admin"
)

type WorkoutGroup struct {
	shared.BaseModel `tstype:",extends"`
	Name             string `json:"name"`
	WorkoutID        uint   `json:"workoutId"`
	Owner            string `json:"owner"`
}

type WorkoutGroupMembership struct {
	shared.BaseModel `tstype:",extends"`
	GroupID          uint             `json:"groupId"`
	UserID           string           `json:"userId"`
	Role             WorkoutGroupRole `json:"role"`
}
