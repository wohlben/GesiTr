# Resource Permissions

**Status:** Implemented

## Summary

Per-resource permission endpoints return the current user's allowed actions as an enum list (`READ`, `MODIFY`, `DELETE`). Route-level components call this endpoint and conditionally render action buttons based on the response. Edit/delete buttons only appear after the permissions query successfully resolves with the appropriate permission.

As a cleanup, the `slug` field was removed from the backend Exercise model — the frontend generates slugs client-side via `SlugifyPipe`. A startup migration drops the `idx_owner_slug` index and `slug` column from existing databases.

## Motivation

With the unified data model, the compendium is a view over all public exercises/equipment. Any user can *see* another user's public entities, but only the owner should be allowed to modify or delete them. The frontend previously rendered edit and delete buttons unconditionally — the backend rejected unauthorized mutations, but the UI shouldn't show options the user can't perform.

A dedicated permission endpoint is cleaner than embedding ownership logic in the frontend (comparing `owner` to current user ID), because:

- The backend is the single source of truth for authorization rules.
- Future permission models (e.g., group-level edit access, admin overrides) are handled in one place — the frontend just reads the list.
- No user identity needs to be exposed to or reasoned about in the frontend.

## Requirements

- **[Unified data model](unify-compendium-and-user-data.md)** — entities must have an `owner` field so the backend can determine permissions.

## Implementation

### Permission Endpoint

Each resource type has a `GET /:id/permissions` sub-route:

```
GET /api/exercises/:id/permissions        -> ["READ"]              (public, not owner)
GET /api/exercises/:id/permissions        -> ["READ", "MODIFY", "DELETE"]  (owner)
GET /api/equipment/:id/permissions        -> ["READ", "MODIFY", "DELETE"]  (owner)
GET /api/exercise-groups/:id/permissions  -> ["READ", "MODIFY", "DELETE"]  (owner)
```

Response body:

```json
{
  "permissions": ["READ", "MODIFY", "DELETE"]
}
```

**Permission enum values:**

| Value    | Meaning                        |
|----------|--------------------------------|
| `READ`   | Can view the resource          |
| `MODIFY` | Can edit the resource          |
| `DELETE` | Can delete the resource        |

**Resolution rules** (owner-based):

- If the requesting user is the entity's `owner` → `READ`, `MODIFY`, `DELETE`
- If the entity is `public = true` and the user is not the owner → `READ`
- If the entity is `public = false` and the user is not the owner → 404 (entity not visible)

These rules live in a shared helper (`internal/shared/permissions.go`) so all resource types use the same logic. Future extensions (group editors, admin roles) modify the helper — individual handlers don't change.

### Slug Removal

1. Removed `Slug` field from `ExerciseEntity` and the `Exercise` DTO.
2. Removed the `idx_owner_slug` unique index from the `Owner` field.
3. Removed slug auto-generation logic in the `CreateExercise` handler.
4. Deleted the `internal/slug/` package.
5. Added a startup migration in `main.go` (`runMigrations()`) that drops the index and column from existing databases.
6. Frontend keeps `:slug` in route URLs — generated client-side from the entity name via the existing `SlugifyPipe`.

### Frontend Consumption

Detail pages (exercise, equipment, exercise-group) fetch permissions alongside entity data. Buttons only render after the permissions query has **successfully** resolved (prevents flash from stale cache):

```typescript
permissionsQuery = injectQuery(() => ({
  queryKey: exerciseKeys.permissions(this.id()),
  queryFn: () => this.api.fetchExercisePermissions(this.id()),
  enabled: !!this.id(),
}));

canModify = computed(() =>
  this.permissionsQuery.isSuccess() &&
  (this.permissionsQuery.data()?.permissions?.includes('MODIFY') ?? false),
);
```

Edit pages check permissions and redirect back to the detail view if the user lacks `MODIFY`.

## Considerations

- **Caching**: permission responses are cached with the same staleness as the entity itself. Invalidate on mutation (if ownership transfer becomes a thing later).
- **Batch endpoint**: for list views, fetching permissions per-item is N+1. A future `POST /api/exercises/permissions` batch endpoint could accept an array of IDs and return a map. Not needed — list views don't show per-item action buttons today.
