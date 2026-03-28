package models

import (
	"time"

	"gesitr/internal/shared"
)

type WorkoutLog struct {
	shared.BaseModel `tstype:",extends"`
	Owner            string              `json:"owner"`
	WorkoutID        *uint               `json:"workoutId"`
	Name             string              `json:"name"`
	Notes            *string             `json:"notes"`
	Date             *time.Time          `json:"date"`
	Status           WorkoutLogStatus    `json:"status"`
	StatusChangedAt  *time.Time          `json:"statusChangedAt"`
	ScheduleID       *uint               `json:"scheduleId"`
	DueStart         *time.Time          `json:"dueStart"`
	DueEnd           *time.Time          `json:"dueEnd"`
	Sections         []WorkoutLogSection `json:"sections" gorm:"-"`
}
