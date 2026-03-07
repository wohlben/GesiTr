# GesiTr - Resistance Training Tracker

**Portmanteau of Golang Resistance Training [Tracker]**

A resistance training tracker with a Go backend for persistence and an Angular frontend for the GUI.

## Project Goals

### Compendium (Public)

A curated, publicly accessible library of training-related entities:

- **Exercises** - catalog of resistance training exercises (e.g. bench press, squat, treadmill running)
- **Equipment** - available equipment (e.g. barbells, dumbbells, cable machines, treadmills)
- **Workouts** - predefined workout templates composed of exercises

The compendium serves as a starting point to make it easy for users to set up their personal training configuration without building everything from scratch.

### User Entities (Private, per-user)

Each user maintains their own private set of entities, only readable by the owning user. These mirror the compendium structure and can be easily imported from it:

- **Exercises** - the user's personal exercise catalog (imported from compendium or custom)
- **Equipment** - equipment the user has available
- **Workouts** - the user's configured workout routines
  - **Goals** - target parameters per exercise within a workout (e.g. "exercise Y: 7 kg, 10 reps")
  - **Logs** - actual achieved results per session (e.g. "exercise Y: 5 kg, 7 reps")
  - **Scheduling** - workouts can be scheduled in flexible ways:
    - **Fixed days** - e.g. "on Mondays and Thursdays"
    - **Frequency-based** - e.g. "three times per week"
- **Personal Records** - automatic tracking of highest scores per exercise, covering metrics like:
  - Maximum weight per rep count
  - Time spent
  - Speed / incline (for running, walking, etc.)

## Tech Stack

- **Backend:** Go
- **Frontend:** Angular (in `web/`)

## Development

### Prerequisites

- Go
- Node.js / npm
- Dolt (`pacman -S dolt`)

### Makefile Targets

| Target | Description |
|---|---|
| `make build` | Full production build: Angular frontend, then Go binary with embedded SPA |
| `make build-web` | Build the Angular frontend only |
| `make build-go` | Build the Go binary (depends on `build-web`) |
| `make dev-api` | Run the Go backend in development mode (`DEV=true`) |
| `make dev-web` | Run the Angular dev server (`ng serve`) |
| `make dolt` | Start the Dolt SQL server for beads issue tracking |
| `make docker` | Build a Docker image |
| `make clean` | Remove built artifacts (`gesitr`, `gesitr.db`, `web/dist/`) |

For development, run `make dev-api` and `make dev-web` in separate terminals.

### Frontend Scripts (`web/`)

| Script | Description |
|---|---|
| `npm start` | Start the Angular dev server |
| `npm run build` | Production build |
| `npm run watch` | Build in watch mode (development) |
| `npm test` | Run tests |
