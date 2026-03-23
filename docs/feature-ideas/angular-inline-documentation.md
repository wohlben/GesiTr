# Angular Inline Documentation

## Summary

Adopt Compodoc to generate browsable HTML documentation from TSDoc comments in Angular components, services, and API clients. Focus on documenting user flows, component responsibilities, and the relationships between route-level components, API clients, and backend endpoints — directly in the code.

## Motivation

The Angular frontend has growing complexity: route-level components orchestrate queries, mutations, and permission checks; API clients map to backend endpoints; UI components handle display logic. Understanding how a user flow threads through these layers requires reading multiple files and mentally connecting them.

Compodoc reads TSDoc comments and Angular decorators to generate a browsable site with component trees, dependency graphs, and route maps — all derived from the code. When the code changes, regenerating the docs picks up the changes automatically. This pairs with the [Go inline documentation](go-inline-documentation.md) effort to keep both halves of the stack self-documenting.

## Requirements

- None — Compodoc works with any Angular project out of the box.

## Design Sketch

### What to document

The highest-value targets are **user flows** — the path from a route through its component, API client, and backend endpoint. These are hard to discover from code alone because they span multiple files.

**Route-level components** — explain which data the component fetches, what actions it exposes, and how it relates to sibling routes:

```typescript
/**
 * Exercise detail page at `/compendium/exercises/:id/:slug`.
 *
 * ## Data
 *
 * Fetches the exercise entity, its version history, and the user's
 * exercise list (to detect "already added" state) in parallel via
 * TanStack Query.
 *
 * ## Actions
 *
 * - **Add to mine**: clones the public exercise into the user's collection
 *   (strips id, slug, owner; sets public=false). Navigates to the new
 *   user exercise on success.
 * - **Edit**: navigates to {@link ExerciseEdit} (requires MODIFY permission).
 * - **Delete**: confirmation dialog, then deletes and navigates to list.
 *
 * ## Related routes
 *
 * - {@link ExerciseEdit} — edit form for this exercise
 * - {@link ExerciseHistory} — version history view
 * - {@link ExerciseList} — parent list page
 */
@Component({ ... })
export class ExerciseDetail { }
```

**API clients** — document the mapping between client methods and backend endpoints:

```typescript
/**
 * API client for compendium (public) exercise and equipment endpoints.
 *
 * All methods target `/api/exercises` or `/api/equipment`.
 * For user-scoped equivalents, see {@link UserApiClient}.
 */
@Injectable({ providedIn: 'root' })
export class CompendiumApiClient { }
```

**UI components** — document the contract (inputs/outputs) and when to use them vs alternatives:

```typescript
/**
 * Generic confirmation dialog with pending state.
 *
 * Used by detail pages for destructive actions (delete).
 * Emits `confirmed` or `cancelled` — the parent owns the mutation.
 */
@Component({ ... })
export class ConfirmDialog { }
```

### Setup

Install and add a make target:

```bash
npm install --save-dev @compodoc/compodoc
```

```makefile
docs-web:
	cd web && npx @compodoc/compodoc -p tsconfig.json -s --port 8081
```

### Where documentation lives

Compodoc pulls documentation from three sources, each serving a different purpose:

1. **TSDoc on components/services** — the primary documentation surface. Compodoc renders these as full pages with the class API. This is where flow descriptions, action lists, and `{@link}` cross-references to related components go. Every route-level component, API client, and shared UI component should have a TSDoc block.

2. **Folder-level `README.md` files** — Compodoc picks up `README.md` in any source folder as "additional documentation" pages. Use these for narrative overviews that tie multiple routes together (e.g., `src/app/pages/compendium/README.md` explaining the compendium exercise flow from list → detail → edit/history). These are the closest equivalent to ExDoc guides.

3. **Route graph** — auto-discovered from `app.routes.ts`, but **visual only**. You cannot attach prose to individual routes in the graph. It's a supplement for orientation, not the primary docs. The real route documentation lives in the component TSDoc (point 1) and folder READMEs (point 2).

### Other Compodoc features

- **Component tree**: shows parent/child relationships and dependency injection.
- **Search**: full-text search across all documented symbols.
- **Coverage report**: `npx @compodoc/compodoc -p tsconfig.json --coverageTest 50` — can enforce minimum documentation coverage in CI.

### Rollout order

1. Add Compodoc to `devDependencies` and the `make docs-web` target.
2. Document the exercise detail flow end-to-end (ExerciseDetail, ExerciseEdit, ExerciseHistory, CompendiumApiClient) as the template for other flows.
3. Document API clients — these are the bridge between frontend and backend, highest leverage for understanding the system.
4. Expand to other route-level components, prioritizing flows with non-obvious behavior (workout logging, exercise group membership).
5. Optionally add coverage threshold to CI once enough components are documented.

## Considerations

- **TSDoc vs JSDoc**: Compodoc supports both. TSDoc (`@link`, `@see`, `@param`) is preferred for consistency with TypeScript tooling.
- **No executable examples**: unlike Go's `Example` functions or Phoenix doctests, Compodoc comments are purely descriptive. The existing screenshot tests and unit tests fill the "executable documentation" role — Compodoc links to test files when they follow Angular naming conventions.
- **Maintenance cost**: comments can drift from code. The coverage report helps catch undocumented new components, but stale comments require code review discipline. Keeping docs focused on *why* and *flow* rather than *what* (which the code already shows) reduces drift.
- **Unified docs target**: a top-level `make docs` could run both pkgsite (Go) and Compodoc (Angular) and link them together via a simple index page.
