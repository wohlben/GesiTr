package models

import "time"

type WorkoutLogEntity struct {
	BaseModel
	Owner     string    `gorm:"not null;index"`
	WorkoutID *uint
	Name      string    `gorm:"not null"`
	Notes     *string
	Date      time.Time `gorm:"not null;index"`
	Sections  []WorkoutLogSectionEntity `gorm:"foreignKey:WorkoutLogID"`
}

func (WorkoutLogEntity) TableName() string { return "workout_logs" }

func (e *WorkoutLogEntity) ToDTO() WorkoutLog {
	dto := WorkoutLog{
		BaseModel: e.BaseModel,
		Owner:     e.Owner,
		WorkoutID: e.WorkoutID,
		Name:      e.Name,
		Notes:     e.Notes,
		Date:      e.Date,
	}
	for _, s := range e.Sections {
		dto.Sections = append(dto.Sections, s.ToDTO())
	}
	return dto
}

func WorkoutLogFromDTO(dto WorkoutLog) WorkoutLogEntity {
	return WorkoutLogEntity{
		BaseModel: dto.BaseModel,
		Owner:     dto.Owner,
		WorkoutID: dto.WorkoutID,
		Name:      dto.Name,
		Notes:     dto.Notes,
		Date:      dto.Date,
	}
}
