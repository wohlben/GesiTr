package handlers

import (
	"context"
	"time"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	workoutmodels "gesitr/internal/user/workout/models"
	workoutlogmodels "gesitr/internal/user/workoutlog/models"
	"gesitr/internal/user/workoutschedule/models"

	"github.com/danielgtaylor/huma/v2"
)

func requireOwner(ctx context.Context, owner string) error {
	if humaconfig.GetUserID(ctx) != owner {
		return huma.Error403Forbidden("not the owner of this resource")
	}
	return nil
}

// ListWorkoutSchedules returns schedules owned by the current user.
// GET /api/user/workout-schedules
func ListWorkoutSchedules(ctx context.Context, input *ListWorkoutSchedulesInput) (*ListWorkoutSchedulesOutput, error) {
	db := database.DB.Model(&models.WorkoutScheduleEntity{}).
		Where("owner = ?", humaconfig.GetUserID(ctx))

	if input.WorkoutID != "" {
		db = db.Where("workout_id = ?", input.WorkoutID)
	}

	var entities []models.WorkoutScheduleEntity
	if err := db.Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	now := time.Now()
	dtos := make([]models.WorkoutSchedule, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO(now)
	}
	return &ListWorkoutSchedulesOutput{Body: dtos}, nil
}

// CreateWorkoutSchedule creates a new workout schedule.
// POST /api/user/workout-schedules
func CreateWorkoutSchedule(ctx context.Context, input *CreateWorkoutScheduleInput) (*CreateWorkoutScheduleOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	// Verify workout exists and is owned by user
	var workout workoutmodels.WorkoutEntity
	if err := database.DB.First(&workout, input.Body.WorkoutID).Error; err != nil {
		return nil, huma.Error404NotFound("Workout not found")
	}
	if workout.Owner != userID {
		return nil, huma.Error403Forbidden("not the owner of this workout")
	}

	initialStatus := "committed"
	if input.Body.InitialStatus != nil {
		switch *input.Body.InitialStatus {
		case "committed", "proposed":
			initialStatus = *input.Body.InitialStatus
		default:
			return nil, huma.Error400BadRequest("initialStatus must be 'committed' or 'proposed'")
		}
	}

	// Default start date: tomorrow
	startDate := startOfDay(time.Now().AddDate(0, 0, 1))
	if input.Body.StartDate != nil {
		startDate = *input.Body.StartDate
	}

	entity := models.WorkoutScheduleEntity{
		Owner:         userID,
		WorkoutID:     input.Body.WorkoutID,
		StartDate:     startDate,
		EndDate:       input.Body.EndDate,
		InitialStatus: initialStatus,
	}

	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dto := entity.ToDTO(time.Now())
	return &CreateWorkoutScheduleOutput{Body: dto}, nil
}

// GetWorkoutSchedule returns a single schedule.
// GET /api/user/workout-schedules/{id}
func GetWorkoutSchedule(ctx context.Context, input *GetWorkoutScheduleInput) (*GetWorkoutScheduleOutput, error) {
	var entity models.WorkoutScheduleEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Schedule not found")
	}
	if err := requireOwner(ctx, entity.Owner); err != nil {
		return nil, err
	}
	dto := entity.ToDTO(time.Now())
	return &GetWorkoutScheduleOutput{Body: dto}, nil
}

// UpdateWorkoutSchedule partially updates a schedule.
// PATCH /api/user/workout-schedules/{id}
func UpdateWorkoutSchedule(ctx context.Context, input *UpdateWorkoutScheduleInput) (*UpdateWorkoutScheduleOutput, error) {
	var entity models.WorkoutScheduleEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Schedule not found")
	}
	if err := requireOwner(ctx, entity.Owner); err != nil {
		return nil, err
	}

	if input.Body.EndDate != nil {
		entity.EndDate = input.Body.EndDate
	}

	if err := database.DB.Save(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dto := entity.ToDTO(time.Now())
	return &UpdateWorkoutScheduleOutput{Body: dto}, nil
}

// DeleteWorkoutSchedule deletes a schedule and orphans its workout logs.
// DELETE /api/user/workout-schedules/{id}
func DeleteWorkoutSchedule(ctx context.Context, input *DeleteWorkoutScheduleInput) (*DeleteWorkoutScheduleOutput, error) {
	var entity models.WorkoutScheduleEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Schedule not found")
	}
	if err := requireOwner(ctx, entity.Owner); err != nil {
		return nil, err
	}

	// Orphan existing workout logs
	database.DB.Model(&workoutlogmodels.WorkoutLogEntity{}).
		Where("schedule_id = ?", entity.ID).
		Updates(map[string]any{"schedule_id": nil, "period_id": nil})

	// Delete commitments, periods, then schedule
	var periodIDs []uint
	database.DB.Model(&models.SchedulePeriodEntity{}).
		Where("schedule_id = ?", entity.ID).Pluck("id", &periodIDs)
	if len(periodIDs) > 0 {
		database.DB.Where("period_id IN ?", periodIDs).Delete(&models.ScheduleCommitmentEntity{})
	}
	database.DB.Where("schedule_id = ?", entity.ID).Delete(&models.SchedulePeriodEntity{})
	database.DB.Delete(&entity)

	return nil, nil
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
