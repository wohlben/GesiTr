package handlers

import (
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

type CreateExerciseLogInput struct {
	RawBody []byte
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

type UpdateExerciseLogInput struct {
	ID      uint `path:"id"`
	RawBody []byte
}

type UpdateExerciseLogOutput struct {
	Body models.ExerciseLog
}

type DeleteExerciseLogInput struct {
	ID uint `path:"id"`
}

type DeleteExerciseLogOutput struct{}
