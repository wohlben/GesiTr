package handlers

import (
	"gesitr/internal/user/workout/models"
)

// --- Workout handlers ---

type ListWorkoutsInput struct{}

type ListWorkoutsOutput struct {
	Body []models.Workout
}

type WorkoutBody struct {
	Name  string  `json:"name" required:"true"`
	Notes *string `json:"notes,omitempty"`
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
