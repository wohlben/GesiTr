package handlers

import "gesitr/internal/workoutgroup/models"

// --- Workout group types ---

type WorkoutGroupBody struct {
	Name      string `json:"name" required:"true"`
	WorkoutID uint   `json:"workoutId" required:"true"`
}

type ListWorkoutGroupsInput struct{}

type ListWorkoutGroupsOutput struct {
	Body []models.WorkoutGroup
}

type CreateWorkoutGroupInput struct {
	Body WorkoutGroupBody
}

type CreateWorkoutGroupOutput struct {
	Body models.WorkoutGroup
}

type GetWorkoutGroupInput struct {
	ID uint `path:"id"`
}

type GetWorkoutGroupOutput struct {
	Body models.WorkoutGroup
}

type UpdateWorkoutGroupInput struct {
	ID   uint `path:"id"`
	Body struct {
		Name string `json:"name" required:"true"`
	}
}

type UpdateWorkoutGroupOutput struct {
	Body models.WorkoutGroup
}

type DeleteWorkoutGroupInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutGroupOutput struct{}

// --- Workout group membership types ---

type WorkoutGroupMembershipBody struct {
	GroupID uint                    `json:"groupId" required:"true"`
	UserID  string                  `json:"userId" required:"true"`
	Role    models.WorkoutGroupRole `json:"role" required:"true"`
}

type ListWorkoutGroupMembershipsInput struct {
	GroupID string `query:"groupId" doc:"Filter by group ID"`
}

type ListWorkoutGroupMembershipsOutput struct {
	Body []models.WorkoutGroupMembership
}

type CreateWorkoutGroupMembershipInput struct {
	Body WorkoutGroupMembershipBody
}

type CreateWorkoutGroupMembershipOutput struct {
	Body models.WorkoutGroupMembership
}

type UpdateWorkoutGroupMembershipInput struct {
	ID   uint `path:"id"`
	Body struct {
		Role models.WorkoutGroupRole `json:"role" required:"true"`
	}
}

type UpdateWorkoutGroupMembershipOutput struct {
	Body models.WorkoutGroupMembership
}

type DeleteWorkoutGroupMembershipInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutGroupMembershipOutput struct{}
