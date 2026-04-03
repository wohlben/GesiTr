package models

// ExerciseNamePreference stores a user's preferred display name for an exercise.
// References exercise_names.id via FK for referential integrity.
type ExerciseNamePreference struct {
	Owner          string `gorm:"primaryKey"`
	ExerciseID     uint   `gorm:"primaryKey"`
	ExerciseNameID uint   `gorm:"not null"`
}

type ExerciseNamePreferenceDTO struct {
	ExerciseID     uint `json:"exerciseId"`
	ExerciseNameID uint `json:"exerciseNameId"`
}

func (e *ExerciseNamePreference) ToDTO() ExerciseNamePreferenceDTO {
	return ExerciseNamePreferenceDTO{
		ExerciseID:     e.ExerciseID,
		ExerciseNameID: e.ExerciseNameID,
	}
}
