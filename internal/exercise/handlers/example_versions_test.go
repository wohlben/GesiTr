package handlers

import (
	"encoding/json"
	"fmt"

	"gesitr/internal/exercise/models"
	"gesitr/internal/shared"
)

// Owner can retrieve version history. Creating an exercise produces version 0,
// updating it produces version 1. Each version contains a snapshot of the
// exercise at that point in time.
func ExampleGetExerciseVersion_owner() {
	setupExampleDB()
	r := newRouter()

	// Create an exercise (version 0).
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat"
	}`)

	// Verify version 0 exists.
	w0 := doJSON(r, "GET", "/api/exercises/templates/squat/versions/0", nil)
	var v0 shared.VersionEntry
	json.Unmarshal(w0.Body.Bytes(), &v0)
	fmt.Println("v0:", w0.Code, "version", v0.Version)

	// Parse the snapshot to check the exercise name at version 0.
	var snap0 models.Exercise
	json.Unmarshal(v0.Snapshot, &snap0)
	fmt.Println("v0 name:", snap0.Name)

	// Update the exercise (creates version 1).
	doRaw(r, "PUT", "/api/exercises/1", `{
		"name": "Back Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "intermediate",
		"bodyWeightScaling": 0.5,
		"description": "Barbell back squat",
		"force": ["PUSH"],
		"primaryMuscles": ["QUADS"]
	}`)

	// Verify version 1 exists with the updated name.
	w1 := doJSON(r, "GET", "/api/exercises/templates/squat/versions/1", nil)
	var v1 shared.VersionEntry
	json.Unmarshal(w1.Body.Bytes(), &v1)
	fmt.Println("v1:", w1.Code, "version", v1.Version)

	var snap1 models.Exercise
	json.Unmarshal(v1.Snapshot, &snap1)
	fmt.Println("v1 name:", snap1.Name)

	// Output:
	// v0: 200 version 0
	// v0 name: Squat
	// v1: 200 version 1
	// v1 name: Back Squat
}

// Non-owner can retrieve version history for a public exercise.
func ExampleGetExerciseVersion_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	// Create a public exercise (version 0).
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"public": true
	}`)

	// Update it (version 1).
	doRaw(r, "PUT", "/api/exercises/1", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "intermediate",
		"bodyWeightScaling": 1.0,
		"description": "Strict push-up",
		"public": true,
		"force": ["PUSH"],
		"primaryMuscles": ["CHEST"]
	}`)

	// Another user can access both versions.
	w0 := doRawAs(r, "GET", "/api/exercises/templates/push-up/versions/0", "", "other")
	var v0 shared.VersionEntry
	json.Unmarshal(w0.Body.Bytes(), &v0)
	fmt.Println("v0:", w0.Code, "version", v0.Version)

	w1 := doRawAs(r, "GET", "/api/exercises/templates/push-up/versions/1", "", "other")
	var v1 shared.VersionEntry
	json.Unmarshal(w1.Body.Bytes(), &v1)
	fmt.Println("v1:", w1.Code, "version", v1.Version)

	var snap1 models.Exercise
	json.Unmarshal(v1.Snapshot, &snap1)
	fmt.Println("v1 description:", snap1.Description)

	// Output:
	// v0: 200 version 0
	// v1: 200 version 1
	// v1 description: Strict push-up
}

// Non-owner cannot retrieve version history for a private exercise.
func ExampleGetExerciseVersion_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Secret Move",
		"templateId": "secret-move",
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise"
	}`)

	w := doRawAs(r, "GET", "/api/exercises/templates/secret-move/versions/0", "", "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Version history survives exercise deletion. The exercise is hard-deleted
// but the history snapshots remain accessible via the templateId.
func ExampleGetExerciseVersion_afterDelete() {
	setupExampleDB()
	r := newRouter()

	// Create an exercise (version 0).
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Overhead Press",
		"templateId": "ohp",
		"type": "STRENGTH",
		"technicalDifficulty": "intermediate",
		"bodyWeightScaling": 0,
		"description": "Standing overhead press"
	}`)

	// Verify it exists.
	w := doJSON(r, "GET", "/api/exercises/1", nil)
	fmt.Println("before delete:", w.Code)

	// Delete the exercise.
	wd := doJSON(r, "DELETE", "/api/exercises/1", nil)
	fmt.Println("delete:", wd.Code)

	// Exercise is gone.
	wg := doJSON(r, "GET", "/api/exercises/1", nil)
	fmt.Println("after delete:", wg.Code)

	// Version history is still accessible.
	wv := doJSON(r, "GET", "/api/exercises/templates/ohp/versions/0", nil)
	var v0 shared.VersionEntry
	json.Unmarshal(wv.Body.Bytes(), &v0)
	fmt.Println("version:", wv.Code, "version", v0.Version)

	var snap models.Exercise
	json.Unmarshal(v0.Snapshot, &snap)
	fmt.Println("snapshot name:", snap.Name)

	// Output:
	// before delete: 200
	// delete: 204
	// after delete: 404
	// version: 200 version 0
	// snapshot name: Overhead Press
}
