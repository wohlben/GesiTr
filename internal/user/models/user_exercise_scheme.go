package models

type UserExerciseScheme struct {
	BaseModel       `tstype:",extends"`
	UserExerciseID  uint     `json:"userExerciseId"`
	MeasurementType string   `json:"measurementType"`
	Sets            *int     `json:"sets"`
	Reps            *int     `json:"reps"`
	Weight          *float64 `json:"weight"`
	RestBetweenSets *int     `json:"restBetweenSets"`
	TimePerRep      *int     `json:"timePerRep"`
	Duration        *int     `json:"duration"`
	Distance        *float64 `json:"distance"`
	TargetTime      *int     `json:"targetTime"`
}
