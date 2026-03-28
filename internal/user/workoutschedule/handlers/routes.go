package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers all workout schedule endpoints on the huma API.
func RegisterRoutes(api huma.API) {
	// --- Schedules ---

	huma.Register(api, huma.Operation{
		OperationID: "ListWorkoutSchedules",
		Method:      http.MethodGet,
		Path:        "/user/workout-schedules",
		Tags:        []string{"workout-schedules"},
		Summary:     "List workout schedules",
	}, ListWorkoutSchedules)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateWorkoutSchedule",
		Method:        http.MethodPost,
		Path:          "/user/workout-schedules",
		Tags:          []string{"workout-schedules"},
		Summary:       "Create workout schedule",
		DefaultStatus: 201,
	}, CreateWorkoutSchedule)

	huma.Register(api, huma.Operation{
		OperationID: "GetWorkoutSchedule",
		Method:      http.MethodGet,
		Path:        "/user/workout-schedules/{id}",
		Tags:        []string{"workout-schedules"},
		Summary:     "Get workout schedule",
	}, GetWorkoutSchedule)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateWorkoutSchedule",
		Method:      http.MethodPatch,
		Path:        "/user/workout-schedules/{id}",
		Tags:        []string{"workout-schedules"},
		Summary:     "Update workout schedule",
	}, UpdateWorkoutSchedule)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteWorkoutSchedule",
		Method:      http.MethodDelete,
		Path:        "/user/workout-schedules/{id}",
		Tags:        []string{"workout-schedules"},
		Summary:     "Delete workout schedule",
	}, DeleteWorkoutSchedule)

	// --- Periods ---

	huma.Register(api, huma.Operation{
		OperationID: "ListSchedulePeriods",
		Method:      http.MethodGet,
		Path:        "/user/schedule-periods",
		Tags:        []string{"workout-schedules"},
		Summary:     "List schedule periods",
	}, ListSchedulePeriods)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateSchedulePeriod",
		Method:        http.MethodPost,
		Path:          "/user/schedule-periods",
		Tags:          []string{"workout-schedules"},
		Summary:       "Create schedule period",
		DefaultStatus: 201,
	}, CreateSchedulePeriod)

	// --- Commitments ---

	huma.Register(api, huma.Operation{
		OperationID: "ListScheduleCommitments",
		Method:      http.MethodGet,
		Path:        "/user/schedule-commitments",
		Tags:        []string{"workout-schedules"},
		Summary:     "List schedule commitments",
	}, ListScheduleCommitments)

	huma.Register(api, huma.Operation{
		OperationID:   "CreateScheduleCommitment",
		Method:        http.MethodPost,
		Path:          "/user/schedule-commitments",
		Tags:          []string{"workout-schedules"},
		Summary:       "Create schedule commitment",
		DefaultStatus: 201,
	}, CreateScheduleCommitment)

	huma.Register(api, huma.Operation{
		OperationID: "DeleteScheduleCommitment",
		Method:      http.MethodDelete,
		Path:        "/user/schedule-commitments/{id}",
		Tags:        []string{"workout-schedules"},
		Summary:     "Delete schedule commitment",
	}, DeleteScheduleCommitment)
}
