package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/workout/models"
	"gesitr/internal/user/workoutgroup"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func preloadWorkout(db *gorm.DB) *gorm.DB {
	return db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	}).Preload("Sections.Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("position")
	})
}

// ListWorkouts returns all workouts owned by the current user, each
// including its sections and section exercises. GET /api/user/workouts
//
// OpenAPI: /api/docs#/operations/ListWorkouts
func ListWorkouts(ctx context.Context, input *ListWorkoutsInput) (*ListWorkoutsOutput, error) {
	userID := humaconfig.GetUserID(ctx)
	db := database.DB.Model(&models.WorkoutEntity{}).Where(`owner = ? OR id IN (
		SELECT wg.workout_id FROM workout_groups wg
		JOIN workout_group_memberships wgm ON wgm.group_id = wg.id
		WHERE wgm.user_id = ? AND wgm.deleted_at IS NULL AND wg.deleted_at IS NULL)`,
		userID, userID)

	var entities []models.WorkoutEntity
	if err := preloadWorkout(db).Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	var nonOwnedIDs []uint
	for _, e := range entities {
		if e.Owner != userID {
			nonOwnedIDs = append(nonOwnedIDs, e.ID)
		}
	}
	groupInfoMap := workoutgroup.GroupInfoForWorkouts(userID, nonOwnedIDs)

	dtos := make([]models.Workout, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
		if info, ok := groupInfoMap[entities[i].ID]; ok {
			dtos[i].WorkoutGroup = &models.WorkoutGroupInfo{
				GroupName:  info.GroupName,
				Membership: info.MembershipRole,
			}
		}
	}
	return &ListWorkoutsOutput{Body: dtos}, nil
}

// CreateWorkout creates an empty workout. Add sections via
// [CreateWorkoutSection] and items via [CreateWorkoutSectionItem].
// POST /api/user/workouts
//
// OpenAPI: /api/docs#/operations/CreateWorkout
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
// OpenAPI: /api/docs#/operations/GetWorkout
func GetWorkout(ctx context.Context, input *GetWorkoutInput) (*GetWorkoutOutput, error) {
	var entity models.WorkoutEntity
	if err := preloadWorkout(database.DB).First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	access := workoutgroup.CheckWorkoutAccess(humaconfig.GetUserID(ctx), entity.Owner, entity.ID)
	if !access.CanRead() {
		return nil, huma.Error403Forbidden("access denied")
	}
	dto := entity.ToDTO()
	if !access.IsOwner && access.GroupName != "" {
		dto.WorkoutGroup = &models.WorkoutGroupInfo{
			GroupName:  access.GroupName,
			Membership: access.MembershipRole,
		}
	}
	return &GetWorkoutOutput{Body: dto}, nil
}

// UpdateWorkout updates workout metadata (name, notes). Sections and exercises
// are managed via their own endpoints. PUT /api/user/workouts/{id}
//
// OpenAPI: /api/docs#/operations/UpdateWorkout
func UpdateWorkout(ctx context.Context, input *UpdateWorkoutInput) (*UpdateWorkoutOutput, error) {
	var existing models.WorkoutEntity
	if err := database.DB.First(&existing, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	access := workoutgroup.CheckWorkoutAccess(humaconfig.GetUserID(ctx), existing.Owner, existing.ID)
	if !access.CanModify() {
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
	dto := entity.ToDTO()
	if !access.IsOwner && access.GroupName != "" {
		dto.WorkoutGroup = &models.WorkoutGroupInfo{
			GroupName:  access.GroupName,
			Membership: access.MembershipRole,
		}
	}
	return &UpdateWorkoutOutput{Body: dto}, nil
}

// DeleteWorkout deletes a workout. DELETE /api/user/workouts/{id}
//
// OpenAPI: /api/docs#/operations/DeleteWorkout
func DeleteWorkout(ctx context.Context, input *DeleteWorkoutInput) (*DeleteWorkoutOutput, error) {
	var entity models.WorkoutEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	access := workoutgroup.CheckWorkoutAccess(humaconfig.GetUserID(ctx), entity.Owner, entity.ID)
	if !access.CanDelete() {
		return nil, huma.Error403Forbidden("access denied")
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}
