// Package handlers implements the HTTP handlers for user workout management.
//
// # Overview
//
// A workout is a user-owned template that organizes exercises into sections.
// The hierarchy is: Workout → Sections → Section Exercises.
// Each section exercise references an exercise scheme from
// [gesitr/internal/exercise/handlers] (which defines sets, reps, and rest).
//
// All endpoints are scoped to the authenticated user — there is no public
// visibility for workouts.
//
// # Workouts
//
// [ListWorkouts] returns only workouts owned by the current user.
// [CreateWorkout] creates a new empty workout — add sections and exercises
// via the section and section-exercise endpoints below.
// [GetWorkout] returns the full workout tree (sections and their exercises).
// [UpdateWorkout] and [DeleteWorkout] require ownership.
//
// # Sections
//
// A section belongs to a workout and groups exercises (e.g. "Warm-up",
// "Main Sets", "Superset A"). Sections are ordered by position.
// [CreateWorkoutSection] requires a valid workoutId owned by the current user.
// [ListWorkoutSections] can be filtered by workoutId.
//
// # Section Exercises
//
// A section exercise links an exercise scheme to a section. The exercise
// scheme is created via [gesitr/internal/exercise/handlers.CreateExerciseScheme]
// and defines which exercise to perform and how (sets, reps, rest).
// [CreateWorkoutSectionExercise] requires a valid workoutSectionId whose
// parent workout is owned by the current user, and a valid exerciseSchemeId.
//
// # Permissions
//
// Workouts do not have a dedicated /permissions endpoint. Instead, all
// endpoints enforce ownership directly:
//
//   - List endpoints filter by the current user's ID
//   - All other endpoints check entity ownership and return 403 for non-owners
package handlers
