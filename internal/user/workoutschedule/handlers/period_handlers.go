package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/user/workoutschedule/models"

	"github.com/danielgtaylor/huma/v2"
)

// ListSchedulePeriods returns periods for a given schedule.
// GET /api/user/schedule-periods
func ListSchedulePeriods(ctx context.Context, input *ListSchedulePeriodsInput) (*ListSchedulePeriodsOutput, error) {
	var schedule models.WorkoutScheduleEntity
	if err := database.DB.First(&schedule, input.ScheduleID).Error; err != nil {
		return nil, huma.Error404NotFound("Schedule not found")
	}
	if schedule.Owner != humaconfig.GetUserID(ctx) {
		return nil, huma.Error403Forbidden("not the owner of this schedule")
	}

	var entities []models.SchedulePeriodEntity
	if err := database.DB.Where("schedule_id = ?", input.ScheduleID).
		Order("period_start").Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	dtos := make([]models.SchedulePeriod, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
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
	return &CreateSchedulePeriodOutput{Body: entity.ToDTO()}, nil
}
