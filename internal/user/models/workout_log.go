package models

import "time"

type WorkoutLog struct {
	BaseModel       `tstype:",extends"`
	Owner           string              `json:"owner"`
	WorkoutID       *uint               `json:"workoutId"`
	Name            string              `json:"name"`
	Notes           *string             `json:"notes"`
	Date            *time.Time          `json:"date"`
	Status          WorkoutLogStatus    `json:"status"`
	StatusChangedAt *time.Time          `json:"statusChangedAt"`
	Sections        []WorkoutLogSection `json:"sections" gorm:"-"`
}
