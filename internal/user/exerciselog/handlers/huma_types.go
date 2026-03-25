package handlers

import (
	"time"

	"gesitr/internal/user/exerciselog/models"
)

// --- ExerciseLog handlers ---

type ListExerciseLogsInput struct {
	ExerciseID      string `query:"exerciseId" doc:"Filter by exercise ID"`
	MeasurementType string `query:"measurementType" doc:"Filter by measurement type"`
	IsRecord        string `query:"isRecord" doc:"'true' to show only records"`
	From            string `query:"from" doc:"Filter logs performed at or after this date"`
	To              string `query:"to" doc:"Filter logs performed at or before this date"`
}

type ListExerciseLogsOutput struct {
	Body []models.ExerciseLog
}

type ExerciseLogBody struct {
	ExerciseID              uint      `json:"exerciseId" required:"true"`
	MeasurementType         string    `json:"measurementType" required:"true"`
	Reps                    *int      `json:"reps,omitempty"`
	Weight                  *float64  `json:"weight,omitempty"`
	Duration                *int      `json:"duration,omitempty"`
	Distance                *float64  `json:"distance,omitempty"`
	Time                    *int      `json:"time,omitempty"`
	PerformedAt             time.Time `json:"performedAt,omitempty"`
	WorkoutLogExerciseSetID *uint     `json:"workoutLogExerciseSetId,omitempty"`
	SourceExerciseSchemeID  *uint     `json:"sourceExerciseSchemeId,omitempty"`
}

type CreateExerciseLogInput struct {
	Body ExerciseLogBody
}

type CreateExerciseLogOutput struct {
	Body models.ExerciseLog
}

type GetExerciseLogInput struct {
	ID uint `path:"id"`
}

type GetExerciseLogOutput struct {
	Body models.ExerciseLog
}

type UpdateExerciseLogBody struct {
	Reps     *int     `json:"reps,omitempty"`
	Weight   *float64 `json:"weight,omitempty"`
	Duration *int     `json:"duration,omitempty"`
	Distance *float64 `json:"distance,omitempty"`
	Time     *int     `json:"time,omitempty"`
}

type UpdateExerciseLogInput struct {
	ID   uint `path:"id"`
	Body UpdateExerciseLogBody
}

type UpdateExerciseLogOutput struct {
	Body models.ExerciseLog
}

type DeleteExerciseLogInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseLogOutput struct{}
