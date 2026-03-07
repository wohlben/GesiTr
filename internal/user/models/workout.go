package models

import "time"

type Workout struct {
	BaseModel `tstype:",extends"`
	Owner     string           `json:"owner"`
	Name      string           `json:"name"`
	Notes     *string          `json:"notes"`
	Date      time.Time        `json:"date"`
	Sections  []WorkoutSection `json:"sections" gorm:"-"`
}
