package handlers

import (
	"time"

	"gesitr/internal/shared"
	"gesitr/internal/workoutschedule/models"
)

// --- Schedule handlers ---

type ListWorkoutSchedulesInput struct {
	WorkoutID string `query:"workoutId" doc:"Filter by workout ID"`
}

type ListWorkoutSchedulesOutput struct {
	Body []models.WorkoutSchedule
}

type WorkoutScheduleBody struct {
	WorkoutID     uint       `json:"workoutId" required:"true"`
	StartDate     *time.Time `json:"startDate,omitempty" doc:"Default: tomorrow"`
	EndDate       *time.Time `json:"endDate,omitempty"`
	InitialStatus *string    `json:"initialStatus,omitempty" doc:"'committed' (default) or 'proposed'"`
}

type CreateWorkoutScheduleInput struct {
	Body WorkoutScheduleBody
}

type CreateWorkoutScheduleOutput struct {
	Body models.WorkoutSchedule
}

type GetWorkoutScheduleInput struct {
	ID uint `path:"id"`
}

type GetWorkoutScheduleOutput struct {
	Body models.WorkoutSchedule
}

type UpdateWorkoutScheduleBody struct {
	EndDate *time.Time `json:"endDate,omitempty"`
}

type UpdateWorkoutScheduleInput struct {
	ID   uint `path:"id"`
	Body UpdateWorkoutScheduleBody
}

type UpdateWorkoutScheduleOutput struct {
	Body models.WorkoutSchedule
}

type DeleteWorkoutScheduleInput struct {
	ID uint `path:"id"`
}

type DeleteWorkoutScheduleOutput struct{}

// --- Period handlers ---

type ListSchedulePeriodsInput struct {
	ScheduleID string `query:"scheduleId" doc:"Filter by schedule ID (omit to list all periods for the user)"`
}

type ListSchedulePeriodsOutput struct {
	Body []models.SchedulePeriod
}

type CreateSchedulePeriodBody struct {
	ScheduleID  uint      `json:"scheduleId" required:"true"`
	PeriodStart time.Time `json:"periodStart" required:"true"`
	PeriodEnd   time.Time `json:"periodEnd" required:"true"`
	Type        string    `json:"type" required:"true" doc:"'fixed_date' or 'frequency'"`
	Mode        *string   `json:"mode,omitempty" doc:"'normal' (default) or 'monthly'"`
}

type CreateSchedulePeriodInput struct {
	Body CreateSchedulePeriodBody
}

type CreateSchedulePeriodOutput struct {
	Body models.SchedulePeriod
}

type GetSchedulePeriodPermissionsInput struct {
	ID uint `path:"id"`
}

type GetSchedulePeriodPermissionsOutput struct {
	Body shared.PermissionsResponse
}

// --- Commitment handlers ---

type ListScheduleCommitmentsInput struct {
	PeriodID string `query:"periodId" doc:"Filter by period ID (omit to list all commitments for the user)"`
}

type ListScheduleCommitmentsOutput struct {
	Body []models.ScheduleCommitment
}

type CreateScheduleCommitmentBody struct {
	PeriodID uint       `json:"periodId" required:"true"`
	Date     *time.Time `json:"date,omitempty" doc:"Specific date for fixed_date schedules"`
}

type CreateScheduleCommitmentInput struct {
	Body CreateScheduleCommitmentBody
}

type CreateScheduleCommitmentOutput struct {
	Body models.ScheduleCommitment
}

type DeleteScheduleCommitmentInput struct {
	ID uint `path:"id"`
}

type DeleteScheduleCommitmentOutput struct{}
