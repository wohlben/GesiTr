package handlers

import (
	exercisemodels "gesitr/internal/compendium/exercise/models"
	"gesitr/internal/user/mastery/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpsertEquipmentExperience records that a user has used equipment via an exercise.
// Looks up the exercise's equipment IDs and increments total_reps for each.
// Safe to call within a transaction.
func UpsertEquipmentExperience(db *gorm.DB, owner string, exerciseID uint, reps *int) error {
	var equipmentIDs []uint
	if err := db.Model(&exercisemodels.ExerciseEquipment{}).
		Where("exercise_id = ?", exerciseID).
		Pluck("equipment_id", &equipmentIDs).Error; err != nil {
		return err
	}
	if len(equipmentIDs) == 0 {
		return nil
	}

	amount := 1
	if reps != nil {
		amount = *reps
	}

	for _, eqID := range equipmentIDs {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "owner"}, {Name: "equipment_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"total_reps": gorm.Expr("total_reps + ?", amount)}),
		}).Create(&models.EquipmentMasteryExperienceEntity{
			Owner:       owner,
			EquipmentID: eqID,
			TotalReps:   amount,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
