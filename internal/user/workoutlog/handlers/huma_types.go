package handlers

import (
	"gesitr/internal/user/workoutlog/models"
)

// --- Workout log handlers ---

type ListWorkoutLogsInput struct {
	WorkoutID string `query:"workoutId" doc:"Filter by workout ID"`
	Status    string `query:"status" doc:"Filter by status"`
}

type ListWorkoutLogsOutput struct {
	Body []models.WorkoutLog
}

type CreateWorkoutLogInput struct {
	RawBody []byte
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

type UpdateWorkoutLogInput struct {
	ID      uint `path:"id"`
	RawBody []byte
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

// --- Workout log section handlers ---

type ListWorkoutLogSectionsInput struct {
	WorkoutLogID string `query:"workoutLogId" doc:"Filter by workout log ID"`
}

type ListWorkoutLogSectionsOutput struct {
	Body []models.WorkoutLogSection
}

type CreateWorkoutLogSectionInput struct {
	RawBody []byte
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

type UpdateWorkoutLogSectionInput struct {
	ID      uint `path:"id"`
	RawBody []byte
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

type CreateWorkoutLogExerciseInput struct {
	RawBody []byte
}

type CreateWorkoutLogExerciseOutput struct {
	Body models.WorkoutLogExercise
}

type UpdateWorkoutLogExerciseInput struct {
	ID      uint `path:"id"`
	RawBody []byte
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

type CreateWorkoutLogExerciseSetInput struct {
	RawBody []byte
}

type CreateWorkoutLogExerciseSetOutput struct {
	Body models.WorkoutLogExerciseSet
}

type UpdateWorkoutLogExerciseSetInput struct {
	ID      uint `path:"id"`
	RawBody []byte
}

type UpdateWorkoutLogExerciseSetOutput struct {
	Body models.WorkoutLogExerciseSet
}

type DeleteWorkoutLogExerciseSetInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutLogExerciseSetOutput struct{}
