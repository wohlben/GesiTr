package models

type WorkoutSectionExercise struct {
	BaseModel            `tstype:",extends"`
	WorkoutSectionID     uint `json:"workoutSectionId"`
	UserExerciseSchemeID uint `json:"userExerciseSchemeId"`
	Position             int  `json:"position"`
}
