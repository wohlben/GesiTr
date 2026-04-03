# Compendium Workout Groups

**Status:** Implemented

## Summary

Users can form groups around shared workouts and track progress together. A workout owner creates a group, invites members by user ID, and members accept after setting up their exercise schemes. Members see the shared workout in their workout list with a readonly start flow — they can customize set targets (reps, weight) but cannot modify the workout structure. The backend enforces access control via a permissions endpoint, and the frontend derives readonly mode from the permission response.

## Motivation

Resistance training is often done in pairs or small groups following the same program. Without workout groups, each person would need to manually recreate the same workout. Groups let a single workout definition be shared, while each member personalizes their own set targets and tracks their own progress independently.

## Roles

Three membership roles, plus the implicit owner:

| Role | Access | Can modify structure | Can track sets |
|------|--------|---------------------|---------------|
| **Owner** | Full | Yes | Yes |
| **Admin** | Read + Modify | Yes | Yes |
| **Member** | Read only | No | Yes |
| **Invited** | Read only | No | No (must accept first) |

- New memberships always start as `invited`
- Members must create exercise schemes for all workout items before accepting
- Acceptance promotes `invited` → `member`
- Owner can promote members to `admin` or demote back

## Implementation

### Backend

**Models** (`internal/user/workoutgroup/`):
- `WorkoutGroupEntity` — name, workoutID (unique), owner
- `WorkoutGroupMembershipEntity` — groupID, userID, role

**Access control** (`internal/user/workoutgroup/access.go`):
- `CheckWorkoutAccess(userID, owner, workoutID)` → `WorkoutAccess{IsOwner, IsMember, IsAdmin}`
- `CanRead()` — owner, member, or admin
- `CanModify()` — owner or admin
- `CanDelete()` — owner only
- `GroupInfoForWorkouts()` — batch lookup for list views

**Endpoints** (10 total):

| Method | Path | Who | Purpose |
|--------|------|-----|---------|
| GET | `/user/workout-groups` | Owner | List owned groups |
| POST | `/user/workout-groups` | Owner | Create group for a workout |
| GET | `/user/workout-groups/{id}` | Owner/Member | Get group |
| PUT | `/user/workout-groups/{id}` | Owner | Update group name |
| DELETE | `/user/workout-groups/{id}` | Owner | Delete group (cascades memberships) |
| GET | `/user/workout-group-memberships` | Owner | List members |
| POST | `/user/workout-group-memberships` | Owner | Invite user (starts as `invited`) |
| PUT | `/user/workout-group-memberships/{id}` | Owner | Change role |
| DELETE | `/user/workout-group-memberships/{id}` | Owner | Remove member |
| POST | `/user/workouts/{id}/group/accept` | Invitee | Accept invite (validates schemes exist) |

**Permissions endpoint** (`GET /user/workouts/{id}/permissions`):
- Extends the existing resource permissions pattern to workouts
- Owner → `[READ, MODIFY, DELETE]`, Admin → `[READ, MODIFY]`, Member/Invited → `[READ]`

**Workout API integration**:
- `ListWorkouts` includes workouts where user is a group member (SQL join)
- `GetWorkout` / `UpdateWorkout` responses include `workoutGroup: {groupName, membership}` for non-owners
- All workout CRUD handlers enforce `CanRead` / `CanModify` / `CanDelete`

### Frontend

**Workout group management page** (`pages/user/workout-group/`):
- Create/edit/delete group (owner only)
- Add/remove members, change roles
- Accessible via group icon on workout list

**Workout list** (`pages/user/workout-list/`):
- Shared workouts appear alongside owned workouts
- Color-coded membership badge: invited (yellow), member (green), admin (purple)

**Workout-start readonly mode** (`pages/user/workout-start/`):
- `isReadonly` computed from permissions query (`!permissions.includes('MODIFY')`)
- Extracted into child components: `WorkoutStartSection`, `WorkoutStartExerciseItem`
- Readonly hides: drag handles, remove buttons, add section/exercise buttons, type/label editing
- Readonly preserves: set target inputs (reps, weight, duration, distance, rest)
- Pending exercise groups remain functional in readonly mode

**API client & query keys**:
- `fetchWorkoutPermissions(id)` → `GET /user/workouts/{id}/permissions`
- `workoutKeys.permissions(id)` for TanStack Query caching
- Full CRUD methods for groups and memberships

### Tests

**Backend**:
- 5 example tests for permissions endpoint (owner, member, admin, invited, no-access)
- 4 example tests for workoutGroup info in workout responses
- 2 example tests for accept invitation (success, denied)
- 1 integration test for full group member start flow (invite → accept → create log → log sets)

**Frontend component tests** (35 tests across 3 files):
- `workout-start-exercise-item.spec.ts` — drag handle/remove hidden in readonly, set inputs always present
- `workout-start-section.spec.ts` — structural controls hidden, pending groups functional, readonly propagated
- `workout-start.spec.ts` — permissions-driven readonly, add-section hidden, sets editable, name/notes editable

**E2E**:
- `workout-group-member-flow.spec.ts` — full flow: API setup (create workout + group + invite) → UI as bob (workout list → readonly start → start workout → complete sets → auto-finish)

## Open Questions (Deferred)

These were listed in the original spec and intentionally deferred:

- **Link-based invites** — currently username-based only
- **Self-leave** — members cannot leave; owner must remove them
- **Group activity feed** — members cannot see each other's workout logs
- **Privacy controls** — no per-log visibility toggle
- **Workout scheduling** — no schedule/deadline fields on groups

## Constraints

- **One group per workout** — unique constraint on `workout_groups.workout_id`
- **Owner not in membership table** — implicit access via ownership check
- **Cascade delete** — deleting a group removes all memberships
