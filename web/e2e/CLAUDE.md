# E2E Test Conventions

## Screenshot Directory Structure

Screenshots are centralized in `e2e/__screenshots__/` using the `snapshotPathTemplate` in `playwright.config.ts`.

**Do NOT** create co-located `*.spec.ts-snapshots/` directories. All screenshots go through the centralized path.

### `toHaveScreenshot()` arg format

**MUST use array syntax** — passing a string with `/` will get sanitized by Playwright (slashes replaced with dashes). Only array args preserve the directory structure.

```typescript
await expect(page).toHaveScreenshot([viewport.name, 'light', 'compendium', 'exercises.png'], { fullPage: true });
```

**MUST pass `{ fullPage: true }`** as the second argument to every `toHaveScreenshot()` call. This captures the entire scrollable page, not just the viewport. This option cannot be set in `playwright.config.ts` — it must be passed per-call.

- `viewport.name` = `desktop` or `mobile`
- Next element = `light` or `dark` (theme)
- Remaining elements = route segments, with the last one including `.png`
- Use `[id]` as placeholder for dynamic route segments

### Examples

| Route | Screenshot arg |
|-------|---------------|
| `/compendium/exercises` | `[viewport.name, 'light', 'compendium', 'exercises.png']` |
| `/compendium/exercises/:id` | `[viewport.name, 'dark', 'compendium', 'exercises', '[id].png']` |
| `/compendium/exercises/:id/edit` | `[viewport.name, 'light', 'compendium', 'exercises', '[id]', 'edit.png']` |
| `/user/exercises` | `[viewport.name, 'light', 'user', 'exercises.png']` |

### Resulting file path

```
e2e/__screenshots__/{viewport}/{theme}/{route}/{projectName}-{platform}.png
```

Example: `e2e/__screenshots__/desktop/light/compendium/exercises/chromium-linux.png`

## E2E Change Verification Workflow

**Every time you modify e2e test files**, you MUST follow this full workflow:

1. **Update screenshots locally**: `make update-screenshots-e2e` (from project root — starts API server automatically)
2. **Visually verify** the updated screenshots look correct
3. **Run the Docker pipeline**: `docker build -t gesitr .` (from project root) to confirm tests pass in CI-like environment

To run e2e tests without updating screenshots: `make test-e2e` (also starts the API server automatically).

**NEVER update screenshots from Docker.** The Docker pipeline is strictly for verification — it ensures the locally-recorded screenshots are correct and match what CI would produce. Always record screenshots against the local dev server.

**Note:** Both `make test-e2e` and `make update-screenshots-e2e` are fully self-contained and safe to run alongside `make dev`. They use dedicated ports (API on :9876, Angular on :4300) and an isolated database (`gesitr-e2e.db`). You do NOT need to start any servers manually.

## Data Values

When creating test data via API helpers, use the **exact casing** the backend expects. The Go backend uses lowercase enum values (e.g., `'free_weights'`, not `'FREE_WEIGHTS'`). Check `internal/compendium/models/enums.go` for canonical values.
