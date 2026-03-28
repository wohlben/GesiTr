# Workout Commitments & Schedules

**Status:** Implemented (core functionality; period activation and UI polish ongoing)

## Summary

A **workout commitment** is a workout log in the future — a concrete promise to perform a workout within a time window. Commitments extend the existing `WorkoutLog` model with new statuses and fields. **Schedules** are rules attached to workouts that automatically generate commitments via a period-based system.

## Workout Log State Machine

Two independent flows exist:

### Manual flow (unchanged)

```
planning → in_progress → finished / partially_finished / aborted
(adhoc is a variant of in_progress that allows structural edits)
```

### Commitment flow

```
proposed → committed → in_progress → finished / partially_finished / aborted
proposed → skipped  (user declines, or due window elapses)
committed → broken  (due window elapses without starting)
```

Key distinctions:
- **`proposed`** — created by a schedule or group cycle. The workout structure is readonly; the user creates exercise schemes before committing.
- **`committed`** — user accepted the proposal. Everything is readonly. Waiting to start.
- **`skipped`** — user declined, or the due window elapsed without committing. Terminal. Set automatically by a background ticker.
- **`broken`** — a committed workout's due window passed without the user starting it. Terminal. Set automatically by a background ticker (15-min interval).

### Fields on WorkoutLog

| Field | Type | Purpose |
|---|---|---|
| `schedule_id` | `*uint` | Which schedule spawned this log (nullable) |
| `period_id` | `*uint` | Which period this log belongs to (nullable) |
| `due_start` | `*time.Time` | Start of the commitment window |
| `due_end` | `*time.Time` | End of the commitment window |

## Entity Hierarchy

```
Schedule → Period → Commitment → WorkoutLog
```

### Schedule

A schedule defines an active date range and the initial status for generated workout logs. It carries no type-specific configuration — that lives on the periods.

| Field | Type | Purpose |
|---|---|---|
| `owner` | `string` | Who owns this schedule |
| `workout_id` | `uint` | Which workout this applies to |
| `start_date` | `time.Time` | When the schedule becomes active (default: tomorrow) |
| `end_date` | `*time.Time` | When the schedule ends (nullable = indefinite) |
| `initial_status` | `string` | `committed` (default) or `proposed` |

Active is derived: `start_date ≤ today AND (end_date IS NULL OR end_date ≥ today)`.

### Period

A period is a concrete time window. The first period is created by the user; subsequent periods are cloned from the last one (template = last period).

| Field | Type | Purpose |
|---|---|---|
| `schedule_id` | `uint` | Parent schedule |
| `period_start` | `time.Time` | Start of the window |
| `period_end` | `time.Time` | End of the window |
| `type` | `string` | `fixed_date` or `frequency` |
| `mode` | `string` | `normal` (fixed duration in days) or `monthly` (calendar month) |

Types:
- **`fixed_date`** — user picks specific days within the period. Each commitment has a date.
- **`frequency`** — user picks a count (N per period). Commitments have no specific date.

Modes:
- **`normal`** — cloned period has the same duration in days.
- **`monthly`** — cloned period advances one calendar month.

### Commitment

A commitment is a join entity between a period and a workout log. It has a two-phase lifecycle:

1. **Plan phase:** created with `workout_log_id = null`. For fixed_date, `date` is set to the specific day.
2. **Activation phase:** when `period_start ≤ today`, a WorkoutLog is created and linked.

| Field | Type | Purpose |
|---|---|---|
| `period_id` | `uint` | Parent period |
| `date` | `*time.Time` | Target day (null for frequency type) |
| `workout_log_id` | `*uint` | Linked workout log (null until activated) |

## Two-Phase Generation

Generation is lazy — triggered when the user lists their workout logs.

**Phase 1 (Clone):** if the last period has ended and no next period exists, clone it forward with the same duration/mode, type, and commitment pattern.

**Phase 2 (Activate):** for each period where `period_start ≤ now`, find commitments with `workout_log_id = null`, create WorkoutLogs in the schedule's `initial_status`, and link them.

## Background Ticker

A Go goroutine runs every 15 minutes:
- `committed` → `broken` when `due_end < now()`
- `proposed` → `skipped` when `due_end < now()`

## API Endpoints

### Workout Logs (extended)
- `POST /user/workout-logs/{id}/skip` — proposed → skipped
- `POST /user/workout-logs/{id}/commit` — proposed → committed
- `GET /user/workout-logs?periodId=X` — filter by period

### Schedules
- `GET/POST /user/workout-schedules` — CRUD
- `GET/PATCH/DELETE /user/workout-schedules/{id}`

### Periods
- `GET/POST /user/schedule-periods`

### Commitments
- `GET/POST /user/schedule-commitments`
- `DELETE /user/schedule-commitments/{id}`

## Frontend

- **Workout list** — schedule icon button per workout row
- **Schedule list** — `/user/workouts/:id/schedules`
- **Schedule create/edit** — `/user/workouts/:id/schedules/new` and `/:scheduleId/edit`
  - Toggle groups for initial status and period type (via signal forms)
  - Spartan date pickers for period start/end with min constraints
  - Monthly mode detection when period spans exactly one calendar month
  - Custom period-day-picker component for fixed_date day selection
  - Number input for frequency count
  - Immediate activation warning when period starts today or earlier
- **Period detail** — `/user/workouts/:id/schedules/:scheduleId/periods/:periodId`
- **Calendar** — committed/proposed logs shown with distinct colors (purple/orange/gray)

## Remaining Work

- Period activation status (pending → active → completed)
- Save period changes from the edit view
- Due indicators on workout list
- Schedule deletion confirmation dialog
