package handlers

import (
	"gesitr/internal/user/mastery/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpsertExperience records that a user has performed an exercise.
// Increments total_reps by the given reps count (or 1 if reps is nil, for non-rep exercises).
// Safe to call within a transaction.
func UpsertExperience(db *gorm.DB, owner string, exerciseID uint, reps *int) error {
	amount := 1
	if reps != nil {
		amount = *reps
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "owner"}, {Name: "exercise_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"total_reps": gorm.Expr("total_reps + ?", amount)}),
	}).Create(&models.MasteryExperienceEntity{
		Owner:      owner,
		ExerciseID: exerciseID,
		TotalReps:  amount,
	}).Error
}
