package models

type WorkoutSectionExerciseEntity struct {
	BaseModel
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
