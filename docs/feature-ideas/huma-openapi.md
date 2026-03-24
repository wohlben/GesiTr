# Huma OpenAPI Migration

## Summary

Replace Gin's `func(*gin.Context)` handlers with [huma v2](https://github.com/danielgtaylor/huma) typed handlers to auto-generate an OpenAPI spec from Go types. Huma works as a Gin adapter — no full framework switch needed.

## Motivation

Gin handlers erase all type information behind `func(*gin.Context)`. This means:

- No way to auto-generate API documentation from the code
- Request/response contracts live only in handler implementations, not in a machine-readable spec
- Frontend TypeScript types (via tygo) are generated from Go models, but there's no spec describing which endpoints accept/return which types

Huma uses Go generics to reflect input/output structs into a full OpenAPI 3.1 spec. The types *are* the spec — there's nothing to keep in sync.

## Design Sketch

### Handler refactor pattern

Current Gin handler:

```go
func CreateExercise(c *gin.Context) {
    var dto models.Exercise
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // ... business logic ...
    c.JSON(http.StatusCreated, resultDTO)
}
```

Huma equivalent:

```go
type CreateExerciseInput struct {
    Body models.Exercise `json:"body"`
}

type CreateExerciseOutput struct {
    Body models.Exercise `json:"body"`
}

func CreateExercise(ctx context.Context, input *CreateExerciseInput) (*CreateExerciseOutput, error) {
    // ... same business logic ...
    return &CreateExerciseOutput{Body: resultDTO}, nil
}
```

Path/query parameters become struct fields with tags:

```go
type ListExercisesInput struct {
    Owner  string `query:"owner" doc:"Filter by owner (use 'me' for current user)"`
    Limit  int    `query:"limit" default:"20"`
    Offset int    `query:"offset" default:"0"`
}

type GetExerciseInput struct {
    ID uint `path:"id"`
}
```

### Route registration

```go
humaAPI := humaginadapter.New(ginEngine, huma.DefaultConfig("GesiTr API", "1.0.0"))

huma.Register(humaAPI, huma.Operation{
    OperationID: "create-exercise",
    Method:      http.MethodPost,
    Path:        "/api/exercises",
    Tags:        []string{"exercises"},
}, exerciseHandlers.CreateExercise)
```

### Auth middleware

Huma supports middleware via resolvers. The current `auth.GetUserID(c)` pattern would become a resolver that injects the user ID into handler inputs:

```go
type AuthInput struct {
    UserID string `header:"X-User-Id" doc:"Authenticated user ID"`
}

// Embedded in handler inputs that need auth:
type CreateExerciseInput struct {
    AuthInput
    Body models.Exercise `json:"body"`
}
```

### Swagger UI

Huma serves the OpenAPI spec at `/openapi.json` and can serve Swagger UI with no additional setup. Restrict to dev mode with a conditional:

```go
if os.Getenv("DEV") == "true" {
    // Swagger UI served automatically by huma
}
```

### Incremental migration

Huma's Gin adapter allows both Gin handlers and huma handlers on the same router. Migrate one resource group at a time:

1. ~~exercises~~ ✓
2. ~~equipment~~ ✓
3. ~~fulfillments, relationships~~ ✓
4. ~~exercise-groups~~ ✓
5. ~~user/* (workouts, workout-logs, exercise-logs)~~ ✓

### Rollout order

1. ~~Add `huma` dependency, set up the adapter alongside existing Gin router.~~ ✓
2. ~~Migrate `exercises` as the pilot — it's the most representative CRUD resource.~~ ✓
3. ~~Verify the generated OpenAPI spec matches the actual API behavior.~~ ✓
4. ~~Migrate remaining resources group by group.~~ ✓
5. Add `make docs-api` target to serve Swagger UI locally.

## What was learned during the pilot

- **`humagin.NewWithGroup`** lets huma routes share a Gin RouterGroup, so Gin middleware (auth, profile) applies automatically. No need for separate huma middleware.
- **Paths are relative to the group**: huma operations use `/exercises/{id}`, not `/api/exercises/{id}`. The OpenAPI spec uses `servers: [{url: "/api"}]`.
- **Auth context propagation**: the Gin auth middleware stores the user ID in both `*gin.Context` (for non-migrated handlers) and `context.Context` (for huma handlers via `context.WithValue`). Huma handlers read it with `humaconfig.GetUserID(ctx)`.
- **`RawBody []byte`** skips huma's automatic body validation for create/update inputs. The Exercise DTO is shared between request and response and has server-set fields (`id`, `createdAt`, etc.) that aren't present in create requests. Using `RawBody` preserves the current behavior (no validation beyond JSON parsing).
- **Error format**: adopted huma's native RFC 7807 error model. The Angular frontend doesn't parse error response bodies (it uses `HttpErrorResponse.message` from Angular), so this change is transparent.
- **Tests**: existing `httptest`-based tests work unchanged since huma routes still go through Gin's engine. Only test setup functions needed updating (register via `exercisehandlers.RegisterRoutes(humaAPI)` instead of individual Gin routes).
- **Infrastructure in `internal/humaconfig/`**: `config.go` (API setup), `auth.go` (context key + `GetUserID`), `pagination.go` (`PaginationInput` with huma tags), `responses.go` (`PaginatedBody[T]` generic). Reusable for all subsequent migrations.

## Considerations

- **Scope**: All 84 handlers migrated. Each follows the same mechanical pattern (extract params into input struct, return output instead of calling `c.JSON`).
- **Complements doc.go**: huma documents the HTTP API contract (endpoints, schemas). `doc.go` documents domain concepts (entity model, relationships, visibility rules). They serve different audiences and don't overlap.
- **Testing**: handler tests use `httptest` with Gin's test router. This continues to work with huma since routes go through the same Gin engine.
