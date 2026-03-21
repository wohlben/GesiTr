package models

import (
	"time"

	exerciselogmodels "gesitr/internal/user/exerciselog/models"

	"gesitr/internal/shared"
)

type WorkoutLogExerciseSet struct {
	shared.BaseModel     `tstype:",extends"`
	WorkoutLogExerciseID uint                           `json:"workoutLogExerciseId"`
	WorkoutLogID         uint                           `json:"workoutLogId"`
	SetNumber            int                            `json:"setNumber"`
	Status               WorkoutLogItemStatus           `json:"status"`
	StatusChangedAt      *time.Time                     `json:"statusChangedAt"`
	BreakAfterSeconds    *int                           `json:"breakAfterSeconds"`
	TargetReps           *int                           `json:"targetReps"`
	TargetWeight         *float64                       `json:"targetWeight"`
	TargetDuration       *int                           `json:"targetDuration"`
	TargetDistance       *float64                       `json:"targetDistance"`
	TargetTime           *int                           `json:"targetTime"`
	ExerciseLog          *exerciselogmodels.ExerciseLog `json:"exerciseLog,omitempty"`
}
