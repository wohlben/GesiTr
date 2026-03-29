package handlers

import (
	"context"
	"time"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/mastery/models"

	"github.com/danielgtaylor/huma/v2"
)

type xpRow struct {
	ExerciseID uint
	TotalXP    float64
}

type daysRow struct {
	ExerciseID   uint
	DistinctDays int
}

// ListExerciseMastery returns mastery for all exercises the user has logged.
// GET /api/user/mastery
func ListExerciseMastery(ctx context.Context, input *ListMasteryInput) (*ListMasteryOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	recencyStart := time.Now().AddDate(0, -6, 0)

	// 1. All contribution mappings for this user.
	var contributions []models.MasteryContributionEntity
	if err := database.DB.Where("owner = ?", userID).Find(&contributions).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// 2. Total XP (reps or 1 per log) per exercise.
	var xpRows []xpRow
	if err := database.DB.Table("exercise_logs").
		Select("exercise_id, SUM(COALESCE(reps, 1)) as total_xp").
		Where("owner = ? AND deleted_at IS NULL", userID).
		Group("exercise_id").
		Scan(&xpRows).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// 3. Distinct days per exercise in recency window.
	var daysRows []daysRow
	if err := database.DB.Table("exercise_logs").
		Select("exercise_id, COUNT(DISTINCT DATE(performed_at)) as distinct_days").
		Where("owner = ? AND performed_at >= ? AND deleted_at IS NULL", userID, recencyStart).
		Group("exercise_id").
		Scan(&daysRows).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &ListMasteryOutput{Body: computeMasteryList(xpRows, daysRows, contributions)}, nil
}

// GetExerciseMastery returns mastery for a specific exercise.
// GET /api/user/mastery/:exerciseId
func GetExerciseMastery(ctx context.Context, input *GetMasteryInput) (*GetMasteryOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	recencyStart := time.Now().AddDate(0, -6, 0)
	exerciseID := input.ExerciseID

	// 1. Contribution mappings for this exercise.
	var contributions []models.MasteryContributionEntity
	if err := database.DB.Where("owner = ? AND exercise_id = ?", userID, exerciseID).
		Find(&contributions).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// Collect all exercise IDs to query (self + contributors).
	queryIDs := []uint{exerciseID}
	for _, c := range contributions {
		queryIDs = append(queryIDs, c.ContributesFromID)
	}

	// 2. Total XP per exercise (filtered to relevant IDs).
	var xpRows []xpRow
	if err := database.DB.Table("exercise_logs").
		Select("exercise_id, SUM(COALESCE(reps, 1)) as total_xp").
		Where("owner = ? AND exercise_id IN ? AND deleted_at IS NULL", userID, queryIDs).
		Group("exercise_id").
		Scan(&xpRows).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	// 3. Distinct days in recency window (across all relevant exercises).
	var daysRows []daysRow
	if err := database.DB.Table("exercise_logs").
		Select("exercise_id, COUNT(DISTINCT DATE(performed_at)) as distinct_days").
		Where("owner = ? AND exercise_id IN ? AND performed_at >= ? AND deleted_at IS NULL", userID, queryIDs, recencyStart).
		Group("exercise_id").
		Scan(&daysRows).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	mastery := computeSingleMastery(exerciseID, xpRows, daysRows, contributions)
	return &GetMasteryOutput{Body: mastery}, nil
}

// computeMasteryList computes mastery for all exercises that have any logs.
func computeMasteryList(xpRows []xpRow, daysRows []daysRow, contributions []models.MasteryContributionEntity) []models.ExerciseMastery {
	xpMap := make(map[uint]float64)
	for _, r := range xpRows {
		xpMap[r.ExerciseID] = r.TotalXP
	}

	daysMap := make(map[uint]int)
	for _, r := range daysRows {
		daysMap[r.ExerciseID] = r.DistinctDays
	}

	// Build contribution index: exerciseID → []contributions
	contribMap := make(map[uint][]models.MasteryContributionEntity)
	for _, c := range contributions {
		contribMap[c.ExerciseID] = append(contribMap[c.ExerciseID], c)
	}

	// Collect all exercise IDs that have logs or are targets of contributions.
	exerciseIDs := make(map[uint]bool)
	for id := range xpMap {
		exerciseIDs[id] = true
	}
	for _, c := range contributions {
		if xpMap[c.ExerciseID] > 0 || xpMap[c.ContributesFromID] > 0 {
			exerciseIDs[c.ExerciseID] = true
		}
	}

	result := make([]models.ExerciseMastery, 0, len(exerciseIDs))
	for exerciseID := range exerciseIDs {
		contribs := contribMap[exerciseID]
		mastery := computeSingleMastery(exerciseID, xpRows, daysRows, contribs)
		if mastery.TotalXP > 0 {
			result = append(result, mastery)
		}
	}
	return result
}

// computeSingleMastery computes mastery for one exercise given pre-fetched data.
func computeSingleMastery(exerciseID uint, xpRows []xpRow, daysRows []daysRow, contributions []models.MasteryContributionEntity) models.ExerciseMastery {
	xpMap := make(map[uint]float64)
	for _, r := range xpRows {
		xpMap[r.ExerciseID] = r.TotalXP
	}

	daysMap := make(map[uint]int)
	for _, r := range daysRows {
		daysMap[r.ExerciseID] = r.DistinctDays
	}

	// Own XP at 1.0 multiplier.
	totalXP := xpMap[exerciseID]

	// Contributor XP scaled by multiplier.
	for _, c := range contributions {
		totalXP += xpMap[c.ContributesFromID] * c.Multiplier
	}

	// Distinct days across self + contributors.
	totalDays := daysMap[exerciseID]
	for _, c := range contributions {
		totalDays += daysMap[c.ContributesFromID]
	}

	baseLevel := models.ComputeLevel(totalXP)
	multiplier := models.ComputeRecencyMultiplier(totalDays, baseLevel)
	effectiveXP := totalXP * multiplier
	level := models.ComputeLevel(effectiveXP)
	tier := models.ComputeTier(level)
	progress := models.ComputeProgress(effectiveXP, level)

	return models.ExerciseMastery{
		ExerciseID:   exerciseID,
		TotalXP:      totalXP,
		EffectiveXP:  effectiveXP,
		Level:        level,
		Tier:         string(tier),
		Progress:     progress,
		DistinctDays: totalDays,
		Multiplier:   multiplier,
	}
}
