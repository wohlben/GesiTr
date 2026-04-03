# Exercise Mastery

## Summary

Replaced the "import to My Exercises" flow (which created carbon copies of compendium exercises) with a **mastery system** that tracks how experienced a user is with each exercise based on their logged workouts. Users consume exercises directly from the compendium — no copying needed.

The "Fork" button remains as a secondary action for users who want to create a custom variant of an exercise.

## Mastery System Design

### XP Source

Each **rep** counts as 1 base XP toward that exercise's mastery. For non-rep-based exercises (time/distance), each exercise log entry = 1 base XP.

### Leveling Formula

**Flat XP per level**: 100 XP to advance each level (the recency multiplier provides the acceleration curve, not the level formula itself). Max level: 100.

### Recency Multiplier

Effective XP is scaled by training consistency:

```
effective_xp = base_xp * min(0.5 * n_days, max(1, current_level / 2))
```

Where:
- `n_days` = number of **distinct days** the user has performed this exercise within a 6-month window
- `current_level / 2` = multiplier cap (grows with level)
- `max(1, ...)` = floor of 1 so level 0 users can still earn XP

### Mastery Tiers

| Level Range | Tier        |
|-------------|-------------|
| 0 - 10      | Novice      |
| 11 - 30     | Journeyman  |
| 31 - 50     | Adept       |
| 51 - 99     | Master      |
| 100         | Mastered    |

### Relationship-Based XP Contributions

Logs from related exercises contribute to mastery with a multiplier derived from two components:

```
multiplier = (strength * 0.5) + type_bonus
```

| Category | Type Bonus | Relationship Types |
|----------|------------|-------------------|
| Equivalent | 0.5 | `equivalent` |
| Skill transfer | 0.25 | `alternative`, `easier_alternative`, `harder_alternative`, `equipment_variation`, `variant`, `variation`, `bilateral_unilateral`, `progression`, `progresses_to`, `regression`, `regresses_to` |
| No transfer | — | `accessory`, `antagonist`, `complementary`, `preparation`, `prerequisite`, `related`, `similar`, `superset_with`, `supports`, `forked` |

An exercise's own logs always count at **1.0**. If multiple relationships exist between the same pair, the highest multiplier wins.

## Implementation

### Backend

- **Models**: `MasteryContributionEntity` (precomputed relationship-based XP lookup), `MasteryExperienceEntity` (total reps per owner+exercise), `ExerciseMastery` DTO
- **Computation** (`internal/user/mastery/models/compute.go`): `ComputeLevel`, `ComputeTier`, `ComputeRecencyMultiplier`, `ComputeProgress`, `ComputeContributionMultiplier`
- **Endpoints**: `GET /api/user/mastery` (list all), `GET /api/user/mastery/:exerciseId` (single)
- **Lifecycle hooks**: `UpsertExperience` (on exercise log creation), `RecalculateContributions` (on relationship create/delete)
- **Migration backfills**: Contributions from existing relationships, experience from existing exercise_logs
- **Exercise list filter**: `GET /api/exercises?mastery=me` returns exercises the user owns or has mastery in
- **Index**: Composite index on `exercise_logs(owner, exercise_id, performed_at)` for efficient queries

### Frontend

- **Exercise detail page**: Shows mastery card (level, tier, XP, progress bar) when the user has logged the exercise
- **"Track Exercise" button** (primary) on exercise detail — navigates to `/compendium/exercises/:id/track` for quick workout creation
- **"Fork Exercise" button** (secondary) — creates a copy for customization, navigates to the forked exercise in the compendium
- **Exercise list**: Mastery column with tier-colored level badges (novice=gray, journeyman=green, adept=blue, master=purple, mastered=yellow)
- **ExerciseConfig**: Sorts exercises with mastery first in the exercise picker

### Removed

- `/user/exercises/:id` and `/user/exercises/:id/track` routes — exercises are accessed directly from compendium
- `UserExerciseDetail` component — no longer needed
- `userExerciseKeys` query keys
- `fetchUserExercises()`, `deleteUserExercise()` API methods

## Resolved Design Decisions

- **Fork**: Kept as secondary action — allows creating a copy as a blueprint for a custom variant
- **"My Exercises" navigation**: Removed as redundant — users access exercises through the compendium
- **Tier naming**: Novice / Journeyman / Adept / Master / Mastered
- **Recency window**: Fixed at 6 months
- **Non-rep XP**: Flat 1 XP per log entry for non-rep exercises
