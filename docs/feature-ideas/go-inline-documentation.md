# Go Inline Documentation

## Summary

Adopt Go's built-in documentation conventions — `doc.go` package comments and `Example` test functions — to embed topical explanations about entity relationships, API behavior, and domain concepts directly in the code. Rendered as browsable HTML via `pkgsite`.

## Motivation

Architectural knowledge currently lives in separate markdown files (`docs/feature-ideas/`, `CLAUDE.md`). These drift from the code they describe — when a handler changes, the markdown doesn't update automatically. Phoenix/Elixir solves this with ExDoc: documentation lives *in* the code, examples are executable tests, and the whole thing renders into a browsable site.

Go doesn't have ExDoc, but it has two underused mechanisms that get most of the value:

1. **`doc.go` files** — package-level comments with Go 1.19+ formatting (headings, lists, links, code blocks). Rendered by `pkgsite` as HTML with cross-references between packages.
2. **`Example` test functions** — runnable code that appears in the rendered docs. They serve as both documentation and regression tests — if the example breaks, `go test` fails.

Together these keep explanations co-located with the code they describe, and the examples stay honest because they're executed.

## Design Sketch

### `doc.go` per domain package

Each package in `internal/` gets a `doc.go` explaining its domain role, entity model, and relationships to other packages:

```
internal/exercise/doc.go           — entity model, ownership, visibility rules
internal/exerciserelationship/doc.go — relationship types, directionality, visibility derivation
internal/equipment/doc.go          — entity model, how it differs from exercises
internal/equipmentfulfillment/doc.go — fulfillment vs equivalence distinction
internal/equipmentrelationship/doc.go — equipment equivalence model
internal/exercisegroup/doc.go      — grouping model, membership
internal/user/doc.go               — user-scoped entities overview
internal/auth/doc.go               — auth middleware, X-User-Id flow
```

Use Go 1.19 doc comment syntax for structure:

```go
// Package exercise manages exercise entities and their lifecycle.
//
// # Entity Model
//
// An [ExerciseEntity] is owned by a user (the Owner field) and optionally
// marked public. The "compendium" is a filtered view over all exercises
// where Public is true — not a separate data model.
//
// # Visibility
//
// Exercises with Public=true are readable by all users. Non-public exercises
// are visible only to their owner. See [auth.GetUserID] for how the
// requesting user is determined.
//
// # Relationships
//
// Exercises link to other exercises via [exerciserelationship] (progressions,
// alternatives, equivalences) and to equipment via [equipmentfulfillment]
// (which equipment an exercise requires or accepts).
package exercise
```

The `[TypeName]` and `[package.Symbol]` syntax creates cross-references in pkgsite.

### Example test functions

Add `Example` functions in `*_test.go` files to document key API behaviors:

```go
func ExampleHandler_CreateExercise() {
    // Creating an exercise sets the owner from the auth context
    // and auto-generates a slug from the name.
    //
    // POST /api/exercises
    // {"name": "Bench Press", "type": "STRENGTH", ...}
    //
    // Response:
    // {"id": 1, "name": "Bench Press", "owner": "alice", "public": false, ...}
}
```

These run with `go test` — if the output comment doesn't match, the test fails.

### Local rendering

Add a make target to serve docs locally:

```makefile
docs:
	go run golang.org/x/pkgsite/cmd/pkgsite@latest -open .
```

### Rollout order

1. Start with `internal/exercise/doc.go` as the template — it's the most central domain package.
2. Add `Example` functions for create/update/delete flows in the exercise handlers.
3. Expand to other packages, prioritizing those with non-obvious domain logic (fulfillment vs equivalence, relationship visibility derivation).
4. Add the `make docs` target.

## Considerations

- **Cross-package narratives**: Go docs are per-package. For stories that span multiple packages (e.g., "how an exercise flows from creation to workout log"), the top-level `internal/doc.go` or `internal/shared/doc.go` can serve as an overview with links to individual packages.
- **Example maintenance**: examples are real tests, so they break when behavior changes — this is a feature, not a bug. Keep examples focused on stable API contracts, not implementation details.
- **Not a replacement for feature docs**: `doc.go` explains *what the code does now*. Feature idea docs explain *what we want to build and why*. Both serve different purposes.
