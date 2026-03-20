package models

import "gesitr/internal/shared"

type UserExerciseEntity struct {
	shared.BaseModel
	Owner                string `gorm:"not null;index;uniqueIndex:idx_owner_compendium_exercise"`
	CompendiumExerciseID string `gorm:"not null;uniqueIndex:idx_owner_compendium_exercise"`
	CompendiumVersion    int    `gorm:"not null"`
}

func (UserExerciseEntity) TableName() string { return "user_exercises" }

func (e *UserExerciseEntity) ToDTO() UserExercise {
	return UserExercise{
		BaseModel:            e.BaseModel,
		Owner:                e.Owner,
		CompendiumExerciseID: e.CompendiumExerciseID,
		CompendiumVersion:    e.CompendiumVersion,
	}
}

func UserExerciseFromDTO(dto UserExercise) UserExerciseEntity {
	return UserExerciseEntity{
		BaseModel:            dto.BaseModel,
		Owner:                dto.Owner,
		CompendiumExerciseID: dto.CompendiumExerciseID,
		CompendiumVersion:    dto.CompendiumVersion,
	}
}

type UserEquipmentEntity struct {
	shared.BaseModel
	Owner                 string `gorm:"not null;index;uniqueIndex:idx_owner_compendium_equipment"`
	CompendiumEquipmentID string `gorm:"not null;uniqueIndex:idx_owner_compendium_equipment"`
	CompendiumVersion     int    `gorm:"not null"`
}

func (UserEquipmentEntity) TableName() string { return "user_equipment" }

func (e *UserEquipmentEntity) ToDTO() UserEquipment {
	return UserEquipment{
		BaseModel:             e.BaseModel,
		Owner:                 e.Owner,
		CompendiumEquipmentID: e.CompendiumEquipmentID,
		CompendiumVersion:     e.CompendiumVersion,
	}
}

func UserEquipmentFromDTO(dto UserEquipment) UserEquipmentEntity {
	return UserEquipmentEntity{
		BaseModel:             dto.BaseModel,
		Owner:                 dto.Owner,
		CompendiumEquipmentID: dto.CompendiumEquipmentID,
		CompendiumVersion:     dto.CompendiumVersion,
	}
}

type UserExerciseSchemeEntity struct {
	shared.BaseModel
	UserExerciseID  uint   `gorm:"not null;index"`
	MeasurementType string `gorm:"not null"`
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
