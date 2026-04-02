package models

import "testing"

func TestCanTransitionTo(t *testing.T) {
	allowed := []struct {
		from, to WorkoutLogStatus
	}{
		{WorkoutLogStatusPlanning, WorkoutLogStatusInProgress},
		{WorkoutLogStatusInProgress, WorkoutLogStatusFinished},
		{WorkoutLogStatusInProgress, WorkoutLogStatusPartiallyFinished},
		{WorkoutLogStatusInProgress, WorkoutLogStatusAborted},
		{WorkoutLogStatusProposed, WorkoutLogStatusCommitted},
		{WorkoutLogStatusProposed, WorkoutLogStatusSkipped},
		{WorkoutLogStatusCommitted, WorkoutLogStatusInProgress},
		{WorkoutLogStatusCommitted, WorkoutLogStatusBroken},
	}
	for _, tt := range allowed {
		if !tt.from.CanTransitionTo(tt.to) {
			t.Errorf("%s -> %s should be allowed", tt.from, tt.to)
		}
	}

	forbidden := []struct {
		from, to WorkoutLogStatus
	}{
		{WorkoutLogStatusPlanning, WorkoutLogStatusFinished},
		{WorkoutLogStatusPlanning, WorkoutLogStatusAborted},
		{WorkoutLogStatusPlanning, WorkoutLogStatusPartiallyFinished},
		{WorkoutLogStatusInProgress, WorkoutLogStatusPlanning},
		{WorkoutLogStatusFinished, WorkoutLogStatusPlanning},
		{WorkoutLogStatusFinished, WorkoutLogStatusInProgress},
		{WorkoutLogStatusFinished, WorkoutLogStatusAborted},
		{WorkoutLogStatusPartiallyFinished, WorkoutLogStatusPlanning},
		{WorkoutLogStatusPartiallyFinished, WorkoutLogStatusInProgress},
		{WorkoutLogStatusPartiallyFinished, WorkoutLogStatusFinished},
		{WorkoutLogStatusAborted, WorkoutLogStatusPlanning},
		{WorkoutLogStatusAborted, WorkoutLogStatusInProgress},
		{WorkoutLogStatusAborted, WorkoutLogStatusFinished},
		// Commitment-specific forbidden transitions
		{WorkoutLogStatusProposed, WorkoutLogStatusInProgress},
		{WorkoutLogStatusProposed, WorkoutLogStatusPlanning},
		{WorkoutLogStatusCommitted, WorkoutLogStatusSkipped},
		{WorkoutLogStatusCommitted, WorkoutLogStatusPlanning},
		{WorkoutLogStatusSkipped, WorkoutLogStatusCommitted},
		{WorkoutLogStatusSkipped, WorkoutLogStatusInProgress},
		{WorkoutLogStatusBroken, WorkoutLogStatusInProgress},
		{WorkoutLogStatusBroken, WorkoutLogStatusCommitted},
	}
	for _, tt := range forbidden {
		if tt.from.CanTransitionTo(tt.to) {
			t.Errorf("%s -> %s should be forbidden", tt.from, tt.to)
		}
	}
}

func TestTransitionTo(t *testing.T) {
	if err := WorkoutLogStatusPlanning.TransitionTo(WorkoutLogStatusInProgress); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if err := WorkoutLogStatusFinished.TransitionTo(WorkoutLogStatusPlanning); err == nil {
		t.Error("expected error for terminal -> planning transition")
	}
}

func TestIsTerminalConsistency(t *testing.T) {
	for status, targets := range validTransitions {
		isTerminal := status.IsTerminal()
		hasNoTransitions := len(targets) == 0
		if isTerminal != hasNoTransitions {
			t.Errorf("%s: IsTerminal=%v but has %d valid transitions", status, isTerminal, len(targets))
		}
	}
}
