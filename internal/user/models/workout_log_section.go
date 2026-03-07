package models

type WorkoutLogSection struct {
	BaseModel            `tstype:",extends"`
	WorkoutLogID         uint                   `json:"workoutLogId"`
	Type                 WorkoutSectionType     `json:"type"`
	Label                *string                `json:"label"`
	Position             int                    `json:"position"`
	RestBetweenExercises *int                   `json:"restBetweenExercises"`
	Exercises            []WorkoutLogExercise   `json:"exercises" gorm:"-"`
}
