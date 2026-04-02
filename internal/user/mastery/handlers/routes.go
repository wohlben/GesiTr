package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "ListExerciseMastery",
		Method:      http.MethodGet,
		Path:        "/user/mastery",
		Tags:        []string{"mastery"},
		Summary:     "List mastery for all exercises the user has logged",
	}, ListExerciseMastery)

	huma.Register(api, huma.Operation{
		OperationID: "GetExerciseMastery",
		Method:      http.MethodGet,
		Path:        "/user/mastery/{exerciseId}",
		Tags:        []string{"mastery"},
		Summary:     "Get mastery for a specific exercise",
	}, GetExerciseMastery)

	huma.Register(api, huma.Operation{
		OperationID: "ListEquipmentMastery",
		Method:      http.MethodGet,
		Path:        "/user/equipment-mastery",
		Tags:        []string{"mastery"},
		Summary:     "List mastery for all equipment the user has used",
	}, ListEquipmentMastery)

	huma.Register(api, huma.Operation{
		OperationID: "GetEquipmentMastery",
		Method:      http.MethodGet,
		Path:        "/user/equipment-mastery/{equipmentId}",
		Tags:        []string{"mastery"},
		Summary:     "Get mastery for a specific equipment item",
	}, GetEquipmentMastery)
}
