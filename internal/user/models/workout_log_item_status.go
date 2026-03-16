package models

import "fmt"

type WorkoutLogItemStatus string

const (
	WorkoutLogItemStatusPlanning          WorkoutLogItemStatus = "planning"
	WorkoutLogItemStatusInProgress        WorkoutLogItemStatus = "in_progress"
	WorkoutLogItemStatusFinished          WorkoutLogItemStatus = "finished"
	WorkoutLogItemStatusSkipped           WorkoutLogItemStatus = "skipped"
	WorkoutLogItemStatusPartiallyFinished WorkoutLogItemStatus = "partially_finished"
	WorkoutLogItemStatusAborted           WorkoutLogItemStatus = "aborted"
)

var validItemTransitions = map[WorkoutLogItemStatus][]WorkoutLogItemStatus{
	WorkoutLogItemStatusPlanning:          {WorkoutLogItemStatusInProgress},
	WorkoutLogItemStatusInProgress:        {WorkoutLogItemStatusFinished, WorkoutLogItemStatusSkipped, WorkoutLogItemStatusAborted},
	WorkoutLogItemStatusFinished:          {},
	WorkoutLogItemStatusSkipped:           {},
	WorkoutLogItemStatusPartiallyFinished: {},
	WorkoutLogItemStatusAborted:           {},
}

func (s WorkoutLogItemStatus) IsTerminal() bool {
	return s == WorkoutLogItemStatusFinished ||
		s == WorkoutLogItemStatusSkipped ||
		s == WorkoutLogItemStatusPartiallyFinished ||
		s == WorkoutLogItemStatusAborted
}

func (s WorkoutLogItemStatus) CanTransitionTo(target WorkoutLogItemStatus) bool {
	for _, t := range validItemTransitions[s] {
		if t == target {
			return true
		}
	}
	return false
}

func (s WorkoutLogItemStatus) TransitionTo(target WorkoutLogItemStatus) error {
	if !s.CanTransitionTo(target) {
		return fmt.Errorf("cannot transition from %s to %s", s, target)
	}
	return nil
}
