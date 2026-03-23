package models

import (
	"time"

	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type ExerciseLogEntity struct {
	shared.BaseModel
	Owner                   string                           `gorm:"not null;index"`
	OwnerProfile            *profilemodels.UserProfileEntity `gorm:"foreignKey:Owner;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
	ExerciseID              uint                             `gorm:"not null;index"`
	MeasurementType         string                           `gorm:"not null"`
	Reps                    *int
	Weight                  *float64
	Duration                *int
	Distance                *float64
	Time                    *int
	RecordValue             float64   `gorm:"not null"`
	IsRecord                bool      `gorm:"not null;default:false;index"`
	PerformedAt             time.Time `gorm:"not null;index"`
	WorkoutLogExerciseSetID *uint     `gorm:"index"`
	SourceExerciseSchemeID  *uint
}

func (ExerciseLogEntity) TableName() string { return "exercise_logs" }

func (e *ExerciseLogEntity) ToDTO() ExerciseLog {
	return ExerciseLog{
		BaseModel:               e.BaseModel,
		Owner:                   e.Owner,
		ExerciseID:              e.ExerciseID,
		MeasurementType:         e.MeasurementType,
		Reps:                    e.Reps,
		Weight:                  e.Weight,
		Duration:                e.Duration,
		Distance:                e.Distance,
		Time:                    e.Time,
		RecordValue:             e.RecordValue,
		IsRecord:                e.IsRecord,
		PerformedAt:             e.PerformedAt,
		WorkoutLogExerciseSetID: e.WorkoutLogExerciseSetID,
		SourceExerciseSchemeID:  e.SourceExerciseSchemeID,
	}
}

func ExerciseLogFromDTO(dto ExerciseLog) ExerciseLogEntity {
	return ExerciseLogEntity{
		BaseModel:               dto.BaseModel,
		Owner:                   dto.Owner,
		ExerciseID:              dto.ExerciseID,
		MeasurementType:         dto.MeasurementType,
		Reps:                    dto.Reps,
		Weight:                  dto.Weight,
		Duration:                dto.Duration,
		Distance:                dto.Distance,
		Time:                    dto.Time,
		RecordValue:             dto.RecordValue,
		IsRecord:                dto.IsRecord,
		PerformedAt:             dto.PerformedAt,
		WorkoutLogExerciseSetID: dto.WorkoutLogExerciseSetID,
		SourceExerciseSchemeID:  dto.SourceExerciseSchemeID,
	}
}
