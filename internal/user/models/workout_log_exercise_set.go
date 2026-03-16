package models

import "time"

type WorkoutLogExerciseSet struct {
	BaseModel            `tstype:",extends"`
	WorkoutLogExerciseID uint                 `json:"workoutLogExerciseId"`
	WorkoutLogID         uint                 `json:"workoutLogId"`
	SetNumber            int                  `json:"setNumber"`
	Status               WorkoutLogItemStatus `json:"status"`
	StatusChangedAt      *time.Time           `json:"statusChangedAt"`
	BreakAfterSeconds    *int                 `json:"breakAfterSeconds"`
	TargetReps           *int                 `json:"targetReps"`
	TargetWeight         *float64             `json:"targetWeight"`
	TargetDuration       *int                 `json:"targetDuration"`
	TargetDistance       *float64             `json:"targetDistance"`
	TargetTime           *int                 `json:"targetTime"`
	ActualReps           *int                 `json:"actualReps"`
	ActualWeight         *float64             `json:"actualWeight"`
	ActualDuration       *int                 `json:"actualDuration"`
	ActualDistance       *float64             `json:"actualDistance"`
	ActualTime           *int                 `json:"actualTime"`
}
