package handlers

import (
	"gesitr/internal/user/mastery/models"
)

type ListMasteryInput struct{}

type ListMasteryOutput struct {
	Body []models.ExerciseMastery
}

type GetMasteryInput struct {
	ExerciseID uint `path:"exerciseId"`
}

type GetMasteryOutput struct {
	Body models.ExerciseMastery
}
