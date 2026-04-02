// Package handlers implements the HTTP handlers for the equipment API.
//
// # Overview
//
// Equipment represents physical items used in exercises (barbells, dumbbells,
// benches, etc.). Like exercises, equipment can be public (part of the
// compendium) or private (visible only to the owner).
//
// Exercises reference equipment via their equipmentIds field — see
// [gesitr/internal/compendium/exercise/handlers.CreateExercise]. This allows users to
// track which equipment is needed for an exercise.
//
// # Version History
//
// Equipment maintains version history — each [UpdateEquipment] call creates
// a snapshot. Use [ListEquipmentVersions] and [GetEquipmentVersion] to browse
// previous versions.
//
// # Permissions
//
// [GetEquipmentPermissions] delegates to [gesitr/internal/shared.ResolvePermissions].
// Owner gets READ, MODIFY, DELETE. Non-owner gets READ on public equipment,
// empty permissions on private equipment. Mutating operations (PUT, DELETE)
// enforce owner checks directly with 403 Forbidden.
package handlers
