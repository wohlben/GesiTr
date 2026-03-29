package handlers

import (
	exerciserelmodels "gesitr/internal/exerciserelationship/models"
	"gesitr/internal/user/mastery/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RecalculateContributions recomputes mastery_contributions rows for the given owner
// and affected exercise IDs. Called after relationship create/delete.
func RecalculateContributions(db *gorm.DB, owner string, exerciseIDs ...uint) error {
	// Fetch all relationships for this owner that involve any of the affected exercises.
	var relationships []exerciserelmodels.ExerciseRelationshipEntity
	err := db.Where("owner = ? AND (from_exercise_id IN ? OR to_exercise_id IN ?)", owner, exerciseIDs, exerciseIDs).
		Find(&relationships).Error
	if err != nil {
		return err
	}

	// Build new contribution rows (bidirectional).
	// Key: (exerciseID, contributesFromID) → best row
	type contribKey struct {
		exerciseID      uint
		contributesFrom uint
	}
	best := make(map[contribKey]models.MasteryContributionEntity)

	for _, rel := range relationships {
		mult, ok := models.ComputeContributionMultiplier(rel.Strength, string(rel.RelationshipType))
		if !ok {
			continue
		}

		// Forward: exercise gains mastery from contributor
		pairs := [][2]uint{
			{rel.FromExerciseID, rel.ToExerciseID},
			{rel.ToExerciseID, rel.FromExerciseID},
		}
		for _, pair := range pairs {
			key := contribKey{pair[0], pair[1]}
			if existing, exists := best[key]; !exists || mult > existing.Multiplier {
				best[key] = models.MasteryContributionEntity{
					Owner:             owner,
					ExerciseID:        pair[0],
					ContributesFromID: pair[1],
					Multiplier:        mult,
					RelationshipType:  string(rel.RelationshipType),
				}
			}
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Delete existing contributions for affected exercises.
		if err := tx.Where("owner = ? AND (exercise_id IN ? OR contributes_from_id IN ?)", owner, exerciseIDs, exerciseIDs).
			Delete(&models.MasteryContributionEntity{}).Error; err != nil {
			return err
		}

		// Insert new rows.
		if len(best) == 0 {
			return nil
		}
		rows := make([]models.MasteryContributionEntity, 0, len(best))
		for _, row := range best {
			rows = append(rows, row)
		}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&rows).Error
	})
}

// BackfillContributions computes mastery_contributions for all existing relationships.
// Intended to be called once during migration.
func BackfillContributions(db *gorm.DB) error {
	// Find all distinct owners with relationships.
	var owners []string
	if err := db.Model(&exerciserelmodels.ExerciseRelationshipEntity{}).
		Distinct("owner").Pluck("owner", &owners).Error; err != nil {
		return err
	}

	for _, owner := range owners {
		// Get all exercise IDs involved in relationships for this owner.
		var exerciseIDs []uint
		var fromIDs, toIDs []uint
		db.Model(&exerciserelmodels.ExerciseRelationshipEntity{}).
			Where("owner = ?", owner).Pluck("from_exercise_id", &fromIDs)
		db.Model(&exerciserelmodels.ExerciseRelationshipEntity{}).
			Where("owner = ?", owner).Pluck("to_exercise_id", &toIDs)
		seen := make(map[uint]bool)
		for _, id := range append(fromIDs, toIDs...) {
			if !seen[id] {
				seen[id] = true
				exerciseIDs = append(exerciseIDs, id)
			}
		}

		if err := RecalculateContributions(db, owner, exerciseIDs...); err != nil {
			return err
		}
	}
	return nil
}
