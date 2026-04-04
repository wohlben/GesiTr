package handlers

import (
	"context"
	"strings"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/exercisescheme/models"

	"github.com/danielgtaylor/huma/v2"
)

// ListExerciseSchemeSectionItems returns the current user's scheme assignments
// for the given workout section item IDs.
// GET /api/user/exercise-scheme-section-items
func ListExerciseSchemeSectionItems(ctx context.Context, input *ListExerciseSchemeSectionItemsInput) (*ListExerciseSchemeSectionItemsOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	db := database.DB.Model(&models.ExerciseSchemeSectionItemEntity{}).
		Where("owner = ?", userID)

	if input.WorkoutSectionItemIDs != "" {
		ids := strings.Split(input.WorkoutSectionItemIDs, ",")
		db = db.Where("workout_section_item_id IN ?", ids)
	}

	var entities []models.ExerciseSchemeSectionItemEntity
	if err := db.Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.ExerciseSchemeSectionItem, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListExerciseSchemeSectionItemsOutput{Body: dtos}, nil
}

// UpsertExerciseSchemeSectionItem creates or replaces the current user's scheme
// assignment for a workout section item slot.
// PUT /api/user/exercise-scheme-section-items
func UpsertExerciseSchemeSectionItem(ctx context.Context, input *UpsertExerciseSchemeSectionItemInput) (*UpsertExerciseSchemeSectionItemOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	// Validate scheme exists and belongs to the user
	var scheme models.ExerciseSchemeEntity
	if err := database.DB.First(&scheme, input.Body.ExerciseSchemeID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise scheme not found")
	}
	if scheme.Owner != userID {
		return nil, huma.Error403Forbidden("access denied")
	}

	// Upsert: find existing or create
	var existing models.ExerciseSchemeSectionItemEntity
	result := database.DB.Where("workout_section_item_id = ? AND owner = ?",
		input.Body.WorkoutSectionItemID, userID).First(&existing)

	if result.Error == nil {
		// Update existing
		existing.ExerciseSchemeID = input.Body.ExerciseSchemeID
		if err := database.DB.Save(&existing).Error; err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		return &UpsertExerciseSchemeSectionItemOutput{Body: existing.ToDTO()}, nil
	}

	// Create new
	entity := models.ExerciseSchemeSectionItemEntity{
		ExerciseSchemeID:     input.Body.ExerciseSchemeID,
		WorkoutSectionItemID: input.Body.WorkoutSectionItemID,
		Owner:                userID,
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpsertExerciseSchemeSectionItemOutput{Body: entity.ToDTO()}, nil
}

// DeleteExerciseSchemeSectionItem removes a scheme assignment. Owner only.
// DELETE /api/user/exercise-scheme-section-items/{id}
func DeleteExerciseSchemeSectionItem(ctx context.Context, input *DeleteExerciseSchemeSectionItemInput) (*DeleteExerciseSchemeSectionItemOutput, error) {
	var entity models.ExerciseSchemeSectionItemEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise scheme section item not found")
	}
	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}
	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
