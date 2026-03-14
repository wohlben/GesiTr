# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GesiTr is a resistance training tracker with a Go backend (Gin + GORM + SQLite) and Angular frontend. The Go binary embeds the compiled Angular SPA and serves it as a single binary.

Two data domains:
- **Compendium** — public, curated exercise/equipment library at `/api/exercises`, `/api/equipment`, etc.
- **User** — private per-user entities (exercises, equipment, workouts, logs, records) at `/api/user/*`

## Development Commands

```bash
# Start both API (port 8080) and web dev server (port 4200) — single process
make dev        # Starts Go API, waits for ready, then starts ng serve
# Also works from web/: npm start

# Or run them separately:
make dev-api    # DEV=true AUTH_FALLBACK_USER=anon go run .
make dev-web    # cd web && npx ng serve (Angular only, no API)

# Regenerate TypeScript types from Go structs (tygo)
make generate

# Lint everything (Go format check + ESLint + Prettier)
make lint

# Run all tests (lint + Go tests + web unit tests + e2e)
make test

# Individual test targets
make test-go                    # go test ./...
make test-web                   # cd web && npm test (vitest, excludes screenshot tests)
cd web && npx ng run web:test-screenshot  # Unit screenshot tests only
make test-e2e                   # Builds Go, starts API on :9876, runs Playwright on :4300, cleans up

# Seed database with compendium data from data/ directory
make seed       # deletes gesitr.db, re-seeds

# Production Docker build (multi-stage: builds, tests, then creates runtime image)
make docker
```

## Architecture

### Backend (Go)

- **Entry point**: `main.go` — sets up Gin routes, GORM auto-migration, embeds SPA via `//go:embed`
- **`internal/compendium/`** — models (GORM entities + API response types) and handlers for public compendium
- **`internal/user/`** — models and handlers for per-user entities, scoped by `UserID` from auth middleware
- **`internal/auth/`** — middleware extracting UserID (falls back to `AUTH_FALLBACK_USER` env var in dev)
- **`internal/database/`** — SQLite init via GORM
- **`cmd/seed/`** — database seeder loading CSV/JSON from `data/`
- **DEV mode** (`DEV=true`): exposes `POST /api/ci/reset-db` to wipe all tables (used by e2e tests)
- **`DATABASE_PATH`** env var overrides the default `gesitr.db` (used by `make test-e2e` to isolate test data)

### Frontend (Angular)

Path aliases defined in `tsconfig.json`:
- `$core/*` → `src/app/core/*` — API clients, query keys
- `$features/*` → `src/app/features/*` — feature pages (compendium/, user/)
- `$ui/*` → `src/app/ui/*` — reusable components
- `$generated/*` → `src/app/generated/*` — auto-generated TypeScript types from Go

Key patterns:
- **TanStack Angular Query** for server state management (caching, sync)
- **Lazy-loaded routes** via `loadComponent()` in `app.routes.ts`
- **Tailwind CSS v4** for styling, SCSS for component styles
- Components are standalone (no NgModules)

### Type Generation (Tygo)

`tygo.yaml` maps Go struct packages → TypeScript interfaces. Entity files are excluded (only API response/request types are generated). Run `make generate` after changing Go model structs.

Output files:
- `web/src/app/generated/models.ts` — compendium types
- `web/src/app/generated/user-models.ts` — user types

## Testing

### E2E Tests (Playwright)

- **Run**: `make test-e2e` — fully self-contained, safe to run alongside `make dev`
  - Builds Go binary, starts API on :9876 with isolated `gesitr-e2e.db`, Angular on :4300
  - Uses `ng e2e --configuration=isolated` (see `angular.json` serve:e2e + `proxy.e2e.conf.json`)
  - Cleans up all processes and temp files automatically
- Config: `web/playwright.config.ts`
- Workers: 1 (sequential — tests share database state)
- Two projects: `chromium` (compendium), `chromium-user` (user routes, depends on chromium)
- **Screenshot conventions**: see `web/e2e/CLAUDE.md` for full details
  - 4 variants per route: desktop-light, desktop-dark, mobile-light, mobile-dark
  - Desktop: 1280x720, Mobile: 375x667
  - **Must use array syntax**: `[viewport.name, 'light', 'compendium', 'exercises.png']`
  - **Must pass `{ fullPage: true }`** to every `toHaveScreenshot()` call
  - Screenshots centralized in `e2e/__screenshots__/`
- **Update workflow**: `make update-screenshots-e2e` (starts API automatically) → verify visually → `docker build -t gesitr .`
- Never update screenshots from Docker — Docker is verification only

### Unit Tests (Vitest + Browser Mode)

- `*.spec.ts` — standard unit tests, run via `npm test` / `make test-web`
- `*.screenshot.spec.ts` — unit screenshot tests, run separately via `npx ng run web:test-screenshot`
- Screenshot tests use the same 4-variant pattern (desktop/mobile + light/dark)

### Go Tests

- `go test ./...` — standard Go testing

## Conventions

- **Commits**: conventional commits (`feat:`, `fix:`, `chore:`, `refactor:`, `docs:`)
- **Enum values**: Go backend uses lowercase (e.g., `'free_weights'`, not `'FREE_WEIGHTS'`). Check `internal/compendium/models/enums.go` for canonical values.
- **Issue tracking**: uses `bd` (beads) — see `AGENTS.md` for workflow
