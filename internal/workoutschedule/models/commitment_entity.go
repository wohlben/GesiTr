package models

import (
	"time"

	"gesitr/internal/shared"
)

// ScheduleCommitmentEntity is the join between a period and a workout log.
// It has a two-phase lifecycle:
//
//  1. Plan phase: created with WorkoutLogID=null. For fixed_date schedules,
//     Date is set to the specific target day. For frequency schedules,
//     Date is null (the user decides when during the period).
//  2. Activation phase: when the period's start date arrives, a WorkoutLog
//     is created and WorkoutLogID is set. The commitment is now fulfilled.
//
// When a period is cloned forward, each commitment is replicated into the
// new period with WorkoutLogID reset to null. For fixed_date commitments,
// the day offset relative to the period start is preserved.
//
// Commitments can only be deleted while WorkoutLogID is null (before
// activation). Once linked to a workout log, the commitment is immutable.
//
// Relationships:
//   - belongs to a [SchedulePeriodEntity] (PeriodID)
//   - optionally links to a WorkoutLog (WorkoutLogID)
type ScheduleCommitmentEntity struct {
	shared.BaseModel
	PeriodID     uint `gorm:"not null;index"`
	Date         *time.Time
	WorkoutLogID *uint `gorm:"index"`
}

func (ScheduleCommitmentEntity) TableName() string { return "schedule_commitments" }

func (e *ScheduleCommitmentEntity) ToDTO() ScheduleCommitment {
	return ScheduleCommitment{
		BaseModel:    e.BaseModel,
		PeriodID:     e.PeriodID,
		Date:         e.Date,
		WorkoutLogID: e.WorkoutLogID,
	}
}
