# Go Inline Documentation

## Summary

Go's built-in documentation conventions — `doc.go` package comments and `Example` test functions — are used to embed API documentation, entity relationships, and domain concepts directly in the code. Examples are executable tests that verify documented behavior. Documentation is rendered as HTML and served at `/docs`.

## What was implemented

### `doc.go` per handler package

Each handler package has a `doc.go` with a package-level overview explaining entity relationships, endpoint hierarchy, and cross-package links:

- `internal/shared/doc.go` — permission model, pagination, history, base model
- `internal/exercise/handlers/doc.go` — exercises, schemes, versioning, cross-links to equipment and workouts
- `internal/equipment/handlers/doc.go` — equipment, versioning, cross-links to exercises
- `internal/user/workout/handlers/doc.go` — workout → section → section-exercise hierarchy, cross-links to exercise schemes

### Handler doc comments

Every exported handler function has a doc comment explaining:
- What the endpoint does and its HTTP method/path
- Relationships to other endpoints (e.g., "requires a workout created via CreateWorkout")
- Permission/visibility behavior
- Cross-package links using Go's `[package.Symbol]` syntax

### `Example` test functions

Executable examples in `example_*_test.go` files that run with `go test` and verify output via `// Output:` comments. Request payloads use raw JSON strings for readability. Examples cover:

- **Permissions**: owner / non-owner+public / non-owner+private for exercises, equipment, exercise schemes, workouts
- **CRUD flows**: CreateExercise (with/without equipment), CreateExerciseScheme (rep-based/time-based), UpdateExercise (with versioning), ListExercises, ListExerciseSchemes
- **Versioning**: create → update → verify snapshots, version history after delete
- **Cross-API flows**: creating equipment via the equipment API then referencing it from an exercise

### Documentation rendering

- `make docs` runs `scripts/docs.sh` which auto-discovers all `doc.go` files and generates markdown via `gomarkdoc`
- Markdown is embedded into the binary via `//go:embed` and converted to HTML at startup using `goldmark`
- Served at `/docs` with an index page linking to each package's documentation
- The Dockerfile generates docs during the build so the production image includes them

### Behavior fixes driven by documentation

Writing examples exposed missing visibility checks and inconsistent behavior:

- **GetExercise / GetEquipment**: added permission checks (previously returned any exercise by ID regardless of visibility)
- **GetExerciseVersion**: now queries the history table directly via `json_extract` so version history survives exercise deletion
- **GetExerciseScheme**: permissions now derive from the linked exercise's visibility, not scheme ownership
- **ListExerciseSchemes**: now returns schemes for public exercises, not just the current user's
- **Permissions endpoints**: non-owner+private now returns empty permissions array instead of 404

## Considerations

- **pkgsite doesn't work** with non-domain module names like `gesitr` (returns 424 on all pages). We use `gomarkdoc` + `goldmark` instead.
- **Example maintenance**: examples are real tests, so they break when behavior changes — this is a feature, not a bug.
- **Future: OpenAPI integration**: the `doc.go` function links currently point to GitHub source. Once the Huma migration is done (see `huma-openapi.md`), these can link to OpenAPI endpoint documentation instead.
