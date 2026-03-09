# E2E Test Conventions

## Screenshot Directory Structure

Screenshots are centralized in `e2e/__screenshots__/` using the `snapshotPathTemplate` in `playwright.config.ts`.

**Do NOT** create co-located `*.spec.ts-snapshots/` directories. All screenshots go through the centralized path.

### `toHaveScreenshot()` arg format

```typescript
await expect(page).toHaveScreenshot(`${viewport.name}/${theme}/${route}.png`);
```

- `viewport.name` = `desktop` or `mobile`
- `theme` = `light` or `dark`
- `route` = the app route with `[id]` as placeholder for dynamic IDs

### Examples

| Route | Screenshot arg |
|-------|---------------|
| `/compendium/exercises` | `desktop/light/compendium/exercises.png` |
| `/compendium/exercises/:id` | `desktop/dark/compendium/exercises/[id].png` |
| `/compendium/exercises/:id/edit` | `mobile/light/compendium/exercises/[id]/edit.png` |
| `/user/exercises` | `desktop/light/user/exercises.png` |

### Resulting file path

```
e2e/__screenshots__/{viewport}/{theme}/{route}/{projectName}-{platform}.png
```

Example: `e2e/__screenshots__/desktop/light/compendium/exercises/chromium-linux.png`
