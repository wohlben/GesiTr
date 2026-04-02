# Equipment Mastery

## Summary

Replace the "My Equipment" flow (which creates carbon copies of compendium equipment) with a **mastery system** that tracks how experienced a user is with each piece of equipment based on their logged workouts. Users consume equipment directly from the compendium — no copying needed. This mirrors the Exercise Mastery pattern.

The "Fork" button remains as a secondary action for users who want to create a custom variant (e.g., a specific brand or configuration).

## Motivation

"My Equipment" creates private copies of compendium equipment, duplicating data with no real benefit. The user doesn't customize the equipment definition — they just want to indicate they have access to it. Equipment mastery solves this by deriving the relationship from actual training data: when you log exercises that use a barbell, you gain experience with that barbell. That experience is the signal, not a manual import.

The result: a single unified equipment list where relevant equipment surfaces automatically based on training history.

## Equipment Mastery Design

### XP Source

Equipment mastery XP is **derived from exercise logs**. When a user logs reps for an exercise, every piece of equipment referenced by that exercise (via `exercise_equipments`) earns XP:

- Each **rep** logged = 1 base XP for each piece of equipment the exercise uses
- For non-rep-based exercises (time/distance), each log entry = 1 base XP per piece of equipment

XP is **not divided** across equipment. If an exercise uses a barbell and a bench, both earn the full rep count — performing 10 reps of bench press gives 10 reps of barbell experience *and* 10 reps of bench experience.

### Leveling Formula

Reuse the exercise mastery formulas exactly:

- **100 XP per level**, max level 100
- **Recency multiplier**: `effective_xp = base_xp * min(0.5 * n_days, max(1, current_level / 2))`
  - `n_days` = distinct days the user logged any exercise using this equipment within a 6-month window
- **Same tier structure**:

| Level Range | Tier        |
|-------------|-------------|
| 0 - 10      | Novice      |
| 11 - 30     | Journeyman  |
| 31 - 50     | Adept       |
| 51 - 99     | Master      |
| 100         | Mastered    |

The computation functions in `internal/user/mastery/models/compute.go` (`ComputeLevel`, `ComputeTier`, `ComputeRecencyMultiplier`, `ComputeProgress`) are reusable as-is — they operate on raw numeric inputs with no exercise-specific logic.

### Relationship-Based XP Contributions

Equipment relationships currently have two types: `equivalent` and `forked`.

| Relationship Type | Contributes? | Multiplier | Rationale |
|---|---|---|---|
| `equivalent` | Yes | `(strength * 0.5) + 0.5` | Same formula as exercise mastery. Experience with equivalent equipment should transfer. |
| `forked` | No | — | Fork is a copy action, not a skill-transfer signal. |

If contributions are enabled, a `equipment_mastery_contributions` table is needed, mirroring `mastery_contributions`. Recalculated on equipment relationship create/delete.

### Equipment Fulfillments as Contributions

The `fulfillments` table ("Equipment X fulfills Equipment Y") represents substitutability. If you mark that your home cable machine fulfills a gym lat pulldown machine, experience should transfer between them.

Fulfillments could be treated as a contribution source with a skill-transfer-level multiplier (e.g., `0.5 + 0.25 = 0.75`). This would make the system more useful for users with home gym equivalents of commercial equipment.

## Technical Approach

### Precomputed Table: `equipment_mastery_experience`

```
equipment_mastery_experience
  owner         string  (PK)   -- user ID
  equipment_id  uint    (PK)   -- equipment item
  total_reps    int            -- precomputed rep count
```

Mirrors `mastery_experience(owner, exercise_id, total_reps)` exactly.

**Why not compute at query time?** The recency multiplier requires counting distinct training days per equipment item in a 6-month window. Combined with sorting the equipment list by mastery value, this aggregation is too expensive to run on every page load without precomputed data.

### Optional: `equipment_mastery_contributions`

```
equipment_mastery_contributions
  owner                string   (PK)
  equipment_id         uint     (PK)   -- equipment gaining mastery
  contributes_from_id  uint     (PK)   -- equipment whose logs contribute
  multiplier           float64         -- combined (strength * 0.5) + type_bonus
  relationship_type    string          -- source relationship (for debugging)
```

Recalculated on equipment relationship create/delete. Same bidirectional pattern as exercise mastery contributions.

### Lifecycle Hooks

**Exercise log creation** — two call sites where `UpsertExperience` is already called:

1. `internal/user/exerciselog/handlers/exerciselog_handlers.go` (line 86)
2. `internal/user/workoutlog/handlers/set.go` (line 214)

