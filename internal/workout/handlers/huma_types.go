package handlers

import (
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"
	"gesitr/internal/workout/models"
)

// --- Workout handlers ---

type ListWorkoutsInput struct {
	humaconfig.PaginationInput
	Owner  string `query:"owner" doc:"Filter by owner ('me' for current user)"`
	Public string `query:"public" doc:"'true' to show only public workouts"`
	Logged string `query:"logged" doc:"'me' to show workouts you own or have logged"`
	Q      string `query:"q" doc:"Search by name"`
}

type ListWorkoutsOutput struct {
	Body humaconfig.PaginatedBody[models.Workout]
}

type WorkoutBody struct {
	Name            string  `json:"name" required:"true"`
	Notes           *string `json:"notes,omitempty"`
	Public          bool    `json:"public" required:"false"`
	SourceWorkoutID *uint   `json:"sourceWorkoutId,omitempty" doc:"Source workout ID for forks (creates forked+equivalent relationships)"`
}

type CreateWorkoutInput struct {
	Body WorkoutBody
}

type CreateWorkoutOutput struct {
	Body models.Workout
}

type GetWorkoutInput struct {
	ID uint `path:"id"`
}

type GetWorkoutOutput struct {
	Body models.Workout
}

type UpdateWorkoutInput struct {
	ID   uint `path:"id"`
	Body WorkoutBody
}

type UpdateWorkoutOutput struct {
	Body models.Workout
}

type DeleteWorkoutInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutOutput struct{}

type GetWorkoutPermissionsInput struct {
	ID uint `path:"id"`
}

type GetWorkoutPermissionsOutput struct {
	Body shared.PermissionsResponse
}

type ListWorkoutVersionsInput struct {
	ID uint `path:"id"`
}

type ListWorkoutVersionsOutput struct {
	Body []shared.VersionEntry
}

type GetWorkoutVersionInput struct {
	ID      uint `path:"id"`
	Version int  `path:"version"`
}

type GetWorkoutVersionOutput struct {
	Body shared.VersionEntry
}

// --- Workout section handlers ---

type ListWorkoutSectionsInput struct {
	WorkoutID string `query:"workoutId" doc:"Filter by workout ID"`
}

type ListWorkoutSectionsOutput struct {
	Body []models.WorkoutSection
}

type WorkoutSectionBody struct {
	WorkoutID            uint                      `json:"workoutId" required:"true"`
	Type                 models.WorkoutSectionType `json:"type" required:"true"`
	Label                *string                   `json:"label,omitempty"`
	Position             int                       `json:"position"`
	RestBetweenExercises *int                      `json:"restBetweenExercises,omitempty"`
}

type CreateWorkoutSectionInput struct {
	Body WorkoutSectionBody
}

type CreateWorkoutSectionOutput struct {
	Body models.WorkoutSection
}

type GetWorkoutSectionInput struct {
	ID uint `path:"id"`
}

type GetWorkoutSectionOutput struct {
	Body models.WorkoutSection
}

type DeleteWorkoutSectionInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutSectionOutput struct{}

// --- Workout section item handlers ---

type ListWorkoutSectionItemsInput struct {
	WorkoutSectionID string `query:"workoutSectionId" doc:"Filter by workout section ID"`
}

type ListWorkoutSectionItemsOutput struct {
	Body []models.WorkoutSectionItem
}

type WorkoutSectionItemBody struct {
	WorkoutSectionID uint                          `json:"workoutSectionId" required:"true"`
	Type             models.WorkoutSectionItemType `json:"type" required:"true"`
	ExerciseSchemeID *uint                         `json:"exerciseSchemeId,omitempty"`
	ExerciseGroupID  *uint                         `json:"exerciseGroupId,omitempty"`
	Data             *string                       `json:"data,omitempty"`
	Position         int                           `json:"position"`
}

type CreateWorkoutSectionItemInput struct {
	Body WorkoutSectionItemBody
}

type CreateWorkoutSectionItemOutput struct {
	Body models.WorkoutSectionItem
}

type DeleteWorkoutSectionItemInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutSectionItemOutput struct{}
