package models

type UserExerciseEntity struct {
	BaseModel
	Owner              string `gorm:"not null;index;uniqueIndex:idx_owner_exercise_template"`
	ExerciseTemplateID string `gorm:"not null;uniqueIndex:idx_owner_exercise_template"`
	CompendiumVersion  int    `gorm:"not null"`
}

func (UserExerciseEntity) TableName() string { return "user_exercises" }

func (e *UserExerciseEntity) ToDTO() UserExercise {
	return UserExercise{
		BaseModel:          e.BaseModel,
		Owner:              e.Owner,
		ExerciseTemplateID: e.ExerciseTemplateID,
		CompendiumVersion:  e.CompendiumVersion,
	}
}

func UserExerciseFromDTO(dto UserExercise) UserExerciseEntity {
	return UserExerciseEntity{
		BaseModel:          dto.BaseModel,
		Owner:              dto.Owner,
		ExerciseTemplateID: dto.ExerciseTemplateID,
		CompendiumVersion:  dto.CompendiumVersion,
	}
}
