package models

import (
	"time"

	"gesitr/internal/shared"
	"gesitr/internal/user/workout"
)

type WorkoutLogSection struct {
	shared.BaseModel     `tstype:",extends"`
	WorkoutLogID         uint                       `json:"workoutLogId"`
	Type                 workout.WorkoutSectionType `json:"type"`
	Label                *string                    `json:"label"`
	Position             int                        `json:"position"`
	RestBetweenExercises *int                       `json:"restBetweenExercises"`
	Status               WorkoutLogItemStatus       `json:"status"`
	StatusChangedAt      *time.Time                 `json:"statusChangedAt"`
	Exercises            []WorkoutLogExercise       `json:"exercises" gorm:"-"`
}

type WorkoutLogSectionEntity struct {
	shared.BaseModel
	WorkoutLogID         uint                       `gorm:"not null;index"`
	Type                 workout.WorkoutSectionType `gorm:"not null"`
	Label                *string
	Position             int `gorm:"not null"`
	RestBetweenExercises *int
	Status               WorkoutLogItemStatus `gorm:"not null;default:'planning'"`
	StatusChangedAt      *time.Time
	Exercises            []WorkoutLogExerciseEntity `gorm:"foreignKey:WorkoutLogSectionID"`
}

func (WorkoutLogSectionEntity) TableName() string { return "workout_log_sections" }

func (e *WorkoutLogSectionEntity) ToDTO() WorkoutLogSection {
	dto := WorkoutLogSection{
		BaseModel:            e.BaseModel,
		WorkoutLogID:         e.WorkoutLogID,
		Type:                 e.Type,
		Label:                e.Label,
		Position:             e.Position,
		RestBetweenExercises: e.RestBetweenExercises,
		Status:               e.Status,
		StatusChangedAt:      e.StatusChangedAt,
	}
	for _, ex := range e.Exercises {
		dto.Exercises = append(dto.Exercises, ex.ToDTO())
	}
	return dto
}

func WorkoutLogSectionFromDTO(dto WorkoutLogSection) WorkoutLogSectionEntity {
	return WorkoutLogSectionEntity{
		BaseModel:            dto.BaseModel,
		WorkoutLogID:         dto.WorkoutLogID,
		Type:                 dto.Type,
		Label:                dto.Label,
		Position:             dto.Position,
		RestBetweenExercises: dto.RestBetweenExercises,
		Status:               dto.Status,
		StatusChangedAt:      dto.StatusChangedAt,
	}
}
