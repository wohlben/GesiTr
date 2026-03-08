package models

type UserExerciseEntity struct {
	BaseModel
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
