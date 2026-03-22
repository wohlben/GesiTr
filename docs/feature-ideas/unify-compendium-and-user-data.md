# Unify Compendium and User Data

## Summary

Remove the compendium as a separate entity. Instead, user-owned exercises and equipment gain a **public** visibility toggle. The compendium becomes a view over all public user data rather than its own parallel data model. Exercise relationships — which already exist but are compendium-scoped — become user-owned entities in the same way, with a new **equivalent** type solving the cross-user deduplication problem for both exercises (as a relationship type) and equipment (as a distinct equivalence table).

## Motivation

The compendium currently exists only so users can import exercises/equipment into their personal collection. This creates two parallel entity hierarchies (compendium exercises vs. user exercises) that model the same real-world thing. If user data itself could be marked public, the compendium is just a query filter — not a separate domain.

Similarly, exercise relationships today live in the compendium domain with no user ownership. They describe how exercises relate to each other (progressions, alternatives, accessories, etc.) but users can't create, customize, or extend them. Folding relationships into the user domain gives users control over their own exercise graphs while preserving the curated seed data as a public starting point.

## Requirements

- **[User Profiles](user-profiles.md)** — `user_profiles` table with display names, plus FK migration of all `Owner`/`CreatedBy` fields. Must land first so that the unified model has referential integrity from the start.

## Design Sketch

### Exercises & Equipment

- Exercises and equipment get a visibility field: `public: bool`.
- The "compendium" UI becomes a filtered view: all exercises/equipment where `public = true`.
- The current compendium tables graduate into a secondary role — possibly as a curated/canonical seed set, or removed entirely.
- Seed data becomes owned by the system user `"Claude"` and marked public.

**Why a bool, not an enum?** Other visibility scopes (e.g., "visible to members of a workout group") are determined by membership in a separate table — the exercise itself doesn't know which group it belongs to. An enum value like `group` would be redundant with (and less authoritative than) the group membership lookup. The entity-level question is strictly "is this in the public compendium?" — yes or no. Group-scoped access is derived at query time from the relevant feature's own tables.

### Relationships as User Entities

The codebase already has a full exercise relationship system (`internal/compendium/exerciserelationship/`):

- **20 directional types**: accessory, alternative, antagonist, bilateral_unilateral, complementary, easier_alternative, equipment_variation, harder_alternative, preparation, prerequisite, progresses_to, progression, regresses_to, regression, related, similar, superset_with, supports, variant, variation
- **Strength-weighted**: each relationship has a `strength` float (0.0–1.0)
- **4,046 seeded relationships** in `data/compendium_relationships/`
- **`createdBy` field** already exists on the model
- **No frontend UI** — backend and seed data only

In the unified model, relationships become user-owned like exercises:

- `createdBy` becomes the owner (scoped to a user).
- Relationships have **no own visibility field** — visibility is derived from the linked exercises (visible when both endpoints are visible to the viewer).
- The 4,046 seeded relationships become owned by `"Claude"`. Since Claude's exercises are public, these relationships are visible to everyone.
- Users inherit all public relationships and can create their own (visible to anyone who can see both linked exercises).
- A relationship between a user's private exercise and a public exercise is only visible to that user.

### Exercise Equivalence

Instead of a separate equivalence table, **`equivalent`** becomes a new exercise relationship type with special semantics:

- **Directional by ownership**: each user declares their own equivalences ("my pushup ≡ that public pushup"). User A's declaration does not create any record in User B's space. Each user manages only their own outbound equivalence links.
- **Query semantics**: when looking up equivalences for an exercise, the system checks both `from` and `to` directions — so a single record A→B means both "A is equivalent to B" and "B is equivalent to A" from the perspective of whoever can see the relationship.
- **System-derived grouping**: the system follows all visible equivalence declarations to build clusters for deduplication and comparison, without any single declaration leaking across user boundaries.

This solves the cross-user deduplication problem natively within the existing relationship model:

