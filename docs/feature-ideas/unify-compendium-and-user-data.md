# Unify Compendium and User Data

## Summary

Remove the compendium as a separate entity. Instead, user-owned exercises and equipment gain a **public** visibility toggle. The compendium becomes a view over all public user data rather than its own parallel data model.

## Motivation

The compendium currently exists only so users can import exercises/equipment into their personal collection. This creates two parallel entity hierarchies (compendium exercises vs. user exercises) that model the same real-world thing. If user data itself could be marked public, the compendium is just a query filter — not a separate domain.

## Design Sketch

- Exercises and equipment get a visibility field (e.g., `public: bool`).
- The "compendium" UI becomes a filtered view: all exercises/equipment where `public = true`.
- The current compendium tables graduate into a secondary role — possibly as a curated/canonical seed set, or removed entirely.

### The Equivalence Problem

If user A has a "pushup" and user B has a "pushup," these are distinct rows (user-scoped). But they represent the same real-world exercise. We need a way to express that.

A new **equivalence** table (or similar join structure) would link user-specific entities that represent the same thing. This enables:

- Deduplication in the public view ("show one pushup, not 500").
- Cross-user comparison (group features, leaderboards).
- Merging/suggesting when a user creates something that already exists publicly.

## Open Questions

- Who creates and manages equivalence links? Automatic (name matching)? User-proposed? Admin-curated?
- What happens to the existing compendium import flow? Does "import" become "copy a public exercise to my collection"?
- How do we handle conflicts — two public exercises with the same name but different definitions?
- Does this affect the data seeding pipeline (`make seed`)? Seed data would need to be owned by a system user?
- Should equivalence be 1:1 or many-to-one (canonical exercise with many user variants)?
- Performance: querying all public exercises across all users vs. a dedicated table.
