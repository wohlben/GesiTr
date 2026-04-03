package handlers

import (
	"gesitr/internal/compendium/exercise/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"
)

// ExerciseBody contains the client-provided fields for creating or updating an exercise.
type ExerciseBody struct {
	Names                         []string                     `json:"names" required:"true" minItems:"1"`
	Type                          models.ExerciseType          `json:"type" required:"true"`
	Force                         []models.Force               `json:"force,omitempty"`
	PrimaryMuscles                []models.Muscle              `json:"primaryMuscles,omitempty"`
	SecondaryMuscles              []models.Muscle              `json:"secondaryMuscles,omitempty"`
	TechnicalDifficulty           models.TechnicalDifficulty   `json:"technicalDifficulty,omitempty"`
	BodyWeightScaling             float64                      `json:"bodyWeightScaling,omitempty"`
	SuggestedMeasurementParadigms []models.MeasurementParadigm `json:"suggestedMeasurementParadigms,omitempty"`
	Description                   string                       `json:"description,omitempty"`
	Instructions                  []string                     `json:"instructions,omitempty"`
	Images                        []string                     `json:"images,omitempty"`
	AuthorName                    *string                      `json:"authorName,omitempty"`
	AuthorUrl                     *string                      `json:"authorUrl,omitempty"`
	Public                        bool                         `json:"public,omitempty"`
	ParentExerciseID              *uint                        `json:"parentExerciseId,omitempty"`
	EquipmentIDs                  []uint                       `json:"equipmentIds,omitempty"`
	SourceExerciseID              *uint                        `json:"sourceExerciseId,omitempty" doc:"Source exercise ID for imports (creates forked+equivalent relationships)"`
}

// ExerciseSchemeBody contains the client-provided fields for creating or updating an exercise scheme.
type ExerciseSchemeBody struct {
	ExerciseID           uint     `json:"exerciseId" required:"true"`
	MeasurementType      string   `json:"measurementType" required:"true"`
	Sets                 *int     `json:"sets,omitempty"`
	Reps                 *int     `json:"reps,omitempty"`
	Weight               *float64 `json:"weight,omitempty"`
	RestBetweenSets      *int     `json:"restBetweenSets,omitempty"`
	TimePerRep           *int     `json:"timePerRep,omitempty"`
	Duration             *int     `json:"duration,omitempty"`
	Distance             *float64 `json:"distance,omitempty"`
	TargetTime           *int     `json:"targetTime,omitempty"`
	WorkoutSectionItemID *uint    `json:"workoutSectionItemId,omitempty"`
}

// --- Exercise handlers ---

type ListExercisesInput struct {
	humaconfig.PaginationInput
	Owner         string `query:"owner" doc:"Filter by owner ('me' for current user)"`
	Public        string `query:"public" doc:"'true' to show only public exercises"`
	Mastery       string `query:"mastery" doc:"'me' to show exercises you own or have mastery in"`
	Q             string `query:"q" doc:"Search by name"`
	Type          string `query:"type" doc:"Filter by exercise type"`
	Difficulty    string `query:"difficulty" doc:"Filter by technical difficulty"`
	Force         string `query:"force" doc:"Filter by force type"`
	Muscle        string `query:"muscle" doc:"Filter by any muscle"`
	PrimaryMuscle string `query:"primaryMuscle" doc:"Filter by primary muscle"`
}

type ListExercisesOutput struct {
	Body humaconfig.PaginatedBody[models.Exercise]
}

type CreateExerciseInput struct {
	Body ExerciseBody
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
	ID   uint `path:"id"`
	Body ExerciseBody
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
	ID      uint `path:"id"`
	Version int  `path:"version"`
}

type GetExerciseVersionOutput struct {
	Body shared.VersionEntry
}

type DeleteExerciseVersionInput struct {
	ID      uint `path:"id"`
	Version int  `path:"version"`
}

type DeleteExerciseVersionOutput struct{}

type DeleteAllExerciseVersionsInput struct {
	ID uint `path:"id"`
}

type DeleteAllExerciseVersionsOutput struct{}

// --- Exercise scheme handlers ---

type ListExerciseSchemesInput struct {
	ExerciseID           string `query:"exerciseId" doc:"Filter by exercise ID"`
	MeasurementType      string `query:"measurementType" doc:"Filter by measurement type"`
	WorkoutSectionItemID string `query:"workoutSectionItemId" doc:"Filter by workout section item ID"`
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
