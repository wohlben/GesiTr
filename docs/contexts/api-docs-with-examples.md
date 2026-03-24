# API Documentation with Go Example Tests

> **Keep this file in sync.** If conventions change (file naming, test patterns, rendering pipeline, etc.), update this context so future sessions don't follow stale guidance.

## How it works

Each handler package has:
- **`doc.go`** — package-level overview explaining entity relationships, endpoint hierarchy, and cross-package links using Go 1.19+ doc comment syntax (`# Headings`, `[Symbol]` cross-references)
- **Handler doc comments** — every exported handler function gets a comment with: what it does, HTTP method/path, relationships to other endpoints, permission behavior
- **`example_*_test.go`** files — Go `Example` functions that serve as both documentation and executable tests. They run with `go test` and fail if the `// Output:` comment doesn't match

## File structure per handler package

```
internal/<domain>/handlers/
├── doc.go                          # Package overview
├── <handler>.go                    # Handler code with doc comments
├── example_setup_test.go           # Shared: setupExampleDB(), doRawAs(), newExampleRouter()
├── example_<resource>_test.go      # Examples for each resource/handler group
└── setup_test.go                   # Existing test infrastructure (setupTestDB, newRouter, etc.)
```

Split examples by handler/resource — one file per logical group (e.g., exercises, schemes, versions, permissions).

## Writing setupExampleDB()

Example functions can't accept `*testing.T`, so `setupExampleDB()` uses `os.Setenv`/`panic` instead of `t.Setenv`/`t.Fatal`. It mirrors the package's `setupTestDB()` but:
- Uses `os.Setenv("AUTH_FALLBACK_USER", "<default-user>")`
- Creates a second user (e.g., "other" or "bob") for non-owner examples
- Migrates all models needed by the examples (including cross-package models if needed)

If examples need cross-API flows (e.g., creating equipment from exercise examples), add a `newExampleRouter()` that extends `newRouter()` with the additional routes.

## Writing Example functions

### Naming convention

```go
func Example<HandlerName>_<scenario>()
```

Go links `Example<FunctionName>` to that function's docs. The `_suffix` differentiates scenarios:
- `ExampleGetExercise_owner`
- `ExampleGetExercise_nonOwnerPublic`
- `ExampleGetExercise_nonOwnerPrivate`

### Request payloads as JSON strings

Use `doRaw()` with raw JSON strings — NOT `doJSON()` with Go maps. This makes the payloads readable as actual API examples:

```go
doRaw(r, "POST", "/api/exercises", `{
    "name": "Bench Press",
    "templateId": "bench-press",
    "type": "STRENGTH",
    "technicalDifficulty": "beginner",
    "bodyWeightScaling": 0.5,
    "description": "Barbell bench press",
    "version": 0
}`)
```

For non-owner requests, use `doRawAs()`:
```go
w := doRawAs(r, "GET", "/api/exercises/1", "", "other")
```

For GET requests as the default user, `doJSON()` with nil body is fine:
```go
w := doJSON(r, "GET", "/api/exercises/1", nil)
```

### Output assertions

Print only stable, meaningful fields. Avoid timestamps and auto-increment IDs where possible (though IDs are stable within a single example since each starts with a fresh DB).

```go
var exercise models.Exercise
json.Unmarshal(w.Body.Bytes(), &exercise)
fmt.Println(w.Code)
fmt.Println(exercise.Name)
fmt.Println(exercise.Owner)
// Output:
// 200
// Bench Press
// testuser
```

### Standard permission scenarios

For resources with public/private visibility (exercises, equipment):
1. **Owner** — full access (200 with data)
2. **Non-owner + public** — read access (200 with data for GET/List, 403 for mutations)
3. **Non-owner + private** — denied (403 for GET, empty list for List, 403 for mutations)

For user-scoped resources without public visibility (workouts, sections):
1. **Owner** — full access (200)
2. **Non-owner** — denied (403 for direct access, empty list for List)

### Documenting relationships

When an endpoint requires another resource to exist first, show the creation chain:

```go
// Creating a section exercise requires a workout, a section, and an exercise
// scheme. The full hierarchy: Workout → Section → SectionExercise → ExerciseScheme.
func ExampleCreateWorkoutSectionExercise() {
    setupExampleDB()
    r := newRouter()

    // Create the exercise and scheme first.
    doRaw(r, "POST", "/api/exercises", `{...}`)
    doRaw(r, "POST", "/api/exercise-schemes", `{...}`)

    // Create a workout and section.
    doRaw(r, "POST", "/api/user/workouts", `{...}`)
    doRaw(r, "POST", "/api/user/workout-sections", `{...}`)

    // Add the exercise scheme to the section.
    w := doRaw(r, "POST", "/api/user/workout-section-exercises", `{...}`)
    // ...assertions...
}
```

If setup gets repetitive, extract a helper (e.g., `createExerciseSchemeForExample(r)`).

## Rendering

- `make docs` generates markdown into `docs/generated/` via `gomarkdoc` (`scripts/docs.sh` auto-discovers all packages with `doc.go` files)
- Generated markdown is **committed to VCS** — changes to doc comments or examples produce visible diffs in `docs/generated/*.md`
- Generated markdown is **embedded into the binary** via `//go:embed docs/generated/*` in `main.go` and **served as HTML** at `/docs` using goldmark for markdown→HTML conversion
- The Dockerfile runs `make docs` before building, so the production image always serves up-to-date docs
- After changing docs or examples, always run `make docs` and commit the regenerated files alongside the code changes

## When behavior changes break examples

This is the point — if a handler's behavior changes, the Example test fails because the `// Output:` no longer matches. Fix the example to match the new behavior, which updates the documentation simultaneously.

## Checklist for adding docs to a new handler package

1. Create `doc.go` with package overview, entity relationships, cross-package links
2. Add doc comments to every exported handler function (endpoint, method, permissions, relationships)
3. Create `example_setup_test.go` with `setupExampleDB()` and `doRawAs()`
4. Create `example_<resource>_test.go` for each resource group
5. Write permission examples (owner / non-owner scenarios)
6. Write CRUD flow examples showing creation chains and data assertions
7. Run `go test -v -run Example ./<package>/` to verify
8. Run `make test-go` to verify nothing broke
9. Run `make docs` to regenerate — check `/docs` renders correctly
10. Commit the regenerated `docs/generated/*.md` files alongside your code changes
