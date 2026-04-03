# Compendium Locality

## Summary

Per-user curated lists of equipment available at specific training locations. A user defines multiple localities (e.g. "Home Gym", "Commercial Gym") and lists which equipment is present at each. The exercise list can be filtered by locality, showing only exercises whose required equipment is available — respecting fulfillments (e.g. an adjustable bench fulfilling a flat bench).

Conceptually distinct from fulfillments: a fulfillment says "equipment A can stand in for equipment B" (substitution), a locality says "at this location, I have these pieces of equipment" (availability).

## Data Model

### Locality (`localities` table)

Compendium-style entity with owner + public toggle.

| Column     | Type      | Notes                             |
|-----------|-----------|-----------------------------------|
| id        | uint      | PK                                |
| name      | string    | e.g. "Home Gym"                   |
| owner     | string    | user ID                           |
| public    | bool      | default false (sharing out of scope for now) |
| created_at| timestamp |                                   |
| updated_at| timestamp |                                   |

### Locality Availability (`locality_availabilities` table)

Join table between locality and equipment with an availability toggle.

| Column       | Type      | Notes                                            |
|-------------|-----------|--------------------------------------------------|
| id          | uint      | PK                                               |
| locality_id | uint      | FK, unique with equipment_id                     |
| equipment_id| uint      | FK, unique with locality_id                      |
| available   | bool      | default true — toggle off for e.g. maintenance   |
| owner       | string    | denormalized from locality for access checks     |
| created_at  | timestamp |                                                  |
| updated_at  | timestamp |                                                  |

## API

```
GET/POST       /localities
GET/PUT/DELETE /localities/{id}
GET            /localities/{id}/permissions
GET/POST       /locality-availabilities
PUT/DELETE     /locality-availabilities/{id}
GET            /exercises?localityId={id}     — fulfillment-aware filter
```

## Exercise Filtering

`GET /exercises?localityId=X` shows exercises where ALL required equipment is in the locality's effective set. The effective set is:
1. Equipment directly available at the locality (`available = true`)
2. Equipment fulfilled by available equipment (via `fulfillments.fulfills_equipment_id`)

Exercises with no equipment requirements always show.

## Frontend

- Locality list, detail, and edit pages under `/compendium/localities/`
- Locality detail page manages equipment availability: add equipment via search, toggle availability on/off, remove
- Equipment list page links to localities via a button next to "New"
- Exercise list has a locality dropdown filter above the table
