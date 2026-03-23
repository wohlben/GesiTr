package models

import (
	"time"

	"gesitr/internal/shared"
)

type ExerciseLog struct {
	shared.BaseModel        `tstype:",extends"`
	Owner                   string    `json:"owner"`
	ExerciseID              uint      `json:"exerciseId"`
	MeasurementType         string    `json:"measurementType"`
	Reps                    *int      `json:"reps"`
	Weight                  *float64  `json:"weight"`
	Duration                *int      `json:"duration"`
	Distance                *float64  `json:"distance"`
	Time                    *int      `json:"time"`
	RecordValue             float64   `json:"recordValue"`
	IsRecord                bool      `json:"isRecord"`
	PerformedAt             time.Time `json:"performedAt"`
	WorkoutLogExerciseSetID *uint     `json:"workoutLogExerciseSetId"`
	SourceExerciseSchemeID  *uint     `json:"sourceExerciseSchemeId"`
}

// ComputeRecordValue calculates the comparable record value for an exercise log entry.
// REP_BASED uses Brzycki e1RM formula; others use the raw primary value.
func ComputeRecordValue(measurementType string, reps *int, weight *float64, duration *int, distance *float64) (float64, bool) {
	switch measurementType {
	case "REP_BASED", "AMRAP":
		if reps == nil || weight == nil || *weight <= 0 {
			return 0, false
		}
		return *weight * (1 + float64(*reps)/30), true
	case "TIME_BASED", "TIME", "EMOM", "ROUNDS_FOR_TIME":
		if duration == nil {
			return 0, false
		}
		return float64(*duration), true
	case "DISTANCE_BASED", "DISTANCE":
		if distance == nil {
			return 0, false
		}
		return *distance, true
	default:
		return 0, false
	}
}
