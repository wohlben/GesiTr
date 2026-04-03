# Ownership Groups

## Summary

Replace the single `owner` string field on compendium entities with a group-based ownership model. Each entity gets its own **ownership group** — a lightweight, unnamed entity that exists solely to hold memberships. The group is auto-created with the creating user as its sole `owner`-role member, preserving current behavior by default while enabling fine-grained sharing of private entities with select users.

**Depends on:** [Resource Permissions](./implemented/26-03-24_resource-permissions.md), [Unify Compendium and User Data](./implemented/26-03-24_unify-compendium-and-user-data.md)

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

Follows the same pattern as workout group roles:

| Role | Semantics |
|------|-----------|
| **owner** | Full control — maps to current owner privileges (modify, delete, manage members) |
| **admin** | Reserved for future use — initially no additional privileges beyond member |
| **member** | Reserved for future use — initially no additional privileges (read access to private entity) |

Only `owner` maps to current owner privileges initially. `admin` and `member` exist in the schema to avoid a future migration when their semantics are defined.

## Affected Entities

Compendium entities that replace `owner` with `ownership_group_id`:

- Exercise
- ExerciseScheme
- ExerciseRelationship
- Equipment
- EquipmentRelationship
- EquipmentFulfillment
- Workout
- ExerciseGroup / ExerciseGroupMember
- WorkoutGroup
- Locality

**Not affected:** User-scoped entities (WorkoutLog, ExerciseLog, WorkoutSchedule, mastery entities) keep their plain `owner` field.

## Behavior

### Creation
When a user creates any compendium entity:
1. A new ownership group is auto-created
2. A membership is created linking the user to that group with role `owner`
3. The entity stores the `ownership_group_id` instead of an `owner` string

### One group per entity
Each entity instance gets its own ownership group. Exercise A and Exercise B always have separate groups, even if the user wants identical members. Managing members across entities is manual — group sharing/linking is out of scope.

### Permission resolution
The existing resource permissions system (`/permissions` endpoints) resolves access by looking up the ownership group membership instead of comparing the `owner` string:
- `owner` role → `READ`, `MODIFY`, `DELETE`
- `admin` / `member` role → TBD (initially just `READ` for private entities)
- No membership + public entity → `READ`
- No membership + private entity → 404

### Sub-entity ownership
Sub-entities (e.g., ExerciseScheme belongs to Exercise) could either:
- Have their own ownership group (current: each has its own `owner` field)
- Inherit from the parent entity's ownership group

This decision affects whether sub-entities can be independently shared.

## Migration

1. Create `ownership_groups` and `ownership_group_memberships` tables
2. For each existing compendium entity row:
   - Create an ownership group
   - Create a membership with the current `owner` value as `user_id` and role `owner`
   - Set `ownership_group_id` on the entity
3. Drop `owner` column from affected entities

Seed data (owned by system user `Sinon`) follows the same migration — each seeded entity gets its own ownership group with `Sinon` as owner.

## Considerations

- How does this interact with the `public` flag? Public entities are readable by everyone regardless of ownership group — does adding someone to the group of a public entity grant any additional privileges?
- Should the compendium list view surface entities shared with the user (where they're a member but not owner) alongside their own?
- How are ownership group members managed in the UI? Per-entity member management could be tedious if a user has many exercises they want to share with the same person.
- What happens to ownership group memberships when an entity is deleted? Cascade delete the group and its memberships.
- Should there be an invitation flow (like workout groups) or direct add?
