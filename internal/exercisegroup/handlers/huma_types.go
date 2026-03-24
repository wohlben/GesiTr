package handlers

import (
	"gesitr/internal/exercisegroup/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"
)

// --- Exercise group handlers ---

type ListExerciseGroupsInput struct {
	humaconfig.PaginationInput
	Q string `query:"q" doc:"Search by name"`
}

type ListExerciseGroupsOutput struct {
	Body humaconfig.PaginatedBody[models.ExerciseGroup]
}

// RawBody skips huma's automatic validation — the ExerciseGroup DTO is shared
// between request and response and has server-set fields (id, createdAt, etc.)
// that aren't present in create requests.
type CreateExerciseGroupInput struct {
	RawBody []byte
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
	ID      uint `path:"id"`
	RawBody []byte
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

// RawBody skips huma's automatic validation — the ExerciseGroupMember DTO is shared
// between request and response and has server-set fields (id, createdAt, etc.)
// that aren't present in create requests.
type CreateExerciseGroupMemberInput struct {
	RawBody []byte
}

type CreateExerciseGroupMemberOutput struct {
	Body models.ExerciseGroupMember
}

type DeleteExerciseGroupMemberInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseGroupMemberOutput struct{}
