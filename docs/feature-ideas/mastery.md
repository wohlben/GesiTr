# Exercise Mastery

## Summary

Replace the "import to My Exercises" flow (which creates carbon copies of compendium exercises) with a **mastery system** that tracks how experienced a user is with each exercise based on their logged workouts. Users consume exercises directly from the compendium — no copying needed.

## Motivation

Importing exercises creates pointless duplication. The user doesn't actually customize the exercise definition — they just want to use it. Mastery gives that "ownership" feeling through progression instead of copying.

## Mastery System Design

### XP Source

Each **rep** counts as 1 base XP toward that exercise's mastery. For non-rep-based exercises (time/distance), each exercise log entry = 1 base XP.

### Leveling Formula

**Flat XP per level**: 100 XP to advance each level (the recency multiplier provides the acceleration curve, not the level formula itself).

### Recency Multiplier

Effective XP is scaled by training consistency:

```
effective_xp = base_xp * min(0.5 * n_days, max(1, current_level / 2))
```

Where:
- `n_days` = number of **distinct days** the user has performed this exercise within a configurable window (default: 6 months)
- `current_level / 2` = multiplier cap (grows with level)
- `max(1, ...)` = floor of 1 so level 0 users can still earn XP

**How it plays out**:
- First ever session (n_days=1): multiplier = min(0.5, 1) = **0.5x** — half XP, proving commitment
- After 4 sessions on separate days (n_days=4): multiplier = min(2.0, cap) — accelerating
- Cap grows with level: level 1 caps at 1x, level 10 caps at 5x, level 30 caps at 15x
- If a user stops training, `n_days` in the window drops over time, naturally slowing progression

### Mastery Tiers

| Level Range | Tier        | Target Timeline |
|-------------|-------------|-----------------|
| 0 - 10      | Novice      | ~3 months       |
| 11 - 30     | Journeyman  | ~12 months cumulative |
| 31 - 50     | Adept       | ~3 years cumulative |
| 51 - 99     | Master      | long-term       |
| 100         | Mastered    | lifetime goal   |

Each tier roughly takes 2x the cumulative time of all previous tiers.

### Example Progression

Assuming 2 workouts/week, 3 sets of 10 reps (60 reps/week base):

| Phase | n_days (6mo window) | Multiplier (at level) | Effective XP/week |
|-------|--------------------|-----------------------|-------------------|
| Week 1 | 2 | min(1.0, max(1, 0/2)) = **1.0x** | 60 |
| Week 4 | 8 | min(4.0, max(1, 1/2)) = **1.0x** | 60 |
| Week 12 (level ~7) | 24 | min(12.0, 7/2) = **3.5x** | 210 |
| Week 24 (level ~15) | 48 | min(24.0, 15/2) = **7.5x** | 450 |
| Week 48 (level ~28) | 96 | min(48.0, 28/2) = **14.0x** | 840 |

Early levels are a slow grind (cap limits you). As you level up, the cap rises and consistency pays off — creating natural acceleration through the Journeyman tier.

### Relationship-Based XP Contributions

Logs from related exercises contribute to mastery with a multiplier derived from two components:

```
multiplier = (strength * 0.5) + type_bonus
```

**Type bonuses** by relationship category:

| Category | Type Bonus | Relationship Types |
|----------|------------|-------------------|
| Equivalent | 0.5 | `equivalent` |
| Skill transfer | 0.25 | `alternative`, `easier_alternative`, `harder_alternative`, `equipment_variation`, `variant`, `variation`, `bilateral_unilateral`, `progression`, `progresses_to`, `regression`, `regresses_to` |
| No transfer | — | `accessory`, `antagonist`, `complementary`, `preparation`, `prerequisite`, `related`, `similar`, `superset_with`, `supports`, `forked` |

"No transfer" relationships are ignored entirely (no contribution row created).

**Examples**:

| Relationship | Strength | Calculation | Multiplier |
|-------------|----------|-------------|------------|
| `equivalent`, strength=1.0 | 1.0 | 0.5 + 0.5 | **1.0** |
| `equipment_variation`, strength=1.0 | 1.0 | 0.5 + 0.25 | **0.75** |
| `variant`, strength=0.8 | 0.8 | 0.4 + 0.25 | **0.65** |
| `alternative`, strength=0.5 | 0.5 | 0.25 + 0.25 | **0.50** |

An exercise's own logs always count at **1.0** (implicit, no row needed).

If multiple relationships exist between the same pair of exercises, the **highest** multiplier wins.

## Technical Approach

### Precomputed table: `mastery_contributions`

```
mastery_contributions
  exercise_id              uint    -- the exercise gaining mastery
  contributes_from_id      uint    -- the exercise whose logs contribute
  multiplier               float64 -- combined (strength * 0.5) + type_bonus
  relationship_type        string  -- source relationship (for debugging)
```

This table is a **precomputed lookup** — not a cache that can go stale, but a derived view that's recalculated on two lifecycle events:
- **Relationship created**: Compute multiplier, insert/update row (keep highest if multiple relationships exist)
- **Relationship removed**: Recompute for the affected pair (may remove row or fall back to next-best relationship)

### Query Flow

To compute mastery for a given exercise:

1. **Look up contributors**: `SELECT contributes_from_id, multiplier FROM mastery_contributions WHERE exercise_id = ?`
2. **Count reps per contributor**: `SELECT exercise_id, SUM(reps) FROM exercise_logs WHERE exercise_id IN (?, self) AND owner = ? GROUP BY exercise_id`
3. **Apply multipliers**: `total_xp = own_reps * 1.0 + SUM(contributor_reps * multiplier)`
4. **Count distinct days** (for recency multiplier): across self + all contributors
5. **Compute** level, tier, progress

### Index

A composite index on `exercise_logs(owner, exercise_id, performed_at)` makes the rep sum and distinct-day count index-only scans.

### Backend

- **Mastery DAO/service**:
  - Queries `mastery_contributions` for contributor exercises + multipliers
  - Aggregates reps + distinct days across contributors via indexed queries
  - Computes effective XP, level, tier, and progress to next level
- **Relationship lifecycle hooks**: Recalculate `mastery_contributions` rows on relationship create/delete
- **Endpoints**: `GET /api/user/mastery` (all exercises with any logs) and/or `GET /api/user/mastery/:exerciseId`

### Frontend

Display mastery badge/level on exercise cards; replace or augment the import flow so users work with compendium exercises directly.

## Open Questions

- **My Exercises**: Remove imports entirely, or keep them for users who genuinely want to customize an exercise definition?
- **Visual treatment**: Colors, icons, progress bars for mastery tiers?
- **Tier naming**: Journeyman / Adept feel right? Other options?
- **Recency window**: 6 months default — should this be user-configurable?
- **Non-rep XP**: For time/distance exercises, should XP scale with duration/distance instead of flat 1 per log?
