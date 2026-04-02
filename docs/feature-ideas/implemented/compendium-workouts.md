# Compendium Workouts

## Summary

Converted user-scoped workouts into compendium workouts — the same pattern used for exercises and equipment. Workouts now have an `owner` field, `public` toggle (default false), version tracking, and a fork mechanic via WorkoutRelationship. "My Workouts" is removed as a separate concept.

## What Changed

### Backend

- **Model**: `WorkoutEntity` gained `Public bool` (default false, indexed) and `Version int` (default 0)
- **Version tracking**: `WorkoutHistoryEntity` stores JSON snapshots on every update, with `WorkoutChanged()` no-op detection
- **Fork mechanic**: `WorkoutRelationship` entity (types: `forked`, `equivalent`) at `internal/workoutrelationship/`. Creating a workout with `SourceWorkoutID` auto-creates forked+equivalent relationships
- **Routes moved**: `/api/user/workouts` → `/api/workouts`, `/api/user/workout-sections` → `/api/workout-sections`, `/api/user/workout-section-items` → `/api/workout-section-items`
- **Listing**: `GET /api/workouts` supports `owner`, `public`, `logged`, `q` filters with pagination. Default shows own + public + group-shared workouts. `logged=me` shows workouts the user has workout logs for
- **Access control**: Public workouts are readable by anyone. Group membership still grants access to private workouts. Owner/admin can modify. Only owner can delete
- **Workout schedules**: Any user who can see a workout can create their own schedule (owner check removed)
- **Workout groups**: Multiple groups per workout now allowed (unique index on `workout_id` dropped)
- **Package location**: Moved from `internal/user/workout/` to `internal/workout/` (same for workoutlog, workoutgroup, workoutschedule)

### Frontend

- **Routes**: Moved from `/user/workouts/...` to `/compendium/workouts/...`
- **Navigation**: "My Workouts" removed from user nav, "Workouts" added to compendium nav
- **API client**: All workout endpoints updated to `/api/workouts`, response parsing handles paginated body
- **Workout list page**: Title changed from "My Workouts" to "Workouts"

### New Endpoints

- `GET /api/workouts/{id}/versions` — list version history
- `GET /api/workouts/{id}/versions/{version}` — get specific version
- `GET /api/workout-relationships` — list workout relationships
- `POST /api/workout-relationships` — create relationship
- `DELETE /api/workout-relationships/{id}` — delete relationship

## Resolved Design Decisions

- **Versioning**: Added (version field + WorkoutHistoryEntity with JSON snapshots)
- **Fork tracking**: WorkoutRelationship entity (same pattern as ExerciseRelationship)
- **Route prefix**: `/compendium/workouts` on frontend, `/api/workouts` on backend
- **Public default**: false — existing workouts remain private, users can explicitly publish
- **Who can create**: Any user (owner-based, same as exercises/equipment)
- **Schedule access**: Any user who can read a workout can create their own schedule
- **Multiple groups**: Allowed per workout (unique constraint dropped)
