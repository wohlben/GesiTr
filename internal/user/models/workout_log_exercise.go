package models

type WorkoutLogExercise struct {
	BaseModel             `tstype:",extends"`
	WorkoutLogSectionID   uint   `json:"workoutLogSectionId"`
	UserExerciseSchemeID  uint   `json:"userExerciseSchemeId"`
	Position              int    `json:"position"`

	// Target fields (snapshotted from scheme on creation)
	TargetMeasurementType string `json:"targetMeasurementType"`
	TargetRestBetweenSets *int   `json:"targetRestBetweenSets"`
	TargetTimePerRep      *int   `json:"targetTimePerRep"`

	Sets []WorkoutLogExerciseSet `json:"sets" gorm:"-"`
}
