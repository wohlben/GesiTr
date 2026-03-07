package models

type UserRecord struct {
	BaseModel                `tstype:",extends"`
	UserExerciseID           uint     `json:"userExerciseId"`
	MeasurementType          string   `json:"measurementType"`
	RecordValue              float64  `json:"recordValue"`
	ActualReps               *int     `json:"actualReps"`
	ActualWeight             *float64 `json:"actualWeight"`
	ActualDuration           *int     `json:"actualDuration"`
	ActualDistance            *float64 `json:"actualDistance"`
	ActualTime               *int     `json:"actualTime"`
	WorkoutLogExerciseSetID  uint     `json:"workoutLogExerciseSetId"`
}
