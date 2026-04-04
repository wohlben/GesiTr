package handlers

import (
	"gesitr/internal/user/exercisescheme/models"
)

// ExerciseSchemeBody contains the client-provided fields for creating or updating an exercise scheme.
type ExerciseSchemeBody struct {
	ExerciseID      uint     `json:"exerciseId" required:"true"`
	MeasurementType string   `json:"measurementType" required:"true"`
	Sets            *int     `json:"sets,omitempty"`
	Reps            *int     `json:"reps,omitempty"`
	Weight          *float64 `json:"weight,omitempty"`
	RestBetweenSets *int     `json:"restBetweenSets,omitempty"`
	TimePerRep      *int     `json:"timePerRep,omitempty"`
	Duration        *int     `json:"duration,omitempty"`
	Distance        *float64 `json:"distance,omitempty"`
	TargetTime      *int     `json:"targetTime,omitempty"`
}

type ListExerciseSchemesInput struct {
	ExerciseID      string `query:"exerciseId" doc:"Filter by exercise ID"`
	MeasurementType string `query:"measurementType" doc:"Filter by measurement type"`
}

type ListExerciseSchemesOutput struct {
	Body []models.ExerciseScheme
}

type CreateExerciseSchemeInput struct {
	Body ExerciseSchemeBody
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
	ID   uint `path:"id"`
	Body ExerciseSchemeBody
}

type UpdateExerciseSchemeOutput struct {
	Body models.ExerciseScheme
}

type DeleteExerciseSchemeInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseSchemeOutput struct{}

// --- Exercise Scheme Section Item (join table) types ---

type ExerciseSchemeSectionItemBody struct {
	ExerciseSchemeID     uint `json:"exerciseSchemeId" required:"true"`
	WorkoutSectionItemID uint `json:"workoutSectionItemId" required:"true"`
}

type ListExerciseSchemeSectionItemsInput struct {
	WorkoutSectionItemIDs string `query:"workoutSectionItemIds" doc:"Comma-separated workout section item IDs"`
}

type ListExerciseSchemeSectionItemsOutput struct {
	Body []models.ExerciseSchemeSectionItem
}

type UpsertExerciseSchemeSectionItemInput struct {
	Body ExerciseSchemeSectionItemBody
}

type UpsertExerciseSchemeSectionItemOutput struct {
	Body models.ExerciseSchemeSectionItem
}

type DeleteExerciseSchemeSectionItemInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseSchemeSectionItemOutput struct{}
