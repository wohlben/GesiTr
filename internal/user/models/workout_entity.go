package models

import "time"

type WorkoutEntity struct {
	BaseModel
	Owner    string     `gorm:"not null;index"`
	Name     string     `gorm:"not null"`
	Notes    *string
	Date     time.Time  `gorm:"not null;index"`
	Sections []WorkoutSectionEntity `gorm:"foreignKey:WorkoutID"`
}

func (WorkoutEntity) TableName() string { return "workouts" }

func (e *WorkoutEntity) ToDTO() Workout {
	dto := Workout{
		BaseModel: e.BaseModel,
		Owner:     e.Owner,
		Name:      e.Name,
		Notes:     e.Notes,
		Date:      e.Date,
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
		Date:      dto.Date,
	}
}
