package models

type UserExerciseSchemeEntity struct {
	BaseModel
	UserExerciseID uint     `gorm:"not null;index"`
	MeasurementType string  `gorm:"not null"`
	Sets            *int
	Reps            *int
	Weight          *float64
	RestBetweenSets *int
	TimePerRep      *int
	Duration        *int
	Distance        *float64
	TargetTime      *int
}

func (UserExerciseSchemeEntity) TableName() string { return "user_exercise_schemes" }

func (e *UserExerciseSchemeEntity) ToDTO() UserExerciseScheme {
	return UserExerciseScheme{
		BaseModel:       e.BaseModel,
		UserExerciseID:  e.UserExerciseID,
		MeasurementType: e.MeasurementType,
		Sets:            e.Sets,
		Reps:            e.Reps,
		Weight:          e.Weight,
		RestBetweenSets: e.RestBetweenSets,
		TimePerRep:      e.TimePerRep,
		Duration:        e.Duration,
		Distance:        e.Distance,
		TargetTime:      e.TargetTime,
	}
}

func UserExerciseSchemeFromDTO(dto UserExerciseScheme) UserExerciseSchemeEntity {
	return UserExerciseSchemeEntity{
		BaseModel:       dto.BaseModel,
		UserExerciseID:  dto.UserExerciseID,
		MeasurementType: dto.MeasurementType,
		Sets:            dto.Sets,
		Reps:            dto.Reps,
		Weight:          dto.Weight,
		RestBetweenSets: dto.RestBetweenSets,
		TimePerRep:      dto.TimePerRep,
		Duration:        dto.Duration,
		Distance:        dto.Distance,
		TargetTime:      dto.TargetTime,
	}
}
