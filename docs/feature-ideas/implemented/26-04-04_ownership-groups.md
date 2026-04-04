# Ownership Groups

**Status:** Implemented

## Summary

Replaced the single `owner` string field on compendium entities with a group-based ownership model. Each entity gets its own **ownership group** — a lightweight, unnamed entity that exists solely to hold memberships. The group is auto-created with the creating user as its sole `owner`-role member, preserving current behavior by default while enabling fine-grained sharing of private entities with select users.

**Depends on:** [Resource Permissions](./26-03-24_resource-permissions.md), [Unify Compendium and User Data](./26-03-24_unify-compendium-and-user-data.md)

## Model

### Ownership Group

A minimal referential entity — no name, no description. Just an ID.

| Field | Type | Notes |
|-------|------|-------|
| ID | uint | Primary key |
| CreatedAt | timestamp | |

### Ownership Group Membership

Join table linking users to ownership groups.

| Field | Type | Notes |
|-------|------|-------|
| ID | uint | Primary key |
| GroupID | uint | FK → OwnershipGroup |
| UserID | string | FK → user |
| Role | enum | `owner`, `admin`, `member` |
| CreatedAt | timestamp | |

Unique constraint on `(GroupID, UserID)`.

### Roles

| Role | Semantics |
|------|-----------|
| **owner** | Full control — modify, delete, manage members |
| **admin** | Reserved for future use — no additional privileges beyond member |
| **member** | Read access to private entities |

## Affected Entities

Top-level compendium entities that get their own ownership group on creation:

- Exercise
- Equipment
- Workout
- ExerciseGroup
- WorkoutGroup
- Locality

Sub-entities that inherit their parent's `ownership_group_id`:

- ExerciseRelationship (from Exercise)
- EquipmentRelationship (from Equipment)
- EquipmentFulfillment (from Equipment)
- ExerciseGroupMember (from ExerciseGroup)
- WorkoutRelationship (from Workout)
- LocalityAvailability (from Locality)

**Not affected:** User-scoped entities (WorkoutLog, ExerciseLog, WorkoutSchedule, ExerciseScheme, mastery entities) keep their plain `owner` field.

## Implementation

### Backend

- **Package:** `internal/ownershipgroup/` with `models/`, `handlers/`, `access.go`, `create.go`, `migrate.go`
- **Access resolution:** `CheckAccess(db, userID, ownershipGroupID)` returns `EntityAccess` with `CanRead()`, `CanModify()`, `CanDelete()` methods
- **List filtering:** `VisibleGroupIDs(db, userID)` subquery replaces `WHERE owner = ?`
- **Permission resolution:** `shared.ResolvePermissionsFromAccess(access, isPublic)` replaces the old `ResolvePermissions`
- **Membership CRUD:** `/api/ownership-groups/{id}/memberships` (GET, POST) and `/api/ownership-group-memberships/{id}` (PUT, DELETE)
- **Direct add:** Members are added directly with `member` role (no invitation flow)

### Frontend

- **Generated types:** `ownershipgroup.ts` with `OwnershipGroup`, `OwnershipGroupMembership`
- **Reusable component:** `<app-ownership-group-panel>` shows members and add/remove UI
- **Integrated into:** Exercise, Equipment, and Locality detail pages (shown when user has modify permission)
- **Owner column removed** from all list pages and detail pages

### Migration

- `MigrateExistingOwners(db)` runs on startup, creates ownership groups for existing entities
- Idempotent — skips rows that already have `ownership_group_id` set
- Drops the `owner` column from all affected tables after migration
- Old unique indexes containing `owner` are dropped and recreated with `ownership_group_id`

### Mastery recalculation

- `RecalculateContributions` and `RecalculateEquipmentContributions` query relationships via ownership group membership instead of the former `owner` column

## Design Decisions

- **Sub-entities inherit parent's group** — no separate group per sub-entity. Sharing an exercise automatically shares its relationships.
- **ExerciseScheme keeps plain `owner`** — user-scoped, not affected.
- **Direct add, no invitation flow** — simpler implementation. Can add invitation flow later.
- **One group per entity** — each entity instance has its own ownership group. No group sharing/linking across entities.
- **Workout dual-access** — workouts retain the existing WorkoutGroup sharing mechanism as a secondary access path alongside ownership groups.
- **`owner` column fully removed** — the field is gone from all affected entity structs, DTOs, and database tables. The `owner` query parameter on list endpoints now resolves through ownership group memberships.

## Resolved Considerations

- **Public flag interaction:** Public entities are readable by everyone regardless of ownership group. Adding someone to a public entity's group grants no additional privileges currently (member role = read-only, which public already provides).
- **List view shared entities:** Entities shared with a user (where they're a member) appear in their default list view alongside their own entities.
- **UI member management:** Per-entity via the ownership group panel on detail pages. Bulk management across entities is out of scope.
- **Entity deletion:** Cascade deletes the ownership group and all its memberships (via GORM foreign key constraint).