After each, add a call to `UpsertEquipmentExperience(db, owner, exerciseID, reps)`. This function:
1. Queries `exercise_equipments WHERE exercise_id = ?` to find equipment IDs
2. Upserts each one into `equipment_mastery_experience`

The `exercise_equipments` lookup is cheap (indexed join table, typically 1-3 rows) and runs within the existing transaction.

**Equipment relationship create/delete** — add calls to `RecalculateEquipmentContributions` (if contributions are enabled).

### Query Flow

To compute equipment mastery for a given equipment item:

1. **Look up contributors**: `SELECT contributes_from_id, multiplier FROM equipment_mastery_contributions WHERE equipment_id = ?`
2. **Count reps per contributor**: `SELECT equipment_id, total_reps FROM equipment_mastery_experience WHERE equipment_id IN (?, self) AND owner = ?`
3. **Apply multipliers**: `total_xp = own_reps * 1.0 + SUM(contributor_reps * multiplier)`
4. **Count distinct days**: `SELECT COUNT(DISTINCT DATE(performed_at)) FROM exercise_logs el JOIN exercise_equipments ee ON ... WHERE ee.equipment_id IN (?, contributors) AND el.owner = ? AND performed_at >= recency_start`
5. **Compute** level, tier, progress using existing `compute.go` functions

### Equipment List Handler Changes

New query parameter `mastery=me` on `GET /api/equipment`:

```sql
SELECT e.*
FROM equipment e
LEFT JOIN equipment_mastery_experience eme
  ON eme.equipment_id = e.id AND eme.owner = ?
WHERE e.owner = ?
   OR e.id IN (SELECT equipment_id FROM equipment_mastery_experience WHERE owner = ?)
ORDER BY
  CASE WHEN e.owner = ? THEN 0 ELSE 1 END,
  COALESCE(eme.total_reps, 0) DESC,
  e.display_name ASC
```

Default behavior (no `mastery` param): unchanged — shows own + public equipment.

### Endpoints

- `GET /api/user/equipment-mastery` — list mastery for all equipment the user has logged
- `GET /api/user/equipment-mastery/:equipmentId` — mastery for a single equipment item

### Migration Backfill

```sql
INSERT INTO equipment_mastery_experience (owner, equipment_id, total_reps)
SELECT el.owner, ee.equipment_id, SUM(COALESCE(el.reps, 1))
FROM exercise_logs el
JOIN exercise_equipments ee ON ee.exercise_id = el.exercise_id
WHERE el.deleted_at IS NULL
GROUP BY el.owner, ee.equipment_id
```

### Backend Package Structure

New package `internal/user/equipmentmastery/` with:
- `models/` — `EquipmentMasteryExperienceEntity`, `EquipmentMastery` DTO, `EquipmentMasteryContributionEntity`
- `handlers/` — `ListEquipmentMastery`, `GetEquipmentMastery`, `UpsertEquipmentExperience`, `RecalculateEquipmentContributions`, `BackfillEquipmentExperience`, routes

### Frontend

**Add:**
- Equipment mastery API client methods (`fetchEquipmentMasteryList`, `fetchEquipmentMastery`)
- Equipment mastery query keys (`equipmentMasteryKeys`)
- Mastery column on equipment list (tier-colored level badges)
- Mastery card on equipment detail page (level, tier, XP, progress bar)
- TypeScript types via tygo

**Modify:**
- Equipment detail: demote "Fork" to secondary, add mastery card
- Equipment list: fetch mastery data, pass to list items

**Remove:**
- `/user/equipment` and `/user/equipment/:id` routes
- `UserEquipmentList` and `UserEquipmentDetail` components
- `userEquipmentKeys` from query-keys
- `fetchUserEquipment()`, `fetchUserEquipmentItem()`, `deleteUserEquipment()` from UserApiClient
- "My Equipment" navigation link

## Open Questions

- **Fulfillments as contributions**: Should the `fulfillments` table feed into equipment mastery contributions? If equipment A fulfills equipment B, does experience with A partially count toward B? Adds complexity (watch fulfillment create/delete events) but is more useful for home gym setups.
- **XP attribution when exercises change equipment**: If exercise X had equipment [A, B] when logged but later changes to [A, C], mastery for B is already stored and won't be retroactively corrected. Same as exercise mastery — acceptable?
- **Equipment mastery inline vs. separate fetch**: Should the equipment list endpoint return mastery inline (optional field on Equipment DTO), or should frontend fetch mastery separately and merge client-side? The exercise pattern (separate fetch + merge) is simpler.
