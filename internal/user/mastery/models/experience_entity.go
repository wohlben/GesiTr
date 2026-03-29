package models

type MasteryExperienceEntity struct {
	Owner      string `gorm:"primaryKey"`
	ExerciseID uint   `gorm:"primaryKey"`
	TotalReps  int    `gorm:"not null;default:0"`
}

func (MasteryExperienceEntity) TableName() string { return "mastery_experience" }
