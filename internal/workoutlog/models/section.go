package models

import (
	"time"

	"gesitr/internal/shared"
	workoutmodels "gesitr/internal/workout/models"
)

type WorkoutLogSection struct {
	shared.BaseModel     `tstype:",extends"`
	WorkoutLogID         uint                             `json:"workoutLogId"`
	Type                 workoutmodels.WorkoutSectionType `json:"type"`
	Label                *string                          `json:"label"`
	Position             int                              `json:"position"`
	RestBetweenExercises *int                             `json:"restBetweenExercises"`
	Status               WorkoutLogItemStatus             `json:"status"`
	StatusChangedAt      *time.Time                       `json:"statusChangedAt"`
	Exercises            []WorkoutLogExercise             `json:"exercises" gorm:"-"`
}
