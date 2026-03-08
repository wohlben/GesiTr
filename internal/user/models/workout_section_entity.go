package models

type WorkoutSectionType string

const (
	WorkoutSectionTypeMain          WorkoutSectionType = "main"
	WorkoutSectionTypeSupplementary WorkoutSectionType = "supplementary"
)

type WorkoutSectionEntity struct {
	BaseModel
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
