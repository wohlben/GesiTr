# GesiTr - Resistance Training Tracker

A resistance training tracker with a Go backend and Angular frontend, deployed as a single binary with the SPA embedded.

## Data Model

### Compendium (Shared Entities)

Exercises, equipment, and workouts are global entities with an **owner** and a **public** visibility toggle. The "compendium" is a filtered view over all public entities — a curated, community-accessible library:

- **Exercises** — resistance training exercises (e.g. bench press, squat, treadmill running), with relationships (similar, variation, progression/regression) and exercise schemes
- **Equipment** — available equipment (e.g. barbells, cable machines, treadmills), with fulfillment relationships between equipment types
- **Workouts** — workout templates composed of sections and exercise groups
- **Workout Groups** — shared groupings of workouts with memberships
- **Localities** — per-user equipment locations with availability tracking

Any user can create their own exercises, equipment, or workouts. Private entities are only visible to the owner; public ones appear in the compendium for everyone. This means a user who creates a custom exercise can keep it private or share it — without duplicating data between "personal" and "compendium" copies.

### User-Scoped Entities

These are inherently per-user and only readable by the owning user:

- **Workout Logs** — a logged session instance of a workout, with sections, exercises, and sets with target values
- **Exercise Logs** — records of what the user actually performed (reps, weight, duration, distance). Created from a workout log or ad-hoc
- **Workout Schedules** — flexible scheduling: fixed days (e.g. "Mondays and Thursdays") or frequency-based (e.g. "three times per week")
- **Mastery** — automatic tracking of exercise expertise via contributions and experience
- **Name Preferences** — user's preferred alternative names for exercises
- **Personal Records** — automatically derived from exercise logs. Best performance per exercise and measurement type:
  - **Rep-based / AMRAP**: estimated 1RM via the Brzycki formula
  - **Time-based / EMOM / Rounds for Time**: raw duration
  - **Distance-based**: raw distance

## Tech Stack

- **Backend:** Go (Gin + GORM + SQLite), Huma v2 for API docs
- **Frontend:** Angular, TanStack Query, Tailwind CSS v4, Spartan-NG
- **Type generation:** Tygo (Go structs to TypeScript interfaces)
- **Testing:** Vitest (unit), Playwright (e2e), `go test` (backend)
- **Deployment:** Single binary with embedded SPA, Docker multi-stage build

## Development

### Prerequisites

- Go 1.25+
- Node.js 24+ / npm

### Quick Start

```bash
make dev        # Starts Go API (port 8080) + Angular dev server (port 4200)
```

Or run them separately:

```bash
make dev-api    # Go backend only (DEV=true, fallback user "anon")
make dev-web    # Angular dev server only
```

### Makefile Targets

| Target | Description |
|---|---|
| `make dev` | Start API + Angular dev server (single process) |
| `make dev-api` | Go backend in dev mode |
| `make dev-web` | Angular dev server only |
| `make build` | Full production build (Angular + Go binary with embedded SPA) |
| `make generate` | Regenerate TypeScript types from Go structs (tygo) |
| `make lint` | Lint everything (Go format + ESLint + Prettier) |
| `make test` | All tests (lint + Go + web unit + e2e) |
| `make test-go` | Go tests only |
| `make test-web` | Web unit tests only (vitest) |
| `make test-e2e` | Playwright e2e (isolated ports and database) |
| `make seed` | Reset and seed database with compendium data from `data/` |
| `make docker` | Multi-stage Docker build (build, test, runtime image) |
| `make clean` | Remove built artifacts |

### Auth

The backend reads a user ID from a configurable header (`AUTH_HEADER`, default `X-User-Id`). In dev mode, it falls back to `AUTH_FALLBACK_USER`. No session or JWT management — identity is provided by an upstream proxy or SSO.

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_PATH` | `gesitr.db` | SQLite file path |
| `PORT` | `8080` | Server port |
| `DEV` | — | Enables dev mode (reset-db endpoint, fallback user) |
| `AUTH_HEADER` | `X-User-Id` | Header for user identity |
| `AUTH_FALLBACK_USER` | — | Fallback user when header is missing (dev only) |
