package handlers

import (
	exercisemodels "gesitr/internal/compendium/exercise/models"
	"gesitr/internal/ownershipgroup"
	"gesitr/internal/user/mastery/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RecalculateContributions recomputes mastery_contributions rows for the given user
// and affected exercise IDs. Called after relationship create/delete.
func RecalculateContributions(db *gorm.DB, userID string, exerciseIDs ...uint) error {
	// Fetch all relationships visible to this user that involve any of the affected exercises.
	var relationships []exercisemodels.ExerciseRelationshipEntity
	err := db.Where("ownership_group_id IN (?) AND (from_exercise_id IN ? OR to_exercise_id IN ?)", ownershipgroup.VisibleGroupIDs(db, userID), exerciseIDs, exerciseIDs).
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
					Owner:             userID,
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
		if err := tx.Where("owner = ? AND (exercise_id IN ? OR contributes_from_id IN ?)", userID, exerciseIDs, exerciseIDs).
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
	// Find all distinct users from ownership group memberships.
	var userIDs []string
	if err := db.Table("ownership_group_memberships").
		Where("deleted_at IS NULL").
		Distinct("user_id").Pluck("user_id", &userIDs).Error; err != nil {
		return err
	}

	for _, userID := range userIDs {
		// Get all exercise IDs involved in relationships visible to this user.
		var exerciseIDs []uint
		var fromIDs, toIDs []uint
		visibleGroups := ownershipgroup.VisibleGroupIDs(db, userID)
		db.Model(&exercisemodels.ExerciseRelationshipEntity{}).
			Where("ownership_group_id IN (?)", visibleGroups).Pluck("from_exercise_id", &fromIDs)
		db.Model(&exercisemodels.ExerciseRelationshipEntity{}).
			Where("ownership_group_id IN (?)", visibleGroups).Pluck("to_exercise_id", &toIDs)
		seen := make(map[uint]bool)
		for _, id := range append(fromIDs, toIDs...) {
			if !seen[id] {
				seen[id] = true
				exerciseIDs = append(exerciseIDs, id)
			}
		}

		if err := RecalculateContributions(db, userID, exerciseIDs...); err != nil {
			return err
		}
	}
	return nil
}
