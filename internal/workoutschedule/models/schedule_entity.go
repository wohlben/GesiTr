package models

import (
	"time"

	"gesitr/internal/shared"
)

// WorkoutScheduleEntity defines a recurring schedule for a workout.
//
// The schedule's only job is to clone the last period forward when it ends.
// It carries no type-specific configuration — that lives on the
// [SchedulePeriodEntity] and its [ScheduleCommitmentEntity] records.
//
// Active is derived from StartDate and EndDate, not stored. See [IsActive].
//
// InitialStatus controls what status generated WorkoutLogs start in:
// "committed" (default, locked immediately) or "proposed" (requires user
// review before locking).
//
// Relationships:
//   - belongs to a Workout (WorkoutID)
//   - owns many [SchedulePeriodEntity] records
type WorkoutScheduleEntity struct {
	shared.BaseModel
	Owner         string    `gorm:"not null;index"`
	WorkoutID     uint      `gorm:"not null;index"`
	StartDate     time.Time `gorm:"not null"`
	EndDate       *time.Time
	InitialStatus string `gorm:"not null;default:'committed'"`
}

func (WorkoutScheduleEntity) TableName() string { return "workout_schedules" }

// IsActive returns whether the schedule is active at the given time.
// Active means startDate ≤ now and (endDate is null or endDate ≥ now).
func (e *WorkoutScheduleEntity) IsActive(now time.Time) bool {
	if now.Before(e.StartDate) {
		return false
	}
	if e.EndDate != nil && now.After(*e.EndDate) {
		return false
	}
	return true
}

func (e *WorkoutScheduleEntity) ToDTO(now time.Time) WorkoutSchedule {
	return WorkoutSchedule{
		BaseModel:     e.BaseModel,
		Owner:         e.Owner,
		WorkoutID:     e.WorkoutID,
		StartDate:     e.StartDate,
		EndDate:       e.EndDate,
		InitialStatus: e.InitialStatus,
		Active:        e.IsActive(now),
	}
}
