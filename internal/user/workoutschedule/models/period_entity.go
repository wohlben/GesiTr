package models

import (
	"time"

	"gesitr/internal/shared"
)

// SchedulePeriodEntity represents a concrete time window within a schedule.
//
// The first period is created by the user via the API (picking start and end
// dates, a type, and a mode). Subsequent periods are generated automatically
// by cloning the last period.
//
// Type determines the behaviour of commitments within this period:
//   - fixed_date: each commitment has a specific Date
//   - frequency: commitments have no Date; the count of commitments is the
//     number of workouts required
//
// Mode determines how the period duration is computed when cloning:
//   - normal: same duration in days, next period starts day after previous ended
//   - monthly: advances one calendar month, preserving the day-of-month offset
//
// The last period for a schedule serves as the template — its type, mode,
// and commitment pattern are replicated into each new period.
//
// Relationships:
//   - belongs to a [WorkoutScheduleEntity] (ScheduleID)
//   - owns many [ScheduleCommitmentEntity] records
type SchedulePeriodEntity struct {
	shared.BaseModel
	ScheduleID  uint         `gorm:"not null;index;uniqueIndex:idx_schedule_period"`
	PeriodStart time.Time    `gorm:"not null;uniqueIndex:idx_schedule_period"`
	PeriodEnd   time.Time    `gorm:"not null;uniqueIndex:idx_schedule_period"`
	Type        ScheduleType `gorm:"not null"`
	Mode        PeriodMode   `gorm:"not null;default:'normal'"`
}

func (SchedulePeriodEntity) TableName() string { return "schedule_periods" }

func (e *SchedulePeriodEntity) ToDTO(now time.Time) SchedulePeriod {
	// PeriodEnd is stored as 23:59:59 of the last inclusive day,
	// so now.Before(PeriodEnd) naturally covers the full end date.
	var status PeriodStatus
	switch {
	case now.Before(e.PeriodStart):
		status = PeriodStatusPlanned
	case now.Before(e.PeriodEnd):
		status = PeriodStatusActive
	default:
		status = PeriodStatusArchived
	}
	return SchedulePeriod{
		BaseModel:   e.BaseModel,
		ScheduleID:  e.ScheduleID,
		PeriodStart: e.PeriodStart,
		PeriodEnd:   e.PeriodEnd,
		Type:        e.Type,
		Mode:        e.Mode,
		Status:      status,
	}
}
