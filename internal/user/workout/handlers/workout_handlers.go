package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/workout/models"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func preloadWorkout(db *gorm.DB) *gorm.DB {
	return db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Sections.Exercises", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	})
}

// ListWorkouts returns all workouts owned by the current user, each
// including its sections and section exercises. GET /api/user/workouts
//
// OpenAPI: /api/docs#/operations/list-workouts
func ListWorkouts(ctx context.Context, input *ListWorkoutsInput) (*ListWorkoutsOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	db := database.DB.Model(&models.WorkoutEntity{}).Where("owner = ?", userID)

	var entities []models.WorkoutEntity
	if err := preloadWorkout(db).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.Workout, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListWorkoutsOutput{Body: dtos}, nil
}

// CreateWorkout creates an empty workout. Add sections via
// [CreateWorkoutSection] and exercises via [CreateWorkoutSectionExercise].
// POST /api/user/workouts
//
// OpenAPI: /api/docs#/operations/create-workout
func CreateWorkout(ctx context.Context, input *CreateWorkoutInput) (*CreateWorkoutOutput, error) {
	entity := models.WorkoutEntity{
		Name:  input.Body.Name,
		Notes: input.Body.Notes,
	}
	entity.Owner = humaconfig.GetUserID(ctx)
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateWorkoutOutput{Body: entity.ToDTO()}, nil
}

// GetWorkout returns a workout with its full section and exercise tree.
// Returns 403 if the caller is not the owner. GET /api/user/workouts/{id}
//
// OpenAPI: /api/docs#/operations/get-workout
func GetWorkout(ctx context.Context, input *GetWorkoutInput) (*GetWorkoutOutput, error) {
	var entity models.WorkoutEntity
	if err := preloadWorkout(database.DB).First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}
	return &GetWorkoutOutput{Body: entity.ToDTO()}, nil
}

// UpdateWorkout updates workout metadata (name, notes). Sections and exercises
// are managed via their own endpoints. PUT /api/user/workouts/{id}
//
// OpenAPI: /api/docs#/operations/update-workout
func UpdateWorkout(ctx context.Context, input *UpdateWorkoutInput) (*UpdateWorkoutOutput, error) {
	var existing models.WorkoutEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	if existing.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}

	entity := models.WorkoutEntity{
		Name:  input.Body.Name,
		Notes: input.Body.Notes,
	}
	entity.ID = existing.ID
	entity.Owner = existing.Owner

	if err := database.DB.Save(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	if err := preloadWorkout(database.DB).First(&entity, entity.ID).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &UpdateWorkoutOutput{Body: entity.ToDTO()}, nil
}

// DeleteWorkout deletes a workout. DELETE /api/user/workouts/{id}
//
// OpenAPI: /api/docs#/operations/delete-workout
func DeleteWorkout(ctx context.Context, input *DeleteWorkoutInput) (*DeleteWorkoutOutput, error) {
	var entity models.WorkoutEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	if entity.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("access denied")
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
