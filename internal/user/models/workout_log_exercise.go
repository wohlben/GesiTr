package models

import "time"

type WorkoutLogExercise struct {
	BaseModel              `tstype:",extends"`
	WorkoutLogSectionID    uint             `json:"workoutLogSectionId"`
	SourceExerciseSchemeID uint             `json:"sourceExerciseSchemeId"`
	Position               int              `json:"position"`
	Status                 WorkoutLogStatus `json:"status"`
	StatusChangedAt        *time.Time       `json:"statusChangedAt"`
	BreakAfterSeconds      *int             `json:"breakAfterSeconds"`

	// Target fields (snapshotted from scheme on creation)
	TargetMeasurementType string `json:"targetMeasurementType"`
	TargetTimePerRep      *int   `json:"targetTimePerRep"`

	Sets []WorkoutLogExerciseSet `json:"sets" gorm:"-"`
}
