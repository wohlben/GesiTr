package models

type EquipmentMasteryContributionEntity struct {
	Owner             string  `gorm:"primaryKey"`
	EquipmentID       uint    `gorm:"primaryKey"`
	ContributesFromID uint    `gorm:"primaryKey"`
	Multiplier        float64 `gorm:"not null"`
	RelationshipType  string  `gorm:"not null"`
}

func (EquipmentMasteryContributionEntity) TableName() string {
	return "equipment_mastery_contributions"
}
