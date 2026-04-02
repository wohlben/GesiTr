package handlers

import (
	"encoding/json"
	"fmt"

	"gesitr/internal/workout/models"
)

// Creating a section requires a workout. The workout must be created first
// via [CreateWorkout].
func ExampleCreateWorkoutSection() {
	setupExampleDB()
	r := newRouter()

	// Create a workout first.
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)

	// Add a section to it.
	w := doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1,
		"type": "main",
		"label": "Compound Lifts",
		"position": 0,
		"restBetweenExercises": 120
	}`)

	var section models.WorkoutSection
	json.Unmarshal(w.Body.Bytes(), &section)
	fmt.Println(w.Code)
	fmt.Println(section.Type)
	fmt.Println(*section.Label)
	fmt.Println(section.WorkoutID)
	// Output:
	// 201
	// main
	// Compound Lifts
	// 1
}

// Non-owner cannot add a section to another user's workout.
func ExampleCreateWorkoutSection_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)

	w := doRawAs(r, "POST", "/api/workout-sections", `{
		"workoutId": 1,
		"type": "main",
		"position": 0
	}`, "bob")
	fmt.Println(w.Code)
	// Output: 403
}

// ListWorkoutSections returns sections owned by the current user. Filter by
// workoutId to get sections for a specific workout.
func ExampleListWorkoutSections_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "label": "Compound", "position": 0
	}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "supplementary", "label": "Isolation", "position": 1
	}`)

	w := doRaw(r, "GET", "/api/workout-sections?workoutId=1", "")

	var sections []models.WorkoutSection
	json.Unmarshal(w.Body.Bytes(), &sections)
	fmt.Println(w.Code)
	fmt.Println(len(sections))
	fmt.Println(*sections[0].Label)
	fmt.Println(*sections[1].Label)
	// Output:
	// 200
	// 2
	// Compound
	// Isolation
}

// Non-owner sees an empty list for another user's sections.
func ExampleListWorkoutSections_nonOwner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)

	w := doRawAs(r, "GET", "/api/workout-sections?workoutId=1", "", "bob")

	var sections []models.WorkoutSection
	json.Unmarshal(w.Body.Bytes(), &sections)
	fmt.Println(w.Code)
	fmt.Println(len(sections))
	// Output:
	// 200
	// 0
}

// Owner can retrieve a single section with its exercises.
func ExampleGetWorkoutSection_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "label": "Compound", "position": 0
	}`)

	w := doJSON(r, "GET", "/api/workout-sections/1", nil)

	var section models.WorkoutSection
	json.Unmarshal(w.Body.Bytes(), &section)
	fmt.Println(w.Code)
	fmt.Println(section.Type)
	// Output:
	// 200
	// main
}

// Non-owner cannot access another user's section.
func ExampleGetWorkoutSection_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)

	w := doRawAs(r, "GET", "/api/workout-sections/1", "", "bob")
	fmt.Println(w.Code)
	// Output: 403
}

// Owner can delete a section from their workout.
func ExampleDeleteWorkoutSection_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)

	w := doJSON(r, "DELETE", "/api/workout-sections/1", nil)
	fmt.Println(w.Code)

	wg := doJSON(r, "GET", "/api/workout-sections/1", nil)
	fmt.Println(wg.Code)
	// Output:
	// 204
	// 404
}

// Non-owner cannot delete another user's section.
func ExampleDeleteWorkoutSection_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)

	w := doRawAs(r, "DELETE", "/api/workout-sections/1", "", "bob")
	fmt.Println(w.Code)
	// Output: 403
}
