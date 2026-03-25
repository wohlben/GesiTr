package models

import (
	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type WorkoutEntity struct {
	shared.BaseModel
	Owner        string                           `gorm:"not null;index"`
	OwnerProfile *profilemodels.UserProfileEntity `gorm:"foreignKey:Owner;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
	Name         string                           `gorm:"not null"`
	Notes        *string
	Sections     []WorkoutSectionEntity `gorm:"foreignKey:WorkoutID"`
}

func (WorkoutEntity) TableName() string { return "workouts" }

func (e *WorkoutEntity) ToDTO() Workout {
	dto := Workout{
		BaseModel: e.BaseModel,
		Owner:     e.Owner,
		Name:      e.Name,
		Notes:     e.Notes,
	}
	for _, s := range e.Sections {
		dto.Sections = append(dto.Sections, s.ToDTO())
	}
	return dto
}

func WorkoutFromDTO(dto Workout) WorkoutEntity {
	return WorkoutEntity{
		BaseModel: dto.BaseModel,
		Owner:     dto.Owner,
		Name:      dto.Name,
		Notes:     dto.Notes,
	}
}

type WorkoutSectionEntity struct {
	shared.BaseModel
	WorkoutID            uint               `gorm:"not null;index"`
	Type                 WorkoutSectionType `gorm:"not null"`
	Label                *string
	Position             int `gorm:"not null"`
	RestBetweenExercises *int
	Items                []WorkoutSectionItemEntity `gorm:"foreignKey:WorkoutSectionID"`
}

func (WorkoutSectionEntity) TableName() string { return "workout_sections" }

func (e *WorkoutSectionEntity) ToDTO() WorkoutSection {
	dto := WorkoutSection{
		BaseModel:            e.BaseModel,
		WorkoutID:            e.WorkoutID,
		Type:                 e.Type,
		Label:                e.Label,
		Position:             e.Position,
		RestBetweenExercises: e.RestBetweenExercises,
	}
	for _, item := range e.Items {
		dto.Items = append(dto.Items, item.ToDTO())
	}
	return dto
}

func WorkoutSectionFromDTO(dto WorkoutSection) WorkoutSectionEntity {
	return WorkoutSectionEntity{
		BaseModel:            dto.BaseModel,
		WorkoutID:            dto.WorkoutID,
		Type:                 dto.Type,
		Label:                dto.Label,
		Position:             dto.Position,
		RestBetweenExercises: dto.RestBetweenExercises,
	}
}

type WorkoutSectionItemEntity struct {
	shared.BaseModel
	WorkoutSectionID uint                   `gorm:"not null;index"`
	Type             WorkoutSectionItemType `gorm:"not null;default:'exercise'"`
	ExerciseSchemeID *uint
	ExerciseGroupID  *uint
	Data             *string
	Position         int `gorm:"not null"`
}

func (WorkoutSectionItemEntity) TableName() string { return "workout_section_items" }

func (e *WorkoutSectionItemEntity) ToDTO() WorkoutSectionItem {
	return WorkoutSectionItem{
		BaseModel:        e.BaseModel,
		WorkoutSectionID: e.WorkoutSectionID,
		Type:             e.Type,
		ExerciseSchemeID: e.ExerciseSchemeID,
		ExerciseGroupID:  e.ExerciseGroupID,
		Data:             e.Data,
		Position:         e.Position,
	}
}

func WorkoutSectionItemFromDTO(dto WorkoutSectionItem) WorkoutSectionItemEntity {
	return WorkoutSectionItemEntity{
		BaseModel:        dto.BaseModel,
		WorkoutSectionID: dto.WorkoutSectionID,
		Type:             dto.Type,
		ExerciseSchemeID: dto.ExerciseSchemeID,
		ExerciseGroupID:  dto.ExerciseGroupID,
		Data:             dto.Data,
		Position:         dto.Position,
	}
}
