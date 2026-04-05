package models

import (
	"time"

	"gesitr/internal/shared"
)

// WorkoutSchedule is the API-facing DTO for a workout schedule.
type WorkoutSchedule struct {
	shared.BaseModel `tstype:",extends"`
	Owner            string     `json:"owner"`
	WorkoutID        uint       `json:"workoutId"`
	StartDate        time.Time  `json:"startDate"`
	EndDate          *time.Time `json:"endDate"`
	InitialStatus    string     `json:"initialStatus"`
	Timezone         string     `json:"timezone"`
	Active           bool       `json:"active"` // derived, not stored
}
