package models

import "fmt"

// WorkoutLogStatus represents the lifecycle state of a workout log.
//
// There are two independent flows:
//
// Manual flow:
//
//	planning → in_progress → finished / partially_finished / aborted
//	(adhoc is a variant of in_progress that allows structural edits)
//
// Commitment flow:
//
//	proposed → committed → in_progress → finished / partially_finished / aborted
//	proposed → skipped  (user declines, or due window elapses)
//	committed → broken  (due window elapses without starting)
type WorkoutLogStatus string

const (
	// WorkoutLogStatusPlanning means the user is planning a workout log and can still
	// change it. Fully editable: sections, exercises, sets, and exercise schemes can
	// all be added, removed, or modified.
	// Transitions to: in_progress.
	WorkoutLogStatusPlanning WorkoutLogStatus = "planning"

	// WorkoutLogStatusInProgress means the workout is actively being performed.
	// The workout structure is readonly — only exercise logs can be created and
	// edited (recording actual performance against targets).
	// Transitions to: finished, partially_finished, aborted.
	WorkoutLogStatusInProgress WorkoutLogStatus = "in_progress"

	// WorkoutLogStatusAdhoc means the workout was started without prior planning.
	// Exercises (but not exercise groups) can be added to the log, but once an
	// exercise log is created for an exercise, that exercise becomes readonly.
	// Transitions to: finished, partially_finished, aborted.
	WorkoutLogStatusAdhoc WorkoutLogStatus = "adhoc"

	// WorkoutLogStatusFinished means all exercises and sets were completed.
	// Terminal state — fully readonly, no changes allowed.
	WorkoutLogStatusFinished WorkoutLogStatus = "finished"

	// WorkoutLogStatusPartiallyFinished means some exercises or sets were completed
	// while others were skipped.
	// Terminal state — fully readonly, no changes allowed.
	WorkoutLogStatusPartiallyFinished WorkoutLogStatus = "partially_finished"

	// WorkoutLogStatusAborted means the workout was abandoned mid-execution.
	// Terminal state — fully readonly, no changes allowed.
	WorkoutLogStatusAborted WorkoutLogStatus = "aborted"

	// WorkoutLogStatusProposed means a schedule or group cycle created this log.
	// The workout structure (exercises, exercise groups) is readonly — the user
	// cannot add or remove exercises. The user's task is to create exercise schemes
	// (sets, reps, weight) for each exercise before committing.
	// Transitions to: committed (user accepts), skipped (user declines or due window elapses).
	WorkoutLogStatusProposed WorkoutLogStatus = "proposed"

	// WorkoutLogStatusCommitted means the user accepted a proposed workout and locked
	// in their exercise schemes. Everything is readonly — structure and schemes.
	// The log is waiting to be started within the due window.
	// Transitions to: in_progress (user starts), broken (due window elapses).
	WorkoutLogStatusCommitted WorkoutLogStatus = "committed"

	// WorkoutLogStatusSkipped means the user declined a proposed workout, or the
	// due window elapsed without the user committing.
	// Terminal state — fully readonly, no changes allowed.
	WorkoutLogStatusSkipped WorkoutLogStatus = "skipped"

	// WorkoutLogStatusBroken means a committed workout's due window passed without
	// the user starting it. Set automatically by a background ticker.
	// Terminal state — fully readonly, no changes allowed.
	WorkoutLogStatusBroken WorkoutLogStatus = "broken"
)

// validTransitions defines the allowed state transitions.
// This is the single source of truth for the workout log state machine.
var validTransitions = map[WorkoutLogStatus][]WorkoutLogStatus{
	WorkoutLogStatusPlanning:          {WorkoutLogStatusInProgress},
	WorkoutLogStatusInProgress:        {WorkoutLogStatusFinished, WorkoutLogStatusPartiallyFinished, WorkoutLogStatusAborted},
	WorkoutLogStatusAdhoc:             {WorkoutLogStatusFinished, WorkoutLogStatusPartiallyFinished, WorkoutLogStatusAborted},
	WorkoutLogStatusFinished:          {},
	WorkoutLogStatusPartiallyFinished: {},
	WorkoutLogStatusAborted:           {},
	WorkoutLogStatusProposed:          {WorkoutLogStatusCommitted, WorkoutLogStatusSkipped},
	WorkoutLogStatusCommitted:         {WorkoutLogStatusInProgress, WorkoutLogStatusBroken},
	WorkoutLogStatusSkipped:           {},
	WorkoutLogStatusBroken:            {},
}

func (s WorkoutLogStatus) IsTerminal() bool {
	return s == WorkoutLogStatusFinished || s == WorkoutLogStatusPartiallyFinished || s == WorkoutLogStatusAborted || s == WorkoutLogStatusSkipped || s == WorkoutLogStatusBroken
}

// CanTransitionTo reports whether transitioning from s to target is valid.
func (s WorkoutLogStatus) CanTransitionTo(target WorkoutLogStatus) bool {
	for _, t := range validTransitions[s] {
		if t == target {
			return true
		}
	}
	return false
}

// TransitionTo validates the transition and returns an error if it is invalid.
func (s WorkoutLogStatus) TransitionTo(target WorkoutLogStatus) error {
	if !s.CanTransitionTo(target) {
		return fmt.Errorf("cannot transition from %s to %s", s, target)
	}
	return nil
}
