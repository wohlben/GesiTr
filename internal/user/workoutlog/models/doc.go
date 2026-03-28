// Package models defines the data types for workout logs.
//
// A workout log records a single workout session — either planned from a
// workout template or started ad-hoc. The log contains sections, exercises,
// and sets that mirror the workout structure but with actual performance data.
//
// # State Machine
//
// [WorkoutLogStatus] governs the lifecycle of a workout log. There are two
// independent flows:
//
// Manual flow (user-initiated):
//
//	planning → in_progress → finished / partially_finished / aborted
//
// Commitment flow (schedule/group-initiated):
//
//	proposed → committed → in_progress → finished / partially_finished / aborted
//	proposed → skipped
//	committed → broken
//
// See the individual status constants for permission models and editability
// rules in each state.
package models
