# Remove Exercise Imports

## Summary

With mastery tracking in place, the "import to My Exercises" concept becomes redundant. This feature removes the import flow and replaces it with a unified exercise view where any compendium exercise can be used directly in workouts.

## Changes

### 1. "Import" becomes "Fork"

The import flow still exists but is rebranded as **"Fork"**. Forking creates a user-owned copy of a compendium exercise — useful when a user genuinely wants to customize the exercise definition (different instructions, different equipment, etc.). This should be the exception, not the default way to use exercises.

### 2. "My Exercises" becomes "Exercises"

The current "My Exercises" page is replaced with a unified **"Exercises"** view. The default filter shows:
- Exercises the user has **any mastery in** (i.e., has logged at least once)
- Exercises the user is the **owner** of (forked/created exercises)

This requires a precomputed **`mastery_experience`** table to efficiently filter exercises by "has mastery":

```
mastery_experience
  owner        string  -- user ID
  exercise_id  uint    -- exercise with any logged experience
  total_reps   int     -- precomputed rep count (for sorting/display)
```

Updated on exercise log creation. This is the lightweight counterpart to the full mastery computation — just tracks "has this user ever done this exercise" for fast filtering.

### 3. Remove `user/exercises` endpoint

The `GET /api/user/exercises` endpoint becomes entirely redundant. Replace it with the standard `GET /api/exercises` endpoint, enhanced with mastery-aware query parameters:

- `mastery=me` — filter to exercises the current user has mastery in OR owns
- Default sort when `mastery=me`: owned first, then by mastery level descending

All frontend queries that currently hit `user/exercises` switch to `exercises?mastery=me`.

### 4. Exercise picker in workouts

When selecting exercises for a workout, **every exercise is valid** — not just owned/imported ones. The picker should:
- Show mastered + owned exercises first (the "primary selection")
- Show all remaining compendium exercises below, separated visually
- All exercises are selectable regardless of section

This means the workout exercise picker query becomes: `GET /api/exercises?sort=mastery` (or similar), which returns all exercises ordered with the user's mastered/owned ones at the top.

### 5. Workouts reference any exercise

Workout section items / exercise schemes must accept any exercise ID, not just user-owned exercises. Check and remove any ownership validation on exercise references in:
- Workout creation/editing
- Exercise scheme creation
- Workout log exercise references

## Migration

- Existing "imported" exercises remain as user-owned (they become "forked" exercises retroactively)
- No data migration needed — the ownership model doesn't change, only the UI flow
- The `user/exercises` endpoint can be deprecated first, then removed

## Dependencies

- Mastery feature (backend API) — already implemented
- `mastery_experience` precomputed table — new, needed for efficient filtering

## Open Questions

- Should the exercise picker have a search/filter that spans the full compendium, or only load more on demand?
- Should "forked" exercises show their source exercise (linked via `forked` relationship type)?

## Resolved

- **Exercise versioning**: Workout logs reference a specific exercise version, so historical accuracy is preserved even if the compendium exercise is updated later. No snapshotting needed.
