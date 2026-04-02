// Package models defines the data types for workout schedules.
//
// # Overview
//
// A workout schedule automates the creation of workout logs by defining a
// recurrence pattern. The user configures a template period with commitments,
// and the system clones that pattern forward in time. When a period becomes
// active (its start date arrives), the commitments are fulfilled by creating
// workout logs.
//
// # Entity Hierarchy
//
//	Schedule → Period → Commitment → WorkoutLog
//
// [WorkoutScheduleEntity] defines what type of recurrence (fixed_date or
// frequency), the active date range, and which status generated logs start
// in (committed or proposed).
//
// [SchedulePeriodEntity] is a concrete time window. The first period is
// created manually by the user. Subsequent periods are cloned from the last
// one — same duration, same commitment pattern shifted forward in time.
//
// [ScheduleCommitmentEntity] represents a planned workout within a period.
// It is created first with WorkoutLogID=null (the plan), then linked to a
// WorkoutLog when the period activates (the execution).
//
// # Schedule Types
//
//   - fixed_date: the user picks specific days within the period by setting
//     Date on each commitment. When cloned, day offsets relative to the
//     period start are preserved.
//   - frequency: the user sets a count (N commitments per period) without
//     picking specific days. Commitments have Date=null.
//
// # Two-Phase Generation
//
// Generation is lazy — triggered when the user lists their workout logs.
//
// Phase 1 (Clone): if the last period has ended and no next period exists,
// clone it forward. A new period is created with the same duration, starting
// the day after the previous one ended. Commitments are cloned with
// WorkoutLogID=null.
//
// Phase 2 (Activate): for each period where PeriodStart ≤ now, find
// commitments with WorkoutLogID=null. For each, create a WorkoutLog in the
// schedule's InitialStatus and link the commitment to it.
//
// # Active Derivation
//
// A schedule has no stored "active" flag. It is derived:
// active = StartDate ≤ now AND (EndDate is null OR EndDate ≥ now).
// StartDate defaults to tomorrow when not specified, so the first period
// is never immediately activated on creation.
package models
