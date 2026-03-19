package models

import (
	"time"

	"gesitr/internal/shared"
)

type WorkoutLog struct {
	shared.BaseModel `tstype:",extends"`
	Owner            string              `json:"owner"`
	WorkoutID        *uint               `json:"workoutId"`
	Name             string              `json:"name"`
	Notes            *string             `json:"notes"`
	Date             *time.Time          `json:"date"`
	Status           WorkoutLogStatus    `json:"status"`
	StatusChangedAt  *time.Time          `json:"statusChangedAt"`
	Sections         []WorkoutLogSection `json:"sections" gorm:"-"`
}

type WorkoutLogEntity struct {
	shared.BaseModel
	Owner           string `gorm:"not null;index"`
	WorkoutID       *uint
	Name            string `gorm:"not null"`
	Notes           *string
	Date            *time.Time       `gorm:"index"`
	Status          WorkoutLogStatus `gorm:"not null;default:'planning'"`
	StatusChangedAt *time.Time
	Sections        []WorkoutLogSectionEntity `gorm:"foreignKey:WorkoutLogID"`
}

func (WorkoutLogEntity) TableName() string { return "workout_logs" }

func (e *WorkoutLogEntity) ToDTO() WorkoutLog {
	dto := WorkoutLog{
		BaseModel:       e.BaseModel,
		Owner:           e.Owner,
		WorkoutID:       e.WorkoutID,
		Name:            e.Name,
		Notes:           e.Notes,
		Date:            e.Date,
		Status:          e.Status,
		StatusChangedAt: e.StatusChangedAt,
	}
	for _, s := range e.Sections {
		dto.Sections = append(dto.Sections, s.ToDTO())
	}
	return dto
}

func WorkoutLogFromDTO(dto WorkoutLog) WorkoutLogEntity {
	return WorkoutLogEntity{
		BaseModel:       dto.BaseModel,
		Owner:           dto.Owner,
		WorkoutID:       dto.WorkoutID,
		Name:            dto.Name,
		Notes:           dto.Notes,
		Date:            dto.Date,
		Status:          dto.Status,
		StatusChangedAt: dto.StatusChangedAt,
	}
}
