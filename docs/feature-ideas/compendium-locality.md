# Compendium Locality

## Summary

Curated, per-user lists of equipment available at specific training locations. A user can define multiple localities (e.g. "Home Gym", "Commercial Gym") and list which equipment is present at each. The primary goal is to filter exercises down to what's actually available at the user's current location.

This is conceptually distinct from **fulfillments** — a fulfillment says "an adjustable bench can stand in for a flat bench and an incline bench" (equipment substitution). A locality says "when I'm at this location, I have these pieces of equipment available."

## Data Model

### Locality

User-scoped entity representing a training location.

| Column       | Type      | Notes                            |
|-------------|-----------|----------------------------------|
| id          | uint      | PK                               |
| user_id     | string    | owner                            |
| name        | string    | e.g. "Home Gym", "Commercial Gym"|
| created_at  | timestamp |                                  |
| updated_at  | timestamp |                                  |

### Locality Availability

Join table between a locality and equipment, with an availability toggle.

| Column        | Type      | Notes                                          |
|--------------|-----------|------------------------------------------------|
| id           | uint      | PK                                             |
| locality_id  | uint      | FK → Locality                                  |
| equipment_id | uint      | FK → Equipment                                 |
| available    | bool      | default `true` — toggle off for e.g. maintenance |
| created_at   | timestamp | tracks when equipment was added                |
| updated_at   | timestamp | tracks when availability was last toggled      |

A user might have overlapping equipment across localities — e.g. dumbbells at both the home gym and commercial gym, but a cable machine only at the commercial gym.

## Behavior

- One user → many localities, each locality → many equipment entries via availability.
- The `available` flag defaults to `true`. Toggling to `false` marks equipment as temporarily unavailable (e.g. machine under maintenance) without deleting the relationship. This is easier than deleting and recreating the connection.
- Selecting a locality filters the exercise list to exercises whose required equipment is available at that locality.

## Considerations

- How does this interact with fulfillments during filtering? If an exercise requires a flat bench and the user doesn't have one, but has an adjustable bench that fulfills it, should the exercise still show? -> yes
- Should localities support a "default" or "active" flag so the user doesn't have to pick one every session? -> unecessary
- Could localities eventually be shared (e.g. a gym's equipment list maintained collaboratively)? -> out of scope, but they should again have the 'public' toggle as every other compendium entity - so maintained by one person, but visible by many. we will flesh taht feature out later on.
- Should the filter be strict (all required equipment must be available) or lenient (at least one variant available via fulfillments)? -> fulfillments usecase is exactly that, its not lenient because thats literally what its for (so yes, it should show things which are fulfilled by the fullfilment)
