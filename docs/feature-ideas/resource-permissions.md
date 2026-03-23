# Resource Permissions

## Summary

Add a per-resource permission endpoint that returns the current user's allowed actions as an enum list (`READ`, `MODIFY`, `DELETE`). Route-level components call this endpoint and conditionally render action buttons based on the response. This replaces the current unconditional display of edit/delete buttons on all detail pages.

As a cleanup, remove the `slug` field from the backend models and API responses — the frontend keeps `:slug` in its routes for URL readability, but the backend no longer stores or returns it.

## Motivation

With the unified data model, the compendium is a view over all public exercises/equipment. Any user can *see* another user's public entities, but only the owner should be allowed to modify or delete them. Today the frontend renders edit and delete buttons unconditionally — the backend rejects unauthorized mutations, but the UI shouldn't show options the user can't perform.

A dedicated permission endpoint is cleaner than embedding ownership logic in the frontend (comparing `owner` to current user ID), because:

- The backend is the single source of truth for authorization rules.
- Future permission models (e.g., group-level edit access, admin overrides) are handled in one place — the frontend just reads the list.
- No user identity needs to be exposed to or reasoned about in the frontend.

### Slug removal rationale

The `slug` field on exercises exists in the backend model (`ExerciseEntity.Slug`, unique index `idx_owner_slug`) and is returned in API responses. It's never used for backend lookups — all routes use `:id`. The frontend includes `:slug` in its URL patterns (`/compendium/exercises/:id/:slug`) for readability, but extracts only `:id` from route params. The slug can be generated client-side from the name for URL construction — no backend storage or transport needed.

## Requirements

- **[Unified data model](implemented/unify-compendium-and-user-data.md)** — entities must have an `owner` field so the backend can determine permissions.

## Design Sketch

### Permission Endpoint

Each resource type gets a `GET /:id/permissions` sub-route:

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

**Resolution rules** (initial, owner-based):

- If the requesting user is the entity's `owner` -> `READ`, `MODIFY`, `DELETE`
- If the entity is `public = true` and the user is not the owner -> `READ`
- If the entity is `public = false` and the user is not the owner -> 404 (entity not visible)

These rules live in a shared helper so all resource types use the same logic. Future extensions (group editors, admin roles) modify the helper — individual handlers don't change.

### Backend Route Changes

Register the new routes alongside existing CRUD routes in `main.go`:

```go
exercises.GET("/:id/permissions", exerciseHandlers.GetPermissions)
equipment.GET("/:id/permissions", equipmentHandlers.GetPermissions)
exerciseGroups.GET("/:id/permissions", exerciseGroupHandlers.GetPermissions)
```

### Slug Removal (Backend Only)

1. Remove `Slug` field from `ExerciseEntity` (and the `idx_owner_slug` unique index).
2. Remove `Slug` field from the `Exercise` DTO.
3. Remove the slug auto-generation logic in `CreateExercise` handler.
4. Remove the `internal/slug/` package.
5. Update seed data pipeline if it references slug.
6. Frontend keeps `:slug` in route URLs — generate it client-side from the entity name for link construction (a simple `slugify` utility in the frontend).

### Frontend Consumption

Route-level components (exercise-detail, equipment-detail, etc.) fetch permissions alongside the entity data:

```typescript
permissionsQuery = injectQuery(() => ({
  queryKey: exerciseKeys.permissions(this.id()),
  queryFn: () => this.api.fetchExercisePermissions(this.id()),
  enabled: !!this.id(),
}));

canModify = computed(() => this.permissionsQuery.data()?.includes('MODIFY') ?? false);
canDelete = computed(() => this.permissionsQuery.data()?.includes('DELETE') ?? false);
```

Action buttons are gated:

```html
@if (canModify()) {
  <a routerLink="./edit" ...>Edit</a>
}
@if (canDelete()) {
  <button (click)="showDeleteDialog.set(true)" ...>Delete</button>
}
```

The permissions query runs in parallel with the entity query — no waterfall.

## Considerations

- **Caching**: permission responses can be cached with the same staleness as the entity itself. Invalidate on mutation (if ownership transfer becomes a thing later).
- **Batch endpoint**: for list views, fetching permissions per-item is N+1. A future `POST /api/exercises/permissions` batch endpoint could accept an array of IDs and return a map. Not needed initially — list views don't show per-item action buttons today.
- **Edit page guard**: the edit route should also check permissions and redirect if the user navigates directly to `/exercises/:id/:slug/edit` without `MODIFY`. This can be a route guard or handled in the edit component itself.
