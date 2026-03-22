# Compendium Workout Groups

Enable users to form groups around shared workouts and track progress together.

**Depends on:** [Compendium Workouts](./compendium-workouts.md)

## Workout Groups

A new entity: **user-workout-group** — a named group that ties users together around shared workouts.

Each group has **participants**, where each participant has a **role**:

- **Planner** — can modify group membership (invite/remove members), assign workouts, and manage the group. Shared ownership: all planners have equal control.
- **Member** — can view group workouts, import them, and track progress.

## Workflow

1. A planner creates a workout group and invites participants.
2. Planners attach compendium workouts to the group.
3. Each participant imports the workout into their "my workouts."
4. Participants complete the workout and log their results as usual (exercise-logs tied to their imported my-exercises).
5. All group members can see each other's results — the exercise-logs inserted in relation to each imported my-exercise are visible to the group.

## Open Questions

- How are invites handled? Link-based? Username search? Approval flow?
- Can members leave a group on their own?
- Should there be a feed/timeline view of group activity?
- Privacy controls — can a member hide specific logs from the group?
- Should planners be able to schedule workouts (e.g., "do this on Monday")?
