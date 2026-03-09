# E2E Test Conventions

## Screenshot Directory Structure

Screenshots are centralized in `e2e/__screenshots__/` using the `snapshotPathTemplate` in `playwright.config.ts`.

**Do NOT** create co-located `*.spec.ts-snapshots/` directories. All screenshots go through the centralized path.

### `toHaveScreenshot()` arg format

**MUST use array syntax** — passing a string with `/` will get sanitized by Playwright (slashes replaced with dashes). Only array args preserve the directory structure.

```typescript
await expect(page).toHaveScreenshot([viewport.name, 'light', 'compendium', 'exercises.png']);
```

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

1. **Start the API server**: `DEV=true AUTH_FALLBACK_USER=anon go run .` (from project root), or use `make dev` which also starts the web server
2. **Update screenshots locally**: `cd web && npx ng e2e --update-snapshots` (the `ng e2e` command starts the Angular dev server automatically)
3. **Visually verify** the updated screenshots look correct
4. **Run the Docker pipeline**: `docker build -t gesitr .` (from project root) to confirm tests pass in CI-like environment

**NEVER update screenshots from Docker.** The Docker pipeline is strictly for verification — it ensures the locally-recorded screenshots are correct and match what CI would produce. Always record screenshots against the local dev server.

## Data Values

When creating test data via API helpers, use the **exact casing** the backend expects. The Go backend uses lowercase enum values (e.g., `'free_weights'`, not `'FREE_WEIGHTS'`). Check `internal/compendium/models/enums.go` for canonical values.
