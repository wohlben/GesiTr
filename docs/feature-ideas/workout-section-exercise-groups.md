# Workout Section Exercise Groups

## Summary

Allow a workout section to reference either a **specific exercise** (current behavior) or an **exercise group**. When a section uses an exercise group, the user picks which specific exercise from that group to perform each time they log the workout — introducing per-session variety.

## Motivation

Many training programs prescribe a slot like "pick a horizontal pull" or "any quad-dominant movement" rather than a fixed exercise. Currently users must either pick one exercise and stick with it, or duplicate sections for each variant. Exercise groups let the workout template express this flexibility directly.

## Design

A workout section currently has:
- A reference to one exercise

The change: a section references **either** a single exercise **or** an exercise group (mutually exclusive).

An **exercise group** is a named collection of exercises (e.g., "Horizontal Pulls" containing barbell row, cable row, dumbbell row).

### Logging Flow

1. User starts logging a workout.
2. For sections with a specific exercise — unchanged, exercise is pre-filled.
3. For sections with an exercise group — the user is prompted to pick one exercise from the group before logging sets.
4. The exercise-log records which specific exercise was chosen for that session.

### Where Do Exercise Groups Live?

Options:
- **Compendium-level** — shared groups curated alongside the exercise library. Any user can reference them.
- **User-level** — personal groups a user defines for their own workouts.
- **Both** — compendium provides defaults, users can create their own.

## Considerations

- How does this interact with progress tracking? If the user alternates between barbell row and cable row, are their logs shown together (grouped) or separately?
- Should the UI suggest the exercise the user picked last time, or rotate?
- Can a user convert an existing single-exercise section to a group section (and vice versa) without losing log history?
