package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/namepreference/models"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm/clause"
)

// ListExerciseNamePreferences returns all name preferences for the current user.
// GET /api/user/exercise-name-preferences
func ListExerciseNamePreferences(ctx context.Context, input *struct{}) (*ListExerciseNamePreferencesOutput, error) {
	var entities []models.ExerciseNamePreference
	if err := database.DB.Where("owner = ?", humaconfig.GetUserID(ctx)).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	dtos := make([]models.ExerciseNamePreferenceDTO, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListExerciseNamePreferencesOutput{Body: dtos}, nil
}

// SetExerciseNamePreference upserts the user's preferred name for an exercise.
// PUT /api/user/exercise-name-preferences/:exerciseId
func SetExerciseNamePreference(ctx context.Context, input *SetExerciseNamePreferenceInput) (*SetExerciseNamePreferenceOutput, error) {
	entity := models.ExerciseNamePreference{
		Owner:          humaconfig.GetUserID(ctx),
		ExerciseID:     input.ExerciseID,
		ExerciseNameID: input.Body.ExerciseNameID,
	}
	if err := database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "owner"}, {Name: "exercise_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"exercise_name_id"}),
	}).Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &SetExerciseNamePreferenceOutput{Body: entity.ToDTO()}, nil
}

// DeleteExerciseNamePreference removes the user's preferred name for an exercise.
// DELETE /api/user/exercise-name-preferences/:exerciseId
func DeleteExerciseNamePreference(ctx context.Context, input *DeleteExerciseNamePreferenceInput) (*DeleteExerciseNamePreferenceOutput, error) {
	result := database.DB.Where("owner = ? AND exercise_id = ?", humaconfig.GetUserID(ctx), input.ExerciseID).
		Delete(&models.ExerciseNamePreference{})
	if result.Error != nil {
		return nil, huma.Error500InternalServerError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, huma.Error404NotFound("preference not found")
	}
	return nil, nil
}
