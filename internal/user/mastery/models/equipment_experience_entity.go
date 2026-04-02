package models

type EquipmentMasteryExperienceEntity struct {
	Owner       string `gorm:"primaryKey"`
	EquipmentID uint   `gorm:"primaryKey"`
	TotalReps   int    `gorm:"not null;default:0"`
}

func (EquipmentMasteryExperienceEntity) TableName() string { return "equipment_mastery_experience" }
