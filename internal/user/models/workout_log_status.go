package models

import "fmt"

type WorkoutLogStatus string

const (
	WorkoutLogStatusPlanning   WorkoutLogStatus = "planning"
	WorkoutLogStatusInProgress WorkoutLogStatus = "in_progress"
	WorkoutLogStatusFinished   WorkoutLogStatus = "finished"
	WorkoutLogStatusAborted    WorkoutLogStatus = "aborted"
)

// validTransitions defines the allowed state transitions.
// This is the single source of truth for the workout log state machine.
var validTransitions = map[WorkoutLogStatus][]WorkoutLogStatus{
	WorkoutLogStatusPlanning:   {WorkoutLogStatusInProgress},
	WorkoutLogStatusInProgress: {WorkoutLogStatusFinished, WorkoutLogStatusAborted},
	WorkoutLogStatusFinished:   {},
	WorkoutLogStatusAborted:    {},
}

func (s WorkoutLogStatus) IsTerminal() bool {
	return s == WorkoutLogStatusFinished || s == WorkoutLogStatusAborted
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
