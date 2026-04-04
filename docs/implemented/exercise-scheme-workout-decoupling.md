# Exercise Scheme / Workout Decoupling

## Status: Implemented

## Summary

Resolved the dissonance between public workouts and user-owned exercise schemes. Workouts now only reference exercises (public entities), not exercise schemes (user entities). A new join table `exercise_scheme_section_items` connects user exercise schemes to workout section items, and the backend joins this data for scheme-dependent flows.

## Changes Made

### WorkoutSectionItem: replaced ExerciseSchemeID with ExerciseID

`WorkoutSectionItemEntity` now has an `ExerciseID` field (replaces `ExerciseSchemeID`). The workout only references the exercise to perform at that slot, with no dependency on user data.

### ExerciseScheme: removed WorkoutSectionItemID

`ExerciseSchemeEntity` in `internal/user/exercisescheme/` keeps `ExerciseID` and all configuration fields but dropped `WorkoutSectionItemID`. It is a pure user-owned entity: "my personal config for this exercise."

### New join table: exercise_scheme_section_items

A join entity linking `ExerciseSchemeID` to `WorkoutSectionItemID`, with a unique constraint on `(WorkoutSectionItemID, Owner)` — one scheme per user per item slot.

```
ExerciseSchemeSectionItemEntity:
├── BaseModel
├── ExerciseSchemeID (uint, not null, indexed)
├── WorkoutSectionItemID (uint, not null, uniqueIndex with Owner)
├── Owner (string, not null, uniqueIndex with WorkoutSectionItemID)
```

### New API endpoints

- `GET /api/user/exercise-scheme-section-items?workoutSectionItemIds=1,2,3` — batch lookup
- `PUT /api/user/exercise-scheme-section-items` — upsert scheme assignment
- `DELETE /api/user/exercise-scheme-section-items/{id}` — remove assignment

### Backend changes

- `CreateWorkoutSectionItem` validates `ExerciseID` instead of `ExerciseSchemeID`
- `AcceptWorkoutGroupInvitation` checks the join table for scheme assignments
- Workout schedule generation (`generate.go`) looks up schemes via the join table
- Seed data creates section items with `ExerciseID` and join-table entries

### Frontend changes

- Workout editor creates section items with `exerciseId`, schemes separately, then links via join table
- Workout start page fetches scheme assignments via join table for planning log creation
- Workout start store loads exercise display using join-table assignments
- API client has new methods: `fetchSchemeSectionItems`, `upsertSchemeSectionItem`, `deleteSchemeSectionItem`

## Notes

- Exercise groups (`type = "exercise_group"`) are unaffected
- `WorkoutLogExercise.SourceExerciseSchemeID` remains as a historical reference
- GORM AutoMigrate does not drop old columns — `exercise_scheme_id` on `workout_section_items` and `workout_section_item_id` on `exercise_schemes` remain in the DB but are unused
