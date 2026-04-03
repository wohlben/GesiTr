package handlers

import (
	"encoding/json"
	"fmt"

	"gesitr/internal/compendium/exercise/models"
	"gesitr/internal/shared"
)

// Owner sees their own exercises (both public and private) in the list.
func ExampleListExercises_owner() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Secret Move"],
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise"
	}`)

	// Create a public exercise.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Push-up"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"public": true
	}`)

	w := doJSON(r, "GET", "/api/exercises", nil)

	var page paginatedJSON
	json.Unmarshal(w.Body.Bytes(), &page)
	fmt.Println(w.Code)
	fmt.Println(page.Total)
	// Output:
	// 200
	// 2
}

// Non-owner sees only public exercises in the list.
func ExampleListExercises_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	// Create a public exercise.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Push-up"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"public": true
	}`)

	// Another user sees the public exercise.
	w := doRawAs(r, "GET", "/api/exercises", "", "other")

	var page paginatedJSON
	json.Unmarshal(w.Body.Bytes(), &page)
	var exercises []models.Exercise
	json.Unmarshal(page.Items, &exercises)
	fmt.Println(w.Code)
	fmt.Println(page.Total)
	fmt.Println(exercises[0].Names[0].Name)
	// Output:
	// 200
	// 1
	// Push-up
}

// Non-owner does not see private exercises in the list.
func ExampleListExercises_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Secret Move"],
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise"
	}`)

	// Another user sees an empty list.
	w := doRawAs(r, "GET", "/api/exercises", "", "other")

	var page paginatedJSON
	json.Unmarshal(w.Body.Bytes(), &page)
	fmt.Println(w.Code)
	fmt.Println(page.Total)
	// Output:
	// 200
	// 0
}

// Creating an exercise without equipment. The owner is set from the
// authenticated user, not from the request body.
func ExampleCreateExercise_simple() {
	setupExampleDB()
	r := newRouter()

	w := doRaw(r, "POST", "/api/exercises", `{
		"names": ["Push-up"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up"
	}`)

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println(w.Code)
	fmt.Println(exercise.Names[0].Name)
	fmt.Println(exercise.Owner)
	fmt.Println(exercise.EquipmentIDs)
	// Output:
	// 201
	// Push-up
	// testuser
	// []
}

// Creating an exercise that requires equipment. The equipment must be created
// first via the equipment API (see [gesitr/internal/compendium/equipment/handlers]).
func ExampleCreateExercise_withEquipment() {
	setupExampleDB()
	r := newExampleRouter()

	// First, create the equipment via the equipment API.
	doRaw(r, "POST", "/api/equipment", `{
		"name": "barbell",
		"displayName": "Barbell",
		"description": "Standard barbell",
		"category": "free_weights",
		"templateId": "barbell"
	}`)

	// Create an exercise referencing the equipment by ID.
	w := doRaw(r, "POST", "/api/exercises", `{
		"names": ["Bench Press"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell bench press",
		"equipmentIds": [1]
	}`)

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println(w.Code)
	fmt.Println(exercise.Names[0].Name)
	fmt.Println(exercise.EquipmentIDs)
	// Output:
	// 201
	// Bench Press
	// [1]
}

// Owner can retrieve their own exercise with full details.
func ExampleGetExercise_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Bench Press"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell bench press"
	}`)

	w := doJSON(r, "GET", "/api/exercises/1", nil)

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println(w.Code)
	fmt.Println(exercise.Names[0].Name)
	fmt.Println(exercise.Owner)
	// Output:
	// 200
	// Bench Press
	// testuser
}

// Non-owner can retrieve a public exercise.
func ExampleGetExercise_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Squat"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat",
		"public": true
	}`)

	w := doRawAs(r, "GET", "/api/exercises/1", "", "other")

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println(w.Code)
	fmt.Println(exercise.Names[0].Name)
	// Output:
	// 200
	// Squat
}

// Non-owner cannot access a private exercise — returns 403 Forbidden.
func ExampleGetExercise_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Secret Move"],
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise"
	}`)

	w := doRawAs(r, "GET", "/api/exercises/1", "", "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Owner can update their exercise. Returns the updated exercise.
func ExampleUpdateExercise_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Squat"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat"
	}`)

	w := doRaw(r, "PUT", "/api/exercises/1", `{
		"names": ["Back Squat"],
		"type": "STRENGTH",
		"technicalDifficulty": "intermediate",
		"bodyWeightScaling": 0.5,
		"description": "Barbell back squat",
		"force": ["PUSH"],
		"primaryMuscles": ["QUADS"]
	}`)

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println(w.Code)
	fmt.Println(exercise.Names[0].Name)
	fmt.Println(exercise.Owner)
	// Output:
	// 200
	// Back Squat
	// testuser
}

// Each update bumps the version and creates a history snapshot. Version 0
// preserves the original state, version 1 the updated state.
func ExampleUpdateExercise_versioning() {
	setupExampleDB()
	r := newRouter()

	// Create an exercise (version 0).
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Squat"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat"
	}`)

	// Update it — version bumps to 1.
	w := doRaw(r, "PUT", "/api/exercises/1", `{
		"names": ["Back Squat"],
		"type": "STRENGTH",
		"technicalDifficulty": "intermediate",
		"bodyWeightScaling": 0.5,
		"description": "Barbell back squat",
		"force": ["PUSH"],
		"primaryMuscles": ["QUADS"]
	}`)

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println("current version:", exercise.Version)

	// Version 0 preserves the original name.
	wv0 := doJSON(r, "GET", "/api/exercises/1/versions/0", nil)
	var v0 shared.VersionEntry
	json.Unmarshal(wv0.Body.Bytes(), &v0)
	var snap0 models.Exercise
	json.Unmarshal(v0.Snapshot, &snap0)
	fmt.Println("v0 name:", snap0.Names[0].Name)

	// Version 1 has the updated name.
	wv1 := doJSON(r, "GET", "/api/exercises/1/versions/1", nil)
	var v1 shared.VersionEntry
	json.Unmarshal(wv1.Body.Bytes(), &v1)
	var snap1 models.Exercise
	json.Unmarshal(v1.Snapshot, &snap1)
	fmt.Println("v1 name:", snap1.Names[0].Name)

	// Output:
	// current version: 1
	// v0 name: Squat
	// v1 name: Back Squat
}

// Non-owner cannot update a public exercise — returns 403.
func ExampleUpdateExercise_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Push-up"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"public": true
	}`)

	w := doRawAs(r, "PUT", "/api/exercises/1", `{
		"names": ["Modified Push-up"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Tampered"
	}`, "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Non-owner cannot update a private exercise — returns 403.
func ExampleUpdateExercise_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Secret Move"],
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise"
	}`)

	w := doRawAs(r, "PUT", "/api/exercises/1", `{
		"names": ["Stolen Move"],
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "Tampered"
	}`, "other")
	fmt.Println(w.Code)
	// Output: 403
}
