# User Profiles

## Summary

Introduce a `user_profiles` table that maps the opaque user ID (UUIDv4 from `x-user-id` header) to a human-readable display name. This is a prerequisite for the compendium/user unification — once exercises are user-owned and publicly visible, the UI needs to show who created them.

## Motivation

Today, user identity is a bare string (`Owner` on user entities, `CreatedBy` on compendium entities) with no profile record behind it. There's no way to display a readable author name, and no referential integrity — any string is accepted as an owner.

## Design

### Table: `user_profiles`

| Column | Type   | Notes                                      |
|--------|--------|--------------------------------------------|
| `id`   | string | Primary key. The value from `x-user-id`.   |
| `name` | string | Human-readable display name, user-changeable. |

### Foreign Key Migration

All existing `Owner` and `CreatedBy` string fields become foreign keys referencing `user_profiles.id`:

**User domain** (`Owner` field):
- `user_exercises.owner`
- `user_equipment.owner`
- `workouts.owner`
- `exercise_logs.owner`
- `workout_logs.owner`

**Compendium domain** (`CreatedBy` field):
- `exercises.created_by`
- `equipment.created_by`
- `exercise_relationships.created_by`
- `fulfillments.created_by`
- `exercise_groups.created_by`

This gives referential integrity — no entity can reference a non-existent user.

### Seed Data

`make seed` creates a profile for the system user:

```json
{ "id": "claude", "name": "Claude" }
```

All seeded entities use `owner: "claude"` / `createdBy: "claude"` referencing this profile. No migration needed — the database is wiped and re-seeded with the new schema.

### Auto-Creation

On first request from an unknown `x-user-id`, the auth middleware (or a dedicated middleware) creates a profile with a default name (the UUID itself or "Anonymous"). The user can change their name later via `PUT /api/user/profile`.

### API

- `GET /api/user/profile` — returns the current user's profile.
- `PUT /api/user/profile` — updates the current user's display name.
- `GET /api/profiles/:id` — public endpoint, returns a profile by ID (for displaying author names on public entities).

## Open Questions

- **Name uniqueness**: should display names be unique? Probably not — multiple users named "Alex" is fine.
- **Profile deletion**: what happens if a user is deleted? Cascade-delete all their entities, or soft-delete the profile and orphan the data?
- **Additional fields**: avatar, bio, join date? Keep it minimal for now — just `id` and `name`.
- **Migration strategy**: not needed. No user data exists yet — wipe and re-seed.
