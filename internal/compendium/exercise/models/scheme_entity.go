package models

import (
	profilemodels "gesitr/internal/profile/models"
	"gesitr/internal/shared"
)

type ExerciseSchemeEntity struct {
	shared.BaseModel
	Owner                string                           `gorm:"not null;index"`
	OwnerProfile         *profilemodels.UserProfileEntity `gorm:"foreignKey:Owner;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
	ExerciseID           uint                             `gorm:"not null;index"`
	MeasurementType      string                           `gorm:"not null"`
	Sets                 *int
	Reps                 *int
	Weight               *float64
	RestBetweenSets      *int
	TimePerRep           *int
	Duration             *int
	Distance             *float64
	TargetTime           *int
	WorkoutSectionItemID *uint `gorm:"index"`
}

func (ExerciseSchemeEntity) TableName() string { return "exercise_schemes" }

func (e *ExerciseSchemeEntity) ToDTO() ExerciseScheme {
	return ExerciseScheme{
		BaseModel:            e.BaseModel,
		Owner:                e.Owner,
		ExerciseID:           e.ExerciseID,
		MeasurementType:      e.MeasurementType,
		Sets:                 e.Sets,
		Reps:                 e.Reps,
		Weight:               e.Weight,
		RestBetweenSets:      e.RestBetweenSets,
		TimePerRep:           e.TimePerRep,
		Duration:             e.Duration,
		Distance:             e.Distance,
		TargetTime:           e.TargetTime,
		WorkoutSectionItemID: e.WorkoutSectionItemID,
	}
}

func ExerciseSchemeFromDTO(dto ExerciseScheme) ExerciseSchemeEntity {
	return ExerciseSchemeEntity{
		BaseModel:            dto.BaseModel,
		Owner:                dto.Owner,
		ExerciseID:           dto.ExerciseID,
		MeasurementType:      dto.MeasurementType,
		Sets:                 dto.Sets,
		Reps:                 dto.Reps,
		Weight:               dto.Weight,
		RestBetweenSets:      dto.RestBetweenSets,
		TimePerRep:           dto.TimePerRep,
		Duration:             dto.Duration,
		Distance:             dto.Distance,
		TargetTime:           dto.TargetTime,
		WorkoutSectionItemID: dto.WorkoutSectionItemID,
	}
}
