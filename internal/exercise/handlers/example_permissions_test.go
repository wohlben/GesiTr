package handlers

import (
	"encoding/json"
	"fmt"

	"gesitr/internal/shared"
)

// Owner gets full permissions (READ, MODIFY, DELETE) on their own exercise.
func ExampleGetExercisePermissions_owner() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise owned by testuser.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Bench Press",
		"templateId": "bench-press",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell bench press",
		"version": 0
	}`)

	// Query permissions as the owner.
	w := doJSON(r, "GET", "/api/exercises/1/permissions", nil)

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(resp.Permissions)
	// Output: [READ MODIFY DELETE]
}

// Non-owner can read a public exercise but cannot modify or delete it.
func ExampleGetExercisePermissions_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	// Create a public exercise owned by testuser.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat",
		"version": 0,
		"public": true
	}`)

	// Query permissions as a different user.
	w := doRawAs(r, "GET", "/api/exercises/1/permissions", "", "other")

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(resp.Permissions)
	// Output: [READ]
}

// Non-owner has no permissions on a private exercise.
func ExampleGetExercisePermissions_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise owned by testuser.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Deadlift",
		"templateId": "deadlift",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell deadlift",
		"version": 0
	}`)

	// Query permissions as a different user.
	w := doRawAs(r, "GET", "/api/exercises/1/permissions", "", "other")

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(resp.Permissions)
	// Output: []
}
