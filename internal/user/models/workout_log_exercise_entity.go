package models

type WorkoutLogExerciseEntity struct {
	BaseModel
	WorkoutLogSectionID  uint   `gorm:"not null;index"`
	UserExerciseSchemeID uint   `gorm:"not null"`
	Position             int    `gorm:"not null"`

	// Target fields (snapshotted from scheme on creation)
	TargetMeasurementType string `gorm:"not null"`
	TargetRestBetweenSets *int
	TargetTimePerRep      *int

	Sets []WorkoutLogExerciseSetEntity `gorm:"foreignKey:WorkoutLogExerciseID"`
}

func (WorkoutLogExerciseEntity) TableName() string { return "workout_log_exercises" }

func (e *WorkoutLogExerciseEntity) ToDTO() WorkoutLogExercise {
	dto := WorkoutLogExercise{
		BaseModel:             e.BaseModel,
		WorkoutLogSectionID:   e.WorkoutLogSectionID,
		UserExerciseSchemeID:  e.UserExerciseSchemeID,
		Position:              e.Position,
		TargetMeasurementType: e.TargetMeasurementType,
		TargetRestBetweenSets: e.TargetRestBetweenSets,
		TargetTimePerRep:      e.TargetTimePerRep,
	}
	for _, s := range e.Sets {
		dto.Sets = append(dto.Sets, s.ToDTO())
	}
	return dto
}

func WorkoutLogExerciseFromDTO(dto WorkoutLogExercise) WorkoutLogExerciseEntity {
	return WorkoutLogExerciseEntity{
		BaseModel:             dto.BaseModel,
		WorkoutLogSectionID:   dto.WorkoutLogSectionID,
		UserExerciseSchemeID:  dto.UserExerciseSchemeID,
		Position:              dto.Position,
		TargetMeasurementType: dto.TargetMeasurementType,
		TargetRestBetweenSets: dto.TargetRestBetweenSets,
		TargetTimePerRep:      dto.TargetTimePerRep,
	}
}
