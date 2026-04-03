# Exercise Multiple Names

## Summary

Exercises should have 1:n names instead of a single primary name. There's already a concept of alternative names — this formalizes it into a proper relation where all names are equal peers (no primary).

## Model Changes

- Remove the single `name` field from exercises
- Introduce a separate `exercise_names` relation (exercise_id, name)
- All names are equal — no primary/secondary distinction

## UI Impact

- **List views** show the same exercise once per name — an exercise with 3 names appears 3 times in search results and compendium lists
- Search should match against all names
- Exercise detail view shows all names

## Considerations

- How does this affect exercise logs and workout references that currently store/display the exercise name?
- Should there be a "display name" preference per user, or is duplication in lists sufficient?
- Migration path for existing exercises: current `name` becomes the first entry in the names relation
