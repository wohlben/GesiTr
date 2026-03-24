package handlers

import (
	"gesitr/internal/user/workout/models"
)

// --- Workout handlers ---

type ListWorkoutsInput struct{}

type ListWorkoutsOutput struct {
	Body []models.Workout
}

type CreateWorkoutInput struct {
	RawBody []byte
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
	ID      uint `path:"id"`
	RawBody []byte
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

type CreateWorkoutSectionInput struct {
	RawBody []byte
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

// --- Workout section exercise handlers ---

type ListWorkoutSectionExercisesInput struct {
	WorkoutSectionID string `query:"workoutSectionId" doc:"Filter by workout section ID"`
}

type ListWorkoutSectionExercisesOutput struct {
	Body []models.WorkoutSectionExercise
}

type CreateWorkoutSectionExerciseInput struct {
	RawBody []byte
}

type CreateWorkoutSectionExerciseOutput struct {
	Body models.WorkoutSectionExercise
}

type DeleteWorkoutSectionExerciseInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutSectionExerciseOutput struct{}
