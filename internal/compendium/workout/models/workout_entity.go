package models

import (
	"gesitr/internal/shared"
)

type WorkoutEntity struct {
	shared.BaseModel
	OwnershipGroupID uint   `gorm:"index"`
	Name             string `gorm:"not null"`
	Notes            *string
	Public           bool                   `gorm:"not null;default:false;index"`
	Version          int                    `gorm:"not null;default:0"`
	Sections         []WorkoutSectionEntity `gorm:"foreignKey:WorkoutID"`
}

func (WorkoutEntity) TableName() string { return "workouts" }

func (e *WorkoutEntity) ToDTO() Workout {
	dto := Workout{
		BaseModel:        e.BaseModel,
		OwnershipGroupID: e.OwnershipGroupID,
		Name:             e.Name,
		Notes:            e.Notes,
		Public:           e.Public,
		Version:          e.Version,
	}
	for _, s := range e.Sections {
		dto.Sections = append(dto.Sections, s.ToDTO())
	}
	return dto
}

func WorkoutFromDTO(dto Workout) WorkoutEntity {
	return WorkoutEntity{
		BaseModel:        dto.BaseModel,
		OwnershipGroupID: dto.OwnershipGroupID,
		Name:             dto.Name,
		Notes:            dto.Notes,
		Public:           dto.Public,
		Version:          dto.Version,
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
	ExerciseID       *uint
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
		ExerciseID:       e.ExerciseID,
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
		ExerciseID:       dto.ExerciseID,
		ExerciseGroupID:  dto.ExerciseGroupID,
		Data:             dto.Data,
		Position:         dto.Position,
	}
}
