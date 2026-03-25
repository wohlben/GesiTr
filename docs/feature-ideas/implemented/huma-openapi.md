# Huma OpenAPI Migration

## Summary

Replaced Gin's `func(*gin.Context)` handlers with [huma v2](https://github.com/danielgtaylor/huma) typed handlers to auto-generate an OpenAPI spec from Go types. Huma works as a Gin adapter — no full framework switch was needed.

## Motivation

Gin handlers erased all type information behind `func(*gin.Context)`. This meant:

- No way to auto-generate API documentation from the code
- Request/response contracts lived only in handler implementations, not in a machine-readable spec
- Frontend TypeScript types (via tygo) were generated from Go models, but there was no spec describing which endpoints accept/return which types

Huma uses Go generics to reflect input/output structs into a full OpenAPI 3.1 spec. The types *are* the spec — there's nothing to keep in sync.

## What was implemented

### Handler pattern

All handlers use typed input/output structs:

```go
type CreateExerciseInput struct {
    Body models.Exercise
}

type CreateExerciseOutput struct {
    Body models.Exercise
}

func CreateExercise(ctx context.Context, input *CreateExerciseInput) (*CreateExerciseOutput, error) {
    // ... business logic ...
    return &CreateExerciseOutput{Body: resultDTO}, nil
}
```

Path/query parameters are struct fields with tags:

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

Each resource has a `routes.go` with huma operations:

```go
huma.Register(humaAPI, huma.Operation{
    OperationID: "CreateExercise",
    Method:      http.MethodPost,
    Path:        "/exercises",
    Tags:        []string{"exercises"},
}, exerciseHandlers.CreateExercise)
```

### OpenAPI docs UI

Huma serves the OpenAPI spec at `/api/openapi.json` and a docs UI (Stoplight Elements) at `/api/docs`. The `x-user-id` header is exposed as a security scheme so it appears as a fillable field in the docs UI.

### Auth middleware

`humagin.NewWithGroup` lets huma routes share a Gin `RouterGroup`, so Gin middleware (auth, profile) applies automatically. The Gin auth middleware stores the user ID in both `*gin.Context` and `context.Context`. Huma handlers read it with `humaconfig.GetUserID(ctx)`.

### Infrastructure

Reusable infrastructure in `internal/humaconfig/`:
- `config.go` — API setup, server URL, security schemes
- `auth.go` — context key + `GetUserID`
- `pagination.go` — `PaginationInput` with huma tags
- `responses.go` — `PaginatedBody[T]` generic

## What was learned

- **Paths are relative to the group**: huma operations use `/exercises/{id}`, not `/api/exercises/{id}`. The OpenAPI spec uses `servers: [{url: "/api"}]`.
- **Error format**: adopted huma's native RFC 7807 error model. The Angular frontend doesn't parse error response bodies (it uses `HttpErrorResponse.message`), so this was transparent.
- **Tests**: existing `httptest`-based tests work unchanged since huma routes still go through Gin's engine. Only test setup functions needed updating.

## Scope

All 84 handlers across every resource group were migrated: exercises, equipment, fulfillments, relationships, exercise-groups, profiles, workouts, workout-logs, and exercise-logs.
