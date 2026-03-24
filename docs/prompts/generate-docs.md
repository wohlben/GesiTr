# Prompt: Generate API Documentation Tests

Copy-paste the following prompt, replacing the placeholders, to generate executable documentation for a handler package.

---

I want to add Go Example test documentation for the `<HANDLER_PACKAGE>` handlers at `internal/<PATH>/handlers/`.

Follow the conventions in `docs/contexts/api-docs-with-examples.md`. Use the exercise handler examples as reference (`internal/exercise/handlers/example_*_test.go`).

**What to create:**

1. `doc.go` — package overview explaining entity relationships, endpoint hierarchy, and cross-package links
2. Doc comments on every exported handler function (endpoint, HTTP method/path, permissions, relationships)
3. `example_setup_test.go` — `setupExampleDB()`, `doRawAs()`, and `newExampleRouter()` if cross-API flows are needed
4. `example_<resource>_test.go` files — one per logical resource group, with:
   - Permission examples (owner / non-owner scenarios appropriate for this resource)
   - CRUD flow examples showing creation chains and meaningful data assertions
   - Request payloads as raw JSON strings via `doRaw()`

**The handlers to document are:**

- `<LIST THE ENDPOINTS, e.g.:`
- `GET /api/... — list`
- `POST /api/... — create`
- `GET /api/.../:id — get`
- `PUT /api/.../:id — update`
- `DELETE /api/.../:id — delete>`

**Permission model for this resource:**

- `<DESCRIBE, e.g.: "owner-only, no public visibility" or "owner gets full access, non-owner gets READ on public, 403 on private">`

**Relationships to other resources:**

- `<DESCRIBE, e.g.: "requires an exercise to exist (exerciseId field)" or "sections belong to workouts">`

After writing the examples, run `go test -v -run Example ./<package>/` and `make test-go`. Then run `make docs` and commit the regenerated `docs/generated/*.md` files alongside the code. Run `docker build -t gesitr .` before committing.
