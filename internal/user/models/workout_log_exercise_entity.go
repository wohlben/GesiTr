package models

import "time"

type WorkoutLogExerciseEntity struct {
	BaseModel
	WorkoutLogSectionID    uint                 `gorm:"not null;index"`
	WorkoutLogID           uint                 `gorm:"not null;index"`
	SourceExerciseSchemeID uint                 `gorm:"not null"`
	Position               int                  `gorm:"not null"`
	Status                 WorkoutLogItemStatus `gorm:"not null;default:'planning'"`
	StatusChangedAt        *time.Time
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
		WorkoutLogID:           e.WorkoutLogID,
		SourceExerciseSchemeID: e.SourceExerciseSchemeID,
		Position:               e.Position,
		Status:                 e.Status,
		StatusChangedAt:        e.StatusChangedAt,
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
		WorkoutLogID:           dto.WorkoutLogID,
		SourceExerciseSchemeID: dto.SourceExerciseSchemeID,
		Position:               dto.Position,
		Status:                 dto.Status,
		StatusChangedAt:        dto.StatusChangedAt,
		BreakAfterSeconds:      dto.BreakAfterSeconds,
		TargetMeasurementType:  dto.TargetMeasurementType,
		TargetTimePerRep:       dto.TargetTimePerRep,
	}
}
