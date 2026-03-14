package models

type WorkoutLogExerciseSet struct {
	BaseModel            `tstype:",extends"`
	WorkoutLogExerciseID uint     `json:"workoutLogExerciseId"`
	SetNumber            int      `json:"setNumber"`
	Completed            bool     `json:"completed"`
	BreakAfterSeconds    *int     `json:"breakAfterSeconds"`
	TargetReps           *int     `json:"targetReps"`
	TargetWeight         *float64 `json:"targetWeight"`
	TargetDuration       *int     `json:"targetDuration"`
	TargetDistance       *float64 `json:"targetDistance"`
	TargetTime           *int     `json:"targetTime"`
	ActualReps           *int     `json:"actualReps"`
	ActualWeight         *float64 `json:"actualWeight"`
	ActualDuration       *int     `json:"actualDuration"`
	ActualDistance       *float64 `json:"actualDistance"`
	ActualTime           *int     `json:"actualTime"`
}
