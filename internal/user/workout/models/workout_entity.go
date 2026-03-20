package models

import "gesitr/internal/shared"

type WorkoutEntity struct {
	shared.BaseModel
	Owner    string `gorm:"not null;index"`
	Name     string `gorm:"not null"`
	Notes    *string
	Sections []WorkoutSectionEntity `gorm:"foreignKey:WorkoutID"`
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
	Exercises            []WorkoutSectionExerciseEntity `gorm:"foreignKey:WorkoutSectionID"`
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
	for _, ex := range e.Exercises {
		dto.Exercises = append(dto.Exercises, ex.ToDTO())
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

type WorkoutSectionExerciseEntity struct {
	shared.BaseModel
	WorkoutSectionID     uint `gorm:"not null;index"`
	UserExerciseSchemeID uint `gorm:"not null"`
	Position             int  `gorm:"not null"`
}

func (WorkoutSectionExerciseEntity) TableName() string { return "workout_section_exercises" }

func (e *WorkoutSectionExerciseEntity) ToDTO() WorkoutSectionExercise {
	return WorkoutSectionExercise{
		BaseModel:            e.BaseModel,
		WorkoutSectionID:     e.WorkoutSectionID,
		UserExerciseSchemeID: e.UserExerciseSchemeID,
		Position:             e.Position,
	}
}

func WorkoutSectionExerciseFromDTO(dto WorkoutSectionExercise) WorkoutSectionExerciseEntity {
	return WorkoutSectionExerciseEntity{
		BaseModel:            dto.BaseModel,
		WorkoutSectionID:     dto.WorkoutSectionID,
		UserExerciseSchemeID: dto.UserExerciseSchemeID,
		Position:             dto.Position,
	}
}
