package models

import "time"

type WorkoutLog struct {
	BaseModel `tstype:",extends"`
	Owner     string              `json:"owner"`
	WorkoutID *uint               `json:"workoutId"`
	Name      string              `json:"name"`
	Notes     *string             `json:"notes"`
	Date      time.Time           `json:"date"`
	Sections  []WorkoutLogSection `json:"sections" gorm:"-"`
}
