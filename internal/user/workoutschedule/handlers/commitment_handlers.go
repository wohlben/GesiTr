package handlers

import (
	"context"

	"gesitr/internal/database"
	"gesitr/internal/user/workoutschedule/models"

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

// ListScheduleCommitments returns commitments for a given period.
// GET /api/user/schedule-commitments
func ListScheduleCommitments(ctx context.Context, input *ListScheduleCommitmentsInput) (*ListScheduleCommitmentsOutput, error) {
	periodID := input.PeriodID
	if _, err := requirePeriodOwner(ctx, parseUint(periodID)); err != nil {
		return nil, err
	}

	var entities []models.ScheduleCommitmentEntity
	if err := database.DB.Where("period_id = ?", periodID).
		Order("date").Find(&entities).Error; err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
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
