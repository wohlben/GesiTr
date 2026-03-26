package models

import "gesitr/internal/shared"

// ExerciseScheme is the API DTO for exercise schemes
type ExerciseScheme struct {
	shared.BaseModel     `tstype:",extends"`
	Owner                string   `json:"owner"`
	ExerciseID           uint     `json:"exerciseId"`
	MeasurementType      string   `json:"measurementType"`
	Sets                 *int     `json:"sets"`
	Reps                 *int     `json:"reps"`
	Weight               *float64 `json:"weight"`
	RestBetweenSets      *int     `json:"restBetweenSets"`
	TimePerRep           *int     `json:"timePerRep"`
	Duration             *int     `json:"duration"`
	Distance             *float64 `json:"distance"`
	TargetTime           *int     `json:"targetTime"`
	WorkoutSectionItemID *uint    `json:"workoutSectionItemId"`
}
