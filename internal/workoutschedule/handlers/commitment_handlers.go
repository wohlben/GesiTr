package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/humaconfig"
	"gesitr/internal/workoutschedule/models"

	"github.com/danielgtaylor/huma/v2"
)

// requirePeriodOwner checks that the period belongs to a schedule owned by the current user.
func requirePeriodOwner(ctx context.Context, periodID uint) (models.SchedulePeriodEntity, error) {
	var period models.SchedulePeriodEntity
	if err := database.DB.First(&period, periodID).Error; err != nil {
		return period, huma.Error404NotFound("Period not found")
	}
	var schedule models.WorkoutScheduleEntity
	if err := database.DB.First(&schedule, period.ScheduleID).Error; err != nil {
		return period, huma.Error404NotFound("Schedule not found")
	}
	if err := requireOwner(ctx, schedule.Owner); err != nil {
		return period, err
	}
	return period, nil
}

// ListScheduleCommitments returns commitments for a given period, or all
// commitments for the user when periodId is omitted.
// GET /api/user/schedule-commitments
func ListScheduleCommitments(ctx context.Context, input *ListScheduleCommitmentsInput) (*ListScheduleCommitmentsOutput, error) {
	userID := humaconfig.GetUserID(ctx)

	var entities []models.ScheduleCommitmentEntity
	if input.PeriodID != "" {
		// Filter by specific period
		if _, err := requirePeriodOwner(ctx, parseUint(input.PeriodID)); err != nil {
			return nil, err
		}
		if err := database.DB.Where("period_id = ?", input.PeriodID).
			Order("date").Find(&entities).Error; err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
	} else {
		// Return all commitments for the user's schedules → periods
		var periodIDs []uint
		if err := database.DB.Model(&models.SchedulePeriodEntity{}).
			Joins("JOIN workout_schedules ON workout_schedules.id = schedule_periods.schedule_id").
			Where("workout_schedules.owner = ?", userID).
			Pluck("schedule_periods.id", &periodIDs).Error; err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if len(periodIDs) > 0 {
			if err := database.DB.Where("period_id IN ?", periodIDs).
				Order("date").Find(&entities).Error; err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
		}
	}

	dtos := make([]models.ScheduleCommitment, len(entities))
	for i := range entities {
		dtos[i] = entities[i].ToDTO()
	}
	return &ListScheduleCommitmentsOutput{Body: dtos}, nil
}

// CreateScheduleCommitment creates a commitment within a period.
// POST /api/user/schedule-commitments
func CreateScheduleCommitment(ctx context.Context, input *CreateScheduleCommitmentInput) (*CreateScheduleCommitmentOutput, error) {
	if _, err := requirePeriodOwner(ctx, input.Body.PeriodID); err != nil {
		return nil, err
	}

	entity := models.ScheduleCommitmentEntity{
		PeriodID: input.Body.PeriodID,
		Date:     input.Body.Date,
	}
	if err := database.DB.Create(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &CreateScheduleCommitmentOutput{Body: entity.ToDTO()}, nil
}

// DeleteScheduleCommitment deletes a commitment (only if not yet linked to a workout log).
// DELETE /api/user/schedule-commitments/{id}
func DeleteScheduleCommitment(ctx context.Context, input *DeleteScheduleCommitmentInput) (*DeleteScheduleCommitmentOutput, error) {
	var entity models.ScheduleCommitmentEntity
	if err := database.DB.First(&entity, input.ID).Error; err != nil {
		return nil, huma.Error404NotFound("Commitment not found")
	}
	if _, err := requirePeriodOwner(ctx, entity.PeriodID); err != nil {
		return nil, err
	}
	if entity.WorkoutLogID != nil {
		return nil, huma.Error409Conflict("cannot delete a commitment that is already linked to a workout log")
	}
	if err := database.DB.Delete(&entity).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}

func parseUint(s string) uint {
	var n uint
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + uint(c-'0')
		}
	}
	return n
}
