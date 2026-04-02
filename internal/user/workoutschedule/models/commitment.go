package models

import (
	"time"

	"gesitr/internal/shared"
)

// ScheduleCommitment is the API-facing DTO for a schedule commitment.
// A commitment links a period to a workout log. It is created first
// (with workoutLogId=null) and linked to a WorkoutLog when the period activates.
type ScheduleCommitment struct {
	shared.BaseModel `tstype:",extends"`
	PeriodID         uint       `json:"periodId"`
	Date             *time.Time `json:"date"`
	WorkoutLogID     *uint      `json:"workoutLogId"`
}
