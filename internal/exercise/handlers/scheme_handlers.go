package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/exercise/models"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"

	"github.com/danielgtaylor/huma/v2"
)

func exerciseSchemeDTOFromBody(b ExerciseSchemeBody) models.ExerciseScheme {
	return models.ExerciseScheme{
		ExerciseID:      b.ExerciseID,
		MeasurementType: b.MeasurementType,
		Sets:            b.Sets,
		Reps:            b.Reps,
		Weight:          b.Weight,
		RestBetweenSets: b.RestBetweenSets,
		TimePerRep:      b.TimePerRep,
		Duration:        b.Duration,
		Distance:        b.Distance,
		TargetTime:      b.TargetTime,
	}
}

// ListExerciseSchemes returns schemes the current user has access to: their
// own schemes plus schemes linked to public exercises. Filter by exerciseId
// or measurementType query params. GET /api/exercise-schemes
//
// OpenAPI: /api/docs#/operations/ListExerciseSchemes
func ListExerciseSchemes(ctx context.Context, input *ListExerciseSchemesInput) (*ListExerciseSchemesOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	// FIXME: subquery doesn't scale — replace with a join or denormalize visibility
	db := database.DB.Model(&models.ExerciseSchemeEntity{}).
		Where("owner = ? OR exercise_id IN (SELECT id FROM exercises WHERE public = ? AND deleted_at IS NULL)", userID, true)

	if input.ExerciseID != "" {
		db = db.Where("exercise_id = ?", input.ExerciseID)
	}
	if input.MeasurementType != "" {
		db = db.Where("measurement_type = ?", input.MeasurementType)
	}

	var entities []models.ExerciseSchemeEntity
	if err := db.Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.ExerciseScheme, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListExerciseSchemesOutput{Body: dtos}, nil
}

// CreateExerciseScheme creates an exercise scheme — a user-specific
// configuration of an exercise (sets, reps, measurement type). Requires
// an exerciseId referencing an existing exercise (see [CreateExercise]).
// Schemes are referenced when adding exercises to workouts via
// [gesitr/internal/user/workout/handlers.CreateWorkoutSectionItem].
// POST /api/exercise-schemes
//
// OpenAPI: /api/docs#/operations/CreateExerciseScheme
func CreateExerciseScheme(ctx context.Context, input *CreateExerciseSchemeInput) (*CreateExerciseSchemeOutput, error) {
	dto := exerciseSchemeDTOFromBody(input.Body)

	var exercise models.ExerciseEntity
	if err := database.DB.First(&exercise, dto.ExerciseID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise not found")
	}

	entity := models.ExerciseSchemeFromDTO(dto)
	entity.Owner = humaconfig.GetUserID(ctx)
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateExerciseSchemeOutput{Body: entity.ToDTO()}, nil
}

// GetExerciseScheme returns a single exercise scheme. Access is determined by
// the linked exercise's visibility — if the user can see the exercise, they
// can see its schemes. GET /api/exercise-schemes/:id
//
// OpenAPI: /api/docs#/operations/GetExerciseScheme
func GetExerciseScheme(ctx context.Context, input *GetExerciseSchemeInput) (*GetExerciseSchemeOutput, error) {
	var entity models.ExerciseSchemeEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise scheme not found")
	}
	var exercise models.ExerciseEntity
	if err := database.DB.First(&exercise, entity.ExerciseID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise not found")
	}
	userID := humaconfig.GetUserID(ctx)
	perms, _ := shared.ResolvePermissions(userID, exercise.Owner, exercise.Public)
	if len(perms) == 0 {
		return nil, huma.Error403Forbidden("access denied")
	}
	return &GetExerciseSchemeOutput{Body: entity.ToDTO()}, nil
}

// UpdateExerciseScheme updates a scheme's configuration. The exerciseId
// cannot be changed. PUT /api/exercise-schemes/:id
//
// OpenAPI: /api/docs#/operations/UpdateExerciseScheme
func UpdateExerciseScheme(ctx context.Context, input *UpdateExerciseSchemeInput) (*UpdateExerciseSchemeOutput, error) {
	var existing models.ExerciseSchemeEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise scheme not found")
	}
	if existing.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}

	dto := exerciseSchemeDTOFromBody(input.Body)
	entity := models.ExerciseSchemeFromDTO(dto)
	entity.ID = existing.ID
	entity.Owner = existing.Owner
	entity.ExerciseID = existing.ExerciseID

	if err := database.DB.Save(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateExerciseSchemeOutput{Body: entity.ToDTO()}, nil
}

// DeleteExerciseScheme deletes an exercise scheme. Owner only.
// DELETE /api/exercise-schemes/:id
//
// OpenAPI: /api/docs#/operations/DeleteExerciseScheme
func DeleteExerciseScheme(ctx context.Context, input *DeleteExerciseSchemeInput) (*DeleteExerciseSchemeOutput, error) {
	var entity models.ExerciseSchemeEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Exercise scheme not found")
	}
	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}
	if err := database.DB.Unscoped().Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
