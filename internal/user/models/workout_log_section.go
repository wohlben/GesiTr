package models

import "time"

type WorkoutLogSection struct {
	BaseModel            `tstype:",extends"`
	WorkoutLogID         uint                 `json:"workoutLogId"`
	Type                 WorkoutSectionType   `json:"type"`
	Label                *string              `json:"label"`
	Position             int                  `json:"position"`
	RestBetweenExercises *int                 `json:"restBetweenExercises"`
	Status               WorkoutLogItemStatus `json:"status"`
	StatusChangedAt      *time.Time           `json:"statusChangedAt"`
	Exercises            []WorkoutLogExercise `json:"exercises" gorm:"-"`
}
