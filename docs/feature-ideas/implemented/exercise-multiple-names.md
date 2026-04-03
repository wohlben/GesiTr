# Exercise Multiple Names

## Summary

Exercises have 1:n names instead of a single primary name. All names are equal peers stored in a dedicated `exercise_names` join table with auto-increment IDs, enabling referential integrity for user preferences.

## Model

- `exercise_names` table: `id` (PK, auto-increment), `exercise_id`, `position`, `name`
- Replaces the old single `name` column on exercises and the `exercise_alternative_names` table
- `exercise_name_preferences` table: `owner` + `exercise_id` (composite PK), `exercise_name_id` (FK to `exercise_names.id`) — stores each user's preferred display name
- Names are served ordered by popularity (number of users who selected each name as their preference), with original position as tiebreaker

## API

- Exercise DTO exposes `names: [{id, name}, ...]` — structured objects so clients can reference names by ID
- Exercise create/update bodies accept `names: string[]` (plain strings; IDs are server-assigned)
- `GET /api/user/exercise-name-preferences` — list user's preferred names
- `PUT /api/user/exercise-name-preferences/:exerciseId` with `{exerciseNameId}` — set preference
- `DELETE /api/user/exercise-name-preferences/:exerciseId` — remove preference
- Search (`?q=`) queries the `exercise_names` table

## UI

- **List view (no search)**: one row per exercise, displays user's preferred name or the most popular name (`names[0]`)
- **List view (with search)**: primary name is the best match (startsWith > contains); other matching names shown as muted sub-rows below
- **Clicking any name** in search results navigates to the exercise detail AND saves it as the user's preferred default (via the preference API)
- **Detail view**: shows all names in a "Names" section
- **Edit form**: unified names list with add/remove (at least one required)
- **Exercise config combobox**: searches across all names, displays preferred/first name

## Migration

- Idempotent `runMigrations()` block: moves `exercises.name` to `exercise_names` at position 0, then `exercise_alternative_names` entries at position 1+ (deduplicated), drops old column and table
- Seed code keeps existing JSON format (`name` + `alternativeNames`) and merges them at seed time
