// Package handlers implements the HTTP handlers for user workout management.
//
// # Overview
//
// A workout is a user-owned template that organizes exercises into sections.
// The hierarchy is: Workout → Sections → Section Items.
// Each section item has a type discriminator: "exercise" items reference an
// exercise scheme, while "exercise_group" items reference an exercise group.
//
// All endpoints are scoped to the authenticated user — there is no public
// visibility for workouts.
//
// # Workouts
//
// [ListWorkouts] returns only workouts owned by the current user.
// [CreateWorkout] creates a new empty workout — add sections and items
// via the section and section-item endpoints below.
// [GetWorkout] returns the full workout tree (sections and their items).
// [UpdateWorkout] and [DeleteWorkout] require ownership.
//
// # Sections
//
// A section belongs to a workout and groups items (e.g. "Warm-up",
// "Main Sets", "Superset A"). Sections are ordered by position.
// [CreateWorkoutSection] requires a valid workoutId owned by the current user.
// [ListWorkoutSections] can be filtered by workoutId.
//
// # Section Items
//
// A section item links either an exercise scheme or an exercise group to a
// section. For "exercise" type items, the scheme defines which exercise to
// perform and how (sets, reps, rest). For "exercise_group" type items, the
// user selects a specific exercise when starting the workout.
// [CreateWorkoutSectionItem] requires a valid workoutSectionId whose parent
// workout is owned by the current user, plus a valid exerciseSchemeId or
// exerciseGroupId depending on the item type.
//
// # Permissions
//
// Workouts do not have a dedicated /permissions endpoint. Instead, all
// endpoints enforce ownership directly:
//
//   - List endpoints filter by the current user's ID
//   - All other endpoints check entity ownership and return 403 for non-owners
package handlers
