package models

type WorkoutLogSectionEntity struct {
	BaseModel
	WorkoutLogID         uint               `gorm:"not null;index"`
	Type                 WorkoutSectionType `gorm:"not null"`
	Label                *string
	Position             int `gorm:"not null"`
	RestBetweenExercises *int
	Completed            bool                       `gorm:"not null;default:false"`
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
		Completed:            e.Completed,
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
		Completed:            dto.Completed,
	}
}
