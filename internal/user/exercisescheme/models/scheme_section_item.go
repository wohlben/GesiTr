package models

import "gesitr/internal/shared"

// ExerciseSchemeSectionItem is the API DTO for the join between exercise schemes and workout section items.
type ExerciseSchemeSectionItem struct {
	shared.BaseModel     `tstype:",extends"`
	ExerciseSchemeID     uint   `json:"exerciseSchemeId"`
	WorkoutSectionItemID uint   `json:"workoutSectionItemId"`
	Owner                string `json:"owner"`
}
