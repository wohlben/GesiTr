package handlers

import (
	"time"

	workoutmodels "gesitr/internal/user/workout/models"
	"gesitr/internal/user/workoutlog/models"
)

// --- Workout log handlers ---

type ListWorkoutLogsInput struct {
	WorkoutID string `query:"workoutId" doc:"Filter by workout ID"`
	Status    string `query:"status" doc:"Filter by status"`
	PeriodID  string `query:"periodId" doc:"Filter by schedule period ID"`
}

type ListWorkoutLogsOutput struct {
	Body []models.WorkoutLog
}

type WorkoutLogBody struct {
	WorkoutID *uint      `json:"workoutId,omitempty"`
	Name      string     `json:"name" required:"true"`
	Notes     *string    `json:"notes,omitempty"`
	Date      *time.Time `json:"date,omitempty"`
	Status    *string    `json:"status,omitempty" doc:"Initial status: 'proposed' or 'committed'. Defaults to 'planning'."`
	DueStart  *time.Time `json:"dueStart,omitempty" doc:"Start of the commitment window. Required when status is proposed or committed."`
	DueEnd    *time.Time `json:"dueEnd,omitempty" doc:"End of the commitment window. Required when status is proposed or committed."`
}

type CreateWorkoutLogInput struct {
	Body WorkoutLogBody
}

type CreateWorkoutLogOutput struct {
	Body models.WorkoutLog
}

type GetWorkoutLogInput struct {
	ID uint `path:"id"`
}

type GetWorkoutLogOutput struct {
	Body models.WorkoutLog
}

type UpdateWorkoutLogBody struct {
	Name  *string `json:"name,omitempty"`
	Notes *string `json:"notes,omitempty"`
}

type UpdateWorkoutLogInput struct {
	ID   uint `path:"id"`
	Body UpdateWorkoutLogBody
}

type UpdateWorkoutLogOutput struct {
	Body models.WorkoutLog
}

type DeleteWorkoutLogInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutLogOutput struct{}

type StartWorkoutLogInput struct {
	ID uint `path:"id"`
}

type StartWorkoutLogOutput struct {
	Body models.WorkoutLog
}

type StartAdhocWorkoutLogInput struct{}

type StartAdhocWorkoutLogOutput struct {
	Body models.WorkoutLog
}

type FinishWorkoutLogInput struct {
	ID uint `path:"id"`
}

type FinishWorkoutLogOutput struct {
	Body models.WorkoutLog
}

type AbandonWorkoutLogInput struct {
	ID uint `path:"id"`
}

type AbandonWorkoutLogOutput struct {
	Body models.WorkoutLog
}

type SkipWorkoutLogInput struct {
	ID uint `path:"id"`
}

type SkipWorkoutLogOutput struct {
	Body models.WorkoutLog
}

type CommitWorkoutLogInput struct {
	ID uint `path:"id"`
}

type CommitWorkoutLogOutput struct {
	Body models.WorkoutLog
}

// --- Workout log section handlers ---

type ListWorkoutLogSectionsInput struct {
	WorkoutLogID string `query:"workoutLogId" doc:"Filter by workout log ID"`
}

type ListWorkoutLogSectionsOutput struct {
	Body []models.WorkoutLogSection
}

type WorkoutLogSectionBody struct {
	WorkoutLogID         uint                             `json:"workoutLogId" required:"true"`
	Type                 workoutmodels.WorkoutSectionType `json:"type" required:"true"`
	Label                *string                          `json:"label,omitempty"`
	Position             int                              `json:"position"`
	RestBetweenExercises *int                             `json:"restBetweenExercises,omitempty"`
}

type CreateWorkoutLogSectionInput struct {
	Body WorkoutLogSectionBody
}

type CreateWorkoutLogSectionOutput struct {
	Body models.WorkoutLogSection
}

type GetWorkoutLogSectionInput struct {
	ID uint `path:"id"`
}

type GetWorkoutLogSectionOutput struct {
	Body models.WorkoutLogSection
}

type UpdateWorkoutLogSectionBody struct {
	Type                 *workoutmodels.WorkoutSectionType `json:"type,omitempty"`
	Label                *string                           `json:"label,omitempty"`
	RestBetweenExercises *int                              `json:"restBetweenExercises,omitempty"`
	Position             *int                              `json:"position,omitempty"`
}

