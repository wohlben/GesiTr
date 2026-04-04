package models

import (
	"gesitr/internal/shared"
)

type ExerciseSchemeSectionItemEntity struct {
	shared.BaseModel
	ExerciseSchemeID     uint   `gorm:"not null;index"`
	WorkoutSectionItemID uint   `gorm:"not null;uniqueIndex:idx_section_item_owner"`
	Owner                string `gorm:"not null;uniqueIndex:idx_section_item_owner"`
}

func (ExerciseSchemeSectionItemEntity) TableName() string {
	return "exercise_scheme_section_items"
}

func (e *ExerciseSchemeSectionItemEntity) ToDTO() ExerciseSchemeSectionItem {
	return ExerciseSchemeSectionItem{
		BaseModel:            e.BaseModel,
		ExerciseSchemeID:     e.ExerciseSchemeID,
		WorkoutSectionItemID: e.WorkoutSectionItemID,
		Owner:                e.Owner,
	}
}
