# E2E Test Conventions

## Running Tests

To run e2e tests: `make test-e2e` (from project root -- starts API server automatically).

**Note:** `make test-e2e` is fully self-contained and safe to run alongside `make dev`. It uses dedicated ports (API on :9876, Angular on :4300) and an isolated in-memory database. You do NOT need to start any servers manually.

## Data Values

When creating test data via API helpers, use the **exact casing** the backend expects. The Go backend uses lowercase enum values for some types (e.g., `'free_weights'`, `'beginner'`) and uppercase for others (e.g., `'STRENGTH'`, `'CHEST'`). Check `internal/exercise/models/exercise.go` and `internal/equipment/models/equipment.go` for canonical values.
