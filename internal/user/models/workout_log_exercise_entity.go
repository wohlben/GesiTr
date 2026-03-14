package models

type WorkoutLogExerciseEntity struct {
	BaseModel
	WorkoutLogSectionID    uint `gorm:"not null;index"`
	SourceExerciseSchemeID uint `gorm:"not null"`
	Position               int  `gorm:"not null"`
	Completed              bool `gorm:"not null;default:false"`
	BreakAfterSeconds      *int

	// Target fields (snapshotted from scheme on creation)
	TargetMeasurementType string `gorm:"not null"`
	TargetTimePerRep      *int

	Sets []WorkoutLogExerciseSetEntity `gorm:"foreignKey:WorkoutLogExerciseID"`
}

func (WorkoutLogExerciseEntity) TableName() string { return "workout_log_exercises" }

func (e *WorkoutLogExerciseEntity) ToDTO() WorkoutLogExercise {
	dto := WorkoutLogExercise{
		BaseModel:              e.BaseModel,
		WorkoutLogSectionID:    e.WorkoutLogSectionID,
		SourceExerciseSchemeID: e.SourceExerciseSchemeID,
		Position:               e.Position,
		Completed:              e.Completed,
		BreakAfterSeconds:      e.BreakAfterSeconds,
		TargetMeasurementType:  e.TargetMeasurementType,
		TargetTimePerRep:       e.TargetTimePerRep,
	}
	for _, s := range e.Sets {
		dto.Sets = append(dto.Sets, s.ToDTO())
	}
	return dto
}

func WorkoutLogExerciseFromDTO(dto WorkoutLogExercise) WorkoutLogExerciseEntity {
	return WorkoutLogExerciseEntity{
		BaseModel:              dto.BaseModel,
		WorkoutLogSectionID:    dto.WorkoutLogSectionID,
		SourceExerciseSchemeID: dto.SourceExerciseSchemeID,
		Position:               dto.Position,
		Completed:              dto.Completed,
		BreakAfterSeconds:      dto.BreakAfterSeconds,
		TargetMeasurementType:  dto.TargetMeasurementType,
		TargetTimePerRep:       dto.TargetTimePerRep,
	}
}
