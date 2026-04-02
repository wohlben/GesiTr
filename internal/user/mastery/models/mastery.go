package models

type ExerciseMastery struct {
	ExerciseID   uint    `json:"exerciseId"`
	TotalXP      float64 `json:"totalXp"`
	EffectiveXP  float64 `json:"effectiveXp"`
	Level        int     `json:"level"`
	Tier         string  `json:"tier"`
	Progress     float64 `json:"progress"`
	DistinctDays int     `json:"distinctDays"`
	Multiplier   float64 `json:"multiplier"`
}

type EquipmentMastery struct {
	EquipmentID  uint    `json:"equipmentId"`
	TotalXP      float64 `json:"totalXp"`
	EffectiveXP  float64 `json:"effectiveXp"`
	Level        int     `json:"level"`
	Tier         string  `json:"tier"`
	Progress     float64 `json:"progress"`
	DistinctDays int     `json:"distinctDays"`
	Multiplier   float64 `json:"multiplier"`
}

type MasteryTier string

const (
	TierNovice     MasteryTier = "novice"
	TierJourneyman MasteryTier = "journeyman"
	TierAdept      MasteryTier = "adept"
	TierMaster     MasteryTier = "master"
	TierMastered   MasteryTier = "mastered"
)
