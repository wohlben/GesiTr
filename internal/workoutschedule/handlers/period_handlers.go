package handlers

import (
	"context"
	"time"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/shared"
	"gesitr/internal/workoutschedule/models"

	"github.com/danielgtaylor/huma/v2"
)

// ListSchedulePeriods returns periods for a given schedule, or all periods
// for the user when scheduleId is omitted.
// GET /api/user/schedule-periods
func ListSchedulePeriods(ctx context.Context, input *ListSchedulePeriodsInput) (*ListSchedulePeriodsOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	var entities []models.SchedulePeriodEntity
	if input.ScheduleID != "" {
		// Filter by specific schedule
		var schedule models.WorkoutScheduleEntity
		if err := database.DB.First(&schedule, input.ScheduleID).Error; err != nil {
			return nil, huma.Error404NotFound("Schedule not found")
		}
		if schedule.Owner != userID {
			return nil, huma.Error403Forbidden("not the owner of this schedule")
		}
		if err := database.DB.Where("schedule_id = ?", input.ScheduleID).
			Order("period_start").Find(&entities).Error; err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
	} else {
		// Return all periods for the user's schedules
		var scheduleIDs []uint
		if err := database.DB.Model(&models.WorkoutScheduleEntity{}).
			Where("owner = ?", userID).Pluck("id", &scheduleIDs).Error; err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if len(scheduleIDs) > 0 {
			if err := database.DB.Where("schedule_id IN ?", scheduleIDs).
				Order("period_start").Find(&entities).Error; err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
		}
	}

	now := time.Now()
	dtos := make([]models.SchedulePeriod, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO(now)
	}
	return &ListSchedulePeriodsOutput{Body: dtos}, nil
}

// CreateSchedulePeriod creates the first period for a schedule.
// POST /api/user/schedule-periods
func CreateSchedulePeriod(ctx context.Context, input *CreateSchedulePeriodInput) (*CreateSchedulePeriodOutput, error) {
	var schedule models.WorkoutScheduleEntity
	if err := database.DB.First(&schedule, input.Body.ScheduleID).Error; err != nil {
		return nil, huma.Error404NotFound("Schedule not found")
	}
	if schedule.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("not the owner of this schedule")
	}

	if !input.Body.PeriodEnd.After(input.Body.PeriodStart) {
		return nil, huma.Error400BadRequest("periodEnd must be after periodStart")
	}

	scheduleType := models.ScheduleType(input.Body.Type)
	switch scheduleType {
	case models.ScheduleTypeFixedDate, models.ScheduleTypeFrequency:
		// valid
	default:
		return nil, huma.Error400BadRequest("type must be 'fixed_date' or 'frequency'")
	}

	mode := models.PeriodModeNormal
	if input.Body.Mode != nil {
		switch models.PeriodMode(*input.Body.Mode) {
		case models.PeriodModeNormal, models.PeriodModeMonthly:
			mode = models.PeriodMode(*input.Body.Mode)
		default:
			return nil, huma.Error400BadRequest("mode must be 'normal' or 'monthly'")
		}
	}

	entity := models.SchedulePeriodEntity{
		ScheduleID:  input.Body.ScheduleID,
		PeriodStart: input.Body.PeriodStart,
		PeriodEnd:   input.Body.PeriodEnd,
		Type:        scheduleType,
		Mode:        mode,
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateSchedulePeriodOutput{Body: entity.ToDTO(time.Now())}, nil
}

// GetSchedulePeriodPermissions returns the permissions for a period based on its status.
// GET /api/user/schedule-periods/{id}/permissions
func GetSchedulePeriodPermissions(ctx context.Context, input *GetSchedulePeriodPermissionsInput) (*GetSchedulePeriodPermissionsOutput, error) {
	period, err := requirePeriodOwner(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	dto := period.ToDTO(time.Now())

	var perms []shared.Permission
	if dto.Status == models.PeriodStatusPlanned {
		perms = []shared.Permission{shared.PermissionRead, shared.PermissionModify, shared.PermissionDelete}
	} else {
		perms = []shared.Permission{shared.PermissionRead}
	}

	return &GetSchedulePeriodPermissionsOutput{
		Body: shared.PermissionsResponse{Permissions: perms},
	}, nil
}
