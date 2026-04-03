# User Profiles

## Summary

Introduce a `user_profiles` table that maps the opaque user ID (UUIDv4 from `x-user-id` header) to a human-readable display name. This is a prerequisite for the compendium/user unification — once exercises are user-owned and publicly visible, the UI needs to show who created them.

## Motivation

Today, user identity is a bare string (`Owner` on user entities, `CreatedBy` on compendium entities) with no profile record behind it. There's no way to display a readable author name, and no referential integrity — any string is accepted as an owner.

## Implementation

### Table: `user_profiles`

| Column | Type   | Notes                                      |
|--------|--------|--------------------------------------------|
| `id`   | string | Primary key. The value from `x-user-id`.   |
| `name` | string | Human-readable display name, user-changeable. |

Uses a string primary key (not `shared.BaseModel` which has `uint` ID). Has its own `created_at`, `updated_at`, `deleted_at` fields.

**Package**: `internal/profile/models/`

### Foreign Key Constraints

All `Owner`, `CreatedBy`, `AddedBy`, and `ChangedBy` string fields reference `user_profiles.id` via GORM association fields (`gorm:"foreignKey:...;references:ID;constraint:OnDelete:RESTRICT"`). The association fields are `json:"-"` and never preloaded — zero query cost.

SQLite foreign key enforcement enabled globally via `?_foreign_keys=on` on the connection string (`internal/database/database.go`).

**User domain** (`Owner` FK — 5 entities):
- `user_exercises.owner` → `OwnerProfile`
- `user_equipment.owner` → `OwnerProfile`
- `workouts.owner` → `OwnerProfile`
- `exercise_logs.owner` → `OwnerProfile`
- `workout_logs.owner` → `OwnerProfile`

**Compendium domain** (`CreatedBy` FK — 5 entities):
- `exercises.created_by` → `CreatedByProfile`
- `equipment.created_by` → `CreatedByProfile`
- `exercise_relationships.created_by` → `CreatedByProfile`
- `fulfillments.created_by` → `CreatedByProfile`
- `exercise_groups.created_by` → `CreatedByProfile`

**Special cases**:
- `exercise_group_members.added_by` → `AddedByProfile`
- `exercise_history.changed_by` → `ChangedByProfile`
- `equipment_history.changed_by` → `ChangedByProfile`

### Server-Side Owner Enforcement

All create handlers now set `Owner`/`CreatedBy`/`AddedBy` from `auth.GetUserID(c)` — the value from the request body is ignored. Update handlers preserve the existing `CreatedBy` (immutable after creation).

### Auto-Creation

`EnsureProfile()` middleware (`internal/profile/middleware.go`) runs after `auth.UserID()` on the `/api` group. Uses a `sync.Map` to cache known profile IDs — on cache miss, calls `db.FirstOrCreate` with `Name` defaulting to the user ID itself.

### Seed Data

`make seed` creates a `"claude"` profile as the first seed step. All seed JSON files use `"createdBy": "claude"` / `"addedBy": "claude"`. Equipment seed code also hardcodes `CreatedBy: "claude"`.

### API

- `GET /api/user/profile` — returns the current user's profile.
- `PUT /api/user/profile` — updates the current user's display name (body: `{ "name": "..." }`).
- `GET /api/profiles/:id` — public endpoint, returns a profile by ID.

### Frontend

- TypeScript types generated via tygo → `web/src/app/generated/profile.ts`
- API client methods added to `UserApiClient`: `fetchProfile()`, `updateProfile()`, `fetchPublicProfile()`
- Query keys added: `profileKeys.mine()`, `profileKeys.public(id)`
- No profile page UI yet — just the API client for future use.

## Resolved Questions

- **Name uniqueness**: not enforced — multiple users named "Alex" is fine.
- **Profile deletion**: FK constraint uses `OnDelete:RESTRICT` — cannot delete a profile that owns entities.
- **Additional fields**: minimal for now — just `id` and `name`. `created_at` and `updated_at` are automatic.
- **Migration strategy**: not needed. No user data exists yet — wipe and re-seed.
