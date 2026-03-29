package models

type MasteryContributionEntity struct {
	Owner             string  `gorm:"primaryKey"`
	ExerciseID        uint    `gorm:"primaryKey"`
	ContributesFromID uint    `gorm:"primaryKey"`
	Multiplier        float64 `gorm:"not null"`
	RelationshipType  string  `gorm:"not null"`
}

func (MasteryContributionEntity) TableName() string { return "mastery_contributions" }
