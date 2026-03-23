# Frontend Migration: Unified Data Model

The backend has unified compendium and user data into a single model. This document summarizes the frontend changes needed to complete the migration.

## API Endpoint Changes

| Old Endpoint | New Endpoint | Notes |
|-------------|-------------|-------|
| `GET /api/user/exercises` | `GET /api/exercises?owner=me` | Returns user's own exercises |
| `POST /api/user/exercises` | `POST /api/exercises` | Creates a full exercise (not a thin wrapper) |
| `GET /api/user/exercises/:id` | `GET /api/exercises/:id` | Same endpoint, richer response |
| `DELETE /api/user/exercises/:id` | `DELETE /api/exercises/:id` | Ownership enforced server-side |
| `GET /api/user/equipment` | `GET /api/equipment?owner=me` | Same pattern as exercises |
| `POST /api/user/equipment` | `POST /api/equipment` | Full equipment entity |
| `GET /api/user/equipment/:id` | `GET /api/equipment/:id` | |
| `DELETE /api/user/equipment/:id` | `DELETE /api/equipment/:id` | |
| `* /api/user/exercise-schemes` | `* /api/exercise-schemes` | Moved out of `/user` prefix |
| `GET /api/exercises` | `GET /api/exercises` | Now returns `owner`+`public` fields; default shows own + public |
| `GET /api/equipment` | `GET /api/equipment` | Same |
| — | `* /api/equipment-relationships` | New endpoint |

## TypeScript Type Changes

### Removed Types
- `UserExercise` — replaced by `Exercise` (temporary aliases exist in `user-models.ts`)
- `UserEquipment` — replaced by `Equipment`
- `UserExerciseScheme` — renamed to `ExerciseScheme` (alias exists)

### Renamed Fields
| Old Field | New Field | Types Affected |
|-----------|-----------|---------------|
| `createdBy` | `owner` | Exercise, Equipment, ExerciseGroup |
| `addedBy` | `owner` | ExerciseGroupMember |
| `userExerciseId` | `exerciseId` | ExerciseScheme, ExerciseLog |
| `userExerciseSchemeId` | `exerciseSchemeId` | WorkoutSectionExercise |
| `compendiumExerciseId` | — (use `templateId` on the exercise itself) | Was on UserExercise |
| `compendiumEquipmentId` | — (use `templateId` on the equipment itself) | Was on UserEquipment |
| `compendiumVersion` | — (use `version` on the exercise/equipment itself) | Was on UserExercise/UserEquipment |
| `fromExerciseTemplateId` | `fromExerciseId` (uint) | ExerciseRelationship |
| `toExerciseTemplateId` | `toExerciseId` (uint) | ExerciseRelationship |
| `equipmentTemplateId` | `equipmentId` (uint) | Fulfillment |
| `fulfillsEquipmentTemplateId` | `fulfillsEquipmentId` (uint) | Fulfillment |
| `groupTemplateId` | `groupId` (uint) | ExerciseGroupMember |
| `exerciseTemplateId` | `exerciseId` (uint) | ExerciseGroupMember |

### New Fields
| Field | Type | On |
|-------|------|-----|
| `public` | `boolean` | Exercise, Equipment |
| `owner` | `string` | Exercise, Equipment, ExerciseScheme, ExerciseRelationship, Fulfillment, ExerciseGroup, ExerciseGroupMember |

### Type Change
- `Exercise.equipmentIds`: `string[]` → `number[]` (template IDs → numeric IDs)

## Query Key Updates

| Old Key | New Key |
|---------|---------|
| `userExerciseKeys.*` | Remove — use `exerciseKeys.*` with owner filter |
| `userEquipmentKeys.*` | Remove — use `equipmentKeys.*` with owner filter |
| `exerciseSchemeKeys.*` | Keep, but URL changed |

## API Client Changes

### `CompendiumApiClient`
- All methods still work — the endpoints haven't changed
- Response types now include `owner` and `public` fields

### `UserApiClient`
- Remove: `fetchUserExercises`, `fetchUserExercise`, `createUserExercise`, `deleteUserExercise`
- Remove: `fetchUserEquipment`, `fetchUserEquipmentItem`, `createUserEquipment`, `deleteUserEquipment`
- Update: `fetchExerciseSchemes` URL from `/api/user/exercise-schemes` to `/api/exercise-schemes`
- Add scheme methods also need URL update for create/get/update/delete

### Consider merging `CompendiumApiClient` and `UserApiClient`
Since exercises/equipment are now unified, the two clients could merge into one `ApiClient`. Exercise/equipment methods live in the compendium client; workout/log methods live in the user client. Or unify fully.

## Component Migration

### Key Pattern Change
**Old**: User exercise list fetches thin `UserExercise` wrappers, then for each, queries compendium for the full exercise snapshot via `templateId + version`. Two-step fetch.

**New**: Exercise list with `?owner=me` returns full exercises directly. No second fetch needed.

### Pages to Update

1. **User Exercise List** (`pages/user/user-exercise-list/`) — Query `/api/exercises?owner=me` instead of `/api/user/exercises`. Remove snapshot enrichment logic.

2. **User Exercise Detail** (`pages/user/user-exercise-detail/`) — Query `/api/exercises/:id` directly. Remove compendium version lookup.

3. **User Equipment List/Detail** — Same pattern as exercises.

4. **Compendium Exercise Detail** (`pages/compendium/exercise-detail/`) — The "Add to Mine" button no longer makes sense (exercises are already owned). Consider replacing with a "Copy" or "Fork" action.

5. **Workout Edit** (`pages/user/workout-edit/`) — Exercise scheme references are now via `exerciseId` not `userExerciseId`. The exercise enrichment logic (fetching compendium snapshots) can be removed.

6. **Workout Start** (`pages/user/workout-start/`) — Same simplification.

7. **Exercise Config** (`ui/exercise-config/`) — Uses `compendiumExerciseId`/`compendiumVersion` — replace with direct exercise data.

### Temporary Compatibility
Type aliases (`UserExercise = Exercise`, `UserEquipment = Equipment`, `UserExerciseScheme = ExerciseScheme`) are defined in `user-models.ts` to keep the build passing. Remove these once components are migrated.

## Route Structure
The `/compendium` and `/user` route split may no longer make sense. Consider:
- `/exercises` (unified, filtered by visibility)
- `/exercises/:id` (detail, shows owner info)
- `/my/exercises` or `/exercises?owner=me` (user's own)
- Keep `/user/workouts`, `/user/calendar` under user routes (inherently user-scoped)

## New Entity: Equipment Relationships
- Endpoint: `/api/equipment-relationships`
- Type: `EquipmentRelationship` with `fromEquipmentId`, `toEquipmentId`, `owner`, `type` ("equivalent")
- No UI exists yet — build when needed