- **Deduplication in public view**: when displaying the compendium, public exercises linked by `equivalent` are grouped. One representative is shown instead of N duplicates from N users.
- **Cross-user comparison**: when users mark their exercises public and declare equivalences, this enables features like leaderboards, aggregate stats, and "how does my bench compare to others?"
- **Import suggestions**: when a user creates an exercise, the system can suggest equivalences to existing public exercises based on name/attribute similarity.
- **Merge path**: accepting an equivalence suggestion links the user's exercise to the public graph without losing their personal data or customizations.

### Equipment Equivalence

Equipment needs the same cross-user identity linking, but as a **distinct table** — not part of the exercise relationship system.

**Why not reuse equipment fulfillment?** The existing `fulfillments` table (`internal/compendium/equipmentfulfillment/`) solves a different problem: *substitutability*. "Adjustable-bench fulfills decline-bench" means it can stand in for it in exercise requirements — not that they're the same piece of equipment. Fulfillment is directional and asymmetric (adjustable-bench fulfills decline-bench, but not vice versa). Overloading it for identity would conflate "can substitute for" with "is the same as," making frontend display confusing.

Equipment equivalence is a separate concept:

- **Distinct table**: `equipment_relationships` — mirrors the exercise relationship model structurally (`type`, `strength`, `owner`, `from`/`to`). Initially only `equivalent` is defined as a type; the type field allows future equipment-specific relationship types to be added without schema changes.
- **No visibility field**: derived from endpoint visibility, same as exercise relationships.
- **Directional by ownership**: same model as exercise equivalence — each user declares "my dumbbell ≡ that public dumbbell." No cross-user record creation. Strength is always 1.0 for equivalence.
- **Same use cases**: deduplication in public view, cross-user matching, import suggestions.
- **Coexists with fulfillment**: fulfillment answers "what can I use instead?" while equivalence answers "is this the same thing?" A user's "adjustable bench" is *equivalent* to another user's "adjustable bench," and separately *fulfills* exercises requiring a "bench."

### Relationship Lifecycle

1. **Seed**: `make seed` loads 4,046 exercise relationships owned by `"claude"` (progressions, alternatives, etc.). Equivalence links between seed exercises are not needed — all seed exercises are already owned by the same user (`"claude"`), so there's nothing to deduplicate.
2. **Inherit**: all users see public relationships when browsing exercises. These provide a curated exercise graph out of the box.
3. **Create**: users add private relationships for their own exercise collections (e.g., "for me, incline dumbbell press is an easier_alternative to flat barbell bench").
4. **Propose equivalence**: when a user creates a new exercise, the system suggests equivalences. Accepting creates an `equivalent` relationship linking the user's exercise to the public graph.
5. **Go public**: when a user marks their exercise public, all relationships involving it become visible to anyone who can also see the other endpoint — no separate publish step for relationships.

## Resolved Decisions

- **Equivalence governance**: no trust or approval needed. Each user only manages their own equivalence declarations — there's nothing to moderate.
- **Equivalence transitivity**: doesn't apply. Equivalences are per-user declarations, not a shared graph. User A's equivalences don't chain with User B's.
- **Conflict resolution**: not an issue. Public exercises are namespaced by author in the UI (`@claude/Pushup`, `@alice/Pushup`). Same-named exercises from different users are distinct and unambiguous.
- **Seeding pipeline**: `make seed` creates all entities under owner `"claude"` (the user profile ID) with `public = true`.
- **Relationship visibility**: relationships don't have their own visibility field. Visibility is derived from the linked exercises — if both endpoints are visible to a viewer, the relationship is too. If either endpoint is non-public, the relationship is only visible to users who can see both exercises.
- **Migration**: not needed. No user data exists yet — just wipe the database and re-seed with the new schema (`make seed`).
- **Strength on equivalence**: always 1.0. Equivalence means identity, not similarity. "Similar but debatable" is what the `alternative` relationship type is for.
- **Performance**: index on `(public, relationship_type)` for efficient compendium queries.

## Open Questions

None remaining — all design questions have been resolved above.
