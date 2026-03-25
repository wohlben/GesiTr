package models

import (
	"time"

	"gesitr/internal/shared"
)

type WorkoutLogExercise struct {
	shared.BaseModel       `tstype:",extends"`
	WorkoutLogSectionID    uint                 `json:"workoutLogSectionId"`
	WorkoutLogID           uint                 `json:"workoutLogId"`
	SourceExerciseSchemeID uint                 `json:"sourceExerciseSchemeId"`
	SourceExerciseGroupID  *uint                `json:"sourceExerciseGroupId"`
	Position               int                  `json:"position"`
	Status                 WorkoutLogItemStatus `json:"status"`
	StatusChangedAt        *time.Time           `json:"statusChangedAt"`
	BreakAfterSeconds      *int                 `json:"breakAfterSeconds"`

	// Target fields (snapshotted from scheme on creation)
	TargetMeasurementType string `json:"targetMeasurementType"`
	TargetTimePerRep      *int   `json:"targetTimePerRep"`

	Sets []WorkoutLogExerciseSet `json:"sets" gorm:"-"`
}
