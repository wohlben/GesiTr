# Workout Section Exercise Groups

**Status:** Implemented

## Summary

Workout sections now support two item types via a union discriminator: **exercise** (existing behavior — references an exercise scheme) and **exercise_group** (references an exercise group whose specific exercise is chosen each session). Exercise groups were simplified from standalone compendium entities to lightweight workout-scoped building blocks: optionally named collections of exercises managed inline on the workout edit page.

## Motivation

Many training programs prescribe a slot like "pick a horizontal pull" or "any quad-dominant movement" rather than a fixed exercise. Previously users had to pick one exercise and stick with it, or duplicate sections for each variant. Exercise groups let the workout template express this flexibility directly.

## Implementation

### Data Model Changes

**`WorkoutSectionExercise` → `WorkoutSectionItem`** (table: `workout_section_items`):
- Added `type` discriminator: `exercise` | `exercise_group`
- `exerciseSchemeId` made nullable (only set for `exercise` type)
- Added `exerciseGroupId` (only set for `exercise_group` type)
- Added `data` JSON column for future extensibility (not used yet)
- Migration copies data from old `workout_section_exercises` table with `type='exercise'`

**`ExerciseGroup` simplified**:
- Removed `description` field
- Made `name` optional (groups are unnamed by default)
- Removed from SPA top navigation (managed inline from workout edit page)

**`WorkoutLogExercise`**:
- Added `sourceExerciseGroupId` to track which group the exercise was chosen from

### API Changes

- `POST/GET/DELETE /api/user/workout-section-exercises` → `/api/user/workout-section-items`
- `CreateWorkoutSectionItem` validates per type: exercise requires `exerciseSchemeId`, exercise_group requires `exerciseGroupId`
- Added `GET/POST/DELETE /api/exercise-group-members` to frontend API client

### Workout Edit Page

- Section items have a type selector (Exercise vs Exercise Group)
- **Exercise type**: existing flow (exercise picker + scheme configuration)
- **Exercise Group type**: inline `ExerciseGroupConfig` component with:
  - Group picker dropdown (existing groups or "New Group")
  - Group name input (optional)
  - Exercise member picker (add exercises to the group)
  - Member list with remove buttons
- Save flow creates new groups + members, or syncs existing groups (updates name, adds/removes members)
- Edit mode loads existing group name + members from API

### Workout Start Page

- Exercise group items appear as pending group cards with an exercise picker
- User picks an exercise from the group → opens AddExerciseDialog with exercise pre-selected
- After configuring the scheme, the log exercise is created and the pending group is resolved
- Start Workout button is disabled until all groups are resolved
- Once started, group exercises appear as regular exercises in the workout log

### E2E Testing Infrastructure

- Switched e2e tests to in-memory SQLite (`DATABASE_PATH=:memory:`)
- Added per-test DB reset via `base-test.ts` fixture (calls `POST /api/ci/reset-db` before each test)
- Fixed profile cache invalidation on DB reset (`ResetProfileCache()`)

## Open Questions (Deferred)

- **Progress tracking**: when alternating exercises from a group, logs are recorded per-exercise (separately). Grouped progress views could be a follow-up.
- **Last-used suggestion**: the UI doesn't yet suggest which exercise was picked last time.
- **Converting between types**: changing a section item from exercise to exercise_group (or vice versa) is supported by deleting and re-adding the item.
