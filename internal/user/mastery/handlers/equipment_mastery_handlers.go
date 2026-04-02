package handlers

import (
	"context"
	"time"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/mastery/models"

	"github.com/danielgtaylor/huma/v2"
)

type equipmentXPRow struct {
	EquipmentID uint
	TotalXP     float64
}

type equipmentDateRow struct {
	EquipmentID uint
	LogDate     string
}

// ListEquipmentMastery returns mastery for all equipment the user has used.
// GET /api/user/equipment-mastery
func ListEquipmentMastery(ctx context.Context, input *ListEquipmentMasteryInput) (*ListEquipmentMasteryOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	recencyStart := time.Now().AddDate(0, -6, 0)

	// 1. All equipment contribution mappings for this user.
	var contributions []models.EquipmentMasteryContributionEntity
	if err := database.DB.Where("owner = ?", userID).Find(&contributions).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// 2. Total XP per equipment (from exercise logs via exercise_equipments).
	var xpRows []equipmentXPRow
	if err := database.DB.Table("exercise_logs").
		Select("ee.equipment_id, SUM(COALESCE(exercise_logs.reps, 1)) as total_xp").
		Joins("JOIN exercise_equipments ee ON ee.exercise_id = exercise_logs.exercise_id").
		Where("exercise_logs.owner = ? AND exercise_logs.deleted_at IS NULL", userID).
		Group("ee.equipment_id").
		Scan(&xpRows).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// 3. Distinct (equipment_id, date) pairs in recency window.
	var dateRows []equipmentDateRow
	if err := database.DB.Table("exercise_logs").
		Select("ee.equipment_id, DATE(exercise_logs.performed_at) as log_date").
		Joins("JOIN exercise_equipments ee ON ee.exercise_id = exercise_logs.exercise_id").
		Where("exercise_logs.owner = ? AND exercise_logs.performed_at >= ? AND exercise_logs.deleted_at IS NULL", userID, recencyStart).
		Group("ee.equipment_id, DATE(exercise_logs.performed_at)").
		Scan(&dateRows).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &ListEquipmentMasteryOutput{Body: computeEquipmentMasteryList(xpRows, dateRows, contributions)}, nil
}

// GetEquipmentMastery returns mastery for a specific equipment item.
// GET /api/user/equipment-mastery/:equipmentId
func GetEquipmentMastery(ctx context.Context, input *GetEquipmentMasteryInput) (*GetEquipmentMasteryOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	recencyStart := time.Now().AddDate(0, -6, 0)
	equipmentID := input.EquipmentID

	// 1. Contribution mappings for this equipment.
	var contributions []models.EquipmentMasteryContributionEntity
	if err := database.DB.Where("owner = ? AND equipment_id = ?", userID, equipmentID).
		Find(&contributions).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Collect all equipment IDs to query (self + contributors).
	queryIDs := []uint{equipmentID}
	for _, c := range contributions {
		queryIDs = append(queryIDs, c.ContributesFromID)
	}

	// 2. Total XP per equipment (filtered to relevant IDs).
	var xpRows []equipmentXPRow
	if err := database.DB.Table("exercise_logs").
		Select("ee.equipment_id, SUM(COALESCE(exercise_logs.reps, 1)) as total_xp").
		Joins("JOIN exercise_equipments ee ON ee.exercise_id = exercise_logs.exercise_id").
		Where("exercise_logs.owner = ? AND ee.equipment_id IN ? AND exercise_logs.deleted_at IS NULL", userID, queryIDs).
		Group("ee.equipment_id").
		Scan(&xpRows).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// 3. Distinct days across all relevant equipment (union, not per-equipment sum).
	var unionDays int64
	if err := database.DB.Table("exercise_logs").
		Select("COUNT(DISTINCT DATE(exercise_logs.performed_at))").
		Joins("JOIN exercise_equipments ee ON ee.exercise_id = exercise_logs.exercise_id").
		Where("exercise_logs.owner = ? AND ee.equipment_id IN ? AND exercise_logs.performed_at >= ? AND exercise_logs.deleted_at IS NULL", userID, queryIDs, recencyStart).
		Scan(&unionDays).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	mastery := computeSingleEquipmentMastery(equipmentID, xpRows, int(unionDays), contributions)
	return &GetEquipmentMasteryOutput{Body: mastery}, nil
}

func computeEquipmentMasteryList(xpRows []equipmentXPRow, dateRows []equipmentDateRow, contributions []models.EquipmentMasteryContributionEntity) []models.EquipmentMastery {
	xpMap := make(map[uint]float64)
	for _, r := range xpRows {
		xpMap[r.EquipmentID] = r.TotalXP
	}

	dateSets := make(map[uint]map[string]bool)
	for _, r := range dateRows {
		if dateSets[r.EquipmentID] == nil {
			dateSets[r.EquipmentID] = make(map[string]bool)
		}
		dateSets[r.EquipmentID][r.LogDate] = true
	}

	contribMap := make(map[uint][]models.EquipmentMasteryContributionEntity)
	for _, c := range contributions {
		contribMap[c.EquipmentID] = append(contribMap[c.EquipmentID], c)
	}

	equipmentIDs := make(map[uint]bool)
	for id := range xpMap {
		equipmentIDs[id] = true
	}
	for _, c := range contributions {
		if xpMap[c.EquipmentID] > 0 || xpMap[c.ContributesFromID] > 0 {
			equipmentIDs[c.EquipmentID] = true
		}
	}

	result := make([]models.EquipmentMastery, 0, len(equipmentIDs))
	for eqID := range equipmentIDs {
		contribs := contribMap[eqID]

		unionDates := make(map[string]bool)
		for d := range dateSets[eqID] {
			unionDates[d] = true
		}
		for _, c := range contribs {
			for d := range dateSets[c.ContributesFromID] {
				unionDates[d] = true
			}
		}

		mastery := computeSingleEquipmentMastery(eqID, xpRows, len(unionDates), contribs)
		if mastery.TotalXP > 0 {
			result = append(result, mastery)
		}
	}
	return result
}

func computeSingleEquipmentMastery(equipmentID uint, xpRows []equipmentXPRow, unionDays int, contributions []models.EquipmentMasteryContributionEntity) models.EquipmentMastery {
	xpMap := make(map[uint]float64)
	for _, r := range xpRows {
		xpMap[r.EquipmentID] = r.TotalXP
	}

	totalXP := xpMap[equipmentID]
	for _, c := range contributions {
		totalXP += xpMap[c.ContributesFromID] * c.Multiplier
	}

	baseLevel := models.ComputeLevel(totalXP)
	multiplier := models.ComputeRecencyMultiplier(unionDays, baseLevel)
	effectiveXP := totalXP * multiplier
	level := models.ComputeLevel(effectiveXP)
	tier := models.ComputeTier(level)
	progress := models.ComputeProgress(effectiveXP, level)

	return models.EquipmentMastery{
		EquipmentID:  equipmentID,
		TotalXP:      totalXP,
		EffectiveXP:  effectiveXP,
		Level:        level,
		Tier:         string(tier),
		Progress:     progress,
		DistinctDays: unionDays,
		Multiplier:   multiplier,
	}
}
