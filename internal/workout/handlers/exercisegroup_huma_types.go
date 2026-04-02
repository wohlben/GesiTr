package handlers

import (
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"
	"gesitr/internal/workout/models"
)

// ExerciseGroupBody contains the client-provided fields for creating or updating an exercise group.
type ExerciseGroupBody struct {
	Name *string `json:"name,omitempty"`
}

// ExerciseGroupMemberBody contains the client-provided fields for creating an exercise group member.
type ExerciseGroupMemberBody struct {
	GroupID    uint `json:"groupId" required:"true"`
	ExerciseID uint `json:"exerciseId" required:"true"`
}

// --- Exercise group handlers ---

type ListExerciseGroupsInput struct {
	humaconfig.PaginationInput
	Q string `query:"q" doc:"Search by name"`
}

type ListExerciseGroupsOutput struct {
	Body humaconfig.PaginatedBody[models.ExerciseGroup]
}

type CreateExerciseGroupInput struct {
	Body ExerciseGroupBody
}

type CreateExerciseGroupOutput struct {
	Body models.ExerciseGroup
}

type GetExerciseGroupPermissionsInput struct {
	ID uint `path:"id"`
}

type GetExerciseGroupPermissionsOutput struct {
	Body shared.PermissionsResponse
}

type GetExerciseGroupInput struct {
	ID uint `path:"id"`
}

type GetExerciseGroupOutput struct {
	Body models.ExerciseGroup
}

type UpdateExerciseGroupInput struct {
	ID   uint `path:"id"`
	Body ExerciseGroupBody
}

type UpdateExerciseGroupOutput struct {
	Body models.ExerciseGroup
}

type DeleteExerciseGroupInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseGroupOutput struct{}

// --- Exercise group member handlers ---

type ListExerciseGroupMembersInput struct {
	GroupID    string `query:"groupId" doc:"Filter by group ID"`
	ExerciseID string `query:"exerciseId" doc:"Filter by exercise ID"`
}

type ListExerciseGroupMembersOutput struct {
	Body []models.ExerciseGroupMember
}

type CreateExerciseGroupMemberInput struct {
	Body ExerciseGroupMemberBody
}

type CreateExerciseGroupMemberOutput struct {
	Body models.ExerciseGroupMember
}

type DeleteExerciseGroupMemberInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseGroupMemberOutput struct{}
