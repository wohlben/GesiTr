# Workout Commitments

## Summary

A **workout commitment** is a promise to perform a specific workout (and hence log it) on a given day. Commitments are the missing layer between workout schedules and actual workout logs — schedules define *patterns*, commitments materialize those patterns into concrete, date-bound obligations that can be tracked and fulfilled.

**Depends on:** [Workout Schedules](./workout-schedules.md)

## Commitment Model

A commitment ties a user to a workout on a specific date.

Key fields:
- **User** — the user who is committed.
- **Workout** — the workout to be performed.
- **Date** — the day the workout should be done.
- **Status** — derived via a state machine (see below).

## Status State Machine

Commitment status is **not stored directly** — it is evaluated by a state machine from underlying boolean/data fields. This is necessary because more states will be added in the future and the evaluation logic needs to remain centralized.

Initial states:
- **PROPOSED** — the commitment has been generated but the user has not yet confirmed it.
- **COMMITTED** — the user has accepted the commitment.

The `COMMITTED` status evaluates from a `committed` boolean on the commitment record. Future states (e.g., `COMPLETED`, `SKIPPED`, `MISSED`) will follow the same pattern: stored fields feeding into the state machine evaluation.

## Creation Flows

### 1. From Workout Schedules

A workout schedule can automatically spawn commitments based on a configurable trigger:

- The **spawn point** is configured on the schedule itself.
- Example: a weekly schedule generates a commitment each week, with the day of week being configurable on the schedule.
- The generated commitment starts in `PROPOSED` status.

### 2. From Workout Groups

Workout group planners can generate commitments for group members:

- Commitments are created per user in the group.
- The generated commitment starts in `PROPOSED` status.
- Workout groups will need their own scheduling concept to drive commitment generation, but the requirements differ enough from workout schedules that this should be designed as a separate feature. For now, the commitment model should only be aware that its source can be a workout group — the group scheduling mechanism itself is out of scope.

## Considerations

- How does a user transition a commitment from `PROPOSED` to `COMMITTED`? Explicit accept button? Bulk accept?
- What happens when a commitment date passes without a log — auto-transition to a `MISSED` state (future)?
- Should commitments link to the resulting workout log once completed?
- UI placement: calendar view? List on the workout detail page? Dashboard widget?
- Can a user manually create a commitment without a schedule or group?
- Notification/reminder system for upcoming commitments?
