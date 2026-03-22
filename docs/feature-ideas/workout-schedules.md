# Workout Schedules

## Summary

Allow users to attach a **schedule** to their workouts, defining how often they intend to train. When a workout is "due," a red exclamation mark (!) appears next to its name — both in the top navigation "my workouts" link and on the specific workout line in the my workouts list.

## Schedule Types

### Fixed Day (e.g., "Mondays")

- Workout is due on the specified weekday(s).
- Logic: check if today matches any of the configured days.

### Frequency per Week (e.g., "3x/week")

- Looks at a rolling 7-day window.
- Workout is due if the number of workout-log entries referencing that workout within the window is below the specified count.

### Every X Days (e.g., "every 3 days")

- Checks the most recent workout-log entry for that workout.
- Due if the gap since the last log meets or exceeds X days.

## UI

- **Due indicator**: red `(!)` before the workout name.
- Shown in two places:
  1. The "my workouts" navigation item (if any workout is due).
  2. The individual workout line in the my workouts list.

## Considerations

- Can a workout have multiple schedule rules combined (e.g., "Mondays and Thursdays" = fixed days, or is that 2x/week)?
- What happens if a user has no logs yet — is a new scheduled workout immediately due?
- Should there be a "rest day" override or snooze?
- Calendar/week view showing upcoming scheduled workouts?
