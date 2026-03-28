package models

import (
	"time"

	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type WorkoutLogEntity struct {
	shared.BaseModel
	Owner           string                           `gorm:"not null;index"`
	OwnerProfile    *profilemodels.UserProfileEntity `gorm:"foreignKey:Owner;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
	WorkoutID       *uint
	Name            string `gorm:"not null"`
	Notes           *string
	Date            *time.Time       `gorm:"index"`
	Status          WorkoutLogStatus `gorm:"not null;default:'planning'"`
	StatusChangedAt *time.Time
	ScheduleID      *uint `gorm:"index"`
	DueStart        *time.Time
	DueEnd          *time.Time
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
		ScheduleID:      e.ScheduleID,
		DueStart:        e.DueStart,
		DueEnd:          e.DueEnd,
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
		ScheduleID:      dto.ScheduleID,
		DueStart:        dto.DueStart,
		DueEnd:          dto.DueEnd,
	}
}
