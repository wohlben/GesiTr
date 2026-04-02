package handlers

import (
	"encoding/json"
	"fmt"

	"gesitr/internal/compendium/workout/models"
)

// ListWorkouts returns only the current user's workouts.
func ExampleListWorkouts_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workouts", `{"name": "Pull Day"}`)

	w := doJSON(r, "GET", "/api/workouts", nil)

	var page struct {
		Items []models.Workout `json:"items"`
	}
	json.Unmarshal(w.Body.Bytes(), &page)
	fmt.Println(w.Code)
	fmt.Println(len(page.Items))
	fmt.Println(page.Items[0].Name)
	fmt.Println(page.Items[1].Name)
	// Output:
	// 200
	// 2
	// Push Day
	// Pull Day
}

// Non-owner sees an empty list — workouts are always private.
func ExampleListWorkouts_nonOwner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)

	w := doRawAs(r, "GET", "/api/workouts", "", "bob")

	var page struct {
		Items []models.Workout `json:"items"`
	}
	json.Unmarshal(w.Body.Bytes(), &page)
	fmt.Println(w.Code)
	fmt.Println(len(page.Items))
	// Output:
	// 200
	// 0
}

// CreateWorkout creates an empty workout. The owner is set from the
// authenticated user. Add sections via [CreateWorkoutSection].
func ExampleCreateWorkout() {
	setupExampleDB()
	r := newRouter()

	w := doRaw(r, "POST", "/api/workouts", `{
		"name": "Push Day",
		"notes": "Chest, shoulders, triceps"
	}`)

	var workout models.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)
	fmt.Println(w.Code)
	fmt.Println(workout.Name)
	fmt.Println(workout.Owner)
	fmt.Println(*workout.Notes)
	// Output:
	// 201
	// Push Day
	// alice
	// Chest, shoulders, triceps
}

// Owner can retrieve their own workout with full section and exercise tree.
func ExampleGetWorkout_ownerAccess() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)

	w := doJSON(r, "GET", "/api/workouts/1", nil)

	var workout models.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)
	fmt.Println(w.Code)
	fmt.Println(workout.Name)
	fmt.Println(workout.Owner)
	// Output:
	// 200
	// Push Day
	// alice
}

// Non-owner is denied access to another user's workout with 403 Forbidden.
func ExampleGetWorkout_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)

	w := doRawAs(r, "GET", "/api/workouts/1", "", "bob")
	fmt.Println(w.Code)
	// Output: 403
}

// Owner can update workout metadata (name, notes).
func ExampleUpdateWorkout_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)

	w := doRaw(r, "PUT", "/api/workouts/1", `{
		"name": "Upper Body Push",
		"notes": "Updated focus"
	}`)

	var workout models.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)
	fmt.Println(w.Code)
	fmt.Println(workout.Name)
	// Output:
	// 200
	// Upper Body Push
}

// Non-owner cannot update another user's workout.
func ExampleUpdateWorkout_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)

	w := doRawAs(r, "PUT", "/api/workouts/1", `{"name": "Stolen"}`, "bob")
	fmt.Println(w.Code)
	// Output: 403
}

// Owner can delete their workout.
func ExampleDeleteWorkout_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)

	w := doJSON(r, "DELETE", "/api/workouts/1", nil)
	fmt.Println(w.Code)

	// Workout is gone.
	wg := doJSON(r, "GET", "/api/workouts/1", nil)
	fmt.Println(wg.Code)
	// Output:
	// 204
	// 404
}

// Non-owner cannot delete another user's workout.
func ExampleDeleteWorkout_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)

	w := doRawAs(r, "DELETE", "/api/workouts/1", "", "bob")
	fmt.Println(w.Code)
	// Output: 403
}
