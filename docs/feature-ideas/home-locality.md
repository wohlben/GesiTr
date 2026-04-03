# Home Locality

## Summary

A "Home" locality is always visible in the UI as a quick-access entry point, even before the user has created any locality. Clicking it either navigates to the user's existing home locality or transparently creates one and then navigates to it.

## Identifying the Home Locality

Since localities can be renamed, the name cannot be used as an identifier. Instead, the home locality is defined as:

- Any locality where the user is the **owner** AND it is **private** (non-public).
- If at least one such locality exists, the first one (e.g. oldest by creation date) is treated as "home", regardless of its current name.
- If none exists, clicking "Home" creates a new private locality named "Home" and navigates to it.

This means a user who renames their home locality to "Garage Gym" still sees it surfaced as the home entry point.

## UI Behavior

The "Home" button appears in the equipment list page (next to the existing "Localities" button). On click:

1. Check if the user has any owned, private locality.
2. **If yes** — navigate to that locality's detail page.
3. **If no** — create a new locality with name "Home" and `public: false`, then navigate to its detail page.

The creation should feel seamless — no intermediate form or confirmation.

## Considerations

- What if the user has multiple private localities? Use the oldest (first created) as home? Or should there be an explicit `isHome` flag on the model?
- Should the home locality be prevented from being made public, or is that the user's choice (at which point it stops being "home")?
- Should the home entry point show the actual name of the locality (e.g. "Garage Gym") or always display "Home"?
