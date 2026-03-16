package models

import "testing"

func TestItemCanTransitionTo(t *testing.T) {
	allowed := []struct {
		from, to WorkoutLogItemStatus
	}{
		{WorkoutLogItemStatusPlanning, WorkoutLogItemStatusInProgress},
		{WorkoutLogItemStatusInProgress, WorkoutLogItemStatusFinished},
		{WorkoutLogItemStatusInProgress, WorkoutLogItemStatusSkipped},
		{WorkoutLogItemStatusInProgress, WorkoutLogItemStatusAborted},
	}
	for _, tt := range allowed {
		if !tt.from.CanTransitionTo(tt.to) {
			t.Errorf("%s -> %s should be allowed", tt.from, tt.to)
		}
	}

	forbidden := []struct {
		from, to WorkoutLogItemStatus
	}{
		{WorkoutLogItemStatusPlanning, WorkoutLogItemStatusFinished},
		{WorkoutLogItemStatusPlanning, WorkoutLogItemStatusSkipped},
		{WorkoutLogItemStatusPlanning, WorkoutLogItemStatusAborted},
		{WorkoutLogItemStatusInProgress, WorkoutLogItemStatusPlanning},
		{WorkoutLogItemStatusFinished, WorkoutLogItemStatusPlanning},
		{WorkoutLogItemStatusFinished, WorkoutLogItemStatusInProgress},
		{WorkoutLogItemStatusSkipped, WorkoutLogItemStatusPlanning},
		{WorkoutLogItemStatusSkipped, WorkoutLogItemStatusInProgress},
		{WorkoutLogItemStatusPartiallyFinished, WorkoutLogItemStatusPlanning},
		{WorkoutLogItemStatusPartiallyFinished, WorkoutLogItemStatusInProgress},
		{WorkoutLogItemStatusAborted, WorkoutLogItemStatusPlanning},
		{WorkoutLogItemStatusAborted, WorkoutLogItemStatusInProgress},
	}
	for _, tt := range forbidden {
		if tt.from.CanTransitionTo(tt.to) {
			t.Errorf("%s -> %s should be forbidden", tt.from, tt.to)
		}
	}
}

func TestItemIsTerminalConsistency(t *testing.T) {
	for status, targets := range validItemTransitions {
		isTerminal := status.IsTerminal()
		hasNoTransitions := len(targets) == 0
		if isTerminal != hasNoTransitions {
			t.Errorf("%s: IsTerminal=%v but has %d valid transitions", status, isTerminal, len(targets))
		}
	}
}
