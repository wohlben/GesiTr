package models

type WorkoutSection struct {
	BaseModel            `tstype:",extends"`
	WorkoutID            uint                     `json:"workoutId"`
	Type                 WorkoutSectionType       `json:"type"`
	Label                *string                  `json:"label"`
	Position             int                      `json:"position"`
	RestBetweenExercises *int                     `json:"restBetweenExercises"`
	Exercises            []WorkoutSectionExercise `json:"exercises" gorm:"-"`
}
