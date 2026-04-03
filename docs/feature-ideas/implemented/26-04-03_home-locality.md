# Home Locality

## Summary

A "Home" button on the equipment list page provides quick access to the user's home locality. If no home locality exists, clicking the button transparently creates one named "Home" and navigates to it.

## Identifying the Home Locality

The home locality is the user's first owned, private (non-public) locality by creation date. Name is not used as identifier — renaming the locality does not break the association.

The frontend queries `GET /localities?owner=me&public=false&limit=1` to find it.

## UI

The "Home" button sits in the equipment list page actions bar alongside "Localities" and "New". On click:

1. If a private locality exists → navigate to its detail page.
2. If none exists → create a locality with name "Home" and `public: false`, then navigate to it.

## Backend

Added `public` query parameter to `GET /localities` — supports `true` (public only) and `false` (private only) filtering.
