# Calendar Commitment Redesign

## Summary

Redesign the calendar view to show commitments as horizontal bars instead of name bubbles, de-emphasize the active schedule, and add a quick-move action on the current date.

## Changes

### Commitments as bars
Replace the current bubble-with-name display for commitments with horizontal bars. This gives a clearer visual representation of planned workouts across the calendar.

### De-emphasize active schedule
When an active workout schedule is present, reduce its visual prominence on the calendar so commitments and actual logs stand out more.

### Quick-move button on current date
Show a button on today's date, styled like a commitment bar, that prompts the user to move the nearest committed or planned workout-log/commitment. The system automatically picks the nearest one (past-due or upcoming) and offers to reschedule it to today.

## Considerations

- What defines "nearest" — closest by date regardless of direction, or prefer overdue commitments first?
- Should the move button appear if there's already a commitment on today?
- How does moving interact with group cycle commitments — can those be moved independently?
