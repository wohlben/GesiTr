# Workout Group Cycles

## Summary

Workout groups operate on a rolling **cycle** system (typically weekly). Each cycle, every participant must commit to a plan — concretely, they configure their exercise schemes for the group's assigned workouts. This solves a key UX problem: when a user imports a group workout, they don't go through the normal workout creation/editing flow, so there's no natural point where they configure sets, reps, weight, etc. The cycle commitment process forces that configuration step while also producing trackable commitments.

**Depends on:** [Compendium Workout Groups](./implemented/compendium-workout-groups.md), [Workout Commitments](./workout-commitments.md), [Workout Schedules](./workout-schedules.md)

## Cycle Model

A group operates on a **configurable cycle length** (default: 1 week). The group always works one cycle forward — while the current cycle is in progress, participants are committing to the *next* cycle.

Key fields on the group:
- **Cycle length** — duration of one cycle (e.g., 1 week).
- **Cycle start day** — when cycles begin (e.g., Monday).
- **Planning window** — how far ahead the next cycle opens for commitment (e.g., the cycle becomes plannable mid-week).

## Commitment Flow

Each cycle, every participant must:

1. **Review** the workouts assigned to them for the upcoming cycle.
2. **Configure exercise schemes** — set their target sets, reps, weight, rest times, etc. for each exercise in each workout. This is the core purpose: it forces the user to make concrete plans rather than arriving at a workout with unconfigured exercises.
3. **Commit** — confirm their plan for the cycle.

The result of this process is effectively a **workout schedule** for that cycle — the participant has declared what they will do and when. This naturally feeds into the [Workout Commitments](./workout-commitments.md) system: each committed workout becomes a commitment with `COMMITTED` status.

## Cycle States

A cycle progresses through states:

- **PLANNING** — the cycle is open for commitments. Participants are configuring their schemes and committing.
- **ACTIVE** — all participants have committed (or the planning window has closed) and the cycle is in progress. Commitments are locked.
- **COMPLETED** — the cycle's date range has passed.

Transition from `PLANNING` → `ACTIVE` happens when **all members have committed** for that cycle. This is a group-level gate — the cycle doesn't activate until everyone has planned.

## Relationship to Workout Schedules

The output of a cycle commitment is conceptually a workout schedule, but scoped to a single cycle and driven by group context rather than individual configuration. Rather than creating a persistent workout schedule entity, the cycle commitment produces:

- Exercise scheme configurations on the user's imported workout (the concrete sets/reps/weight).
- Workout commitments for each planned session (date + workout).

If a user also has a personal workout schedule, the group cycle commitments coexist alongside it — they're additive.

## Considerations

- What happens if a participant doesn't commit before the planning window closes? Auto-carry-forward from last cycle? Mark as uncommitted?
- Can planners override or suggest schemes for participants?
- Should there be a "propose" step where planners suggest a plan that participants then accept/modify?
- How does this interact with workout modifications — if the group workout template changes mid-cycle, do existing commitments update?
- Visibility: can participants see each other's committed schemes before the cycle activates?
- What if a new member joins mid-cycle?
