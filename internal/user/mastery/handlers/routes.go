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
}
