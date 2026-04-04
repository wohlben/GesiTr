// Package handlers implements the HTTP handlers for exercises.
//
// # Overview
//
// Exercises are the core building blocks of workouts. They can be public
// (visible to all users as part of the compendium) or private (visible only
// to their owner). Exercises can reference equipment via equipmentIds — see
// [gesitr/internal/compendium/equipment/handlers].
//
// Exercise schemes (user-specific configurations of exercises) are in
// [gesitr/internal/user/exercisescheme/handlers].
//
// # Version History
//
// Exercises maintain version history — each [UpdateExercise] call creates a
// snapshot. Use [ListExerciseVersions] and [GetExerciseVersion] to browse
// previous versions.
//
// # Permissions
//
// [GetExercisePermissions] delegates to [gesitr/internal/shared.ResolvePermissions].
// Owner gets READ, MODIFY, DELETE. Non-owner gets READ on public exercises,
// empty permissions on private exercises. Mutating operations (PUT, DELETE)
// enforce owner checks directly with 403 Forbidden.
package handlers
