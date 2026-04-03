package handlers

import "gesitr/internal/user/namepreference/models"

type ListExerciseNamePreferencesOutput struct {
	Body []models.ExerciseNamePreferenceDTO
}

type SetExerciseNamePreferenceInput struct {
	ExerciseID uint `path:"exerciseId"`
	Body       struct {
		ExerciseNameID uint `json:"exerciseNameId" required:"true"`
	}
}

type SetExerciseNamePreferenceOutput struct {
	Body models.ExerciseNamePreferenceDTO
}

type DeleteExerciseNamePreferenceInput struct {
	ExerciseID uint `path:"exerciseId"`
}

type DeleteExerciseNamePreferenceOutput struct{}
