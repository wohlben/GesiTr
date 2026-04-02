package models

import (
	"math"

	equipmentrelmodels "gesitr/internal/equipmentrelationship/models"
	exerciserelmodels "gesitr/internal/exerciserelationship/models"
)

const (
	xpPerLevel       = 100.0
	maxLevel         = 100
	recencyWindowStr = "6 months"
)

// ComputeLevel returns the level for a given total XP, capped at maxLevel.
func ComputeLevel(totalXP float64) int {
	level := int(totalXP / xpPerLevel)
	if level > maxLevel {
		return maxLevel
	}
	return level
}

// ComputeTier returns the mastery tier for a given level.
func ComputeTier(level int) MasteryTier {
	switch {
	case level >= 100:
		return TierMastered
	case level >= 51:
		return TierMaster
	case level >= 31:
		return TierAdept
	case level >= 11:
		return TierJourneyman
	default:
		return TierNovice
	}
}

// ComputeRecencyMultiplier computes the recency-based XP multiplier.
// nDays is the number of distinct days the exercise was performed in the recency window.
// baseLevel is the level computed from raw XP (without recency), used as the cap to avoid circular dependency.
func ComputeRecencyMultiplier(nDays, baseLevel int) float64 {
	cap := math.Max(1, float64(baseLevel)/2)
	mult := 0.5 * float64(nDays)
	return math.Min(mult, cap)
}

// ComputeProgress returns a 0.0-1.0 value representing progress within the current level.
func ComputeProgress(effectiveXP float64, level int) float64 {
	if level >= maxLevel {
		return 1.0
	}
	xpIntoLevel := effectiveXP - float64(level)*xpPerLevel
	return xpIntoLevel / xpPerLevel
}

// RelationshipTypeBonus returns the type bonus for a relationship type and whether it contributes to mastery at all.
func RelationshipTypeBonus(relType string) (float64, bool) {
	switch exerciserelmodels.ExerciseRelationshipType(relType) {
	case exerciserelmodels.ExerciseRelationshipTypeEquivalent:
		return 0.5, true

	case exerciserelmodels.ExerciseRelationshipTypeAlternative,
		exerciserelmodels.ExerciseRelationshipTypeEasierAlternative,
		exerciserelmodels.ExerciseRelationshipTypeHarderAlternative,
		exerciserelmodels.ExerciseRelationshipTypeEquipmentVariation,
		exerciserelmodels.ExerciseRelationshipTypeVariant,
		exerciserelmodels.ExerciseRelationshipTypeVariation,
		exerciserelmodels.ExerciseRelationshipTypeBilateralUnilateral,
		exerciserelmodels.ExerciseRelationshipTypeProgression,
		exerciserelmodels.ExerciseRelationshipTypeProgressesTo,
		exerciserelmodels.ExerciseRelationshipTypeRegression,
		exerciserelmodels.ExerciseRelationshipTypeRegressesTo:
		return 0.25, true

	default:
		return 0, false
	}
}

// ComputeContributionMultiplier calculates the combined multiplier from relationship strength and type.
func ComputeContributionMultiplier(strength float64, relType string) (float64, bool) {
	typeBonus, contributes := RelationshipTypeBonus(relType)
	if !contributes {
		return 0, false
	}
	return (strength * 0.5) + typeBonus, true
}

// EquipmentRelationshipTypeBonus returns the type bonus for an equipment relationship type.
func EquipmentRelationshipTypeBonus(relType string) (float64, bool) {
	switch equipmentrelmodels.EquipmentRelationshipType(relType) {
	case equipmentrelmodels.EquipmentRelationshipTypeEquivalent:
		return 0.5, true
	default:
		return 0, false
	}
}

// ComputeEquipmentContributionMultiplier calculates the combined multiplier for equipment relationships.
func ComputeEquipmentContributionMultiplier(strength float64, relType string) (float64, bool) {
	typeBonus, contributes := EquipmentRelationshipTypeBonus(relType)
	if !contributes {
		return 0, false
	}
	return (strength * 0.5) + typeBonus, true
}

// FulfillmentContributionMultiplier returns the fixed multiplier for equipment fulfillments.
// Fulfillments represent substitutability — experience with a substitute partially transfers.
func FulfillmentContributionMultiplier() float64 {
	return 0.75 // (1.0 * 0.5) + 0.25
}
