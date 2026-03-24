package handlers

import (
	"gesitr/internal/exercise/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"
)

// --- Exercise handlers ---

type ListExercisesInput struct {
	humaconfig.PaginationInput
	Owner         string `query:"owner" doc:"Filter by owner ('me' for current user)"`
	Public        string `query:"public" doc:"'true' to show only public exercises"`
	Q             string `query:"q" doc:"Search by name or alternative name"`
	Type          string `query:"type" doc:"Filter by exercise type"`
	Difficulty    string `query:"difficulty" doc:"Filter by technical difficulty"`
	Force         string `query:"force" doc:"Filter by force type"`
	Muscle        string `query:"muscle" doc:"Filter by any muscle"`
	PrimaryMuscle string `query:"primaryMuscle" doc:"Filter by primary muscle"`
}

type ListExercisesOutput struct {
	Body humaconfig.PaginatedBody[models.Exercise]
}

// RawBody skips huma's automatic validation — the Exercise DTO is shared
// between request and response and has server-set fields (id, createdAt, etc.)
// that aren't present in create requests.
type CreateExerciseInput struct {
	RawBody []byte
}

type CreateExerciseOutput struct {
	Body models.Exercise
}

type GetExerciseInput struct {
	ID uint `path:"id"`
}

type GetExerciseOutput struct {
	Body models.Exercise
}

type UpdateExerciseInput struct {
	ID      uint `path:"id"`
	RawBody []byte
}

type UpdateExerciseOutput struct {
	Body models.Exercise
}

type DeleteExerciseInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseOutput struct{}

type GetExercisePermissionsInput struct {
	ID uint `path:"id"`
}

type GetExercisePermissionsOutput struct {
	Body shared.PermissionsResponse
}

type ListExerciseVersionsInput struct {
	ID uint `path:"id"`
}

type ListExerciseVersionsOutput struct {
	Body []shared.VersionEntry
}

type GetExerciseVersionInput struct {
	TemplateID string `path:"templateId"`
	Version    int    `path:"version"`
}

type GetExerciseVersionOutput struct {
	Body shared.VersionEntry
}

// --- Exercise scheme handlers ---

type ListExerciseSchemesInput struct {
	ExerciseID      string `query:"exerciseId" doc:"Filter by exercise ID"`
	MeasurementType string `query:"measurementType" doc:"Filter by measurement type"`
}

type ListExerciseSchemesOutput struct {
	Body []models.ExerciseScheme
}

type CreateExerciseSchemeInput struct {
	RawBody []byte
}

type CreateExerciseSchemeOutput struct {
	Body models.ExerciseScheme
}

type GetExerciseSchemeInput struct {
	ID uint `path:"id"`
}

type GetExerciseSchemeOutput struct {
	Body models.ExerciseScheme
}

type UpdateExerciseSchemeInput struct {
	ID      uint `path:"id"`
	RawBody []byte
}

type UpdateExerciseSchemeOutput struct {
	Body models.ExerciseScheme
}

type DeleteExerciseSchemeInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseSchemeOutput struct{}