type UpdateWorkoutLogSectionInput struct {
	ID   uint `path:"id"`
	Body UpdateWorkoutLogSectionBody
}

type UpdateWorkoutLogSectionOutput struct {
	Body models.WorkoutLogSection
}

type DeleteWorkoutLogSectionInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutLogSectionOutput struct{}

// --- Workout log exercise handlers ---

type ListWorkoutLogExercisesInput struct {
	WorkoutLogSectionID string `query:"workoutLogSectionId" doc:"Filter by workout log section ID"`
}

type ListWorkoutLogExercisesOutput struct {
	Body []models.WorkoutLogExercise
}

type WorkoutLogExerciseBody struct {
	WorkoutLogSectionID    uint `json:"workoutLogSectionId" required:"true"`
	SourceExerciseSchemeID uint `json:"sourceExerciseSchemeId" required:"true"`
	Position               int  `json:"position"`
	BreakAfterSeconds      *int `json:"breakAfterSeconds,omitempty"`
}

type CreateWorkoutLogExerciseInput struct {
	Body WorkoutLogExerciseBody
}

type CreateWorkoutLogExerciseOutput struct {
	Body models.WorkoutLogExercise
}

type UpdateWorkoutLogExerciseBody struct {
	Position          *int `json:"position,omitempty"`
	BreakAfterSeconds *int `json:"breakAfterSeconds,omitempty"`
}

type UpdateWorkoutLogExerciseInput struct {
	ID   uint `path:"id"`
	Body UpdateWorkoutLogExerciseBody
}

type UpdateWorkoutLogExerciseOutput struct {
	Body models.WorkoutLogExercise
}

type DeleteWorkoutLogExerciseInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutLogExerciseOutput struct{}

// --- Workout log exercise set handlers ---

type ListWorkoutLogExerciseSetsInput struct {
	WorkoutLogExerciseID string `query:"workoutLogExerciseId" doc:"Filter by workout log exercise ID"`
}

type ListWorkoutLogExerciseSetsOutput struct {
	Body []models.WorkoutLogExerciseSet
}

type WorkoutLogExerciseSetBody struct {
	WorkoutLogExerciseID uint     `json:"workoutLogExerciseId" required:"true"`
	SetNumber            int      `json:"setNumber" required:"true"`
	BreakAfterSeconds    *int     `json:"breakAfterSeconds,omitempty"`
	TargetReps           *int     `json:"targetReps,omitempty"`
	TargetWeight         *float64 `json:"targetWeight,omitempty"`
	TargetDuration       *int     `json:"targetDuration,omitempty"`
	TargetDistance       *float64 `json:"targetDistance,omitempty"`
	TargetTime           *int     `json:"targetTime,omitempty"`
}

type CreateWorkoutLogExerciseSetInput struct {
	Body WorkoutLogExerciseSetBody
}

type CreateWorkoutLogExerciseSetOutput struct {
	Body models.WorkoutLogExerciseSet
}

type UpdateWorkoutLogExerciseSetBody struct {
	Status            models.WorkoutLogItemStatus `json:"status,omitempty"`
	BreakAfterSeconds *int                        `json:"breakAfterSeconds,omitempty"`
	TargetReps        *int                        `json:"targetReps,omitempty"`
	TargetWeight      *float64                    `json:"targetWeight,omitempty"`
	TargetDuration    *int                        `json:"targetDuration,omitempty"`
	TargetDistance    *float64                    `json:"targetDistance,omitempty"`
	TargetTime        *int                        `json:"targetTime,omitempty"`
	ActualReps        *int                        `json:"actualReps,omitempty"`
	ActualWeight      *float64                    `json:"actualWeight,omitempty"`
	ActualDuration    *int                        `json:"actualDuration,omitempty"`
	ActualDistance    *float64                    `json:"actualDistance,omitempty"`
	ActualTime        *int                        `json:"actualTime,omitempty"`
}

type UpdateWorkoutLogExerciseSetInput struct {
	ID   uint `path:"id"`
	Body UpdateWorkoutLogExerciseSetBody
}

type UpdateWorkoutLogExerciseSetOutput struct {
	Body models.WorkoutLogExerciseSet
}

type DeleteWorkoutLogExerciseSetInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutLogExerciseSetOutput struct{}
